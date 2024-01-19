package tilemap

import (
	"fmt"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)



func NewTileMapEmpty(depth int, tileSize int, worldW, worldH int) *TileMap {
	tic := time.Now()
	maxDepth := depth

	defer func() {
		fmt.Println("done new empty tilemap from zip in. Mem usage:", time.Since(tic).Seconds())
		PrintMemUsage()
	}()
	tm := &TileMap{
		tileSize: tileSize,
		worldW:   worldW,
		worldH:   worldH,
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
		// name := fmt.Sprintf("%d/%d/%d.png", start.z, start.x, start.y)
		start.decoded = ebiten.NewImage(tileSize, tileSize)
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
	tm.Root = root
	return tm
}



// func (tm *TileMap) Redraw(drawers []Drawer) {
// 	ww, wh := tm.WorldSize()
// 	var dfs func(start *Tile)
// 	dfs = func(start *Tile) {
// 		scalar := float64(int(1) << (tm.Depth() - start.z))
// 		overlap := false
// 		for _, d := range drawers {
// 			dBounds := d.InWorldBounds(ww, wh)
// 			if start.bounds.Overlaps(dBounds) {
// 				overlap = true
// 				corner := dBounds.Sub(start.bounds.Min).Min
// 				// fmt.Println(corner)
// 				sprite, scale, tx, ty := d.SpriteOp()
// 				// fmt.Println(sprite.Bounds())
// 				// op := &ebiten.DrawImageOptions{}
// 				op := &ebiten.DrawImageOptions{}
// 				op.GeoM.Scale(scale/scalar, scale/scalar)
// 				op.GeoM.Translate((float64(corner.X)+tx)/scalar, (float64(corner.Y)+ty)/scalar)
// 				start.decoded.DrawImage(sprite, op)
// 			}
// 		}
// 		if !overlap {
// 			return
// 		}
// 		for _, c := range start.Children() {
// 			dfs(c)
// 		}
// 	}
// 	dfs(tm.Root)
// }
