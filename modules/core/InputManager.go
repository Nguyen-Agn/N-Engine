package core

import (
	"autoworld/domain"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ─── Key Mapping Table ───
//
// domain.Key --> ebiten.Key
var ebitenKeyMap = map[domain.Key]ebiten.Key{
	domain.KeyA:            ebiten.KeyA,
	domain.KeyB:            ebiten.KeyB,
	domain.KeyC:            ebiten.KeyC,
	domain.KeyD:            ebiten.KeyD,
	domain.KeyE:            ebiten.KeyE,
	domain.KeyF:            ebiten.KeyF,
	domain.KeyG:            ebiten.KeyG,
	domain.KeyH:            ebiten.KeyH,
	domain.KeyI:            ebiten.KeyI,
	domain.KeyJ:            ebiten.KeyJ,
	domain.KeyK:            ebiten.KeyK,
	domain.KeyL:            ebiten.KeyL,
	domain.KeyM:            ebiten.KeyM,
	domain.KeyN:            ebiten.KeyN,
	domain.KeyO:            ebiten.KeyO,
	domain.KeyP:            ebiten.KeyP,
	domain.KeyQ:            ebiten.KeyQ,
	domain.KeyR:            ebiten.KeyR,
	domain.KeyS:            ebiten.KeyS,
	domain.KeyT:            ebiten.KeyT,
	domain.KeyU:            ebiten.KeyU,
	domain.KeyV:            ebiten.KeyV,
	domain.KeyW:            ebiten.KeyW,
	domain.KeyX:            ebiten.KeyX,
	domain.KeyY:            ebiten.KeyY,
	domain.KeyZ:            ebiten.KeyZ,
	domain.Key0:            ebiten.Key0,
	domain.Key1:            ebiten.Key1,
	domain.Key2:            ebiten.Key2,
	domain.Key3:            ebiten.Key3,
	domain.Key4:            ebiten.Key4,
	domain.Key5:            ebiten.Key5,
	domain.Key6:            ebiten.Key6,
	domain.Key7:            ebiten.Key7,
	domain.Key8:            ebiten.Key8,
	domain.Key9:            ebiten.Key9,
	domain.KeyApostrophe:   ebiten.KeyApostrophe,
	domain.KeyBackslash:    ebiten.KeyBackslash,
	domain.KeyBackspace:    ebiten.KeyBackspace,
	domain.KeyCapsLock:     ebiten.KeyCapsLock,
	domain.KeyComma:        ebiten.KeyComma,
	domain.KeyDelete:       ebiten.KeyDelete,
	domain.KeyDown:         ebiten.KeyDown,
	domain.KeyEnd:          ebiten.KeyEnd,
	domain.KeyEnter:        ebiten.KeyEnter,
	domain.KeyEqual:        ebiten.KeyEqual,
	domain.KeyEscape:       ebiten.KeyEscape,
	domain.KeyF1:           ebiten.KeyF1,
	domain.KeyF2:           ebiten.KeyF2,
	domain.KeyF3:           ebiten.KeyF3,
	domain.KeyF4:           ebiten.KeyF4,
	domain.KeyF5:           ebiten.KeyF5,
	domain.KeyF6:           ebiten.KeyF6,
	domain.KeyF7:           ebiten.KeyF7,
	domain.KeyF8:           ebiten.KeyF8,
	domain.KeyF9:           ebiten.KeyF9,
	domain.KeyF10:          ebiten.KeyF10,
	domain.KeyF11:          ebiten.KeyF11,
	domain.KeyF12:          ebiten.KeyF12,
	domain.KeyGraveAccent:  ebiten.KeyGraveAccent,
	domain.KeyHome:         ebiten.KeyHome,
	domain.KeyInsert:       ebiten.KeyInsert,
	domain.KeyKP0:          ebiten.KeyKP0,
	domain.KeyKP1:          ebiten.KeyKP1,
	domain.KeyKP2:          ebiten.KeyKP2,
	domain.KeyKP3:          ebiten.KeyKP3,
	domain.KeyKP4:          ebiten.KeyKP4,
	domain.KeyKP5:          ebiten.KeyKP5,
	domain.KeyKP6:          ebiten.KeyKP6,
	domain.KeyKP7:          ebiten.KeyKP7,
	domain.KeyKP8:          ebiten.KeyKP8,
	domain.KeyKP9:          ebiten.KeyKP9,
	domain.KeyKPAdd:        ebiten.KeyKPAdd,
	domain.KeyKPDecimal:    ebiten.KeyKPDecimal,
	domain.KeyKPDivide:     ebiten.KeyKPDivide,
	domain.KeyKPEnter:      ebiten.KeyKPEnter,
	domain.KeyKPEqual:      ebiten.KeyKPEqual,
	domain.KeyKPMultiply:   ebiten.KeyKPMultiply,
	domain.KeyKPSubtract:   ebiten.KeyKPSubtract,
	domain.KeyLeft:         ebiten.KeyLeft,
	domain.KeyLeftAlt:      ebiten.KeyAlt,
	domain.KeyLeftBracket:  ebiten.KeyLeftBracket,
	domain.KeyLeftControl:  ebiten.KeyControl,
	domain.KeyLeftShift:    ebiten.KeyShift,
	domain.KeyMenu:         ebiten.KeyMenu,
	domain.KeyMinus:        ebiten.KeyMinus,
	domain.KeyNumLock:      ebiten.KeyNumLock,
	domain.KeyPageDown:     ebiten.KeyPageDown,
	domain.KeyPageUp:       ebiten.KeyPageUp,
	domain.KeyPause:        ebiten.KeyPause,
	domain.KeyPeriod:       ebiten.KeyPeriod,
	domain.KeyPrintScreen:  ebiten.KeyPrintScreen,
	domain.KeyRight:        ebiten.KeyRight,
	domain.KeyRightAlt:     ebiten.KeyAlt,
	domain.KeyRightBracket: ebiten.KeyRightBracket,
	domain.KeyRightControl: ebiten.KeyControl,
	domain.KeyRightShift:   ebiten.KeyShift,
	domain.KeyScrollLock:   ebiten.KeyScrollLock,
	domain.KeySemicolon:    ebiten.KeySemicolon,
	domain.KeySlash:        ebiten.KeySlash,
	domain.KeySpace:        ebiten.KeySpace,
	domain.KeyTab:          ebiten.KeyTab,
	domain.KeyUp:           ebiten.KeyUp,
}

// ─── Mouse Mapping Table ───
//
// domain.Mouse --> ebiten.Mouse
var ebitenMouseMap = map[domain.MouseButton]ebiten.MouseButton{
	domain.MouseButtonLeft:   ebiten.MouseButtonLeft,
	domain.MouseButtonRight:  ebiten.MouseButtonRight,
	domain.MouseButtonMiddle: ebiten.MouseButtonMiddle,
}

// ─── InputManager ───
//
// InputManager: "string" -> key code
type InputManager struct {
	actions map[string][]domain.Key
}

// NewInputManager creates a new InputManager instance.
//
// Purpose: Initializes an InputManager that handles physical key/mouse tracking and virtual action mapping.
//
// Outputs:
// - *InputManager: A newly initialized InputManager.
func NewInputManager() *InputManager {
	return &InputManager{
		actions: make(map[string][]domain.Key),
	}
}

// ─── Keyboard ───
//
// Update is called once per frame.
//
// Purpose: Provided to satisfy interface requirements. Ebiten processes inputs automatically behind the scenes, so no internal logic is needed here.
func (im *InputManager) Update() {}

// IsKeyPressed checks if a specific key is currently being held down.
//
// Purpose: Determines continuous key press states.
//
// Inputs:
// - key (domain.Key): The physical key to check.
//
// Outputs:
// - bool: True if the key is currently pressed, false otherwise.
func (im *InputManager) IsKeyPressed(key domain.Key) bool {
	if eKey, ok := ebitenKeyMap[key]; ok {
		return ebiten.IsKeyPressed(eKey)
	}
	return false
}

// IsKeyJustPressed checks if a specific key was pressed in the current frame.
//
// Purpose: Detects the initial press event (edge-trigger) of a key.
//
// Inputs:
// - key (domain.Key): The physical key to check.
//
// Outputs:
// - bool: True if the key transitioned from unpressed to pressed this frame.
func (im *InputManager) IsKeyJustPressed(key domain.Key) bool {
	if eKey, ok := ebitenKeyMap[key]; ok {
		return inpututil.IsKeyJustPressed(eKey)
	}
	return false
}

// IsKeyJustReleased checks if a specific key was released in the current frame.
//
// Purpose: Detects the release event (edge-trigger) of a key.
//
// Inputs:
// - key (domain.Key): The physical key to check.
//
// Outputs:
// - bool: True if the key transitioned from pressed to unpressed this frame.
func (im *InputManager) IsKeyJustReleased(key domain.Key) bool {
	if eKey, ok := ebitenKeyMap[key]; ok {
		return inpututil.IsKeyJustReleased(eKey)
	}
	return false
}

// ─── Mouse ────────────────────────────────────────────────────────────────────

// CursorPosition retrieves the current mouse cursor coordinates.
//
// Purpose: Returns the position of the mouse relative to the top-left corner of the window.
//
// Outputs:
// - (int, int): The X and Y pixel coordinates of the cursor.
func (im *InputManager) CursorPosition() (int, int) {
	return ebiten.CursorPosition()
}

// IsMouseButtonPressed checks if a specific mouse button is currently being held down.
//
// Purpose: Determines continuous mouse button press states.
//
// Inputs:
// - button (domain.MouseButton): The mouse button to check.
//
// Outputs:
// - bool: True if the button is currently pressed, false otherwise.
func (im *InputManager) IsMouseButtonPressed(button domain.MouseButton) bool {
	if eBtn, ok := ebitenMouseMap[button]; ok {
		return ebiten.IsMouseButtonPressed(eBtn)
	}
	return false
}

// IsMouseButtonJustPressed checks if a specific mouse button was pressed in the current frame.
//
// Purpose: Detects the initial press event (edge-trigger) of a mouse button.
//
// Inputs:
// - button (domain.MouseButton): The mouse button to check.
//
// Outputs:
// - bool: True if the button transitioned from unpressed to pressed this frame.
func (im *InputManager) IsMouseButtonJustPressed(button domain.MouseButton) bool {
	if eBtn, ok := ebitenMouseMap[button]; ok {
		return inpututil.IsMouseButtonJustPressed(eBtn)
	}
	return false
}

// IsMouseButtonJustReleased checks if a specific mouse button was released in the current frame.
//
// Purpose: Detects the release event (edge-trigger) of a mouse button.
//
// Inputs:
// - button (domain.MouseButton): The mouse button to check.
//
// Outputs:
// - bool: True if the button transitioned from pressed to unpressed this frame.
func (im *InputManager) IsMouseButtonJustReleased(button domain.MouseButton) bool {
	if eBtn, ok := ebitenMouseMap[button]; ok {
		return inpututil.IsMouseButtonJustReleased(eBtn)
	}
	return false
}

// WheelOffset retrieves the scroll offsets for the mouse wheel or touchpad.
//
// Purpose: Gets the X and Y movement of the scrolling device. A wrapper for ebiten.Wheel.
//
// Outputs:
// - (float64, float64): The X and Y scroll offsets. Returns (0, 0) if no scrolling occurred.
func (im *InputManager) WheelOffset() (float64, float64) {
	return ebiten.Wheel()
}

// ─── Virtual Action Mapping ───────────────────────────────────────────────────

// BindAction assigns one or more physical keys to a virtual action string.
//
// Purpose: Maps an action name to hardware keys so that game logic can check for the action rather than specific keys.
//
// Inputs:
// - action (string): The name of the virtual action (e.g., "jump").
// - keys (...domain.Key): A variadic list of physical keys that trigger this action.
func (im *InputManager) BindAction(action string, keys ...domain.Key) {
	im.actions[action] = keys
}

// IsActionPressed checks if any key bound to the specified action is currently being held down.
//
// Purpose: Evaluates continuous state for a virtual action.
//
// Inputs:
// - action (string): The name of the virtual action to check.
//
// Outputs:
// - bool: True if any associated key is pressed.
func (im *InputManager) IsActionPressed(action string) bool {
	keys, ok := im.actions[action]
	if !ok {
		return false
	}
	return slices.ContainsFunc(keys, im.IsKeyPressed)
}

// IsActionJustPressed checks if any key bound to the specified action was pressed in the current frame.
//
// Purpose: Evaluates the initial press event (edge-trigger) for a virtual action.
//
// Inputs:
// - action (string): The name of the virtual action to check.
//
// Outputs:
// - bool: True if any associated key transitioned from unpressed to pressed this frame.
func (im *InputManager) IsActionJustPressed(action string) bool {
	keys, ok := im.actions[action]
	if !ok {
		return false
	}
	return slices.ContainsFunc(keys, im.IsKeyJustPressed)
}

// IsActionJustReleased checks if any key bound to the specified action was released in the current frame.
//
// Purpose: Evaluates the release event (edge-trigger) for a virtual action.
//
// Inputs:
// - action (string): The name of the virtual action to check.
//
// Outputs:
// - bool: True if any associated key transitioned from pressed to unpressed this frame.
func (im *InputManager) IsActionJustReleased(action string) bool {
	keys, ok := im.actions[action]
	if !ok {
		return false
	}
	return slices.ContainsFunc(keys, im.IsKeyJustReleased)
}
