package body

import (
	"image/png"
	"io"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/pkg/vector"
)

type BodyConfig struct {
	Id       int
	Name     string  `yaml:"name"`
	Parent   string  `yaml:"parent"`
	Mass     float64 `yaml:"mass"`
	Diameter float64 `yaml:"diameter"`
	X        float64 `yaml:"x"`
	Y        float64 `yaml:"y"`
	Dx       float64 `yaml:"dx"`
	Dy       float64 `yaml:"dy"`
}

func FromConfig(c *BodyConfig, imgBuf io.Reader) (*Body, error) {
	b := &Body{

		Name: c.Name,
		Mass: c.Mass,
	}
	// startAngle := rand.Float64() * math.Pi * 2
	x := c.X * DistScale
	y := c.Y * DistScale
	dx := c.Dx * DistScale
	dy := c.Dy * DistScale

	b.Pos = vector.FromFloats(x, y)
	b.Vel = vector.FromFloats(dx, dy)

	dia := math.Log2(c.Diameter+1) / 3
	b.Diameter = dia
	b.trailHue = rand.Float64() * math.Pi * 2
	decoded, err := png.Decode(imgBuf)
	if err != nil {
		return nil, err
	}
	b.image = ebiten.NewImageFromImage(decoded)
	return b, nil
}

func NewBody(c *BodyConfig, img *ebiten.Image) *Body {
	b := &Body{

		Name: c.Name,
		Mass: c.Mass,
	}
	// startAngle := rand.Float64() * math.Pi * 2
	x := c.X * DistScale
	y := c.Y * DistScale
	dx := c.Dx * DistScale
	dy := c.Dy * DistScale

	b.Pos = vector.FromFloats(x, y)
	b.Vel = vector.FromFloats(dx, dy)

	dia := math.Log2(c.Diameter+1) / 3
	b.Diameter = dia
	b.trailHue = rand.Float64() * math.Pi * 2
	b.image = img
	return b
}
