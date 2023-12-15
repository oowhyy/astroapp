package assets

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed starfield_alpha.png
	BackgroundPng []byte

	//go:embed red_gem.png
	RedGemPng []byte
)

var (
	BackgroundImage *ebiten.Image
	RedGemImage     *ebiten.Image
)

func init() {
	var err error
	img, _, err := image.Decode(bytes.NewReader(BackgroundPng))
	if err != nil {
		log.Fatal(err)
	}
	BackgroundImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(RedGemPng))
	if err != nil {
		log.Fatal(err)
	}
	RedGemImage = ebiten.NewImageFromImage(img)
}
