package components

import (
	"image/color"

	"autoworld/modules/enginetype"

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
}

// packageDefaultFont is the engine-wide default font for DrawComponent.Text().
// Override with napi.SetDefaultFont() or per-instance with DrawComponent.SetFont().
var packageDefaultFont font.Face = basicfont.Face7x13

// SetPackageDefaultFont replaces the engine-wide default font.
// Intended to be called only by napi.SetDefaultFont().
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
func (d *DrawComponent) BindComponent(base IObject) {
	d.IObject = base
	d.data = enginetype.GetComponent(base, Draw)
}

// SetFont overrides the default font for this DrawComponent instance only.
func (d *DrawComponent) SetFont(f font.Face) {
	d.defaultFont = f
}

// --- internal helpers ---

func (d *DrawComponent) screen() *ebiten.Image {
	if d.data == nil {
		return nil
	}
	return d.data.Screen
}

// toScreen converts map-space (x, y) to screen-space by subtracting camera offset.
func (d *DrawComponent) toScreen(x, y float32) (float32, float32) {
	if d.data == nil {
		return x, y
	}
	return x - d.data.CamX, y - d.data.CamY
}

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
// Useful when building a custom vector.Path before passing it to PathFill/PathStroke.
func (d *DrawComponent) ScreenX(mapX float32) float32 {
	if d.data == nil {
		return mapX
	}
	return mapX - d.data.CamX
}

// ScreenY converts a map-space Y coordinate to screen-space Y.
// Useful when building a custom vector.Path before passing it to PathFill/PathStroke.
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
func (d *DrawComponent) Rect(x, y, w, h float32, c color.RGBA) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.FillRect(s, sx, sy, w, h, c, false)
}

// Circle draws a filled circle centered at map-space (x, y) with radius r.
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
// strokeWidth controls the border thickness in pixels.
func (d *DrawComponent) RectStroke(x, y, w, h float32, c color.RGBA, strokeWidth float32) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.StrokeRect(s, sx, sy, w, h, strokeWidth, c, false)
}

// CircleStroke draws a circle outline centered at map-space (x, y) with radius r.
// strokeWidth controls the border thickness in pixels.
func (d *DrawComponent) CircleStroke(x, y, r float32, c color.RGBA, strokeWidth float32) {
	s := d.screen()
	if s == nil {
		return
	}
	sx, sy := d.toScreen(x, y)
	vector.StrokeCircle(s, sx, sy, r, strokeWidth, c, true)
}

// Line draws a straight line from map-space (x0, y0) to (x1, y1).
// strokeWidth controls the line thickness in pixels.
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
// Build the path using screen-space coordinates via ScreenX() / ScreenY().
//
// Example:
//
//	var p vector.Path
//	p.MoveTo(o.ScreenX(100), o.ScreenY(50))
//	p.LineTo(o.ScreenX(150), o.ScreenY(100))
//	p.Close()
//	o.PathFill(&p, color.RGBA{255, 0, 0, 255}, true)
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
// Build the path using screen-space coordinates via ScreenX() / ScreenY().
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

// Text draws a string at map-space (x, y) using the active default font.
// Color c is applied to the rendered glyphs.
func (d *DrawComponent) Text(text string, x, y float32, c color.RGBA) {
	s := d.screen()
	if s == nil || text == "" {
		return
	}
	sx, sy := d.toScreen(x, y)
	face := textv2.NewGoXFace(d.activeFont())
	op := &textv2.DrawOptions{}
	op.GeoM.Translate(float64(sx), float64(sy))
	op.ColorScale.ScaleWithColor(c)
	textv2.Draw(s, text, face, op)
}

// TextEx draws a string at map-space (x, y) with a uniform scale applied to the font.
// scale 1.0 = default size, 2.0 = double size. Uses the same font as Text().
func (d *DrawComponent) TextEx(text string, x, y float32, c color.RGBA, scale float64) {
	s := d.screen()
	if s == nil || text == "" {
		return
	}
	sx, sy := d.toScreen(x, y)
	face := textv2.NewGoXFace(d.activeFont())
	op := &textv2.DrawOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(sx), float64(sy))
	op.ColorScale.ScaleWithColor(c)
	textv2.Draw(s, text, face, op)
}

// =============================================================================
// Image
// =============================================================================

// Image draws frame idx of the given ISpriteLW at map-space (x, y).
// Allows manual sprite rendering independent of SpriteComponent.
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
