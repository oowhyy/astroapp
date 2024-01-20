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
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/tilemap"
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
	worldSize vector.Vector
	Camera    *camera.Camera
	UI        webui.UserInterface

	background *tilemap.TileMap
	trailLayer *tilemap.TileMap

	// for adding new body
	addStage  AddStage
	rockPath  string
	blueArrow *ebiten.Image
	newId     int

	Bodies    []*body.Body
	GConstant float64

	simSpeed  int
	showTrail bool
}

func (g *Game) WorldSize() (float64, float64) {
	return g.worldSize.X, g.worldSize.Y
}

func (g *Game) Update() error {
	// addMode := g.UI.IsAddMode()
	// w, h := g.WorldSize()
	// if addMode {
	// 	//add mode
	// 	switch g.addStage {
	// 	case AddStage(""):
	// 		g.addStage = AddStageMove
	// 		newId := len(g.Bodies) + 1
	// 		cfg := &body.BodyConfig{
	// 			Id:       newId,
	// 			Name:     "rock",
	// 			Image:    g.rockPath,
	// 			Mass:     1,
	// 			Diameter: 40,
	// 			X:        1000,
	// 			Y:        1000,
	// 			Dx:       0,
	// 			Dy:       0,
	// 		}
	// 		newb, err := body.FromConfig(cfg)
	// 		if err != nil {
	// 			panic("failed to add new body")
	// 		}
	// 		g.newId = newId
	// 		newb.Frozen = true
	// 		g.Bodies[newId] = newb
	// 	case AddStageMove:
	// 		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
	// 			g.addStage = AddStageDraw
	// 		}
	// 		mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
	// 		newb := g.Bodies[g.newId]
	// 		newb.Pos.X, newb.Pos.Y = body.ToLocal(mx, my, w, h)
	// 	case AddStageDraw:
	// 		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
	// 			g.addStage = AddStageDraw
	// 			g.addStage = AddStage("")
	// 			mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
	// 			mouseVec := vector.FromFloats(mx, my)
	// 			rockVec := g.Bodies[g.newId].WorldCoords(w, h)
	// 			velRaw := vector.Diff(rockVec, mouseVec)
	// 			scaled := vector.Scaled(velRaw, 1.0/200.0)
	// 			newb := g.Bodies[g.newId]
	// 			newb.Vel = scaled
	// 			newb.Frozen = false
	// 			g.newId = 0
	// 		}
	// 	}
	// } else {
	// 	g.Camera.Update()
	// 	if g.newId != 0 {
	// 		delete(g.Bodies, g.newId)
	// 		g.addStage = AddStage("")
	// 		g.newId = 0
	// 	}
	// }
	g.Camera.Update()
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

	// dd := make([]tilemap.Drawer, len(g.b0))
	// for i, b := range g.b0 {
	// 	dd[i] = tilemap.Drawer(b)
	// }
	// g.bodiesMap.Redraw(dd)
	// g.Camera.RenderTilemap(screen, g.bodiesMap)
	// sub := g.world.SubImage(rect).(*ebiten.Image)
	// screen.DrawImage(sub, nil)
	// g.World.Clear()
	// g.World.Fill(color.RGBA{50, 50, 50, 255})

	// worldW, worldH := g.WorldSize()

	// if g.addStage == AddStageDraw {
	// 	mx, my := g.Camera.ScreenToWorld(ebiten.CursorPosition())
	// 	mouseVec := vector.FromFloats(mx, my)
	// 	rockVec := g.Bodies[g.newId].WorldCoords(worldW, worldH)
	// 	velRaw := vector.Diff(rockVec, mouseVec)
	// 	g.DrawPlanetVector(g.World, rockVec, velRaw)
	// 	toDraw := vector.Scaled(velRaw, 1.0/500.0)
	// 	txt := toDraw.String()
	// 	op := &text.DrawOptions{}
	// 	op.GeoM.Translate(rockVec.X+50, rockVec.Y+50)
	// 	text.Draw(g.World, txt, mplusNormalFace, op)
	// }
	// bodies
	// for _, body := range g.Bodies {
	// 	body.Draw(g.World, g.trailLayer)
	// }
	// if g.showTrail {
	// 	g.World.DrawImage(g.trailLayer, &ebiten.DrawImageOptions{})
	// }
	// g.Camera.Render(g.World, screen)
	// g.Camera.Render(g.trailLayer, screen)

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

func (g *Game) CenterToWorld(x, y float64) vector.Vector {
	w, h := g.WorldSize()
	return vector.Vector{X: w/2 + x, Y: h/2 + y}
}

func (g *Game) PlanetToWorld(pos vector.Vector) vector.Vector {
	w, h := g.WorldSize()
	return vector.Vector{X: w/2 + pos.X, Y: h/2 + pos.Y}
}
