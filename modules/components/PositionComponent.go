package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// Khai báo token cho Position Component
var Position = enginetype.Position

func init() {
	enginetype.RegisterComponentInitializer("pos", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Position, PositionData{X: 0, Y: 0})
	})
}

// PositionComponent là Mixin để nhúng vào Custom Object.
// Yêu cầu Object phải có base IObject để gọi Entry().
type PositionComponent struct {
	IObject
	data *PositionData
}

func (p *PositionComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Position)
}

func (p PositionComponent) X() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.X
}

func (p PositionComponent) Y() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.Y
}

func (p PositionComponent) SetX(x float32) {

	if p.data != nil {
		p.data.X = x
	}
}

func (p PositionComponent) SetY(y float32) {

	if p.data != nil {
		p.data.Y = y
	}
}
