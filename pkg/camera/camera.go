package camera

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/oowhyy/astroapp/pkg/tilemap"
	"github.com/oowhyy/astroapp/pkg/vector"
)

type Camera struct {
	PerfectSize int
	WorldSize   image.Point
	ScreenSize  image.Point
	MousePan    vector.Vector

	geom ebiten.GeoM
}

func NewCamera(wordlSize image.Point, screenSize image.Point) *Camera {
	c := &Camera{
		ScreenSize: screenSize,
		WorldSize:  wordlSize,
		geom:       ebiten.GeoM{},
	}
	return c
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"Scale: %f",
		c.ZoomFactor(),
	)
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	return c.geom

}

func (c *Camera) GetRect() image.Rectangle {
	mat := c.worldMatrix()
	x1, y1 := mat.Apply(0, 0)
	x2, y2 := mat.Apply(float64(c.ScreenSize.X), float64(c.ScreenSize.Y))
	return image.Rect(int(x1), int(y1), int(x2), int(y2))
}

func (c *Camera) RenderDrawers(screen *ebiten.Image, drawers []Drawer) {
	cameraRect := c.GetRect()
	scale0 := c.ZoomFactor()
	for _, dr := range drawers {
		dBounds := dr.InWorldBounds(c.WorldSize.X, c.WorldSize.Y)
		if !dBounds.Overlaps(cameraRect) {
			continue
		}
		bodyX := dBounds.Min.X + dBounds.Dx()/2
		bodyY := dBounds.Min.Y + dBounds.Dy()/2
		screenX, screenY := c.WorldToScreen(float64(bodyX), float64(bodyY))
		img, scale, tx, ty := dr.SpriteOp()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale*scale0, scale*scale0)
		op.GeoM.Translate(tx*scale0+float64(screenX), ty*scale0+float64(screenY))
		screen.DrawImage(img, op)
	}
}

func (c *Camera) RenderTilemap(screen *ebiten.Image, tm *tilemap.TileMap) {
	camera := c.GetRect()
	// runtime.GC()
	// PrintMemUsage()
	zoomPow, scale0 := getZooms(screen.Bounds().Dx(), camera.Bounds().Dx(), tm.Depth())
	zoomLevel := tm.Depth() - zoomPow
	scalar := float64(int(1) << (zoomPow))
	// snappedX := camera.Min.X / minTileSize * minTileSize
	// snappedY := camera.Min.Y / minTileSize * minTileSize
	// snappedCameraCorner := image.Point{snappedX, snappedY}
	var dfs func(start *tilemap.Tile)
	dfs = func(start *tilemap.Tile) {
		if !start.Bounds().Overlaps(camera) {
			return
		}
		if start.Z() == zoomLevel {
			op := &ebiten.DrawImageOptions{}
			cornerPos := start.Bounds().Sub(camera.Min).Min
			screenx0 := float64(cornerPos.X) / scalar
			screeny0 := float64(cornerPos.Y) / scalar
			op.GeoM.Translate(screenx0, screeny0)
			op.GeoM.Scale(scale0, scale0)
			img := start.Image()
			screen.DrawImage(img, op)
			return
		}
		for _, c := range start.Children() {
			dfs(c)
		}
	}
	dfs(tm.Root)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("current zoom level: %d\nscale0: %f", zoomLevel, scale0), 10, 800)
}

func (c *Camera) CleanUpVisible(screen *ebiten.Image, tm *tilemap.TileMap) {
	camera := c.GetRect()
	// runtime.GC()
	// PrintMemUsage()
	zoomPow, _ := getZooms(screen.Bounds().Dx(), camera.Bounds().Dx(), tm.Depth())
	zoomLevel := tm.Depth() - zoomPow
	// snappedX := camera.Min.X / minTileSize * minTileSize
	// snappedY := camera.Min.Y / minTileSize * minTileSize
	// snappedCameraCorner := image.Point{snappedX, snappedY}
	var dfs func(start *tilemap.Tile)
	dfs = func(start *tilemap.Tile) {
		if !start.Bounds().Overlaps(camera) {
			return
		}
		if start.Z() == zoomLevel {
			start.Clear()
			return
		}
		for _, c := range start.Children() {
			dfs(c)
		}
	}
	dfs(tm.Root)
}

// getZooms returns zoom depth and scalar.  1.0 <= scalar < 2.0 if max depth not reached
func getZooms(screenW, cameraW, maxDepth int) (int, float64) {
	zoomLevels := 0
	// only SCALE UP if needed

	for zoomLevels < maxDepth && screenW < cameraW {
		zoomLevels++
		screenW *= 2
	}
	// fmt.Println(screenW, cameraW, float64(screenW)/float64(cameraW))
	return zoomLevels, float64(screenW) / float64(cameraW)
}

func (c *Camera) Update() {
	// pan
	mousePos := vector.FromInts(ebiten.CursorPosition())
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.MousePan = mousePos
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		diff := vector.Diff(mousePos, c.MousePan)
		diff.Scale(1 / c.ZoomFactor())
		c.geom.Translate(-diff.X, -diff.Y)
		// c.Position.Add(diff)
		c.MousePan = mousePos
	}
	// zoom
	_, dy := ebiten.Wheel()
	// scroll down - zoom out

	scaleFactor := 1.06
	rect := c.GetRect()
	ww := float64(rect.Min.X) + float64(rect.Dx())/2
	wh := float64(rect.Min.Y) + float64(rect.Dy())/2
	// scroll down - zoom out
	if dy < 0 {
		c.geom.Translate(-ww, -wh)
		c.geom.Scale(scaleFactor, scaleFactor)
		c.geom.Translate(ww, wh)
	}
	// scroll up - zoom in
	if dy > 0 {
		c.geom.Translate(-ww, -wh)
		c.geom.Scale(1/scaleFactor, 1/scaleFactor)
		c.geom.Translate(ww, wh)
	}
}

// ZoomFactor is current camera dx divided by screen dx
func (c *Camera) ZoomFactor() float64 {
	return float64(c.ScreenSize.X) / float64(c.GetRect().Dx())
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	// inverseMatrix := c.worldMatrix()
	// if inverseMatrix.IsInvertible() {
	// 	inverseMatrix.Invert()
	// 	return inverseMatrix.Apply(float64(posX), float64(posY))
	// } else {
	// 	// When scaling it can happened that matrix is not invertable
	// 	return math.NaN(), math.NaN()
	// }
	return c.geom.Apply(float64(posX), float64(posY))
}

func (c *Camera) WorldToScreen(posX, posY float64) (int, int) {
	inverse := c.geom
	if inverse.IsInvertible() {
		inverse.Invert()
		x, y := inverse.Apply(posX, posY)
		return int(x), int(y)
	} else {
		return 0, 0
	}
}

func (c *Camera) Reset() {
	c.geom.Reset()
}

func (c *Camera) Translate(tx, ty float64) {
	c.geom.Translate(tx, ty)
}

// func PrintMemUsage() {
// 	var m runtime.MemStats
// 	runtime.ReadMemStats(&m)
// 	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
// 	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
// 	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
// 	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
// 	fmt.Printf("\tNumGC = %v\n", m.NumGC)
// }

// func bToMb(b uint64) uint64 {
// 	return b / 1024 / 1024
// }
