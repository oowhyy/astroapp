package camera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/oowhyy/astroapp/pkg/vector"
)

type Camera struct {
	ViewPort   vector.Vector
	Position   vector.Vector
	MousePan   vector.Vector
	ZoomFactor float64
	Rotation   int
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"Position: %.1f, Rotation: %d, Scale: %f",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() vector.Vector {
	return vector.Vector{
		X: c.ViewPort.X * 0.5,
		Y: c.ViewPort.Y * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position.X, -c.Position.Y)
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter().X, -c.viewportCenter().Y)
	m.Scale(
		c.ZoomFactor,
		c.ZoomFactor,
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter().X, c.viewportCenter().Y)
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})
}

func (c *Camera) Update() {
	// pan

	mousePos := vector.FromInts(ebiten.CursorPosition())
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.MousePan = mousePos
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		c.Position.Sub(vector.Diff(mousePos, c.MousePan))
		c.MousePan = mousePos
	}

	// zoom
	_, dy := ebiten.Wheel()
	if dy > 0 {
		c.ZoomFactor *= 1.07
	}
	if dy < 0 {
		c.ZoomFactor *= 0.93
	}
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position.X = 0
	c.Position.Y = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}
