package components

import (
	"image/color"
	"strconv"
	"strings"

	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// SpriteComponent là Mixin để nhúng vào Custom Object.
type SpriteComponent struct {
	IObject
	data *SpriteData
}

// interface adready defined at domain

// resgiter new Component
var Sprite = enginetype.Sprite

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
func (p *SpriteComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Sprite)
}

func (p SpriteComponent) SpriteIdx() int {
	if p.data == nil {
		return 0
	}
	return p.data.SpriteIdx
}

func (p SpriteComponent) SetSpriteIdx(spriteIdx int) {
	if p.data != nil {
		p.data.SpriteIdx = spriteIdx
	}
}

func (p SpriteComponent) ImageSpeed() float32 {
	if p.data == nil {
		return 1
	}
	return p.data.ImageSpeed
}

func (p SpriteComponent) SetImageSpeed(imageSpeed float32) {
	if p.data != nil {
		p.data.ImageSpeed = imageSpeed
	}
}

func (p SpriteComponent) Rotation() float32 {
	if p.data == nil {
		return 0
	}
	return p.data.Rotation
}

func (p SpriteComponent) SetRotation(rotation float32) {
	if p.data != nil {
		p.data.Rotation = rotation
	}
}

func (p SpriteComponent) OffsetX() float32 {
	if p.data == nil {
		return 0
	}
	return p.data.OffsetX
}

func (p SpriteComponent) SetOffsetX(offsetX float32) {
	if p.data != nil {
		p.data.OffsetX = offsetX
	}
}

func (p SpriteComponent) OffsetY() float32 {
	if p.data == nil {
		return 0
	}
	return p.data.OffsetY
}

func (p SpriteComponent) SetOffsetY(offsetY float32) {
	if p.data != nil {
		p.data.OffsetY = offsetY
	}
}

func (p SpriteComponent) ImageColor() color.RGBA {
	if p.data == nil {
		return color.RGBA{255, 255, 255, 255}
	}
	return p.data.ImageColor
}

func (p SpriteComponent) SetImageColor(imageColor color.RGBA) {
	if p.data != nil {
		p.data.ImageColor = imageColor
	}
}

func (p SpriteComponent) ScaleX() float32 {
	if p.data == nil {
		return 1
	}
	return p.data.ScaleX
}

func (p SpriteComponent) SetScaleX(scaleX float32) {
	if p.data != nil {
		p.data.ScaleX = scaleX
	}
}

func (p SpriteComponent) ScaleY() float32 {
	if p.data == nil {
		return 1
	}
	return p.data.ScaleY
}

func (p SpriteComponent) SetScaleY(scaleY float32) {
	if p.data != nil {
		p.data.ScaleY = scaleY
	}
}

func (p SpriteComponent) Sprite(name string) ISpriteLW {
	if p.data == nil {
		return nil
	}
	return p.data.Sprite[name]
}

func (p SpriteComponent) SetSprite(name string, sprite ISpriteLW) {
	if p.data != nil {
		p.data.Sprite[name] = sprite
	}
}

func (p SpriteComponent) AddSprite(name string, sprite ISpriteLW) bool {
	if p.data != nil {
		p.data.Sprite[name] = sprite
		return true
	}
	return false
}

func (p SpriteComponent) RemoveSprite(name string) bool {
	if p.data != nil {
		delete(p.data.Sprite, name)
		return true
	}
	return false
}

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

func (p SpriteComponent) ImageIndex() int {
	return p.SpriteIdx()
}

func (p SpriteComponent) SetImageIndex(imageIndex int) {
	p.SetSpriteIdx(imageIndex)
}

func (p SpriteComponent) SetCurrentSprite(name string) {
	if p.data != nil {
		p.data.CurrentSprite = name
	}
}
func (p SpriteComponent) GetCurrentSprite() ISpriteLW {
	if p.data != nil {
		return p.data.Sprite[p.data.CurrentSprite]
	}
	return nil
}

func (p SpriteComponent) IsVisible() bool {
	return p.data.IsVisible
}

func (p SpriteComponent) SetVisible(turn bool) {
	p.data.IsVisible = turn
}

// Enable / Disable 9Slice Mode
// String Ex: "5" ->5:all, "5 6" -> 5:top&bottom, 6:right&left, "1 2 3 4" -> each
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

func (p SpriteComponent) ZOrder() int {
	if p.data == nil {
		return 0
	}
	return p.data.ZOrder
}

func (p SpriteComponent) SetZOrder(z int) {
	if p.data != nil && p.data.ZOrder != z {
		p.data.ZOrder = z
		p.data.IsZOrderDirty = true
	}
}
