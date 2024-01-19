package tilemap

import (
	"bytes"
	"image"
	"image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	decoded  *ebiten.Image
	encoded  *ImageReader
	bounds   image.Rectangle
	z        int
	x        int
	y        int
	children []*Tile
}

func (t *Tile) Children() []*Tile {
	return t.children
}

func (t *Tile) Image() *ebiten.Image {
	if t.decoded == nil {
		if t.encoded == nil {
			t.decoded = ebiten.NewImage(100, 100)
		} else {
			t.decoded = ebiten.NewImageFromImage(t.encoded.Decoded())
		}
	}
	return t.decoded
}

func (t *Tile) Clear() {
	if t.decoded != nil {
		t.decoded.Clear()
	}
}

func (t *Tile) X() int {
	return t.x
}

func (t *Tile) Y() int {
	return t.y
}

func (t *Tile) Z() int {
	return t.z
}

func (t *Tile) Bounds() image.Rectangle {
	return t.bounds
}

type ImageReader struct {
	encoded io.Reader
}

func NewImageReader(encoded io.Reader) *ImageReader {
	return &ImageReader{
		encoded: encoded,
	}
}

func (ir *ImageReader) Decoded() image.Image {
	if ir.encoded == nil {
		return image.NewAlpha(image.Rect(0, 0, 100, 100))
	}
	img, tmp := saveDecode(ir.encoded)
	ir.encoded = tmp
	return img
}

func saveDecode(reader io.Reader) (image.Image, io.Reader) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	img, err := png.Decode(tee)
	if err != nil {
		panic("png decode error")
	}
	return img, &buf
}
