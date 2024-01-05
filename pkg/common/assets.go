package common

import (
	"image"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
)

func ImageFromPath(filePath string) (*ebiten.Image, error) {
	url := "https://raw.githubusercontent.com/oowhyy/astroapp/main/" + filePath
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	image, _, err := image.Decode(resp.Body)
	if image == nil {
		panic("HERE EMPTY IMAGE " + url)
	}
	ebitenImg := ebiten.NewImageFromImage(image)
	return ebitenImg, err
}
