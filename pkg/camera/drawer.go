package camera

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawer interface {
	InWorldBounds(worldWidth, worldHeight int) image.Rectangle
	SpriteOp() (sprite *ebiten.Image, scale float64, tx, ty float64)
}
