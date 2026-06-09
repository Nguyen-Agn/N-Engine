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

// BindComponent binds the base object and its corresponding ECS data to the component.
// Purpose: Initializes the component with the base object and fetches the ECS data.
// Inputs:
//   - base: The IObject representing the base entity to bind.
// Outputs: None.
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

// Direction retrieves the current direction of the entity.
// Purpose: Gets the current facing direction in degrees.
// Inputs: None.
// Outputs: The current direction in degrees as a float32. Returns 0 if data is not initialized.
func (d DirectionComponent) Direction() float32 {

	if d.data == nil {
		return 0
	}
	return d.data.Direction
}

// SetDirection sets the absolute direction of the entity.
// Purpose: Sets a specific direction for the entity.
// Inputs:
//   - dir: The target direction in degrees (float32). Special Requirement: It will be normalized to [0, 360).
// Outputs: None.
func (d DirectionComponent) SetDirection(dir float32) {
	if d.data != nil {
		dir := math.Mod(float64(dir), 360)
		if dir < 0 {
			dir += 360
		}
		d.data.Direction = float32(dir)
	}
}

// Rotate modifies the current direction by a given relative angle.
// Purpose: Adds an angle to the current direction.
// Inputs:
//   - dir: The angle in degrees to add to the current direction (float32). Special Requirement: The result is normalized to [0, 360).
// Outputs: None.
func (d DirectionComponent) Rotate(dir float32) {

	if d.data != nil {
		dir := math.Mod(float64(d.data.Direction+dir), 360)
		if dir < 0 {
			dir += 360
		}
		d.data.Direction = float32(dir)
	}
}
