package components

import (
	"image"
	"image/color"
	"strings"

	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/hajimehoshi/ebiten/v2"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// DrawComponent is a Mixin embedded in Custom Objects to provide primitive drawing.
// Token: "drw" — automatically pulls in "pos" via constraint.
//
// The engine binds this mixin automatically via napi.NewObject reflection.
// Call drawing methods inside the object's Draw() method.
//
// All coordinates are in map space — camera offset is applied automatically.
// DrawSystem injects Screen, CamX, CamY into DrawData before calling Draw() each frame.
type DrawComponent struct {
	IObject
	data        *DrawData
	defaultFont font.Face // per-instance font override; nil = use packageDefaultFont

	// Text rendering states
	hAlign       textv2.Align
	vAlign       textv2.Align
	isJustify    bool
	maxWidth     float32
	maxHeight    float32
	overflowMode string // "visible", "hidden", "scale"
}

// packageDefaultFont is the engine-wide default font for DrawComponent.Text().
// Override with napi.SetDefaultFont() or per-instance with DrawComponent.SetFont().
var packageDefaultFont font.Face = basicfont.Face7x13

// SetPackageDefaultFont replaces the engine-wide default font.
// Purpose: Sets the global font. Intended to be called only by napi.SetDefaultFont().
// Inputs:
//   - f: The font.Face to set as default.
//
// Outputs: None.
func SetPackageDefaultFont(f font.Face) {
	packageDefaultFont = f
}

// Draw component type reference.
var Draw = enginetype.Draw

func init() {
	enginetype.RegisterComponentInitializer("drw", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Draw, DrawData{})
	})
}

// BindComponent is called by the engine reflection binder to wire this mixin to its ECS entry.
// Purpose: Binds the component to the ECS entry.
// Inputs:
//   - base: The IObject representing the base entity.
//
// Outputs: None.
func (d *DrawComponent) BindComponent(base IObject) {
	d.IObject = base
	d.data = enginetype.GetComponent(base, Draw)
}

// SetFont overrides the default font for this DrawComponent instance only.
// Purpose: Sets a specific font for this component instance.
// Inputs:
//   - f: The font.Face to use for this component.
//
// Outputs: None.
func (d *DrawComponent) SetFont(f font.Face) {
	d.defaultFont = f
}

// --- internal helpers ---

// screen retrieves the target screen image for drawing.
// Purpose: Gets the underlying ebiten.Image screen.
// Inputs: None.
// Outputs: Pointer to the ebiten.Image screen, or nil if data is unset.
func (d *DrawComponent) screen() *ebiten.Image {
	if d.data == nil {
		return nil
	}
	return d.data.Screen
}

// toScreen converts map-space (x, y) to screen-space by subtracting camera offset.
// Purpose: Translates game world coordinates into screen coordinates.
// Inputs:
//   - x: Map-space X coordinate.
//   - y: Map-space Y coordinate.
//
// Outputs: A tuple (float32, float32) representing screen-space X and Y coordinates.
func (d *DrawComponent) toScreen(x, y float32) (float32, float32) {
	if d.data == nil {
		return x, y
	}
	return x - d.data.CamX, y - d.data.CamY
}

// activeFont retrieves the current active font for text rendering.
// Purpose: Returns the component-specific font or falls back to the package default.
// Inputs: None.
// Outputs: The font.Face to use.
func (d *DrawComponent) activeFont() font.Face {
	if d.defaultFont != nil {
		return d.defaultFont
	}
	return packageDefaultFont
}

// =============================================================================
// Coordinate helpers
// =============================================================================

// ScreenX converts a map-space X coordinate to screen-space X.
// Purpose: Useful when building a custom vector.Path before passing it to PathFill/PathStroke.
// Inputs:
//   - mapX: Map-space X coordinate.
//
// Outputs: Screen-space X coordinate.
func (d *DrawComponent) ScreenX(mapX float32) float32 {
	if d.data == nil {
		return mapX
	}
	return mapX - d.data.CamX
}

// ScreenY converts a map-space Y coordinate to screen-space Y.
// Purpose: Useful when building a custom vector.Path before passing it to PathFill/PathStroke.
// Inputs:
//   - mapY: Map-space Y coordinate.
//
// Outputs: Screen-space Y coordinate.
func (d *DrawComponent) ScreenY(mapY float32) float32 {
	if d.data == nil {
		return mapY
	}
	return mapY - d.data.CamY
}

