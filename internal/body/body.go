package body

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	evector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/oowhyy/astroapp/pkg/vector"
)

const (
	DistScale    = 200.0 // affects physics
	ToScreenMult = 0.33
	BaseImgSize  = 20.0
	OverlapMult  = 0.2
)

type Body struct {
	Id      int
	Name    string
	Pos     vector.Vector
	PrevPos vector.Vector
	Vel     vector.Vector
	Acc     vector.Vector

	Mass     float64
	Diameter float64
	image    *ebiten.Image

	Frozen     bool
	trailColor color.Color
}

func (b *Body) Update() {
	b.Vel.Scale(0.99996) // friction in space xdd
	b.PrevPos = b.Pos
	b.Pos.Add(b.Vel)
	b.Vel.Add(b.Acc)
	b.Acc.Reset()

}

func (b *Body) ApplyForce(force vector.Vector) {
	acc := vector.Scaled(force, 1.0/b.Mass)
	// fmt.Println("BodyMass", b.Mass, "added acc")
	b.Acc.Add(acc)
}

func (b *Body) DrawTrail(screen *ebiten.Image, trailLayer *ebiten.Image) {
	screenSize := screen.Bounds()
	screenDx := float64(screenSize.Dx())
	screenDy := float64(screenSize.Dy())
	worldCoord := b.WorldCoords(screenDx, screenDy)
	screenX := worldCoord.X
	screenY := worldCoord.Y
	if trailLayer != nil {
		worldPrevX, worldPrevY := b.PrevPos.X*ToScreenMult+screenDx/2.0, b.PrevPos.Y*ToScreenMult+screenDy/2.0
		evector.StrokeLine(trailLayer, float32(worldPrevX), float32(worldPrevY), float32(screenX), float32(screenY), 1, b.trailColor, true)
	}
}

func (b *Body) Draw(screen *ebiten.Image, trailLayer *ebiten.Image) {
	bounds := b.image.Bounds().Size()
	minDimScale := BaseImgSize / float64(min(bounds.X, bounds.Y))
	finalScale := minDimScale * b.Diameter
	halfW := finalScale * float64(bounds.X) / 2
	halfH := finalScale * float64(bounds.Y) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(finalScale, finalScale)

	screenSize := screen.Bounds()
	screenDx := float64(screenSize.Dx())
	screenDy := float64(screenSize.Dy())
	worldCoord := b.WorldCoords(screenDx, screenDy)
	screenX := worldCoord.X
	screenY := worldCoord.Y

	// if trailLayer != nil {
	// 	worldPrevX, worldPrevY := b.PrevPos.X*ToScreenMult+screenDx/2.0, b.PrevPos.Y*ToScreenMult+screenDy/2.0
	// 	evector.StrokeLine(trailLayer, float32(worldPrevX), float32(worldPrevY), float32(screenX), float32(screenY), 1, b.trailColor, true)
	// }

	op.GeoM.Translate(screenX-halfW, screenY-halfH)
	screen.DrawImage(b.image, op)
	// screen.DrawImage(b.image, &ebiten.DrawImageOptions{})
}

func (b *Body) WorldCoords(worldX, worldY float64) vector.Vector {
	return vector.FromFloats(b.Pos.X*ToScreenMult+worldX/2.0, b.Pos.Y*ToScreenMult+worldY/2.0)
}

func ToLocal(wordlX, wordlY, wW, wH float64) (float64, float64) {
	return (wordlX - wW/2.0) / ToScreenMult, (wordlY - wH/2.0) / ToScreenMult
}

func (b *Body) Overlap(b2 *Body) bool {
	d1 := b.Diameter * BaseImgSize * OverlapMult
	d2 := b2.Diameter * BaseImgSize * OverlapMult
	return b.DistTo(b2) < (d1+d2)*0.5
}

func (b *Body) DistTo(b2 *Body) float64 {
	return vector.Diff(b.Pos, b2.Pos).Len()
}

func (b *Body) UnitDir(b2 *Body) vector.Vector {
	diff := vector.Diff(b2.Pos, b.Pos)
	mag := diff.Len()
	diff.Scale(1 / mag)
	return diff
}

func (b Body) String() string {
	return fmt.Sprintf("%s %s", b.Name, b.Pos)
}
