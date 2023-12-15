package body

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/pkg/vector"
)

const (
	BaseImgSize = 100
)

type Body struct {
	Pos vector.Vector
	Vel vector.Vector
	Acc vector.Vector
	// default - 1.0
	Mass   float64
	sprite *ebiten.Image
}

func NewBody(sprite *ebiten.Image, options ...BodyOptions) *Body {
	body := &Body{
		sprite: sprite,
		Mass:   1,
	}
	for _, op := range options {
		op(body)
	}
	return body
}

type BodyOptions func(b *Body)

func WithMass(mass float64) BodyOptions {
	return func(b *Body) {
		b.Mass = mass
	}
}

func WithPosVector(pos vector.Vector) BodyOptions {
	return func(b *Body) {
		b.Pos = pos
	}
}

func WithPos(x, y int) BodyOptions {
	return func(b *Body) {
		b.Pos = vector.Vector{
			X: float64(x),
			Y: float64(y),
		}
	}
}

func WithVel(dx, dy float64) BodyOptions {
	return func(b *Body) {
		b.Vel = vector.Vector{
			X: dx,
			Y: dy,
		}
	}
}

func WithVelVector(vel vector.Vector) BodyOptions {
	return func(b *Body) {
		b.Vel = vel
	}
}

func (b *Body) Update() {
	b.Pos.Add(b.Vel)
	b.Vel.Add(b.Acc)
	b.Acc.Reset()
}

func (b *Body) ApplyForce(force vector.Vector) {
	acc := vector.Scaled(force, 1.0/b.Mass)
	// fmt.Println("BodyMass", b.Mass, "added acc")
	b.Acc.Add(acc)
}

func (b *Body) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	bounds := b.sprite.Bounds().Size()
	minDimScale := BaseImgSize / float64(min(bounds.X, bounds.Y))
	finalScale := minDimScale
	// finalScale
	halfW := finalScale * float64(bounds.X) / 2
	halfH := finalScale * float64(bounds.Y) / 2
	op.GeoM.Scale(finalScale, finalScale)
	op.GeoM.Translate(b.Pos.X-halfW, b.Pos.Y-halfH)
	screen.DrawImage(b.sprite, op)
}

func (b *Body) Overlap(b2 *Body) bool {
	return b.DistTo(b2) < BaseImgSize*0.5
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
