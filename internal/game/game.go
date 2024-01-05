package game

import (
	"bytes"
	"fmt"
	"image"
	"log"

	_ "embed"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/vector"
	"github.com/oowhyy/astroapp/pkg/webui"
)

var (
	trailColorScale = ebiten.ColorScale{}
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
)

func init() {
	trailColorScale.ScaleAlpha(0.5)
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	mplusNormalFace = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
}

type AddStage string

const (
	AddStageMove AddStage = "Not placed"
	AddStageDraw AddStage = "Placed"
)

type Game struct {
	World      *ebiten.Image
	Camera     camera.Camera
	UI         webui.UserInterface
	background *ebiten.Image

	// for adding new body
	addStage  AddStage
	rockPath  string
	blueArrow *ebiten.Image
	newId     int

	Bodies    map[int]*body.Body
	GConstant float64

	// parentMap map[string]string
}

func (g *Game) WorldSize() image.Point {
	return g.World.Bounds().Size()
}

func (g *Game) Update() error {
	addMode := g.UI.IsAddMode()
	if addMode {
		//add mode
		switch g.addStage {
		case AddStage(""):
			g.addStage = AddStageMove
			newId := len(g.Bodies) + 1
			cfg := &body.BodyConfig{
				Id:       newId,
				Name:     "rock",
				Image:    g.rockPath,
				Mass:     1,
				Diameter: 40,
				X:        1000,
				Y:        1000,
				Dx:       0,
				Dy:       0,
			}
			newb, err := body.FromConfig(cfg)
			if err != nil {
				panic("failed to add new body")
			}
			g.newId = newId
			newb.Frozen = true
			g.Bodies[newId] = newb
		case AddStageMove:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.addStage = AddStageDraw
			}
			mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
			size := g.WorldSize()
			halfW := float64(size.X / 2)
			halfH := float64(size.Y / 2)
			newb := g.Bodies[g.newId]
			newb.Pos.X = mx - halfW
			newb.Pos.Y = my - halfH
		case AddStageDraw:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.addStage = AddStageDraw
				g.addStage = AddStage("")
				mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
				mouseVec := vector.FromFloats(mx, my)
				rockVec := g.PlanetToWorld(g.Bodies[g.newId].Pos)
				velRaw := vector.Diff(rockVec, mouseVec)
				scaled := vector.Scaled(velRaw, 1.0/200.0)
				newb := g.Bodies[g.newId]
				newb.Vel = scaled
				newb.Frozen = false
				g.newId = 0
			}
		}
	} else {
		g.Camera.Update()
		if g.newId != 0 {
			delete(g.Bodies, g.newId)
			g.addStage = AddStage("")
			g.newId = 0
		}
	}
	if g.UI.IsPaused() {
		return nil
	}

	// physics
	for _, body := range g.Bodies {
		if body.Frozen {
			continue
		}
		for _, forceSource := range g.Bodies {
			if body == forceSource || body.Frozen {
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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 10)
	g.World.DrawImage(g.background, op)

	if g.addStage == AddStageDraw {
		mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
		mouseVec := vector.FromFloats(mx, my)
		rockVec := g.PlanetToWorld(g.Bodies[g.newId].Pos)
		velRaw := vector.Diff(rockVec, mouseVec)
		g.DrawPlanetVector(g.World, rockVec, velRaw)
		toDraw := vector.Scaled(velRaw, 1.0/200.0)
		txt := toDraw.String()
		op := &text.DrawOptions{}
		op.GeoM.Translate(rockVec.X+50, rockVec.Y+50)
		text.Draw(g.World, txt, mplusNormalFace, op)

	}

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

	g.Camera.Render(g.World, screen)

	txt := fmt.Sprintf("TPS: %f FPS: %f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), g.Camera.String())
	ebitenutil.DebugPrintAt(screen, txt, 100, 0)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Camera.ViewPort = vector.FromInts(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) CenterToWorld(x, y float64) vector.Vector {
	size := g.WorldSize()
	return vector.Vector{X: float64(size.X)/2 + x, Y: float64(size.Y)/2 + y}
}

func (g *Game) PlanetToWorld(pos vector.Vector) vector.Vector {
	size := g.WorldSize()
	return vector.Vector{X: float64(size.X)/2 + pos.X, Y: float64(size.Y)/2 + pos.Y}
}
