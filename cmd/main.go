package main

import (
	_ "embed"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/oowhyy/astroapp/internal/game"
	"gopkg.in/yaml.v3"
)

var (
	WindowWidth  = 1600
	WindowHeight = 900
)

var RunInBrowser = false

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Render an image")
	cfg := &game.Config{}
	err := readConfig("config.yaml", cfg)
	if err != nil {
		log.Fatal(err)
	}

	g := game.FromConfig(cfg)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func readConfig(path string, cfg *game.Config) error {
	pid := os.Getpid()
	// wasm
	if pid == -1 {
		cfg.IsWasm = true
		resp, err := http.Get("https://raw.githubusercontent.com/oowhyy/astroapp/main/" + path)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		err = yaml.NewDecoder(resp.Body).Decode(cfg)
		if err != nil {
			log.Fatal(err)
		}
		cfg.IsWasm = true
		return nil
	}
	// desktop
	return cleanenv.ReadConfig(path, cfg)
}
