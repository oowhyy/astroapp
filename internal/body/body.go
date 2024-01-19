package body

import (
	"fmt"
	"image"
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


func (b *Body) InWorldBounds(ww, wh int) image.Rectangle {
	bounds := b.image.Bounds().Size()
	// return image.Rect(worldSize/2, worldSize/2, worldSize/2+bounds.X/2, worldSize/2+bounds.Y/2)
	finalScale := b.Diameter * (BaseImgSize / float64(min(bounds.X, bounds.Y)))
	w := finalScale * float64(bounds.X)
	h := finalScale * float64(bounds.Y)
	x0 := float64(ww/2) + b.Pos.X -w
	y0 := float64(wh/2) + b.Pos.Y -h
	return image.Rect(int(x0-100), int(y0-100), int(2*w+x0+100), int(2*h+y0+100))
}

func (b *Body) SpriteOp() (*ebiten.Image, float64, float64, float64) {
	bounds := b.image.Bounds().Size()
	finalScale := b.Diameter * (BaseImgSize / float64(min(bounds.X, bounds.Y)))
	halfW := finalScale * float64(bounds.X) / 2
	halfH := finalScale * float64(bounds.Y) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(finalScale, finalScale)
	op.GeoM.Translate(-halfW, -halfH)
	return b.image, finalScale, -halfW, -halfH
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
