package components

import (
	"image/color"
	"strconv"
	"strings"

	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// SpriteComponent is a mixin to embed into Custom Objects for sprite capabilities.
type SpriteComponent struct {
	IObject
	data *SpriteData
}

// Interface already defined in domain.

// Register new Component
var Sprite = enginetype.Sprite

// init registers the Sprite Component initializer with default values.
func init() {
	enginetype.RegisterComponentInitializer("spr", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Sprite, SpriteData{
			Sprite:        make(map[string]ISpriteLW),
			IsVisible:     true,
			ScaleX:        1.0,
			ScaleY:        1.0,
			ImageColor:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			ImageSpeed:    1,
			Rotation:      0,
			OffsetX:       0,
			OffsetY:       0,
			CurrentSprite: "",
			SpriteIdx:     0,
			ZOrder:        0,
			IsZOrderDirty: true,
		})
	})
}
// BindComponent binds the base object and retrieves the SpriteData from the ECS.
// Inputs: base - the base IObject to bind.
func (p *SpriteComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Sprite)
}

// SpriteIdx returns the current sprite frame index.
// Outputs: int representing the current frame index.
func (p SpriteComponent) SpriteIdx() int {
	if p.data == nil {
		return 0
	}
	return p.data.SpriteIdx
}

// SetSpriteIdx sets the current sprite frame index.
// Inputs: spriteIdx - the new frame index.
func (p SpriteComponent) SetSpriteIdx(spriteIdx int) {
	if p.data != nil {
		p.data.SpriteIdx = spriteIdx
	}
}

// ImageSpeed returns the current image animation speed.
// Outputs: float32 representing the animation speed multiplier.
func (p SpriteComponent) ImageSpeed() float32 {
	if p.data == nil {
		return 1
	}
	return p.data.ImageSpeed
}

// SetImageSpeed sets the current image animation speed.
// Inputs: imageSpeed - the new animation speed multiplier.
func (p SpriteComponent) SetImageSpeed(imageSpeed float32) {
	if p.data != nil {
		p.data.ImageSpeed = imageSpeed
	}
}

// Rotation returns the current rotation angle.
// Outputs: float32 representing the rotation angle.
func (p SpriteComponent) Rotation() float32 {
	if p.data == nil {
		return 0
	}
	return p.data.Rotation
}

// SetRotation sets the current rotation angle.
// Inputs: rotation - the new rotation angle.
func (p SpriteComponent) SetRotation(rotation float32) {
	if p.data != nil {
		p.data.Rotation = rotation
	}
}

// OffsetX returns the current X offset for rendering.
// Outputs: float32 representing the X offset.
func (p SpriteComponent) OffsetX() float32 {
	if p.data == nil {
		return 0
	}
	return p.data.OffsetX
}

// SetOffsetX sets the X offset for rendering.
// Inputs: offsetX - the new X offset.
func (p SpriteComponent) SetOffsetX(offsetX float32) {
	if p.data != nil {
		p.data.OffsetX = offsetX
	}
}

// OffsetY returns the current Y offset for rendering.
// Outputs: float32 representing the Y offset.
func (p SpriteComponent) OffsetY() float32 {
	if p.data == nil {
		return 0
	}
	return p.data.OffsetY
}

// SetOffsetY sets the Y offset for rendering.
// Inputs: offsetY - the new Y offset.
func (p SpriteComponent) SetOffsetY(offsetY float32) {
	if p.data != nil {
		p.data.OffsetY = offsetY
	}
}

// ImageColor returns the current color tint applied to the sprite.
// Outputs: color.RGBA representing the tint color.
func (p SpriteComponent) ImageColor() color.RGBA {
	if p.data == nil {
		return color.RGBA{255, 255, 255, 255}
	}
	return p.data.ImageColor
}

// SetImageColor sets the color tint applied to the sprite.
// Inputs: imageColor - the new color tint.
func (p SpriteComponent) SetImageColor(imageColor color.RGBA) {
	if p.data != nil {
		p.data.ImageColor = imageColor
	}
}

// ScaleX returns the current horizontal scale factor.
// Outputs: float32 representing the X scale factor.
func (p SpriteComponent) ScaleX() float32 {
	if p.data == nil {
		return 1
	}
	return p.data.ScaleX
}

// SetScaleX sets the horizontal scale factor.
// Inputs: scaleX - the new X scale factor.
func (p SpriteComponent) SetScaleX(scaleX float32) {
	if p.data != nil {
		p.data.ScaleX = scaleX
	}
}

// ScaleY returns the current vertical scale factor.
// Outputs: float32 representing the Y scale factor.
func (p SpriteComponent) ScaleY() float32 {
	if p.data == nil {
		return 1
	}
	return p.data.ScaleY
}

// SetScaleY sets the vertical scale factor.
// Inputs: scaleY - the new Y scale factor.
func (p SpriteComponent) SetScaleY(scaleY float32) {
	if p.data != nil {
		p.data.ScaleY = scaleY
	}
}

