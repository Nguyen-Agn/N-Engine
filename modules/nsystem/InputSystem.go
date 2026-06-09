package nsystem

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/components"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"
)

// InputSystem iterates through all Objects with InputData and triggers
// their corresponding Handlers every frame based on EventType (Pressed, JustPressed, JustReleased).
//
// OR logic for key groups: if any key in KeyBinding.Keys matches,
// the Handler is called exactly once with the triggered key's name.
//
// Mouse bindings: InputSystem calls UpdateMouseBindings() on every object
// that embeds MouseComponent (detected via IMouse interface).
//
// InputSystem receives domain.IInputManager via interface to avoid
// direct dependency on the core module (Dependency Inversion Principle).
type InputSystem struct {
	input domain.IInputManager
}

// NewInputSystem creates an InputSystem using the provided IInputManager from the Engine.
// Inputs: input (domain.IInputManager) - The input manager instance.
// Outputs: Returns a pointer to a newly initialized InputSystem.
func NewInputSystem(input domain.IInputManager) *InputSystem {
	return &InputSystem{input: input}
}

// Update iterates through the objectList and processes both keyboard and mouse bindings.
// Inputs: objectList ([]IObject) - The list of objects to be updated.
// Purpose: It checks keyboard events (Pressed, JustPressed, JustReleased) based on bindings defined in each object's InputData, triggering handlers if conditions are met. It also processes mouse events for objects implementing the internal mouse updater interface.
func (s *InputSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		// ── Keyboard ──────────────────────────────────────────────────────────
		data := enginetype.GetComponent(obj, enginetype.Input)
		if data != nil {
			for _, binding := range data.Bindings {
				for _, key := range binding.Keys {
					if s.checkKey(key, binding.EventType) {
						keyName := domain.KeyReverseMap[key]
						binding.Handler(keyName) // OR logic: calls once with the first matching key
						break
					}
				}
			}
		}

		// ── Mouse ─────────────────────────────────────────────────────────────
		// Any object embedding MouseComponent implements iMouseUpdater.
		// InputSystem calls UpdateMouseBindings() to process registered bindings.
		if mu, ok := obj.(iMouseUpdater); ok {
			mu.UpdateMouseBindings()
		}
	}
}

// checkKey evaluates the key state based on the provided EventType.
// Inputs:
//
//	key (domain.Key) - The key to check.
//	evt (domain.EventType) - The type of event (Pressed, JustPressed, JustReleased).
//
// Outputs: Returns true if the key state matches the event type, false otherwise.
func (s *InputSystem) checkKey(key domain.Key, evt domain.EventType) bool {
	switch evt {
	case domain.EventPressed:
		return s.input.IsKeyPressed(key)
	case domain.EventJustPressed:
		return s.input.IsKeyJustPressed(key)
	case domain.EventJustReleased:
		return s.input.IsKeyJustReleased(key)
	}
	return false
}

// iMouseUpdater is an internal interface for InputSystem to detect objects with MouseComponent.
// It is not exported to avoid exposing implementation details.
type iMouseUpdater interface {
	UpdateMouseBindings()
}

// Ensure MouseComponent implements iMouseUpdater at compile time.
var _ iMouseUpdater = (*components.MouseComponent)(nil)
