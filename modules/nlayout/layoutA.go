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

func (a *A) Object() IObject       { return a.object }
func (a *A) SetObject(obj IObject) { a.object = obj }

func (a *A) BoxX() int { return a.x }
func (a *A) BoxY() int { return a.y }
func (a *A) BoxW() int { return a.width }
func (a *A) BoxH() int { return a.height }

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

func (a *A) SetBoxW(w int) { a.width = w }
func (a *A) SetBoxH(h int) { a.height = h }

func (a *A) Name() string        { return a.name }
func (a *A) SetName(name string) { a.name = name }

func (a *A) IsVisible() bool             { return a.isVisible }
func (a *A) SetIsVisible(isVisible bool) { a.isVisible = isVisible }

func (a *A) Childrens() []ILayout         { return a.childrens }
func (a *A) AddChildren(children ILayout) { a.childrens = append(a.childrens, children) }
func (a *A) RemoveChildren(children ILayout) {
	for i, child := range a.childrens {
		if child == children {
			a.childrens = append(a.childrens[:i], a.childrens[i+1:]...)
			break
		}
	}
}
func (a *A) GetChildrenByName(name string) ILayout {
	for _, child := range a.childrens {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

func (a *A) SetLayoutConfig(config int) { a.layoutConfig = config }
func (a *A) LayoutConfig() int          { return a.layoutConfig }

func (a *A) SetGap(gap int) { a.gap = gap }
func (a *A) Gap() int       { return a.gap }

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
