package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/game"
	"gopkg.in/yaml.v3"
)

var (
	configReroURI = "https://raw.githubusercontent.com/oowhyy/astroapp/main/config.yaml"
)

var refreshToken string
var appAuth string

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Render an image")
	cfg := &game.Config{}
	err := readConfig(cfg)
	fmt.Println("mercury:", cfg.Bodies[1])
	if err != nil {
		log.Fatal(err)
	}
	g := game.FromConfig(cfg, refreshToken, appAuth)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func readConfig(cfg *game.Config) error {
	fmt.Println("config URL:", configReroURI)
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
