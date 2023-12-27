package game

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/vector"

	evector "github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	trailColorScale = ebiten.ColorScale{}
)

func init() {
	trailColorScale.ScaleAlpha(0.5)
}

type Game struct {
	World  *ebiten.Image
	Camera camera.Camera

	background *ebiten.Image

	Bodies map[int]*body.Body

	GConstant float64

	// parentMap map[string]string
}

func (g *Game) WorldSize() image.Point {
	return g.World.Bounds().Size()
}

func DrawArrow(img *ebiten.Image, from, to vector.Vector) {
	evector.StrokeLine(img, float32(from.X), float32(from.Y), float32(to.X), float32(to.Y), 5, color.White, true)
	evector.DrawFilledCircle(img, float32(to.X), float32(to.Y), 10, color.RGBA{255, 0, 0, 255}, true)
}

func DrawVector(img *ebiten.Image, baseX, baseY float64, vec vector.Vector) {
	if vec.X > 1000 {
		fmt.Println("too big", vec)
		vec.X = 1000
	}
	if vec.Y > 1000 {
		vec.Y = 1000
	}
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
			if body == forceSource {
				continue
			}
			if body.Overlap(forceSource) {
				fmt.Printf("%s overlapping %s\n", body, forceSource)
				continue
			}
			distTo := body.DistTo(forceSource)
			forceVec := body.UnitDir(forceSource)
			forceMag := g.GConstant * body.Mass * forceSource.Mass / (distTo * distTo)
			forceVec.Scale(forceMag)
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
	g.World.DrawImage(g.background, nil)
	// axis
	size := g.WorldSize()

	halfW := float64(size.X / 2)
	halfH := float64(size.Y / 2)
	// bodies
	for _, body := range g.Bodies {
		if body.TrailLayer != nil {
			g.World.DrawImage(body.TrailLayer, &ebiten.DrawImageOptions{ColorScale: trailColorScale})
		}
		body.Draw(g.World, halfW, halfH)
	}

	// vec

	// for _, body := range g.Bodies {
	// 	for _, forceSource := range g.Bodies {
	// 		if body == forceSource || body.Overlap(forceSource) {
	// 			continue
	// 		}
	// 		distTo := body.DistTo(forceSource)
	// 		forceVec := body.UnitDir(forceSource)
	// 		forceMag := g.GConstant * body.Mass * forceSource.Mass / (distTo * distTo)
	// 		forceVec.Scale(forceMag / body.Mass)
	// 		forceVec.Scale(100)
	// 		DrawVector(g.World, body.Pos.X+halfW, body.Pos.Y+halfH, forceVec)
	// 	}
	// }

	g.Camera.Render(g.World, screen)
	_, screenSizeY := ebiten.WindowSize()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f", ebiten.ActualTPS()))
	ebitenutil.DebugPrintAt(
		screen,

		g.Camera.String(),

		0, screenSizeY-32,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Camera.ViewPort = vector.FromInts(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) CenterRelative(x, y int) vector.Vector {
	size := g.WorldSize()
	return vector.Vector{X: float64(size.X)/2 + float64(x), Y: float64(size.Y)/2 + float64(y)}
}
