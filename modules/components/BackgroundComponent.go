package components

import (
	"autoworld/modules/enginetype"
	"image/color"

	"github.com/yohamta/donburi"
)

// Khai báo token cho Background Component
var Background = enginetype.Background

func init() {
	enginetype.RegisterComponentInitializer("bg", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Background, BackgroundData{
			Color:     color.RGBA{0, 0, 0, 0},
			IsVisible: true,
		})
	})
}

// BackgroundComponent là Mixin để nhúng vào Custom Object.
// Yêu cầu Object phải có base IObject để gọi Entry().
type BackgroundComponent struct {
	IObject
	data *BackgroundData
}

func (b *BackgroundComponent) BindComponent(base IObject) {
	b.IObject = base
	b.data = enginetype.GetComponent(base, Background)
}

func (b BackgroundComponent) Color() color.RGBA {
	if b.data == nil {
		return color.RGBA{}
	}
	return b.data.Color
}

func (b *BackgroundComponent) SetColor(c color.RGBA) {
	if b.data != nil {
		b.data.Color = c
	}
}

func (b BackgroundComponent) Sprite() ISpriteLW {
	if b.data == nil {
		return nil
	}
	return b.data.Sprite
}

func (b *BackgroundComponent) SetSprite(s ISpriteLW) {
	if b.data != nil {
		b.data.Sprite = s
	}
}

func (b BackgroundComponent) RepeatX() bool {
	if b.data == nil {
		return false
	}
	return b.data.RepeatX
}

func (b *BackgroundComponent) SetRepeatX(repeatX bool) {
	if b.data != nil {
		b.data.RepeatX = repeatX
	}
}

func (b BackgroundComponent) RepeatY() bool {
	if b.data == nil {
		return false
	}
	return b.data.RepeatY
}

func (b *BackgroundComponent) SetRepeatY(repeatY bool) {
	if b.data != nil {
		b.data.RepeatY = repeatY
	}
}

func (b BackgroundComponent) Stretch() bool {
	if b.data == nil {
		return false
	}
	return b.data.Stretch
}

func (b *BackgroundComponent) SetStretch(stretch bool) {
	if b.data != nil {
		b.data.Stretch = stretch
	}
}

func (b BackgroundComponent) ScrollSpeedX() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.ScrollSpeedX
}

func (b *BackgroundComponent) SetScrollSpeedX(speed float32) {
	if b.data != nil {
		b.data.ScrollSpeedX = speed
	}
}

func (b BackgroundComponent) ScrollSpeedY() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.ScrollSpeedY
}

func (b *BackgroundComponent) SetScrollSpeedY(speed float32) {
	if b.data != nil {
		b.data.ScrollSpeedY = speed
	}
}

func (b BackgroundComponent) OffsetX() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.OffsetX
}

func (b *BackgroundComponent) SetOffsetX(offset float32) {
	if b.data != nil {
		b.data.OffsetX = offset
	}
}

func (b BackgroundComponent) OffsetY() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.OffsetY
}

func (b *BackgroundComponent) SetOffsetY(offset float32) {
	if b.data != nil {
		b.data.OffsetY = offset
	}
}

func (b BackgroundComponent) IsVisible() bool {
	if b.data == nil {
		return false
	}
	return b.data.IsVisible
}

func (b *BackgroundComponent) SetIsVisible(visible bool) {
	if b.data != nil {
		b.data.IsVisible = visible
	}
}
