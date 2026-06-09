package napi

import (
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

// ─── ECS Entry Helper ─────────────────────────────────────────────────────────

// newEntry creates a raw ECS entity (donburi.Entry) within the scene's World using the provided component types.
//
// Purpose: Enables manual entity creation bypassing the standard NewObject pipeline, useful for highly customized ECS integration.
//
// Inputs:
// - scene (IScene): The active scene providing the ECS World.
// - components (...donburi.IComponentType): A variadic list of component types to instantiate the entity with.
//
// Outputs:
// - *donburi.Entry: The raw ECS entry reference.
func newEntry(scene IScene, components ...donburi.IComponentType) *donburi.Entry {
	return scene.World().Entry(scene.World().Create(components...))
}

// ─── Custom Component Type ────────────────────────────────────────────────────

// newComponentType creates and registers a new Donburi component type for custom game data.
//
// Purpose: Allows developers to define new component types. If a non-empty token is provided, it registers the component in the global registry, permitting its use via shorthand strings in NewObject.
//
// Inputs:
// - token (string): A short string identifier (e.g., "sta") used to refer to this component during object creation.
//
// Outputs:
// - *donburi.ComponentType[T]: The initialized generic component type.
func newComponentType[T any](token string) *donburi.ComponentType[T] {
	return enginetype.NewComponentType[T](token)
}

// getComponentType looks up a registered component type using its string token.
//
// Purpose: Retrieves the underlying ECS component type definition required for dynamic entity construction.
//
// Inputs:
// - token (string): The string identifier for the component.
//
// Outputs:
// - donburi.IComponentType: The corresponding component type, or nil if the token is unregistered.
func getComponentType(token string) donburi.IComponentType {
	return enginetype.GetComponentType(token)
}
