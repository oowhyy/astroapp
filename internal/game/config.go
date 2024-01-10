package game

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/common"
	"github.com/oowhyy/astroapp/pkg/webui"
)

type Config struct {
	Background string             `yaml:"background"`
	Rock       string             `yaml:"rock"`
	Bodies     []*body.BodyConfig `yaml:"bodies"`
	Arrow      string             `yaml:"arrow"`
	GConstant  float64            `yaml:"gConstant"`
}

const (
	WorldWidth  = 4000
	WorldHeight = 4000
)

func FromConfig(c *Config) *Game {
	g := &Game{
		// parentMap: make(map[string]string, len(c.Bodies)),
	}
	g.UI = webui.NewWebInterface()

	g.GConstant = c.GConstant
	g.simSpeed = 1
	g.Camera = camera.Camera{ZoomFactor: 1}
	bgImg, err := common.ImageFromPath(c.Background)
	if err != nil {
		log.Fatalf("background image not found at %s", c.Background)
	}
	fmt.Println(c.Arrow)
	arrow, err := common.ImageFromPath(c.Arrow)
	if err != nil {
		log.Fatalf("background image not found at %s", c.Background)
	}

	g.rockPath = c.Rock
	g.background = bgImg
	g.blueArrow = arrow
	g.World = ebiten.NewImage(WorldWidth, WorldHeight)
	g.trailLayer = ebiten.NewImage(WorldWidth, WorldHeight)
	g.Bodies = make(map[int]*body.Body, 10)
	// center camera
	w, h := g.WorldSize()
	screenW, screenH := g.UI.WindowSize()
	g.Camera.Position.X = (w - float64(screenW)) / 2
	g.Camera.Position.Y = (h - float64(screenH)) / 2

	// setup ui callbacks
	g.UI.OnClearTrail(func() {
		g.trailLayer.Clear()
	})
	g.UI.OnSpeedUp(func() int {
		if g.simSpeed >= 9 {
			return g.simSpeed
		}
		g.simSpeed++
		return g.simSpeed
	})
	g.UI.OnSlowDown(func() int {
		if g.simSpeed <= 1 {
			return g.simSpeed
		}
		g.simSpeed--
		return g.simSpeed
	})

	// bodies

	// auto assign ids and fill config map
	configMap := make(map[string]*body.BodyConfig, len(c.Bodies))
	for i, b := range c.Bodies {
		b.Id = i
		configMap[b.Name] = b
	}
	// translate vectors relative to parent
	for _, bConf := range c.Bodies {
		if par, ok := configMap[bConf.Parent]; ok {
			bConf.X += par.X
			bConf.Y += par.Y
			bConf.Dy += par.Dy
			bConf.Dx += par.Dx
		}
	}

	bodiesMap := make(map[int]*body.Body, len(c.Bodies))
	for _, conf := range c.Bodies {
		body, err := body.FromConfig(conf)
		if err != nil {
			log.Fatalf("failed to load body from config: %s", err)
		}
		bodiesMap[body.Id] = body
	}
	g.Bodies = bodiesMap
	return g
}
