package game

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/oowhyy/astroapp/internal/assets"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/vector"
)

type Game struct {
	World  *ebiten.Image
	Camera camera.Camera

	Bodies []*body.Body
}

func (g *Game) WorldSize() image.Point {
	return g.World.Bounds().Size()
}

func NewGame() *Game {
	g := &Game{}
	g.World = assets.BackgroundImage
	// center camera
	size := g.WorldSize()
	screenW, screenH := ebiten.WindowSize()
	if size.X > screenW {
		g.Camera.Position.X = (float64(size.X) - float64(screenW)) / 2
	}
	if size.Y > screenH {
		g.Camera.Position.Y = (float64(size.Y) - float64(screenH)) / 2
	}

	sun := body.NewBody(assets.RedGemImage, body.WithPos(g.CenterRelative(0, 0)))
	palnet := body.NewBody(assets.RedGemImage, body.WithPos(g.CenterRelative(0, -200)))
	g.Bodies = append(g.Bodies, sun)
	g.Bodies = append(g.Bodies, palnet)
	return g
}

func (g *Game) Update() error {
	g.Camera.Update()

	// physics

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, body := range g.Bodies {
		body.Draw(g.World)
	}

	g.Camera.Render(g.World, screen)
	_, screenSizeY := ebiten.WindowSize()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %d", ebiten.TPS()))
	ebitenutil.DebugPrintAt(
		screen,

		g.Camera.String(),

		0, screenSizeY-32,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Camera.ViewPort = vector.FromInt(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) CenterRelative(x, y int) vector.Vector {
	size := g.WorldSize()
	return vector.Vector{X: float64(size.X)/2 + float64(x), Y: float64(size.Y)/2 + float64(y)}
}
