package components

import (
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

// GenericComponent là mixin generic cho phép game dev tạo Component Mixin riêng
// mà vẫn được engine tự động bind, không cần thao tác ECS thủ công.
// Thường được re-export qua napi — dev chỉ cần import napi để dùng.
//
// # Cách dùng cơ bản (chỉ Get/Set) — qua napi
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//
//	type Hero struct {
//	    napi.IObject
//	    napi.IPosition
//	    napi.GenericComponent[StatsData]
//	}
//
//	func NewHero() *Hero {
//	    h := &Hero{
//	        GenericComponent: napi.NewGenericComponent(StatsComp),
//	    }
//	    napi.NewObject(h, "hero", "pos sta")
//	    napi.Register(h, "")
//	    return h
//	}
//
//	func (h *Hero) OnStep() {
//	    h.Get().Health -= 1
//	}
//
// # Cách dùng nâng cao (thêm method riêng) — qua napi
//
//	type StatsComponent struct {
//	    napi.GenericComponent[StatsData]
//	}
//
//	func (s *StatsComponent) TakeDamage(amount int) {
//	    if s.Get() != nil {
//	        s.Get().Health -= amount
//	    }
//	}
//
//	type Hero struct {
//	    napi.IObject
//	    napi.IPosition
//	    StatsComponent   // method TakeDamage được promoted lên Hero
//	}
//
//	func NewHero() *Hero {
//	    h := &Hero{
//	        StatsComponent: StatsComponent{
//	            GenericComponent: napi.NewGenericComponent(StatsComp),
//	        },
//	    }
//	    napi.NewObject(h, "hero", "pos sta")
//	    napi.Register(h, "")
//	    return h
//	}
type GenericComponent[T any] struct {
	IObject
	comp *donburi.ComponentType[T]
	data *T
}

// NewGenericComponent creates a GenericComponent bound to a ComponentType.
// Purpose: Initializes a generic component before it is processed by the engine.
// Inputs:
//   - comp: The donburi ComponentType to bind.
//
// Outputs: A new GenericComponent instance.
// Special Requirements: Must be called before napi.NewObject for BindComponent to work correctly.
func NewGenericComponent[T any](comp *donburi.ComponentType[T]) GenericComponent[T] {
	return GenericComponent[T]{comp: comp}
}

// BindComponent binds the base object and its corresponding ECS data to the component.
// Purpose: Automatically invoked by the engine to fetch the pointer to the ECS data.
// Inputs:
//   - base: The IObject representing the base entity to bind.
//
// Outputs: None.
func (c *GenericComponent[T]) BindComponent(base IObject) {
	c.IObject = base
	if c.comp != nil {
		c.data = enginetype.GetComponent(base, c.comp)
	}
}

// Data retrieves the direct pointer to the component's data in the ECS.
// Purpose: Allows direct modification of the ECS data.
// Inputs: None.
// Outputs: A pointer to the component data. Returns nil if not bound or not present.
func (c *GenericComponent[T]) Data() *T {
	return c.data
}

// SetData overwrites the entire component data with a new value.
// Purpose: Replaces the current data in the ECS.
// Inputs:
//   - val: The new data to set.
//
// Outputs: None.
func (c *GenericComponent[T]) SetData(val T) {
	if c.data != nil {
		*c.data = val
	}
}
