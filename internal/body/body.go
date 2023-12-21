package body

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	evector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/oowhyy/astroapp/pkg/vector"
)

// ideally, pixel mappings should be general to the whole game and passed in body.Draw
// current implementation saves unnecesary calculation
const (
	PixelsPerAU = 200.0
	BaseImgSize = 20.0
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

	TrailLayer *ebiten.Image
	trailColor color.Color
}

func (b *Body) Update() {
	b.PrevPos = b.Pos
	b.Pos.Add(b.Vel)
	b.Vel.Add(b.Acc)
	b.Acc.Reset()
	if b.TrailLayer != nil {
		evector.StrokeLine(b.TrailLayer, float32(b.PrevPos.X), float32(b.PrevPos.Y), float32(b.Pos.X), float32(b.Pos.Y), 1, b.trailColor, true)
	}

}

func (b *Body) ApplyForce(force vector.Vector) {
	acc := vector.Scaled(force, 1.0/b.Mass)
	// fmt.Println("BodyMass", b.Mass, "added acc")
	b.Acc.Add(acc)
}

func (b *Body) Draw(screen *ebiten.Image, dx, dy float64) {
	bounds := b.image.Bounds().Size()
	minDimScale := BaseImgSize / float64(min(bounds.X, bounds.Y))
	finalScale := minDimScale * b.Diameter
	halfW := finalScale * float64(bounds.X) / 2
	halfH := finalScale * float64(bounds.Y) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(finalScale, finalScale)
	op.GeoM.Translate(b.Pos.X+dx-halfW, b.Pos.Y+dy-halfH)
	screen.DrawImage(b.image, op)
	// screen.DrawImage(b.image, &ebiten.DrawImageOptions{})
}

func (b *Body) Overlap(b2 *Body) bool {
	d1 := b.Diameter * BaseImgSize
	d2 := b2.Diameter * BaseImgSize
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
