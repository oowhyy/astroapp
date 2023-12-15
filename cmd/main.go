package main

import (
	_ "embed"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/game"
)

var (
	WindowWidth  = 1600
	WindowHeight = 900
)

func main() {
	// ebiten.SetWindowSize(WindowWidth, WindowHeight)
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Render an image")
	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
