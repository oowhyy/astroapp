package tilemap

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/oowhyy/astroapp/pkg/vector"
)

var pixel *ebiten.Image

func init() {
	pixel = ebiten.NewImage(1, 1)
	pixel.Fill(color.RGBA{255, 0, 0, 255})
}

func DrawLine(screen *ebiten.Image, hue float64, from, to vector.Vector) {
	diff := vector.Diff(to, from)
	op := &colorm.DrawImageOptions{}
	len := diff.Len()
	op.GeoM.Scale(len, 1)
	angle := math.Atan2(diff.Y, diff.X)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(from.X, from.Y)
	cmatrix := colorm.ColorM{}
	cmatrix.RotateHue(hue)
	colorm.DrawImage(screen, pixel, cmatrix, op)
}

type Trailer interface {
	TrailLine(worldW, worldH int) (prevX, prevY, curX, curY float64)
	HueTheta() float64
}

func (tm *TileMap) DrawTrail(trailer Trailer) {
	x0, y0, x1, y1 := trailer.TrailLine(tm.WorldSize())
	rx0, ry0, rx1, ry1 := x0, y0, x1, y1
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	color := trailer.HueTheta()

	var dfs func(start *Tile)
	dfs = func(start *Tile) {

		if !Overlaps(start.bounds.Bounds(), x0, y0, x1, y1) {
			return
		}
		scalar := float64(int(1) << (tm.Depth() - start.z))
		xx0 := rx0 - float64(start.bounds.Min.X)
		yy0 := ry0 - float64(start.bounds.Min.Y)
		xx1 := rx1 - float64(start.bounds.Min.X)
		yy1 := ry1 - float64(start.bounds.Min.Y)
		// fmt.Println(sprite.Bounds())
		// op := &ebiten.DrawImageOptions{}
		// evector.StrokeLine(start.decoded, float32(xx0/scalar), float32(yy0/scalar), float32(xx1/scalar), float32(yy1/scalar), 1, color.White, true)
		DrawLine(start.decoded, color, vector.FromFloats(xx0/scalar, yy0/scalar), vector.FromFloats(xx1/scalar, yy1/scalar))
		start.isEmpty = false
		for _, c := range start.Children() {
			dfs(c)
		}
	}
	dfs(tm.Root)
}

func Overlaps(r image.Rectangle, x0, y0, x1, y1 float64) bool {
	return float64(r.Min.X) < x1 && x0 < float64(r.Max.X) &&
		float64(r.Min.Y) < y1 && y0 < float64(r.Max.Y)
}
