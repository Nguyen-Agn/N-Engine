package layout

import (
	"github.com/yohamta/donburi"
)

// A acts as an adapter/anchor for IObject in the layout system.
type A struct {
	name      string
	isVisible bool
	width     int
	height    int
	x         int
	y         int
	childrens []ILayout

	layoutConfig int
	gap          int

	object IObject
}

// Purpose: Creates a new layout anchor 'A'.
// Inputs: name (string) - Name of the layout, obj (IObject) - The game object to bind to, width (int) - Layout width, height (int) - Layout height.
// Outputs: (*A) - The initialized layout element.
func NewA(name string, obj IObject, width int, height int) *A {
	return &A{
		name:         name,
		object:       obj,
		width:        width,
		height:       height,
		childrens:    []ILayout{},
		layoutConfig: DirRow | AlignStart | JustifyStart,
		isVisible:    true,
	}
}

// Purpose: Gets the underlying game object bound to this layout.
// Outputs: (IObject) - The bound game object.
func (a *A) Object() IObject       { return a.object }

// Purpose: Sets the underlying game object for this layout.
// Inputs: obj (IObject) - The game object to bind.
func (a *A) SetObject(obj IObject) { a.object = obj }

// Purpose: Gets the X coordinate of the layout box.
func (a *A) BoxX() int { return a.x }

// Purpose: Gets the Y coordinate of the layout box.
func (a *A) BoxY() int { return a.y }

// Purpose: Gets the width of the layout box.
func (a *A) BoxW() int { return a.width }

// Purpose: Gets the height of the layout box.
func (a *A) BoxH() int { return a.height }

// Purpose: Sets the X coordinate of the layout box and updates the bound object's Position and Box components.
// Inputs: x (int) - The new X coordinate.
func (a *A) SetBoxX(x int) {
	a.x = x
	if a.object != nil {
		entry := a.object.Entry()
		if entry.HasComponent(Position) {
			donburi.Get[PositionData](entry, Position).X = float32(x)
		}
		if entry.HasComponent(Box) {
			donburi.Get[BoxData](entry, Box).BoxX = float32(x)
		}
	}
}

// Purpose: Sets the Y coordinate of the layout box and updates the bound object's Position and Box components.
// Inputs: y (int) - The new Y coordinate.
func (a *A) SetBoxY(y int) {
	a.y = y
	if a.object != nil {
		entry := a.object.Entry()
		if entry.HasComponent(Position) {
			donburi.Get[PositionData](entry, Position).Y = float32(y)
		}
		if entry.HasComponent(Box) {
			donburi.Get[BoxData](entry, Box).BoxY = float32(y)
		}
	}
}

// Purpose: Sets the width of the layout box.
func (a *A) SetBoxW(w int) { a.width = w }

// Purpose: Sets the height of the layout box.
func (a *A) SetBoxH(h int) { a.height = h }

// Purpose: Gets the name of the layout element.
func (a *A) Name() string        { return a.name }

// Purpose: Sets the name of the layout element.
func (a *A) SetName(name string) { a.name = name }

// Purpose: Checks if the layout is visible.
func (a *A) IsVisible() bool             { return a.isVisible }

// Purpose: Sets the visibility of the layout.
func (a *A) SetIsVisible(isVisible bool) { a.isVisible = isVisible }

// Purpose: Retrieves all child layouts.
func (a *A) Childrens() []ILayout         { return a.childrens }

// Purpose: Adds a child layout.
func (a *A) AddChildren(children ILayout) { a.childrens = append(a.childrens, children) }

// Purpose: Removes a child layout from the children list.
func (a *A) RemoveChildren(children ILayout) {
	for i, child := range a.childrens {
		if child == children {
			a.childrens = append(a.childrens[:i], a.childrens[i+1:]...)
			break
		}
	}
}
// Purpose: Retrieves a child layout recursively by its name.
// Inputs: name (string) - The name of the child layout.
// Outputs: (ILayout) - The found child layout, or nil if not found.
func (a *A) GetChildrenByName(name string) ILayout {
	for _, child := range a.childrens {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

// Purpose: Sets the layout configuration (direction, align, justify).
func (a *A) SetLayoutConfig(config int) { a.layoutConfig = config }

// Purpose: Gets the current layout configuration.
func (a *A) LayoutConfig() int          { return a.layoutConfig }

// Purpose: Sets the gap between child layouts.
func (a *A) SetGap(gap int) { a.gap = gap }

// Purpose: Gets the gap between child layouts.
func (a *A) Gap() int       { return a.gap }

// Purpose: Computes the layout positions and updates the bound game object's position based on layout rules.
// Inputs: parentX, parentY (int) - Parent coordinates, parentW, parentH (int) - Parent dimensions.
func (a *A) ComputeLayout(parentX, parentY int, parentW, parentH int) {
	if a.object != nil {
		entry := a.object.Entry()
		if entry.HasComponent(Position) {
			donburi.Get[PositionData](entry, Position).X = float32(a.x)
			donburi.Get[PositionData](entry, Position).Y = float32(a.y)
		}
		if entry.HasComponent(Box) {
			donburi.Get[BoxData](entry, Box).BoxX = float32(a.x)
			donburi.Get[BoxData](entry, Box).BoxY = float32(a.y)
		}
	}
}