// =============================================================================
// Filled shapes
// =============================================================================

// Rect draws a filled rectangle at map-space (x, y) with size (w, h).
// Purpose: Draws a solid rectangle.
// Inputs:
//   - x, y: Map-space coordinates for the top-left corner.
//   - w, h: Width and height of the rectangle.
//   - c: The fill color.
//
// Outputs: None.
func (d *DrawComponent) Rect(x, y, w, h float32, c color.RGBA) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.FillRect(s, sx, sy, w, h, c, false)
}

// Circle draws a filled circle centered at map-space (x, y) with radius r.
// Purpose: Draws a solid circle.
// Inputs:
//   - x, y: Map-space coordinates for the center.
//   - r: Radius of the circle.
//   - c: The fill color.
//
// Outputs: None.
func (d *DrawComponent) Circle(x, y, r float32, c color.RGBA) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.FillCircle(s, sx, sy, r, c, true)
}

// =============================================================================
// Stroke (outline) shapes
// =============================================================================

// RectStroke draws a rectangle outline at map-space (x, y) with size (w, h).
// Purpose: Draws a rectangle outline.
// Inputs:
//   - x, y: Map-space coordinates for the top-left corner.
//   - w, h: Width and height of the rectangle.
//   - c: The outline color.
//   - strokeWidth: Thickness of the border in pixels.
//
// Outputs: None.
func (d *DrawComponent) RectStroke(x, y, w, h float32, c color.RGBA, strokeWidth float32) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.StrokeRect(s, sx, sy, w, h, strokeWidth, c, false)
}

// CircleStroke draws a circle outline centered at map-space (x, y) with radius r.
// Purpose: Draws a circle outline.
// Inputs:
//   - x, y: Map-space coordinates for the center.
//   - r: Radius of the circle.
//   - c: The outline color.
//   - strokeWidth: Thickness of the border in pixels.
//
// Outputs: None.
func (d *DrawComponent) CircleStroke(x, y, r float32, c color.RGBA, strokeWidth float32) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.StrokeCircle(s, sx, sy, r, strokeWidth, c, true)
}

// Line draws a straight line from map-space (x0, y0) to (x1, y1).
// Purpose: Draws a line segment.
// Inputs:
//   - x0, y0: Map-space coordinates for the start point.
//   - x1, y1: Map-space coordinates for the end point.
//   - c: The line color.
//   - strokeWidth: Thickness of the line in pixels.
//
// Outputs: None.
func (d *DrawComponent) Line(x0, y0, x1, y1 float32, c color.RGBA, strokeWidth float32) {
	s := d.screen()
	if s == nil {
		return
	}
	sx0, sy0 := d.toScreen(x0, y0)
	sx1, sy1 := d.toScreen(x1, y1)
	vector.StrokeLine(s, sx0, sy0, sx1, sy1, strokeWidth, c, true)
}

// =============================================================================
// Custom path — build path in screen space using ScreenX()/ScreenY()
// =============================================================================

// PathFill fills an arbitrary vector.Path with a color.
// Purpose: Fills a custom shape path.
// Inputs:
//   - p: Pointer to the vector.Path built in screen-space coordinates.
//   - c: The fill color.
//   - antialias: Whether to enable antialiasing.
//
// Outputs: None.
func (d *DrawComponent) PathFill(p *vector.Path, c color.RGBA, antialias bool) {
	s := d.screen()
	if s == nil || p == nil {
		return
	}
	op := &vector.DrawPathOptions{AntiAlias: antialias}
	op.ColorScale.ScaleWithColor(c)
	vector.FillPath(s, p, nil, op)
}

// PathStroke strokes an arbitrary vector.Path with a color.
// Purpose: Draws the outline of a custom shape path.
// Inputs:
//   - p: Pointer to the vector.Path built in screen-space coordinates.
//   - c: The stroke color.
//   - strokeWidth: Thickness of the stroke.
//   - antialias: Whether to enable antialiasing.
//
// Outputs: None.
func (d *DrawComponent) PathStroke(p *vector.Path, c color.RGBA, strokeWidth float32, antialias bool) {
	s := d.screen()
	if s == nil || p == nil {
		return
	}
	strokeOp := &vector.StrokeOptions{Width: strokeWidth}
	drawOp := &vector.DrawPathOptions{AntiAlias: antialias}
	drawOp.ColorScale.ScaleWithColor(c)
	vector.StrokePath(s, p, strokeOp, drawOp)
}

