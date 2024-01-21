package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/game"
	"github.com/oowhyy/astroapp/pkg/dropbox"
	"github.com/oowhyy/astroapp/pkg/webloader"
	"gopkg.in/yaml.v3"
)

var refreshToken string
var appAuth string

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Render an image")
	cfg := &game.Config{}
	err := readConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	accessToken, err := dropbox.GetAccessToken(refreshToken, appAuth)
	if err != nil {
		panic(err)
	}
	fmt.Println("token len", len(accessToken))
	DBcfg := dropbox.NewConfig(accessToken)
	client := dropbox.New(DBcfg)
	g,err := game.FromConfig(cfg, client)
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func readConfig(cfg *game.Config) error {
	yamlString := webloader.LoadFile("config.yaml")
	if yamlString == "" {
		log.Fatal("no yaml string from loader")
	}
	err := yaml.Unmarshal([]byte(yamlString), cfg)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
