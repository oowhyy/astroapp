package tilemap

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"io"
	"runtime"
	"time"

	"github.com/oowhyy/astroapp/pkg/dropbox"
)

type TileMap struct {
	maxDepth int
	tileSize int
	worldW   int
	worldH   int
	Root     *Tile
}

func (tm *TileMap) Depth() int {
	return tm.maxDepth
}

func (tm *TileMap) TileSize() int {
	return tm.tileSize
}

func (tm *TileMap) WorldSize() (int, int) {
	return tm.worldW, tm.worldH
}

func NewTileMapFromDropboxZip(client *dropbox.Client, zipPath string) (*TileMap, error) {
	tic := time.Now()
	f, err := client.Files.Download(&dropbox.DownloadInput{
		Path: zipPath,
	})
	fmt.Println("fetched in", time.Since(tic).Seconds())
	if err != nil {
		return nil, fmt.Errorf("db download error: %w", err)
	}
	buff := new(bytes.Buffer)
	size, err := io.Copy(buff, f.Body)
	if err != nil {
		return nil, fmt.Errorf("io copy error: %w", err)
	}
	f.Body.Close()
	br := bytes.NewReader(buff.Bytes())
	reader, err := zip.NewReader(br, size)
	fmt.Println("built reader in ", time.Since(tic).Seconds())
	// reader, err := zip.OpenReader(zipName)
	if err != nil {
		return nil, fmt.Errorf("open zip reader error: %w", err)
	}
	n := 3*len(reader.File) + 1
	maxDepth := 0
	for n > 4 {
		maxDepth++
		n /= 4
	}
	mem := map[string]*ImageReader{}
	for _, f := range reader.File {
		readCloser, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open error: %w", err)
		}
		mem[f.Name] = NewImageReader(readCloser)
	}
	rootImg, ok := mem["0/0/0.png"]
	if !ok {
		return nil, fmt.Errorf("no root image")
	}
	fmt.Println("file to img done in", time.Since(tic).Seconds())
	defer func() {
		fmt.Println("done new tilemap from zip in", time.Since(tic).Seconds())
	}()
	tm := &TileMap{
		tileSize: rootImg.Decoded().Bounds().Dx(),
	}
	fmt.Println(tm.tileSize)
	perfectSize := tm.tileSize << maxDepth
	fmt.Println("perfect size", perfectSize)
	// get resized layers
	tm.maxDepth = maxDepth
	var build func(start *Tile)
	build = func(start *Tile) {
		if start.bounds.Empty() {
			panic("empty bounds")
		}
		name := fmt.Sprintf("%d/%d/%d.png", start.z, start.x, start.y)
		img, ok := mem[name]
		if !ok {
			panic(fmt.Sprintf("file not found: %s", name))
		}
		start.encoded = img
		// PrintMemUsage()
		// no children on last level
		if start.z == maxDepth {
			return
		}
		start.children = make([]*Tile, 4)
		halfW := start.bounds.Dx() / 2
		topX, topY := start.bounds.Min.X, start.bounds.Min.Y
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				rect0 := image.Rect(halfW*j, halfW*i, halfW*(j+1), halfW*(i+1))
				child := &Tile{
					bounds: rect0.Add(image.Point{topX, topY}),
					z:      start.z + 1,
					x:      start.x*2 + j,
					y:      start.y*2 + i,
				}
				build(child)
				start.children[2*i+j] = child
			}
		}
	}
	root := &Tile{
		bounds: image.Rect(0, 0, perfectSize, perfectSize),
		z:      0,
	}
	build(root)
	PrintMemUsage()
	tm.Root = root
	return tm, nil
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