// Sprite retrieves a specific sprite by name.
// Inputs: name - the string name of the sprite.
// Outputs: ISpriteLW representing the requested sprite, or nil if not found.
func (p SpriteComponent) Sprite(name string) ISpriteLW {
	if p.data == nil {
		return nil
	}
	return p.data.Sprite[name]
}

// SetSprite sets or updates a specific sprite by name.
// Inputs: name - the string name of the sprite, sprite - the ISpriteLW instance.
func (p SpriteComponent) SetSprite(name string, sprite ISpriteLW) {
	if p.data != nil {
		p.data.Sprite[name] = sprite
	}
}

// AddSprite adds a new sprite to the component if it exists.
// Inputs: name - the string name of the sprite, sprite - the ISpriteLW instance.
// Outputs: bool indicating success.
func (p SpriteComponent) AddSprite(name string, sprite ISpriteLW) bool {
	if p.data != nil {
		p.data.Sprite[name] = sprite
		return true
	}
	return false
}

// RemoveSprite removes a specific sprite by name.
// Inputs: name - the string name of the sprite to remove.
// Outputs: bool indicating success.
func (p SpriteComponent) RemoveSprite(name string) bool {
	if p.data != nil {
		delete(p.data.Sprite, name)
		return true
	}
	return false
}

// NextImage advances the animation frame to the next index, looping back to 0 if needed.
func (p SpriteComponent) NextImage() {
	if p.data == nil || p.data.CurrentSprite == "" {
		return
	}
	length := p.data.Sprite[p.data.CurrentSprite].Length()
	if length == 0 {
		return
	}
	p.data.SpriteIdx++
	if p.data.SpriteIdx >= length {
		p.data.SpriteIdx = 0
	}
}

// ImageIndex returns the current animation frame index (alias for SpriteIdx).
// Outputs: int representing the current frame index.
func (p SpriteComponent) ImageIndex() int {
	return p.SpriteIdx()
}

// SetImageIndex sets the current animation frame index (alias for SetSpriteIdx).
// Inputs: imageIndex - the new frame index.
func (p SpriteComponent) SetImageIndex(imageIndex int) {
	p.SetSpriteIdx(imageIndex)
}

// SetCurrentSprite sets the active sprite by name.
// Inputs: name - the string name of the sprite to make active.
func (p SpriteComponent) SetCurrentSprite(name string) {
	if p.data != nil {
		p.data.CurrentSprite = name
	}
}
// GetCurrentSprite retrieves the currently active sprite.
// Outputs: ISpriteLW representing the active sprite, or nil if none is active.
func (p SpriteComponent) GetCurrentSprite() ISpriteLW {
	if p.data != nil {
		return p.data.Sprite[p.data.CurrentSprite]
	}
	return nil
}

// IsVisible returns whether the sprite is currently visible.
// Outputs: bool indicating visibility.
func (p SpriteComponent) IsVisible() bool {
	return p.data.IsVisible
}

// SetVisible sets the visibility of the sprite.
// Inputs: turn - true to show, false to hide.
func (p SpriteComponent) SetVisible(turn bool) {
	p.data.IsVisible = turn
}

// Set9Slice enables or disables 9-Slice mode with specified padding values.
// Inputs: turn - enable/disable flag, TopRightBottomLeft - padding string (e.g., "5" -> all 5, "5 6" -> top/bottom 5, right/left 6).
func (p SpriteComponent) Set9Slice(turn bool, TopRightBottomLeft string) {
	p.data.IsNineSlice = turn
	if !turn {
		return
	}
	tokens := strings.Fields(TopRightBottomLeft)
	tokensLen := len(tokens)
	if tokensLen < 1 {
		return
	}
	values := make([]int, tokensLen)
	for idx, tok := range tokens {
		var err error
		values[idx], err = strconv.Atoi(tok)
		if err != nil {
			values[idx] = 0
		}
	}
	var indexs []int
	switch tokensLen {
	case 3:
		indexs = []int{0, 1, 2, 1}
	case 2:
		indexs = []int{0, 1, 0, 1}
	case 1:
		indexs = []int{0, 0, 0, 0}
	default:
		indexs = []int{0, 1, 2, 3}
	}
	p.data.NineSlice.Top = values[indexs[0]]
	p.data.NineSlice.Right = values[indexs[1]]
	p.data.NineSlice.Bottom = values[indexs[2]]
	p.data.NineSlice.Left = values[indexs[3]]
}

// ZOrder returns the current Z-Order for rendering.
// Outputs: int representing the Z-Order.
func (p SpriteComponent) ZOrder() int {
	if p.data == nil {
		return 0
	}
	return p.data.ZOrder
}

// SetZOrder sets the Z-Order for rendering, updating the dirty flag if changed.
// Inputs: z - the new Z-Order value.
func (p SpriteComponent) SetZOrder(z int) {
	if p.data != nil && p.data.ZOrder != z {
		p.data.ZOrder = z
		p.data.IsZOrderDirty = true
	}
}
