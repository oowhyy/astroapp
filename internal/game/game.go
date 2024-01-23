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
	"github.com/oowhyy/astroapp/pkg/tilemap"
	"github.com/oowhyy/astroapp/pkg/vector"
	"github.com/oowhyy/astroapp/pkg/webui"
)

var (
	mplusFaceSource *text.GoTextFaceSource
	mplusNormalFace *text.GoTextFace
)

func init() {
	// clrSc.ScaleAlpha(0.5)
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
	worldSize image.Point
	Camera    *camera.Camera
	UI        webui.UserInterface

	background *tilemap.TileMap
	trailLayer *tilemap.TileMap

	Bodies    []*body.Body
	GConstant float64

	simSpeed  int
	showTrail bool

	// for adding new body
	addStage  AddStage
	blueArrow *ebiten.Image
	rock      *ebiten.Image
	newBody   *body.Body
}

func (g *Game) WorldSize() (int, int) {
	return g.worldSize.X, g.worldSize.Y
}

func (g *Game) Update() error {
	addMode := g.UI.IsAddMode()
	ww, wh := g.WorldSize()
	if addMode {
		//add mode
		switch g.addStage {
		case AddStage(""):
			g.addStage = AddStageMove
			cfg := &body.BodyConfig{
				Name:     "rock",
				Mass:     1,
				Diameter: 40,
				X:        1000,
				Y:        1000,
				Dx:       0,
				Dy:       0,
			}
			newb := body.NewBody(cfg, g.rock)
			newb.Frozen = true
			g.newBody = newb
		case AddStageMove:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.addStage = AddStageDraw
			}
			fmt.Println(g.newBody.Pos)
			mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
			newb := g.newBody
			newb.SetWorldPos(mx, my, ww, wh)
		case AddStageDraw:
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.addStage = AddStageDraw
				g.addStage = AddStage("")
				mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
				mouseVec := vector.FromFloats(mx, my)
				rockVec := g.newBody.WorldPosVec(ww, wh)
				velRaw := vector.Diff(rockVec, mouseVec)
				scaled := vector.Scaled(velRaw, 1.0/200.0)
				newb := g.newBody
				newb.Vel = scaled
				newb.Frozen = false
				g.Bodies = append(g.Bodies, newb)
				g.newBody = nil
			}
		}
	} else {
		g.Camera.Update()
		if g.newBody != nil {
			g.newBody = nil
			g.addStage = AddStage("")
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
			if g.showTrail {
				g.trailLayer.DrawTrail(body)
			}
		}
	}
	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()
	mx, my := g.Camera.ScreenToWorld(cx, cy)
	g.Camera.RenderTilemap(screen, g.background)
	if g.showTrail {
		g.Camera.RenderTilemap(screen, g.trailLayer)
	}
	for _, b := range g.Bodies {
		g.Camera.RenderDrawer(screen, b)
	}
	if g.newBody != nil {
		g.Camera.RenderDrawer(screen, g.newBody)
	}
	worldW, worldH := g.WorldSize()
	if g.addStage == AddStageDraw {
		screenCursorX, screenCursorY := ebiten.CursorPosition()
		// mx, my := g.Camera.ScreenToWorld(screenCursorX, screenCursorY)
		// mouseVec := vector.FromFloats(mx, my)
		rockWorld := g.newBody.WorldPosVec(worldW, worldH)
		rockScreenX, rockScreenY := g.Camera.WorldToScreen(rockWorld.X, rockWorld.Y)
		// velRaw := vector.Diff(rockWorld, vector.FromInts(screenCursorX, screenCursorY))
		velRaw := vector.FromInts(rockScreenX-screenCursorX, rockScreenY-screenCursorY)
		g.DrawPlanetVector(screen, vector.FromInts(rockScreenX, rockScreenY), velRaw)
		toDraw := vector.Scaled(velRaw, 1.0/200.0)
		txt := toDraw.String()
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(rockScreenX)+50, float64(rockScreenY)+50)
		text.Draw(screen, txt, mplusNormalFace, op)
	}

	ebitenutil.DebugPrintAt(screen,
		fmt.Sprintf("cursor screen: %d %d\ncursor world: %f %f\nCamera Rect: %v\nCameraZoom: %f", cx, cy, mx, my, g.Camera.GetRect(), g.Camera.ZoomFactor()),
		10, 500,
	)
	txt := fmt.Sprintf("TPS: %f FPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(screen, txt, 100, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Camera.ScreenSize = image.Pt(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
