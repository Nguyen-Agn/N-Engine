package napi

import (
	"autoworld/modules/components"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// LoadFont loads a TrueType/OpenType font from a file path and returns a font.Face
// at the given size (in pixels). Use the returned face with TextWithFont().
//
// Example:
//
//	face, err := napi.LoadFont("assets/fonts/myfont.ttf", 16)
//	if err != nil { log.Fatal(err) }
//	// inside Draw():
//	o.TextWithFont("Hello!", x, y, color.White, face)
func LoadFont(path string, size float64) (font.Face, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}
	return face, nil
}

// SetDefaultFont sets the engine-wide default font used by DrawComponent.Text().
// If not called, the engine defaults to basicfont.Face7x13 (a small bitmap font).
// Call this once during initialisation, before the game loop starts.
func SetDefaultFont(face font.Face) {
	components.SetPackageDefaultFont(face)
}
