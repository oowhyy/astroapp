package common

import (
	"image"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
)

func ImageFromPath(fullURL string) (*ebiten.Image, error) {
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	image, _, err := image.Decode(resp.Body)
	if image == nil {
		panic("GOT HERE EMPTY IMAGE " + fullURL)
	}
	ebitenImg := ebiten.NewImageFromImage(image)
	return ebitenImg, err
}
