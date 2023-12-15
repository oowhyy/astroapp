package extracolor

import (
	"image/color"
	"math/rand"
)

func RandomRGB() color.Color {
	r := uint8(rand.Intn(2))
	g := uint8(rand.Intn(2))
	b := uint8(rand.Intn(2))
	if r+g+b == 0 {
		return color.RGBA{255, 255, 255, 255}
	}
	return color.RGBA{255 * r, 255 * g, 255 * b, 255}
}
