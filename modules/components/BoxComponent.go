package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// BoxComponent là Mixin để nhúng vào Custom Object.
type BoxComponent struct {
	IObject
	data *BoxData
}

// interface adready defined at domain

func (p *BoxComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Box)
}

// resgiter new Component
var Box = enginetype.Box

func init() {
	enginetype.RegisterComponentInitializer("box", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Box, BoxData{
			Shape: BSRectangle,
		})
	})
}

func (p BoxComponent) BoxW() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxW
}

func (p BoxComponent) SetBoxW(boxW float32) {

	if p.data != nil {
		p.data.BoxW = boxW
	}
}

func (p BoxComponent) BoxH() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxH
}

func (p BoxComponent) SetBoxH(boxH float32) {

	if p.data != nil {
		p.data.BoxH = boxH
	}
}

func (p BoxComponent) BoxX() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxX
}

func (p BoxComponent) SetBoxX(boxX float32) {

	if p.data != nil {
		p.data.BoxX = boxX
	}
}

func (p BoxComponent) BoxY() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxY
}

func (p BoxComponent) SetBoxY(boxY float32) {

	if p.data != nil {
		p.data.BoxY = boxY
	}
}
func (p BoxComponent) Shape() BoxShape {

	if p.data == nil {
		return BSRectangle
	}
	return p.data.Shape
}

func (p BoxComponent) SetShape(shape BoxShape) {
	if p.data != nil {
		p.data.Shape = shape
	}
}
