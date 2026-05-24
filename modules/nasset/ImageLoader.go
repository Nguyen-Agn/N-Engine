package nasset

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"autoworld/domain/bridge"

	"github.com/hajimehoshi/ebiten/v2"
)

// ReadFileFn is a wrapper for reading files to allow easy mocking or switching to embed.FS.
var ReadFileFn = func(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// IImageDecoder defines the strategy for decoding different image formats.
type IImageDecoder interface {
	Decode(r io.Reader) (image.Image, error)
}

// PNGDecoder implements IImageDecoder for .png
type PNGDecoder struct{}

func (d *PNGDecoder) Decode(r io.Reader) (image.Image, error) { return png.Decode(r) }

// JPGDecoder implements IImageDecoder for .jpg/.jpeg
type JPGDecoder struct{}

func (d *JPGDecoder) Decode(r io.Reader) (image.Image, error) { return jpeg.Decode(r) }

// GIFDecoder implements IImageDecoder for .gif
type GIFDecoder struct{}

func (d *GIFDecoder) Decode(r io.Reader) (image.Image, error) { return gif.Decode(r) }

// SpriteLoader implements domain.ISpriteLoader
type SpriteLoader struct {
	decoders map[string]IImageDecoder
}

// NewSpriteLoader creates a new loader with default decoders.
func NewSpriteLoader() *SpriteLoader {
	return &SpriteLoader{
		decoders: map[string]IImageDecoder{
			".png":  &PNGDecoder{},
			".jpg":  &JPGDecoder{},
			".jpeg": &JPGDecoder{},
			".gif":  &GIFDecoder{},
		},
	}
}

// RegisterDecoder allows adding custom decoders for new formats.
func (l *SpriteLoader) RegisterDecoder(ext string, decoder IImageDecoder) {
	l.decoders[ext] = decoder
}

func (l *SpriteLoader) decodeImage(path string) (image.Image, error) {
	file, err := ReadFileFn(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	decoder, ok := l.decoders[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	img, err := decoder.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// LoadSingle loads a single image file into a SpriteLW.
func (l *SpriteLoader) LoadSingle(path string) (ISpriteLW, error) {
	img, err := l.decodeImage(path)
	if err != nil {
		return nil, err
	}

	eImg := ebiten.NewImageFromImage(img)
	bounds := eImg.Bounds()

	sprite := bridge.NewSpriteLW(bounds.Dx(), bounds.Dy())
	sprite.AddImage(eImg)

	return sprite, nil
}

// LoadSheet loads a spritesheet and cuts it into frames.
func (l *SpriteLoader) LoadSheet(path string, frameW, frameH, cols, rows, gapX, gapY int) (ISpriteLW, error) {
	img, err := l.decodeImage(path)
	if err != nil {
		return nil, err
	}

	eImg := ebiten.NewImageFromImage(img)
	sprite := bridge.NewSpriteLW(frameW, frameH)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			startX := col * (frameW + gapX)
			startY := row * (frameH + gapY)

			rect := image.Rect(startX, startY, startX+frameW, startY+frameH)
			frame := eImg.SubImage(rect).(*ebiten.Image)
			sprite.AddImage(frame)
		}
	}

	return sprite, nil
}
