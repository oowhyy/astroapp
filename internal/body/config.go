package body

import (
	"math"

	"github.com/oowhyy/astroapp/pkg/common"
	"github.com/oowhyy/astroapp/pkg/vector"
)

type BodyConfig struct {
	Id       int
	Name     string  `yaml:"name"`
	Image    string  `yaml:"image"`
	Parent   string  `yaml:"parent"`
	Mass     float64 `yaml:"mass"`
	Diameter float64 `yaml:"diameter"`
	X        float64 `yaml:"x"`
	Y        float64 `yaml:"y"`
	Dx       float64 `yaml:"dx"`
	Dy       float64 `yaml:"dy"`
}

func FromConfig(c *BodyConfig) (*Body, error) {
	if len(c.Image) == 0 {
		panic("no image provided")
	}
	b := &Body{
		Id:   c.Id,
		Name: c.Name,
		Mass: c.Mass,
	}
	x := c.X * PixelsPerAU
	y := c.Y * PixelsPerAU
	dx := c.Dx * PixelsPerAU
	dy := c.Dy * PixelsPerAU
	b.Pos = vector.FromFloats(x, y)
	b.Vel = vector.FromFloats(dx, dy)

	dia := math.Log2(c.Diameter + 1)
	b.Diameter = dia

	image, err := common.ImageFromPath(c.Image)
	if err != nil {
		return nil, err
	}
	b.image = image
	return b, nil
}

// func NewBody(img *ebiten.Image, name string) Body {
// 	return
// }
