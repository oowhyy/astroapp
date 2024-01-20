package camera

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/pkg/vector"
)

var pixel *ebiten.Image

func init() {
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)
}

type Drawer interface {
	InWorldBounds(worldWidth, worldHeight int) image.Rectangle
	WorldCoords(worldW, worldH int) (float64, float64)
	SpriteOp() (sprite *ebiten.Image, scale float64, tx, ty float64)
}

func DrawLine(screen *ebiten.Image, from, to vector.Vector) {
	diff := vector.Diff(to, from)
	op := &ebiten.DrawImageOptions{}
	len := diff.Len()
	width := 30.0
	op.GeoM.Scale(len, 1)
	angle := math.Atan2(diff.Y, diff.X)
	op.GeoM.Rotate(angle)
	yy := width / 2 * math.Cos(angle)
	xx := width / 2 * math.Sin(angle)
	op.GeoM.Translate(from.X+xx, from.Y-yy)

	screen.DrawImage(pixel, op)

}
