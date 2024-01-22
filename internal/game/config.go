package game

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"time"

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
	WorldWidth    int                `yaml:"worldWidth"`
	WorldHeight   int                `yaml:"worldHeight"`
}

// const (
// 	WorldWidth  = 4000
// 	WorldHeight = 4000
// )

func FromConfig(c *Config, client *dropbox.Client) (*Game, error) {
	g := &Game{
		UI:        webui.NewWebInterface(),
		GConstant: c.GConstant,
		simSpeed:  1,
		worldSize: image.Pt(c.WorldWidth, c.WorldHeight),
	}
	// timer := timeit.NewTimer(func(lapName string, timedDuration time.Duration) {
	// 	fmt.Println(lapName, "in", timedDuration)
	// })

	// fetch images form dropbox
	g.UI.SetLoadingMessage("fetching assets")
	assets, err := fetchAssets(client)
	if err != nil {
		return nil, err
	}
	fmt.Println("FETCHED ASSETS")
	// fetch background from dropbox
	g.UI.SetLoadingMessage("fetching background")
	zip, err := client.FetchZip("/bg_tiles.zip")
	if err != nil {
		return nil, err
	}
	g.UI.SetLoadingMessage("building tilemap")
	tiles, err := tilemap.NewTileMapFromZip(zip)
	if err != nil {
		return nil, err
	}
	g.background = tiles
	// timer.Lap("background loaded")

	g.UI.SetLoadingMessage("building trail layer")
	g.trailLayer = tilemap.NewTileMapEmpty(tiles.Depth(), tiles.TileSize(), c.WorldWidth, c.WorldHeight)
	// timer.Lap("empty layer built")
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
	g.Camera = camera.NewCamera(g.worldSize, image.Pt(int(screenW), int(screenH)))
	g.Camera.Translate((float64(c.WorldWidth)-screenW)/2, (float64(c.WorldHeight)-screenH)/2)

	// setup ui callbacks
	g.showTrail = true
	g.UI.OnClearTrail(func() {
		tic := time.Now()
		g.showTrail = !g.showTrail
		if !g.showTrail {
			g.trailLayer.Clear()
			fmt.Println("CLEARED TRAIL IN", time.Since(tic).Seconds())
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
	g.UI.DoneLoading()
	return g, nil
}

func fetchAssets(client *dropbox.Client) (map[string]io.Reader, error) {
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
	return assets, nil
}
