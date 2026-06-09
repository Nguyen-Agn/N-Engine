package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// Position is the token for the Position Component used in the ECS.
var Position = enginetype.Position

// init registers the Position Component initializer with default values.
func init() {
	enginetype.RegisterComponentInitializer("pos", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Position, PositionData{X: 0, Y: 0})
	})
}

// PositionComponent is a mixin to embed into Custom Objects.
// It requires the Object to have an IObject base to call Entry().
type PositionComponent struct {
	IObject
	data *PositionData
}

// BindComponent binds the base object and retrieves the PositionData from the ECS.
// Inputs: base - the base IObject to bind.
func (p *PositionComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Position)
}

// X returns the current X coordinate.
// Outputs: float32 representing the X position.
func (p PositionComponent) X() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.X
}

// Y returns the current Y coordinate.
// Outputs: float32 representing the Y position.
func (p PositionComponent) Y() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.Y
}

// SetX sets the X coordinate.
// Inputs: x - the new X position.
func (p PositionComponent) SetX(x float32) {

	if p.data != nil {
		p.data.X = x
	}
}

// SetY sets the Y coordinate.
// Inputs: y - the new Y position.
func (p PositionComponent) SetY(y float32) {

	if p.data != nil {
		p.data.Y = y
	}
}
