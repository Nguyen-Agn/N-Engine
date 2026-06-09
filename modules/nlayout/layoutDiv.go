package layout

// Div layout struct
type Div struct {
	name         string
	isVisible    bool
	width        int
	height       int
	x            int
	y            int
	childrens    []ILayout
	unit         string
	layoutConfig int
	gap          int
}

// Purpose: Creates a new Div layout.
// Inputs: name (string) - Name of the div, x, y (int) - Coordinates, width, height (int) - Dimensions, unit (string) - Layout unit.
// Outputs: (*Div) - A pointer to the newly created Div.
func NewDivWithNameAndWidthAndHeight(name string, x int, y int, width int, height int, unit string) *Div {
	return &Div{
		name:         name,
		width:        width,
		height:       height,
		x:            x,
		y:            y,
		unit:         unit,
		childrens:    []ILayout{},
		layoutConfig: DirRow | AlignStart | JustifyStart,
		isVisible:    true,
	}
}

// Purpose: Gets the X coordinate of the div.
func (d *Div) BoxX() int { return d.x }

// Purpose: Gets the Y coordinate of the div.
func (d *Div) BoxY() int { return d.y }

// Purpose: Gets the width of the div.
func (d *Div) BoxW() int { return d.width }

// Purpose: Gets the height of the div.
func (d *Div) BoxH() int { return d.height }

// Purpose: Sets the X coordinate.
func (d *Div) SetBoxX(x int) { d.x = x }

// Purpose: Sets the Y coordinate.
func (d *Div) SetBoxY(y int) { d.y = y }

// Purpose: Sets the width.
func (d *Div) SetBoxW(w int) { d.width = w }

// Purpose: Sets the height.
func (d *Div) SetBoxH(h int) { d.height = h }

// Purpose: Gets the layout name.
func (d *Div) Name() string        { return d.name }

// Purpose: Sets the layout name.
func (d *Div) SetName(name string) { d.name = name }

// Purpose: Checks if the div is visible.
func (d *Div) IsVisible() bool             { return d.isVisible }

// Purpose: Sets the visibility of the div.
func (d *Div) SetIsVisible(isVisible bool) { d.isVisible = isVisible }

// Purpose: Retrieves all child layouts.
func (d *Div) Childrens() []ILayout { return d.childrens }

// Purpose: Adds a child layout to this div.
func (d *Div) AddChildren(children ILayout) {
	d.childrens = append(d.childrens, children)
}

// Purpose: Removes a child layout from this div.
func (d *Div) RemoveChildren(children ILayout) {
	for i, child := range d.childrens {
		if child == children {
			d.childrens = append(d.childrens[:i], d.childrens[i+1:]...)
			break
		}
	}
}
// Purpose: Recursively searches for and returns a child layout by name.
// Inputs: name (string) - The name to search for.
// Outputs: (ILayout) - The found child layout, or nil if not found.
func (d *Div) GetChildrenByName(name string) ILayout {
	for _, child := range d.childrens {
		if child.Name() == name {
			return child
		}
		candidate := child.GetChildrenByName(name)
		if candidate != nil {
			return candidate
		}
	}
	return nil
}

// Layout specific
// Purpose: Sets the layout configuration (e.g., row/col, justify, align).
func (d *Div) SetLayoutConfig(config int) { d.layoutConfig = config }

// Purpose: Gets the layout configuration.
func (d *Div) LayoutConfig() int          { return d.layoutConfig }

// Purpose: Sets the gap between child layouts.
func (d *Div) SetGap(gap int) { d.gap = gap }

// Purpose: Gets the gap between child layouts.
func (d *Div) Gap() int       { return d.gap }

// Purpose: Computes the layout logic, setting positions of all child layouts based on the current layout configuration.
// Inputs: parentX, parentY (int) - Parent coordinates, parentW, parentH (int) - Parent dimensions.
func (d *Div) ComputeLayout(parentX, parentY int, parentW, parentH int) {
	if len(d.childrens) == 0 {
		return
	}

	var totalChildWidth, totalChildHeight int
	for _, child := range d.childrens {
		totalChildWidth += child.BoxW()
		totalChildHeight += child.BoxH()
	}

	if len(d.childrens) > 1 {
		totalChildWidth += d.gap * (len(d.childrens) - 1)
		totalChildHeight += d.gap * (len(d.childrens) - 1)
	}

	startX := d.x
	startY := d.y

	dir := d.layoutConfig & DirMask
	justify := d.layoutConfig & JustifyMask
	align := d.layoutConfig & AlignMask

	if dir == DirRow {
		switch justify {
		case JustifyCenter:
			startX += (d.width - totalChildWidth) / 2
		case JustifyEnd:
			startX += d.width - totalChildWidth
		}

		currentX := startX
		for _, child := range d.childrens {
			childY := d.y
			switch align {
			case AlignCenter:
				childY += (d.height - child.BoxH()) / 2
			case AlignEnd:
				childY += d.height - child.BoxH()
			}

			child.SetBoxX(currentX)
			child.SetBoxY(childY)
			child.ComputeLayout(currentX, childY, child.BoxW(), child.BoxH())

			currentX += child.BoxW() + d.gap
		}
	} else { // DirColumn
		switch justify {
		case JustifyCenter:
			startY += (d.height - totalChildHeight) / 2
		case JustifyEnd:
			startY += d.height - totalChildHeight
		}

		currentY := startY
		for _, child := range d.childrens {
			childX := d.x
			switch align {
			case AlignCenter:
				childX += (d.width - child.BoxW()) / 2
			case AlignEnd:
				childX += d.width - child.BoxW()
			}

			child.SetBoxX(childX)
			child.SetBoxY(currentY)
			child.ComputeLayout(childX, currentY, child.BoxW(), child.BoxH())

			currentY += child.BoxH() + d.gap
		}
	}
}
