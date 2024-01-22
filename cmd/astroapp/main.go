package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/oowhyy/astroapp/internal/game"
	"github.com/oowhyy/astroapp/pkg/dropbox"
	"gopkg.in/yaml.v3"
)

var refreshToken string
var appAuth string

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Render an image")

	accessToken, err := dropbox.GetAccessToken(refreshToken, appAuth)
	if err != nil {
		panic(err)
	}
	fmt.Println("token len", len(accessToken))
	DBcfg := dropbox.NewConfig(accessToken)
	client := dropbox.New(DBcfg)
	cfg := &game.Config{}

	err = readConfig(client, cfg)
	if err != nil {
		panic(err)
	}
	g, err := game.FromConfig(cfg, client)
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func readConfig(client *dropbox.Client, cfg *game.Config) error {
	file, err := client.FetchFile("/config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}