// =============================================================================
// Text
// =============================================================================

// parseAlign converts a string token to textv2.Align
// Purpose: Parses a string alignment token into an ebiten text alignment enum.
// Inputs:
//   - token: The string token (e.g., "center", "left").
//
// Outputs: A tuple containing the textv2.Align and a boolean indicating if it is justified.
func parseAlign(token string) (textv2.Align, bool) {
	switch token {
	case "start", "left", "top", "l", "t":
		return textv2.AlignStart, false
	case "center", "middle", "c", "m":
		return textv2.AlignCenter, false
	case "end", "right", "bottom", "r", "b":
		return textv2.AlignEnd, false
	case "justify", "j":
		return textv2.AlignStart, true // We use Start alignment and handle justify manually
	default:
		return textv2.AlignStart, false
	}
}

// SetTextAlign sets the text alignment for subsequent text drawing commands.
// Purpose: Configures text alignment.
// Inputs:
//   - hAlign: Horizontal alignment string token.
//   - vAlign: Vertical alignment string token.
//
// Outputs: None.
func (d *DrawComponent) SetTextAlign(hAlign, vAlign string) {
	d.hAlign, d.isJustify = parseAlign(hAlign)
	d.vAlign, _ = parseAlign(vAlign)
}

// SetTextOverflow sets the bounding box for text and the overflow handling mode.
// Purpose: Configures text overflow behavior.
// Inputs:
//   - maxWidth: Maximum width of the text box.
//   - maxHeight: Maximum height of the text box.
//   - mode: The overflow mode ("visible", "hidden", "scale").
//
// Outputs: None.
func (d *DrawComponent) SetTextOverflow(maxWidth, maxHeight float32, mode string) {
	d.maxWidth = maxWidth
	d.maxHeight = maxHeight

	switch mode {
	case "hidden", "h":
		d.overflowMode = "hidden"
	case "scale", "s":
		d.overflowMode = "scale"
	default:
		d.overflowMode = "visible"
	}
}

// Text draws a string at map-space (x, y) using the active default font.
// Purpose: Renders text to the screen.
// Inputs:
//   - text: The string to draw.
//   - x, y: Map-space coordinates for the text.
//   - c: The color of the text.
//
// Outputs: None.
func (d *DrawComponent) Text(text string, x, y float32, c color.RGBA) {
	d.drawTextInternal(text, x, y, c, 1.0)
}

// TextEx draws a string at map-space (x, y) with a uniform scale applied to the font.
// Purpose: Renders scaled text to the screen.
// Inputs:
//   - text: The string to draw.
//   - x, y: Map-space coordinates for the text.
//   - c: The color of the text.
//   - scale: The scaling factor (1.0 is default size).
//
// Outputs: None.
func (d *DrawComponent) TextEx(text string, x, y float32, c color.RGBA, scale float64) {
	d.drawTextInternal(text, x, y, c, scale)
}

