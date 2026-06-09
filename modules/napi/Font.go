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
// Purpose: Loads external font files so they can be rendered by the draw system.
//
// Inputs:
// - path (string): The file system path to the font file (e.g., "assets/fonts/myfont.ttf").
// - size (float64): The desired size of the font in pixels.
//
// Outputs:
// - font.Face: The parsed font face ready for rendering.
// - error: An error if the file cannot be read or parsed.
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
//
// Purpose: Establishes a fallback font for text rendering across the engine. If not set, the engine defaults to basicfont.Face7x13.
//
// Inputs:
// - face (font.Face): The new default font face to use.
//
// Special Requirements:
// - Should be called once during initialization before the game loop starts.
func SetDefaultFont(face font.Face) {
	components.SetPackageDefaultFont(face)
}
