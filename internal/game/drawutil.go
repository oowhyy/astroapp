package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	evector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/oowhyy/astroapp/pkg/vector"
)

func DrawArrow(img *ebiten.Image, from, to vector.Vector) {
	endY := float32(from.Y + to.Y)
	endX := float32(from.X + to.X)
	evector.StrokeLine(img, float32(from.X), float32(from.Y), float32(endX), float32(endY), 5, color.White, true)
	evector.DrawFilledCircle(img, float32(to.X), float32(to.Y), 10, color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) DrawLine(img *ebiten.Image, from, to vector.Vector) {

	diff := vector.Diff(to, from)
	op := &ebiten.DrawImageOptions{}
	len := diff.Len()
	width := 30.0
	arrowBounds := g.blueArrow.Bounds()
	op.GeoM.Scale(len/float64(arrowBounds.Dx()), width/float64(arrowBounds.Dy()))
	angle := math.Atan2(diff.Y, diff.X)
	op.GeoM.Rotate(angle)
	yy := width / 2 * math.Cos(angle)
	xx := width / 2 * math.Sin(angle)
	op.GeoM.Translate(from.X+xx, from.Y-yy)

	img.DrawImage(g.blueArrow, op)

}

func (g *Game) DrawArrow(screen *ebiten.Image, from, to vector.Vector) {

	diff := vector.Diff(to, from)
	op := &ebiten.DrawImageOptions{}
	len := diff.Len()
	width := 30.0
	arrowBounds := g.blueArrow.Bounds()
	op.GeoM.Scale(len/float64(arrowBounds.Dx()), width/float64(arrowBounds.Dy()))
	angle := math.Atan2(diff.Y, diff.X)
	op.GeoM.Rotate(angle)
	yy := width / 2 * math.Cos(angle)
	xx := width / 2 * math.Sin(angle)
	op.GeoM.Translate(from.X+xx, from.Y-yy)

	screen.DrawImage(g.blueArrow, op)

}

func (g *Game) DrawPlanetVector(img *ebiten.Image, worldPos vector.Vector, vec vector.Vector) {
	limVec := vec.Clone()
	limVec.Limit(-1000, 1000)
	endY := float32(worldPos.Y + vec.Y)
	endX := float32(worldPos.X + vec.X)
	// evector.StrokeLine(img, float32(worldPos.X), float32(worldPos.Y), endX, endY, 5, color.RGBA{255, 0, 0, 255}, true)
	g.DrawLine(img, worldPos, vector.FromFloats(float64(endX), float64(endY)))
	// r := 20.0
	// if vec.Len() < 2*r {
	// 	return
	// }
	// a1 := math.Pi/6 + math.Pi
	// a2 := -math.Pi/6 + math.Pi
	// atan := math.Atan2(limVec.Y, limVec.X)
	// leaf1X := r * math.Cos(atan+a1)
	// leaf1Y := r * math.Sin(atan+a1)
	// leaf2X := r * math.Cos(atan+a2)
	// leaf2Y := r * math.Sin(atan+a2)

	// evector.StrokeLine(img, endX, endY, endX+float32(leaf1X), endY+float32(leaf1Y), 5, color.RGBA{255, 0, 0, 255}, true)
	// evector.StrokeLine(img, endX, endY, endX+float32(leaf2X), endY+float32(leaf2Y), 5, color.RGBA{255, 0, 0, 255}, true)
}