// drawTextInternal handles the actual text rendering logic including alignment, scaling, and clipping.
// Purpose: Internal text drawing implementation.
// Inputs:
//   - textContent: The string to draw.
//   - x, y: Map-space coordinates.
//   - c: The text color.
//   - scale: The scaling factor.
//
// Outputs: None.
func (d *DrawComponent) drawTextInternal(textContent string, x, y float32, c color.RGBA, scale float64) {
	s := d.screen()
	if s == nil || textContent == "" {
		return
	}
	sx, sy := d.toScreen(x, y)
	face := textv2.NewGoXFace(d.activeFont())

	op := &textv2.DrawOptions{}
	op.PrimaryAlign = d.hAlign
	op.SecondaryAlign = d.vAlign
	op.GeoM.Scale(scale, scale)

	// 1. Xử lý Overflow Scale
	var minScale float64 = 1.0
	if (d.maxWidth > 0 || d.maxHeight > 0) && d.overflowMode == "scale" {
		w, h := textv2.Measure(textContent, face, face.Metrics().HAscent)
		w *= scale
		h *= scale

		scaleX, scaleY := 1.0, 1.0
		if d.maxWidth > 0 && w > float64(d.maxWidth) {
			scaleX = float64(d.maxWidth) / w
		}
		if d.maxHeight > 0 && h > float64(d.maxHeight) {
			scaleY = float64(d.maxHeight) / h
		}

		minScale = scaleX
		if scaleY < minScale {
			minScale = scaleY
		}
		if minScale < 1.0 {
			op.GeoM.Scale(minScale, minScale)
		}
	}

	// 2. Xử lý Overflow Hidden (Clipping)
	var target *ebiten.Image = s
	if d.overflowMode == "hidden" && (d.maxWidth > 0 || d.maxHeight > 0) {
		clipX, clipY := float64(sx), float64(sy)

		if d.maxWidth > 0 {
			if d.hAlign == textv2.AlignCenter {
				clipX -= float64(d.maxWidth) / 2
			}
			if d.hAlign == textv2.AlignEnd {
				clipX -= float64(d.maxWidth)
			}
		} else {
			clipX = 0 // Không giới hạn
		}

		if d.maxHeight > 0 {
			if d.vAlign == textv2.AlignCenter {
				clipY -= float64(d.maxHeight) / 2
			}
			if d.vAlign == textv2.AlignEnd {
				clipY -= float64(d.maxHeight)
			}
		} else {
			clipY = 0
		}

		cw := float64(d.maxWidth)
		if cw <= 0 {
			cw = 99999
		} // Rất lớn nếu không giới hạn
		ch := float64(d.maxHeight)
		if ch <= 0 {
			ch = 99999
		}

		rect := image.Rect(int(clipX), int(clipY), int(clipX+cw), int(clipY+ch))
		// Cắt giới hạn vẽ (Clipping)
		target = s.SubImage(rect.Intersect(s.Bounds())).(*ebiten.Image)
	}

	op.GeoM.Translate(float64(sx), float64(sy))
	op.ColorScale.ScaleWithColor(c)

	// 3. Xử lý Justify
	if d.isJustify && d.maxWidth > 0 && d.overflowMode != "scale" {
		words := strings.Fields(textContent)
		if len(words) > 1 {
			totalWordWidth := 0.0
			for _, w := range words {
				ww, _ := textv2.Measure(w, face, face.Metrics().HAscent)
				totalWordWidth += ww * minScale
			}

			spaceWidth := (float64(d.maxWidth) - totalWordWidth) / float64(len(words)-1)
			if spaceWidth < 0 {
				spaceWidth = 0
			} // Tránh đè chữ

			currX := float64(sx)
			if d.hAlign == textv2.AlignCenter {
				currX -= float64(d.maxWidth) / 2
			}
			if d.hAlign == textv2.AlignEnd {
				currX -= float64(d.maxWidth)
			}

			for _, w := range words {
				wOp := &textv2.DrawOptions{}
				wOp.PrimaryAlign = textv2.AlignStart
				wOp.SecondaryAlign = d.vAlign
				wOp.GeoM.Scale(minScale, minScale)
				wOp.GeoM.Translate(currX, float64(sy))
				wOp.ColorScale.ScaleWithColor(c)
				textv2.Draw(target, w, face, wOp)

				ww, _ := textv2.Measure(w, face, face.Metrics().HAscent)
				currX += (ww * minScale) + spaceWidth
			}
			return
		}
	}

	textv2.Draw(target, textContent, face, op)
}

// =============================================================================
// Image
// =============================================================================

// Image draws frame idx of the given ISpriteLW at map-space (x, y).
// Purpose: Allows manual sprite rendering independent of SpriteComponent.
// Inputs:
//   - sprite: The ISpriteLW instance to draw.
//   - idx: The frame index of the sprite to draw.
//   - x, y: Map-space coordinates.
//
// Outputs: None.
func (d *DrawComponent) Image(sprite ISpriteLW, idx int, x, y float32) {
	s := d.screen()
	if s == nil || sprite == nil {
		return
	}
	img := sprite.Image(idx)
	if img == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(sx), float64(sy))
	s.DrawImage(img, opts)
}
