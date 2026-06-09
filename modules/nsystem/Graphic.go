package nsystem

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawOptions is the pre-calculated result, ready for rendering.
// Graphic.go is responsible for creating this struct from SpriteData + PositionData.
type DrawOptions struct {
	// Image is the specific frame image to be drawn.
	Image *ebiten.Image
	// Opts contains all the geometric transformation and color parameters.
	Opts *ebiten.DrawImageOptions
}

// BuildDrawOptions converts SpriteData and PositionData into an array of DrawOptions ready for drawing.
// Inputs: 
//   pos (PositionData) - The position component of the entity.
//   spr (SpriteData) - The sprite component of the entity.
//   camX, camY (float32) - Camera coordinates in map space.
// Outputs: Returns an array of DrawOptions. Returns nil if there is no valid sprite to draw.
// Purpose: Calculates final screen position, handles rotation, scaling, NineSlice (if enabled), and color scaling.
func BuildDrawOptions(pos PositionData, spr SpriteData, camX, camY float32) []*DrawOptions {
	// Retrieve the currently selected sprite
	spriteLW, ok := spr.Sprite[spr.CurrentSprite]
	if !ok || spriteLW == nil {
		return nil
	}

	// Retrieve the current image frame by index
	img := spriteLW.Image(spr.SpriteIdx)
	if img == nil {
		return nil
	}

	if spr.IsNineSlice {
		return buildNineSliceDrawOptions(img, pos, spr, camX, camY)
	}

	opts := &ebiten.DrawImageOptions{}

	// 1. Translate to origin (0,0) so rotation and scaling work correctly around the center
	w := float64(spriteLW.Width())
	h := float64(spriteLW.Height())
	opts.GeoM.Translate(-w/2, -h/2)

	// 2. Apply Scale (Negative ScaleX = flip horizontally, Negative ScaleY = flip vertically)
	opts.GeoM.Scale(float64(spr.ScaleX), float64(spr.ScaleY))

	// 3. Apply Rotation (convert degrees to radians)
	opts.GeoM.Rotate(float64(spr.Rotation) * math.Pi / 180)

	// 4. Translate to final position = (map pos + offset) - camera offset = screen pos
	finalX := float64(pos.X) + float64(spr.OffsetX) - float64(camX)
	finalY := float64(pos.Y) + float64(spr.OffsetY) - float64(camY)
	opts.GeoM.Translate(finalX, finalY)

	// 5. Apply Color Scale
	opts.ColorScale.ScaleWithColor(toFloat32Color(spr.ImageColor))

	return []*DrawOptions{
		{
			Image: img,
			Opts:  opts,
		},
	}
}

// buildNineSliceDrawOptions creates drawing options for a 9-slice scaled sprite.
// Inputs: 
//   img (*ebiten.Image) - The source image to slice.
//   pos (PositionData) - Entity's position.
//   spr (SpriteData) - Entity's sprite data containing nine-slice configuration.
//   camX, camY (float32) - Camera coordinates in map space.
// Outputs: Returns a slice of DrawOptions for each visible patch of the 9-slice.
// Purpose: Divides the image into 9 patches based on NineSlice margins, scales the middle patches appropriately to reach the target scaled size without distorting corners, and applies rotation/positioning to the entire assembled shape.
func buildNineSliceDrawOptions(img *ebiten.Image, pos PositionData, spr SpriteData, camX, camY float32) []*DrawOptions {
	var result []*DrawOptions

	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()

	top := spr.NineSlice.Top
	bottom := spr.NineSlice.Bottom
	left := spr.NineSlice.Left
	right := spr.NineSlice.Right

	// Safety checks for slices
	if left+right > srcW {
		left = srcW / 2
		right = srcW - left
	}
	if top+bottom > srcH {
		top = srcH / 2
		bottom = srcH - top
	}

	targetW := float64(srcW) * float64(spr.ScaleX)
	targetH := float64(srcH) * float64(spr.ScaleY)

	srcX := []int{0, left, srcW - right, srcW}
	srcY := []int{0, top, srcH - bottom, srcH}

	midW := targetW - float64(left+right)
	if midW < 0 {
		midW = 0
	}
	midH := targetH - float64(top+bottom)
	if midH < 0 {
		midH = 0
	}

	dstX := []float64{0, float64(left), float64(left) + midW, targetW}
	dstY := []float64{0, float64(top), float64(top) + midH, targetH}

	finalX := float64(pos.X) + float64(spr.OffsetX) - float64(camX)
	finalY := float64(pos.Y) + float64(spr.OffsetY) - float64(camY)

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			sw := srcX[col+1] - srcX[col]
			sh := srcY[row+1] - srcY[row]

			if sw <= 0 || sh <= 0 {
				continue
			}

			dw := dstX[col+1] - dstX[col]
			dh := dstY[row+1] - dstY[row]

			if dw <= 0 || dh <= 0 {
				continue
			}

			rect := image.Rect(srcX[col], srcY[row], srcX[col+1], srcY[row+1])
			subImg := img.SubImage(rect).(*ebiten.Image)

			opts := &ebiten.DrawImageOptions{}

			// Scale patch
			scaleX := dw / float64(sw)
			scaleY := dh / float64(sh)
			opts.GeoM.Scale(scaleX, scaleY)

			// Translate patch to correct coordinate within the local frame (0, 0)
			opts.GeoM.Translate(dstX[col], dstY[row])

			// Translate center so it can be rotated around the middle of the frame
			opts.GeoM.Translate(-targetW/2, -targetH/2)

			// Rotate the entire block
			opts.GeoM.Rotate(float64(spr.Rotation) * math.Pi / 180)

			// Translate to final position
			opts.GeoM.Translate(finalX, finalY)

			// Apply Color Scale
			opts.ColorScale.ScaleWithColor(toFloat32Color(spr.ImageColor))

			result = append(result, &DrawOptions{
				Image: subImg,
				Opts:  opts,
			})
		}
	}

	return result
}

// toFloat32Color passes through the color.RGBA, as ebiten supports it directly.
// Inputs: c (color.RGBA) - The color to be used.
// Outputs: Returns the exact same color.RGBA object.
func toFloat32Color(c color.RGBA) color.RGBA {
	return c // ebiten.ColorScale.ScaleWithColor accepts color.RGBA directly
}
