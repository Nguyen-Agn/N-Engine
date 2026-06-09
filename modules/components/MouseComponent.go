package components

import (
	"strings"

	"autoworld/domain"
)

// globalInputManager là reference singleton tới IInputManager.
// Được inject từ core khi engine khởi động qua SetGlobalInputManager.
var globalInputManager domain.IInputManager

// SetGlobalInputManager injects the global input manager into the components package.
// Purpose: Initializes the singleton input manager for mouse components.
// Inputs:
//   - m: The IInputManager instance to inject.
// Outputs: None.
// Special Requirements: Must be called by core during engine initialization before any MouseComponent is used.
func SetGlobalInputManager(m domain.IInputManager) {
	globalInputManager = m
}

// MouseComponent là Mixin để nhúng vào Custom Object.
// Không cần token ECS — lưu mouse bindings trực tiếp trong struct.
// Cung cấp truy cập tọa độ chuột, tốc độ cuộn và lắng nghe sự kiện nút chuột.
//
// Cách dùng:
//
//	type Player struct {
//	    napi.IObject
//	    ncom.Mouse   // không cần token "inp" hay token riêng
//	}
//
//	func (p *Player) OnCreate() {
//	    // Nhấn chuột trái 1 lần
//	    p.ListenMouseOn("left", "just_pressed", func(btn string) { p.Shoot() })
//	    // Giữ chuột phải mỗi frame
//	    p.ListenMouseOn("right", "pressed", func(btn string) { p.Aim() })
//	    // Thả nút giữa
//	    p.ListenMouseOn("middle", "just_released", func(btn string) { p.PlaceMarker() })
//	}
//
//	func (p *Player) OnStep() {
//	    mx, my := p.MouseX(), p.MouseY()
//	    wy := p.WheelY()   // cuộn lên/xuống
//	    _ = mx; _ = my; _ = wy
//	}
type MouseComponent struct {
	bindings []domain.MouseBinding
	// MouseComponent tự update bindings mỗi frame (được gọi từ InputSystem hoặc LogicSystem).
}

// BindComponent is a no-op for MouseComponent as it stores bindings directly.
// Purpose: Satisfies the IComponent interface requirement.
// Inputs:
//   - _: The IObject representing the base entity (ignored).
// Outputs: None.
func (c *MouseComponent) BindComponent(_ IObject) {}

// MouseX retrieves the current horizontal cursor position.
// Purpose: Gets the X coordinate of the mouse pointer.
// Inputs: None.
// Outputs: The X coordinate in pixels from the top-left of the screen (int). Returns 0 if input manager is nil.
func (c *MouseComponent) MouseX() int {
	if globalInputManager == nil {
		return 0
	}
	x, _ := globalInputManager.CursorPosition()
	return x
}

// MouseY retrieves the current vertical cursor position.
// Purpose: Gets the Y coordinate of the mouse pointer.
// Inputs: None.
// Outputs: The Y coordinate in pixels from the top-left of the screen (int). Returns 0 if input manager is nil.
func (c *MouseComponent) MouseY() int {
	if globalInputManager == nil {
		return 0
	}
	_, y := globalInputManager.CursorPosition()
	return y
}

// WheelX retrieves the horizontal mouse wheel scroll offset.
// Purpose: Gets the amount of horizontal scrolling in the current frame.
// Inputs: None.
// Outputs: The horizontal scroll offset (float64). Positive means scrolling right, negative means left.
func (c *MouseComponent) WheelX() float64 {
	if globalInputManager == nil {
		return 0
	}
	wx, _ := globalInputManager.WheelOffset()
	return wx
}

// WheelY retrieves the vertical mouse wheel scroll offset.
// Purpose: Gets the amount of vertical scrolling in the current frame.
// Inputs: None.
// Outputs: The vertical scroll offset (float64). Positive means scrolling down, negative means up.
func (c *MouseComponent) WheelY() float64 {
	if globalInputManager == nil {
		return 0
	}
	_, wy := globalInputManager.WheelOffset()
	return wy
}

// ListenMouseOn registers a handler to be called when a mouse button triggers an event.
// Purpose: Subscribes to mouse button input events.
// Inputs:
//   - button: The string name of the mouse button ("left", "right", "middle") or multiple space-separated buttons.
//   - eventType: The type of event ("pressed", "just_pressed", "just_released").
//   - handler: The function to execute when the event occurs. It receives the triggered button name.
// Outputs: None.
func (c *MouseComponent) ListenMouseOn(button string, eventType string, handler func(button string)) {
	evt, ok := domain.EventTypeNameMap[eventType]
	if !ok {
		evt = domain.EventPressed // fallback
	}

	tokens := strings.FieldsSeq(button)
	for token := range tokens {
		btn, ok := domain.MouseButtonNameMap[token]
		if !ok {
			continue
		}
		c.bindings = append(c.bindings, domain.MouseBinding{
			Button:    btn,
			EventType: evt,
			Handler:   handler,
		})
	}
}

// UpdateMouseBindings processes all registered mouse bindings for the current frame.
// Purpose: Triggers handlers for mouse events that occurred this frame.
// Inputs: None.
// Outputs: None.
// Special Requirements: Called internally by the engine (InputSystem); should not be called from game code.
func (c *MouseComponent) UpdateMouseBindings() {
	if globalInputManager == nil {
		return
	}
	for _, binding := range c.bindings {
		if checkMouseButton(globalInputManager, binding.Button, binding.EventType) {
			btnName := domain.MouseButtonReverseMap[binding.Button]
			binding.Handler(btnName)
		}
	}
}

// checkMouseButton evaluates if a specific mouse button event occurred.
// Purpose: Checks the state of a mouse button against a desired event type.
// Inputs:
//   - m: The IInputManager instance to use for checking.
//   - btn: The domain.MouseButton to check.
//   - evt: The domain.EventType to test for.
// Outputs: True if the specified event occurred for the button, false otherwise.
func checkMouseButton(m domain.IInputManager, btn domain.MouseButton, evt domain.EventType) bool {
	switch evt {
	case domain.EventPressed:
		return m.IsMouseButtonPressed(btn)
	case domain.EventJustPressed:
		return m.IsMouseButtonJustPressed(btn)
	case domain.EventJustReleased:
		return m.IsMouseButtonJustReleased(btn)
	}
	return false
}
