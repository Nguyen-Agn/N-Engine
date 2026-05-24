package components

import (
	"math"

	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

type DirectionComponent struct {
	IObject
	data *DirectionData
}

func (p *DirectionComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Direction)
}

var Direction = enginetype.Direction

func init() {
	enginetype.RegisterComponentInitializer("dir", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Direction, DirectionData{
			Direction: 0.0,
		})
	})
}

func (d DirectionComponent) Direction() float32 {

	if d.data == nil {
		return 0
	}
	return d.data.Direction
}

func (d DirectionComponent) SetDirection(dir float32) {
	if d.data != nil {
		dir := math.Mod(float64(dir), 360)
		if dir < 0 {
			dir += 360
		}
		d.data.Direction = float32(dir)
	}
}

func (d DirectionComponent) Rotate(dir float32) {

	if d.data != nil {
		dir := math.Mod(float64(d.data.Direction+dir), 360)
		if dir < 0 {
			dir += 360
		}
		d.data.Direction = float32(dir)
	}
}
