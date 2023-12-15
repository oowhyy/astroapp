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

func WithPos(pos vector.Vector) BodyOptions {
	return func(b *Body) {
		b.Pos = pos
	}
}

func WithVel(vel vector.Vector) BodyOptions {
	return func(b *Body) {
		b.Vel = vel
	}
}

func (b *Body) Update() {
	b.Pos.Add(b.Vel)
	b.Pos.Add(b.Acc)
	b.Acc.Reset()
}

func (b *Body) ApplyForce(force vector.Vector) {
	b.Acc.Add(vector.Scaled(force, 1.0/b.Mass))
}

func (b *Body) Draw(screen *ebiten.Image) {
	op := baseScale(b.sprite)
	op.GeoM.Translate(b.Pos.X, b.Pos.Y)
	op.GeoM.Scale(b.Mass, b.Mass)
	screen.DrawImage(b.sprite, op)
}

func baseScale(img *ebiten.Image) *ebiten.DrawImageOptions {
	bounds := img.Bounds().Size()
	maxDimScale := BaseImgSize / float64(max(bounds.X, bounds.Y))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(maxDimScale, maxDimScale)
	return op
}
