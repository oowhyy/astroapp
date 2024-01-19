package game

import (
	"fmt"
	"image"
	"log"

	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/common"
	"github.com/oowhyy/astroapp/pkg/dropbox"
	"github.com/oowhyy/astroapp/pkg/tilemap"
	"github.com/oowhyy/astroapp/pkg/vector"
	"github.com/oowhyy/astroapp/pkg/webui"
)

type Config struct {
	Background    string             `yaml:"background"`
	BackgroundBig string             `yaml:"backgroundBig"`
	Rock          string             `yaml:"rock"`
	Bodies        []*body.BodyConfig `yaml:"bodies"`
	Arrow         string             `yaml:"arrow"`
	GConstant     float64            `yaml:"gConstant"`
}

// const (
// 	WorldWidth  = 4000
// 	WorldHeight = 4000
// )

func FromConfig(c *Config, refreshToken, appAuth string) *Game {
	g := &Game{
		// parentMap: make(map[string]string, len(c.Bodies)),
	}
	g.UI = webui.NewWebInterface()

	g.GConstant = c.GConstant
	g.simSpeed = 1
	fmt.Println("bg url : ", c.BackgroundBig)
	fmt.Println("bg image loaded")
	arrow, err := common.ImageFromPath(c.Arrow)
	if err != nil {
		log.Fatalf("background image not found at %s", c.Background)
	}

	g.rockPath = c.Rock

	// get bg from drop box
	accessToken, err := dropbox.GetAccessToken(refreshToken, appAuth)
	if err != nil {
		panic(err)
	}
	fmt.Println("token len", len(accessToken))
	cfg := dropbox.NewConfig(accessToken)
	client := dropbox.New(cfg)
	tiles, err := tilemap.NewTileMapFromDropboxZip(client, "/bg_tiles.zip")
	if err != nil {
		panic(err)
	}
	g.background = tiles
	fmt.Println("tilemap ok")

	g.blueArrow = arrow
	// g.trailLayer = ebiten.NewImage(WorldWidth, WorldHeight)

	// center camera
	worldX, worldY := 10000, 5625
	screenW, screenH := g.UI.WindowSize()
	g.Camera = camera.NewCamera(image.Pt(worldX, worldY), image.Pt(int(screenW), int(screenH)))
	g.Camera.Translate((float64(worldX)-screenW)/2, (float64(worldY)-screenH)/2)
	g.worldSize = vector.FromFloats(screenW, screenH)

	// setup ui callbacks
	g.UI.OnClearTrail(func() {
		g.showTrail = !g.showTrail
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

	// bodiesMap := make(map[int]*body.Body, len(c.Bodies))
	bodies := make([]*body.Body, 0, len(c.Bodies))
	for _, conf := range c.Bodies {
		body, err := body.FromConfig(conf)
		if err != nil {
			log.Fatalf("failed to load body from config: %s", err)
		}
		// bodiesMap[body.Id] = body
		bodies = append(bodies, body)
	}
	g.Bodies = bodies

	return g
}
