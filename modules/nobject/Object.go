package nobject

import (
	"github.com/Nguyen-Agn/N-Engine/domain"

	"github.com/yohamta/donburi"
)

type Object struct {
	entry *donburi.Entry
	pool  domain.IPool
}

// NewObject creates and returns a new Object instance using the given Donburi entry.
// Inputs: entry - pointer to the donburi ECS entry for this object.
// Outputs: a pointer to the newly created Object.
func NewObject(entry *donburi.Entry) *Object {
	return &Object{entry: entry}
}

// Entry returns the donburi.Entry inside the object.
// Used by napi.SetComponent and napi.GetComponent to assign custom component data.
// Outputs: pointer to the internal donburi.Entry.
func (this *Object) Entry() *donburi.Entry {
	return this.entry
}

// GetPool retrieves the object's parent pool.
// Outputs: the domain.IPool instance managing this object.
func (this *Object) GetPool() domain.IPool {
	return this.pool
}

// SetPool assigns the object's parent pool.
// Inputs: pool - the domain.IPool instance that will manage this object.
func (this *Object) SetPool(pool domain.IPool) {
	this.pool = pool
}

// #region Event

// OnCreate is triggered when the object is initially created or spawned.
func (this *Object) OnCreate() {}

// OnStep is triggered on every frame update to process the object's logic.
func (this *Object) OnStep() {}

// OnDestroy is triggered right before the object is completely removed from the world.
func (this *Object) OnDestroy() {}

// OnSave captures the object's specific save data for persistence.
// Inputs: data - a map where the object should store its persistable state.
func (this *Object) OnSave(data map[string]any) {}

// OnLoad restores the object's state from a loaded save data payload.
// Inputs: data - the previously saved map holding this object's specific state.
func (this *Object) OnLoad(data map[string]any) {}

// SetTokens configures the ECS component tokens bound to this object.
// Inputs: tokenClasses - a space-separated string of component types.
func (this *Object) SetTokens(tokenClasses string) {}

// Remove marks this object for removal from the scene or pool.
func (this *Object) Remove() {}

// #endregion
