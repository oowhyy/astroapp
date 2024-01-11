package game

import (
	"bytes"
	"fmt"
	"image/color"
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
	trailOps        *ebiten.DrawImageOptions
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
)

func init() {
	clrSc := ebiten.ColorScale{}
	// clrSc.ScaleAlpha(0.5)
	trailOps = &ebiten.DrawImageOptions{ColorScale: clrSc}
	trailOps.GeoM.Scale(1, 1)
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

	trailLayer *ebiten.Image

	// for adding new body
	addStage  AddStage
	rockPath  string
	blueArrow *ebiten.Image
	newId     int

	Bodies    map[int]*body.Body
	GConstant float64

	simSpeed  int
	showTrail bool
}

func (g *Game) WorldSize() (float64, float64) {
	p := g.World.Bounds().Size()
	return float64(p.X), float64(p.Y)
}

func (g *Game) Update() error {
	addMode := g.UI.IsAddMode()
	w, h := g.WorldSize()
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
			newb := g.Bodies[g.newId]
			newb.Pos.X, newb.Pos.Y = body.ToLocal(mx, my, w, h)
		case AddStageDraw:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.addStage = AddStageDraw
				g.addStage = AddStage("")
				mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
				mouseVec := vector.FromFloats(mx, my)
				rockVec := g.Bodies[g.newId].WorldCoords(w, h)
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
	for ii := 0; ii < g.simSpeed; ii++ {
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
			body.DrawTrail(g.World, g.trailLayer)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Clear()
	g.World.Fill(color.RGBA{50, 50, 50, 255})
	op := &ebiten.DrawImageOptions{}
	worldW, worldH := g.WorldSize()
	bounds := g.background.Bounds()
	op.GeoM.Scale(worldW/float64(bounds.Dx()), worldH/float64(bounds.Dy()))
	g.World.DrawImage(g.background, op)

	if g.addStage == AddStageDraw {
		mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
		mouseVec := vector.FromFloats(mx, my)
		rockVec := g.Bodies[g.newId].WorldCoords(worldW, worldH)
		velRaw := vector.Diff(rockVec, mouseVec)
		g.DrawPlanetVector(g.World, rockVec, velRaw)
		toDraw := vector.Scaled(velRaw, 1.0/500.0)
		txt := toDraw.String()
		op := &text.DrawOptions{}
		op.GeoM.Translate(rockVec.X+50, rockVec.Y+50)
		text.Draw(g.World, txt, mplusNormalFace, op)

	}
	// bodies
	for _, body := range g.Bodies {
		// if body.TrailLayer != nil {
		// 	g.World.DrawImage(body.TrailLayer, trailOps)
		// }
		body.Draw(g.World, g.trailLayer)
	}
	if g.showTrail {
		g.World.DrawImage(g.trailLayer, &ebiten.DrawImageOptions{})
	}
	g.Camera.Render(g.World, screen)
	// g.Camera.Render(g.trailLayer, screen)

	txt := fmt.Sprintf("TPS: %f FPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(screen, txt, 100, 0)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Camera.ViewPort = vector.FromInts(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) CenterToWorld(x, y float64) vector.Vector {
	w, h := g.WorldSize()
	return vector.Vector{X: w/2 + x, Y: h/2 + y}
}

func (g *Game) PlanetToWorld(pos vector.Vector) vector.Vector {
	w, h := g.WorldSize()
	return vector.Vector{X: w/2 + pos.X, Y: h/2 + pos.Y}
}
