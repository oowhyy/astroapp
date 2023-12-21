package common

import (
	"image"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func ImageFromPath(filePath string, isWasm bool) (*ebiten.Image, error) {
	if isWasm {
		url := "https://raw.githubusercontent.com/oowhyy/astroapp/main/" + filePath
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		image, _, err := image.Decode(resp.Body)
		ebitenImg := ebiten.NewImageFromImage(image)
		return ebitenImg, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	ebitenImg := ebiten.NewImageFromImage(image)
	return ebitenImg, err
}
