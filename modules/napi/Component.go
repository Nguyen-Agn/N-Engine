package napi

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// setComponent assigns a value to a custom component on an object.
//
// Purpose: Attaches or updates a specific custom data component for an entity using the donburi ECS.
//
// Inputs:
// - obj (IObject): The game object to modify.
// - comp (*donburi.ComponentType[T]): The component type previously created.
// - value (T): The data to assign to the component.
func setComponent[T any](obj IObject, comp *donburi.ComponentType[T], value T) {
	enginetype.SetComponent(obj, comp, value)
}

// getComponent retrieves a pointer to the data of a custom component on an object.
//
// Purpose: Allows reading and modifying a custom component's data directly.
//
// Inputs:
// - obj (IObject): The game object to query.
// - comp (*donburi.ComponentType[T]): The custom component type to look for.
//
// Outputs:
// - *T: A pointer to the component data, or nil if the object does not possess this component.
func getComponent[T any](obj IObject, comp *donburi.ComponentType[T]) *T {
	return enginetype.GetComponent(obj, comp)
}

// addComponentType appends a component type to an object's ECS entry.
//
// Purpose: Dynamically adds a new custom component type to an already instantiated object, enabling runtime capability expansion.
//
// Inputs:
// - obj (IObject): The game object to enhance.
// - comp (*donburi.ComponentType[T]): The component type to add.
func addComponentType[T any](obj IObject, comp *donburi.ComponentType[T]) {
	enginetype.AddComponentType(obj, comp)
}
