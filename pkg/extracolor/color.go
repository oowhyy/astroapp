package extracolor

import (
	"image/color"
	"math/rand"
)

func RandomRGB() color.Color {
	bits := 8
	r := uint8(256 / bits * (rand.Intn(bits) + 1))
	g := uint8(256 / bits * (rand.Intn(bits) + 1))
	b := uint8(256 / bits * (rand.Intn(bits) + 1))
	return color.RGBA{r, g, b, 255}
}
