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

func (d *Div) BoxX() int { return d.x }
func (d *Div) BoxY() int { return d.y }
func (d *Div) BoxW() int { return d.width }
func (d *Div) BoxH() int { return d.height }

func (d *Div) SetBoxX(x int) { d.x = x }
func (d *Div) SetBoxY(y int) { d.y = y }
func (d *Div) SetBoxW(w int) { d.width = w }
func (d *Div) SetBoxH(h int) { d.height = h }

func (d *Div) Name() string        { return d.name }
func (d *Div) SetName(name string) { d.name = name }

func (d *Div) IsVisible() bool             { return d.isVisible }
func (d *Div) SetIsVisible(isVisible bool) { d.isVisible = isVisible }

func (d *Div) Childrens() []ILayout { return d.childrens }
func (d *Div) AddChildren(children ILayout) {
	d.childrens = append(d.childrens, children)
}
func (d *Div) RemoveChildren(children ILayout) {
	for i, child := range d.childrens {
		if child == children {
			d.childrens = append(d.childrens[:i], d.childrens[i+1:]...)
			break
		}
	}
}
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
func (d *Div) SetLayoutConfig(config int) { d.layoutConfig = config }
func (d *Div) LayoutConfig() int          { return d.layoutConfig }

func (d *Div) SetGap(gap int) { d.gap = gap }
func (d *Div) Gap() int       { return d.gap }

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
