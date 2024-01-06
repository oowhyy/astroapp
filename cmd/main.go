package main

import (
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/game"
	"gopkg.in/yaml.v3"
)

var (
	configReroURI = "https://raw.githubusercontent.com/oowhyy/astroapp/main/astroapp/config.yaml"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Render an image")
	cfg := &game.Config{}
	err := readConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	g := game.FromConfig(cfg)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func readConfig(cfg *game.Config) error {
	resp, err := http.Get(configReroURI)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	err = yaml.NewDecoder(resp.Body).Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
