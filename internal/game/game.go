package game

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/oowhyy/astroapp/internal/assets"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/vector"

	evector "github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	World  *ebiten.Image
	Camera camera.Camera

	Bodies map[int]*body.Body
}

func (g *Game) WorldSize() image.Point {
	return g.World.Bounds().Size()
}

func NewGame() *Game {
	g := &Game{}
	g.Camera = camera.Camera{ZoomFactor: 0.5}
	g.World = ebiten.NewImage(assets.BackgroundImage.Bounds().Dx(), assets.BackgroundImage.Bounds().Dy())
	g.Bodies = make(map[int]*body.Body, 10)
	// center camera
	size := g.WorldSize()
	screenW, screenH := ebiten.WindowSize()
	if screenW != 0 {
		g.Camera.Position.X = (float64(size.X) - float64(screenW)) / 2
		g.Camera.Position.Y = (float64(size.Y) - float64(screenH)) / 2
	} else {
		g.Camera.Position.X = float64(size.X) / 4
		g.Camera.Position.Y = float64(size.Y) / 4
	}

	// test bodies
	sun := body.NewBody(assets.RedGemImage, body.WithPosVector(g.CenterRelative(0, 0)), body.WithMass(5000))

	planet1 := body.NewBody(assets.RedGemImage, body.WithPosVector(g.CenterRelative(0, -300)), body.WithVel(3.7, 0))
	planet2 := body.NewBody(assets.RedGemImage, body.WithPosVector(g.CenterRelative(-600, 0)), body.WithVel(0, 2.2))
	g.Bodies[0] = sun
	g.Bodies[1] = planet1
	g.Bodies[2] = planet2
	return g
}

const (
	GConstant = 1.2
)

func DrawArrow(img *ebiten.Image, from, to vector.Vector) {
	evector.StrokeLine(img, float32(from.X), float32(from.Y), float32(to.X), float32(to.Y), 5, color.White, true)
	evector.DrawFilledCircle(img, float32(to.X), float32(to.Y), 10, color.RGBA{255, 0, 0, 255}, true)
}

func DrawVector(img *ebiten.Image, baseX, baseY float64, vec vector.Vector) {
	endY := float32(baseY + vec.Y)
	endX := float32(baseX + vec.X)
	evector.StrokeLine(img, float32(baseX), float32(baseY), endX, endY, 2, color.RGBA{255, 0, 0, 255}, true)
	r := 10.0
	a1 := math.Pi/6 + math.Pi
	a2 := -math.Pi/6 + math.Pi
	atan := math.Atan2(vec.Y, vec.X)
	leaf1X := r * math.Cos(atan+a1)
	leaf1Y := r * math.Sin(atan+a1)
	leaf2X := r * math.Cos(atan+a2)
	leaf2Y := r * math.Sin(atan+a2)

	evector.StrokeLine(img, endX, endY, endX+float32(leaf1X), endY+float32(leaf1Y), 2, color.RGBA{255, 0, 0, 255}, true)
	evector.StrokeLine(img, endX, endY, endX+float32(leaf2X), endY+float32(leaf2Y), 2, color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) Update() error {
	g.Camera.Update()

	// physics
	for _, body := range g.Bodies {
		for _, forceSource := range g.Bodies {
			if body == forceSource || body.Overlap(forceSource) {
				continue
			}
			distTo := body.DistTo(forceSource)
			forceVec := body.UnitDir(forceSource)
			forceMag := GConstant * body.Mass * forceSource.Mass / (distTo * distTo)
			forceVec.Scale(forceMag)
			// DrawVector(g.World, body.Pos.X, body.Pos.Y, forceVec)
			body.ApplyForce(forceVec)
		}
	}

	for _, body := range g.Bodies {
		body.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Clear()
	g.World.DrawImage(assets.BackgroundImage, nil)

	// axis
	size := g.WorldSize()
	evector.StrokeLine(g.World, 0, float32(size.Y/2), float32(size.X), float32(size.Y/2), 2, color.RGBA{255, 255, 255, 50}, false)
	evector.StrokeLine(g.World, float32(size.X/2), 0, float32(size.X/2), float32(size.Y), 2, color.RGBA{255, 255, 255, 50}, false)
	// bodies
	for _, body := range g.Bodies {
		body.Draw(g.World)
	}

	// vec

	for _, body := range g.Bodies {
		for _, forceSource := range g.Bodies {
			if body == forceSource || body.Overlap(forceSource) {
				continue
			}
			distTo := body.DistTo(forceSource)
			forceVec := body.UnitDir(forceSource)
			forceMag := GConstant * body.Mass * forceSource.Mass / (distTo * distTo)
			forceVec.Scale(forceMag)
			forceVec.Scale(1000)
			DrawVector(g.World, body.Pos.X, body.Pos.Y, forceVec)
		}
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
