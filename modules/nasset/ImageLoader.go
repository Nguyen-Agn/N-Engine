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
	// Purpose: Decodes an image from a reader.
	// Inputs: r (io.Reader) - The input stream containing the image data.
	// Outputs: (image.Image) - The decoded image, (error) - Error if decoding fails.
	Decode(r io.Reader) (image.Image, error)
}

// PNGDecoder implements IImageDecoder for .png
type PNGDecoder struct{}

// Purpose: Decodes a PNG image.
// Inputs: r (io.Reader) - The input stream.
// Outputs: (image.Image) - The decoded image, (error) - Error if decoding fails.
func (d *PNGDecoder) Decode(r io.Reader) (image.Image, error) { return png.Decode(r) }

// JPGDecoder implements IImageDecoder for .jpg/.jpeg
type JPGDecoder struct{}

// Purpose: Decodes a JPG/JPEG image.
// Inputs: r (io.Reader) - The input stream.
// Outputs: (image.Image) - The decoded image, (error) - Error if decoding fails.
func (d *JPGDecoder) Decode(r io.Reader) (image.Image, error) { return jpeg.Decode(r) }

// GIFDecoder implements IImageDecoder for .gif
type GIFDecoder struct{}

// Purpose: Decodes a GIF image.
// Inputs: r (io.Reader) - The input stream.
// Outputs: (image.Image) - The decoded image, (error) - Error if decoding fails.
func (d *GIFDecoder) Decode(r io.Reader) (image.Image, error) { return gif.Decode(r) }

// SpriteLoader implements domain.ISpriteLoader
type SpriteLoader struct {
	decoders map[string]IImageDecoder
}

// Purpose: Creates a new SpriteLoader with default decoders for common image formats.
// Outputs: (*SpriteLoader) - The newly created loader.
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

// Purpose: Registers a custom image decoder for a specific file extension.
// Inputs: ext (string) - The file extension (e.g. ".webp"), decoder (IImageDecoder) - The decoder strategy.
func (l *SpriteLoader) RegisterDecoder(ext string, decoder IImageDecoder) {
	l.decoders[ext] = decoder
}

// Purpose: Internal helper to read and decode an image from a given file path.
// Inputs: path (string) - The path to the image file.
// Outputs: (image.Image) - The decoded image, (error) - Error if file reading or decoding fails.
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

// Purpose: Loads a single image file and returns it as a Sprite.
// Inputs: path (string) - The path to the image file.
// Outputs: (ISpriteLW) - The loaded sprite, (error) - Error if loading fails.
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

// Purpose: Loads a spritesheet image and slices it into individual frames.
// Inputs: path (string) - The image path, frameW (int) - Frame width, frameH (int) - Frame height, cols (int) - Number of columns, rows (int) - Number of rows, gapX (int) - Horizontal gap between frames, gapY (int) - Vertical gap.
// Outputs: (ISpriteLW) - The loaded sprite containing all frames, (error) - Error if loading fails.
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
