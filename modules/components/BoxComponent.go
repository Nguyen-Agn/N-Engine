package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// BoxComponent provides bounding box properties for an object.
// It is used to define dimensions and offsets for collisions or spatial queries.
type BoxComponent struct {
	IObject
	data *BoxData
}

// interface adready defined at domain

// BindComponent binds the base object and its ECS data to this component.
// Inputs:
//   - base: The base IObject to bind to.
func (p *BoxComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Box)
}

// resgiter new Component
var Box = enginetype.Box

// init initializes the default data for the box component.
// It registers the "box" component token with a default rectangular shape.
func init() {
	enginetype.RegisterComponentInitializer("box", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Box, BoxData{
			Shape: BSRectangle,
		})
	})
}

// BoxW retrieves the width of the bounding box.
// Outputs: Returns the box width as a float32.
func (p BoxComponent) BoxW() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxW
}

// SetBoxW sets the width of the bounding box.
// Inputs:
//   - boxW: The new width for the box.
func (p BoxComponent) SetBoxW(boxW float32) {

	if p.data != nil {
		p.data.BoxW = boxW
	}
}

// BoxH retrieves the height of the bounding box.
// Outputs: Returns the box height as a float32.
func (p BoxComponent) BoxH() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxH
}

// SetBoxH sets the height of the bounding box.
// Inputs:
//   - boxH: The new height for the box.
func (p BoxComponent) SetBoxH(boxH float32) {
	if p.data != nil {
		p.data.BoxH = boxH
	}
}

// BoxX retrieves the X-coordinate offset of the bounding box relative to the object.
// Outputs: Returns the X offset as a float32.
func (p BoxComponent) BoxX() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxX
}

// SetBoxX sets the X-coordinate offset of the bounding box.
// Inputs:
//   - boxX: The new X offset.
func (p BoxComponent) SetBoxX(boxX float32) {

	if p.data != nil {
		p.data.BoxX = boxX
	}
}

// BoxY retrieves the Y-coordinate offset of the bounding box relative to the object.
// Outputs: Returns the Y offset as a float32.
func (p BoxComponent) BoxY() float32 {

	if p.data == nil {
		return 0
	}
	return p.data.BoxY
}

// SetBoxY sets the Y-coordinate offset of the bounding box.
// Inputs:
//   - boxY: The new Y offset.
func (p BoxComponent) SetBoxY(boxY float32) {

	if p.data != nil {
		p.data.BoxY = boxY
	}
}
// Shape retrieves the shape type of the bounding box.
// Outputs: Returns the BoxShape (e.g., Rectangle, Circle).
func (p BoxComponent) Shape() BoxShape {

	if p.data == nil {
		return BSRectangle
	}
	return p.data.Shape
}

// SetShape sets the shape type of the bounding box.
// Inputs:
//   - shape: The new BoxShape to apply.
func (p BoxComponent) SetShape(shape BoxShape) {
	if p.data != nil {
		p.data.Shape = shape
	}
}
