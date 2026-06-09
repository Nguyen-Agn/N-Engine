package components

import (
	"image/color"

	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Background = enginetype.Background

// init initializes the default data for the background component.
// It registers the "bg" component token with a default invisible, transparent background.
func init() {
	enginetype.RegisterComponentInitializer("bg", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Background, BackgroundData{
			Color:     color.RGBA{0, 0, 0, 0},
			IsVisible: true,
		})
	})
}

// Background Feature
type BackgroundComponent struct {
	IObject
	data *BackgroundData
}

// BindComponent binds the base object and its ECS data to this component.
// Inputs:
//   - base: The base IObject to bind to.
func (b *BackgroundComponent) BindComponent(base IObject) {
	b.IObject = base
	b.data = enginetype.GetComponent(base, Background)
}

// Color retrieves the background color.
// Outputs: Returns the color.RGBA representing the background color.
func (b BackgroundComponent) Color() color.RGBA {
	if b.data == nil {
		return color.RGBA{}
	}
	return b.data.Color
}

// SetColor sets the background color.
// Inputs:
//   - c: The new color.RGBA to set for the background.
func (b *BackgroundComponent) SetColor(c color.RGBA) {
	if b.data != nil {
		b.data.Color = c
	}
}

// Sprite retrieves the background sprite.
// Outputs: Returns the ISpriteLW object used as the background, or nil if none.
func (b BackgroundComponent) Sprite() ISpriteLW {
	if b.data == nil {
		return nil
	}
	return b.data.Sprite
}

// SetSprite sets the background sprite.
// Inputs:
//   - s: The ISpriteLW object to use as the background sprite.
func (b *BackgroundComponent) SetSprite(s ISpriteLW) {
	if b.data != nil {
		b.data.Sprite = s
	}
}

// RepeatX checks if the background should repeat along the X-axis.
// Outputs: Returns true if it repeats on the X-axis, false otherwise.
func (b BackgroundComponent) RepeatX() bool {
	if b.data == nil {
		return false
	}
	return b.data.RepeatX
}

// SetRepeatX sets whether the background should repeat along the X-axis.
// Inputs:
//   - repeatX: Boolean indicating whether to enable X-axis repetition.
func (b *BackgroundComponent) SetRepeatX(repeatX bool) {
	if b.data != nil {
		b.data.RepeatX = repeatX
	}
}

// RepeatY checks if the background should repeat along the Y-axis.
// Outputs: Returns true if it repeats on the Y-axis, false otherwise.
func (b BackgroundComponent) RepeatY() bool {
	if b.data == nil {
		return false
	}
	return b.data.RepeatY
}

// SetRepeatY sets whether the background should repeat along the Y-axis.
// Inputs:
//   - repeatY: Boolean indicating whether to enable Y-axis repetition.
func (b *BackgroundComponent) SetRepeatY(repeatY bool) {
	if b.data != nil {
		b.data.RepeatY = repeatY
	}
}

// Stretch checks if the background is currently in stretch mode.
// Outputs: Returns true if the background is stretched to fit the screen, false otherwise.
func (b BackgroundComponent) Stretch() bool {
	if b.data == nil {
		return false
	}
	return b.data.Stretch
}

// SetStretch sets the stretch mode for the background.
// Inputs:
//   - stretch: Boolean indicating whether to stretch the background.
func (b *BackgroundComponent) SetStretch(stretch bool) {
	if b.data != nil {
		b.data.Stretch = stretch
	}
}

// ScrollSpeedX retrieves the scrolling speed of the background along the X-axis.
// Outputs: Returns the scrolling speed as a float32.
func (b BackgroundComponent) ScrollSpeedX() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.ScrollSpeedX
}

// SetScrollSpeedX sets the scrolling speed of the background along the X-axis.
// Inputs:
//   - speed: The new scrolling speed for the X-axis.
func (b *BackgroundComponent) SetScrollSpeedX(speed float32) {
	if b.data != nil {
		b.data.ScrollSpeedX = speed
	}
}

// ScrollSpeedY retrieves the scrolling speed of the background along the Y-axis.
// Outputs: Returns the scrolling speed as a float32.
func (b BackgroundComponent) ScrollSpeedY() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.ScrollSpeedY
}

// SetScrollSpeedY sets the scrolling speed of the background along the Y-axis.
// Inputs:
//   - speed: The new scrolling speed for the Y-axis.
func (b *BackgroundComponent) SetScrollSpeedY(speed float32) {
	if b.data != nil {
		b.data.ScrollSpeedY = speed
	}
}

// OffsetX retrieves the offset position along the X-axis.
// Outputs: Returns the offset as a float32.
func (b BackgroundComponent) OffsetX() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.OffsetX
}

// SetOffsetX sets the offset position along the X-axis.
// Inputs:
//   - offset: The new offset for the X-axis.
func (b *BackgroundComponent) SetOffsetX(offset float32) {
	if b.data != nil {
		b.data.OffsetX = offset
	}
}

// OffsetY retrieves the offset position along the Y-axis.
// Outputs: Returns the offset as a float32.
func (b BackgroundComponent) OffsetY() float32 {
	if b.data == nil {
		return 0
	}
	return b.data.OffsetY
}

// SetOffsetY sets the offset position along the Y-axis.
// Inputs:
//   - offset: The new offset for the Y-axis.
func (b *BackgroundComponent) SetOffsetY(offset float32) {
	if b.data != nil {
		b.data.OffsetY = offset
	}
}

// IsVisible checks if the background is visible.
// Outputs: Returns true if the background is set to be visible, false otherwise.
func (b BackgroundComponent) IsVisible() bool {
	if b.data == nil {
		return false
	}
	return b.data.IsVisible
}

// SetIsVisible sets the visibility of the background.
// Inputs:
//   - visible: Boolean indicating whether the background should be visible.
func (b *BackgroundComponent) SetIsVisible(visible bool) {
	if b.data != nil {
		b.data.IsVisible = visible
	}
}
