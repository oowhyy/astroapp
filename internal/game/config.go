package game

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/body"
	"github.com/oowhyy/astroapp/pkg/camera"
	"github.com/oowhyy/astroapp/pkg/dropbox"
	"github.com/oowhyy/astroapp/pkg/tilemap"
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

func FromConfig(c *Config, client *dropbox.Client) (*Game, error) {
	g := &Game{
		// parentMap: make(map[string]string, len(c.Bodies)),
	}
	g.UI = webui.NewWebInterface()

	g.GConstant = c.GConstant
	g.simSpeed = 1

	// fetch images form dropbox
	assetsZip, err := client.FetchZip("/assets.zip")
	if err != nil {
		return nil, err
	}
	assets := map[string]io.Reader{}
	for _, f := range assetsZip.File {
		file, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open file error: %w", err)
		}
		buff := new(bytes.Buffer)
		_, err = buff.ReadFrom(file)
		if err != nil {
			return nil, fmt.Errorf("buffer ReadFrom error: %w", err)
		}
		assets[f.Name] = buff
		file.Close()
	}

	// fetch background from dropbox
	tiles, err := tilemap.NewTileMapFromDropboxZip(client, "/bg_tiles.zip")
	if err != nil {
		panic(err)
	}
	g.background = tiles
	worldX, worldY := 10000, 5625
	g.worldSize = image.Pt(worldX, worldY)
	g.trailLayer = tilemap.NewTileMapEmpty(tiles.Depth(), tiles.TileSize(), worldX, worldY)
	fmt.Println("tilemap ok")

	// game other assets
	arrowPng, err := png.Decode(assets["blue_arrow.png"])
	if err != nil {
		return nil, err
	}
	g.blueArrow = ebiten.NewImageFromImage(arrowPng)
	rockPng, err := png.Decode(assets["rock.png"])
	if err != nil {
		return nil, err
	}
	g.rock = ebiten.NewImageFromImage(rockPng)

	// center camera
	screenW, screenH := g.UI.WindowSize()
	g.Camera = camera.NewCamera(image.Pt(worldX, worldY), image.Pt(int(screenW), int(screenH)))
	g.Camera.Translate((float64(worldX)-screenW)/2, (float64(worldY)-screenH)/2)

	// setup ui callbacks
	g.showTrail = true
	g.UI.OnClearTrail(func() {
		g.showTrail = !g.showTrail
		if g.showTrail {
			g.trailLayer.Clear()
		}
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
		imgName := conf.Name + ".png"
		body, err := body.FromConfig(conf, assets[imgName])
		if err != nil {
			return nil, fmt.Errorf("failed to load body from config: %w", err)
		}
		// bodiesMap[body.Id] = body
		bodies = append(bodies, body)
	}
	g.Bodies = bodies

	return g, nil
}
