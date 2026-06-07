package components

import (
	"strings"

	"autoworld/domain"
)

// globalInputManager là reference singleton tới IInputManager.
// Được inject từ core khi engine khởi động qua SetGlobalInputManager.
var globalInputManager domain.IInputManager

// SetGlobalInputManager inject IInputManager vào package components.
// Gọi từ core khi khởi tạo engine, trước khi bất kỳ MouseComponent nào được dùng.
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

// BindComponent được gọi tự động bởi engine khi bind mixin vào object.
// MouseComponent không cần dữ liệu từ IObject nên chỉ no-op.
func (c *MouseComponent) BindComponent(_ IObject) {}

// MouseX trả về tọa độ ngang của con trỏ chuột (pixel, tính từ góc trên-trái màn hình).
func (c *MouseComponent) MouseX() int {
	if globalInputManager == nil {
		return 0
	}
	x, _ := globalInputManager.CursorPosition()
	return x
}

// MouseY trả về tọa độ dọc của con trỏ chuột (pixel, tính từ góc trên-trái màn hình).
func (c *MouseComponent) MouseY() int {
	if globalInputManager == nil {
		return 0
	}
	_, y := globalInputManager.CursorPosition()
	return y
}

// WheelX trả về độ cuộn bánh xe chuột theo trục X trong frame hiện tại.
// Dương = cuộn phải, âm = cuộn trái.
func (c *MouseComponent) WheelX() float64 {
	if globalInputManager == nil {
		return 0
	}
	wx, _ := globalInputManager.WheelOffset()
	return wx
}

// WheelY trả về độ cuộn bánh xe chuột theo trục Y trong frame hiện tại.
// Dương = cuộn xuống, âm = cuộn lên.
func (c *MouseComponent) WheelY() float64 {
	if globalInputManager == nil {
		return 0
	}
	_, wy := globalInputManager.WheelOffset()
	return wy
}

// ListenMouseOn đăng ký một handler được gọi theo loại sự kiện khi nút chuột kích hoạt.
//
// Tham số:
//   - button: tên nút chuột ("left", "right", "middle") hoặc nhiều nút cách nhau bằng dấu cách
//   - eventType: "pressed" (giữ), "just_pressed" (vừa nhấn), "just_released" (vừa thả)
//   - handler: hàm nhận tên nút chuột đã trigger (ví dụ: "left")
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

// UpdateMouseBindings được gọi mỗi frame bởi InputSystem để xử lý mouse bindings.
// Không cần gọi trực tiếp từ game code.
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

// checkMouseButton kiểm tra trạng thái nút chuột theo EventType.
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
