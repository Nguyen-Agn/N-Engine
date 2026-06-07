package nsystem

import (
	"autoworld/domain"
	"autoworld/modules/components"
	"autoworld/modules/enginetype"
)

// InputSystem duyệt tất cả Object có InputData và kích hoạt Handler
// tương ứng mỗi frame theo loại sự kiện (Pressed, JustPressed, JustReleased).
//
// OR logic cho nhóm phím: nếu bất kỳ phím nào trong KeyBinding.Keys khớp,
// Handler sẽ được gọi đúng một lần với tên phím đã trigger.
//
// Mouse bindings: InputSystem gọi UpdateMouseBindings() trên mọi object
// có MouseComponent được nhúng (phát hiện qua interface IMouse).
//
// InputSystem nhận domain.IInputManager qua interface để tránh
// phụ thuộc trực tiếp vào module core (nguyên tắc Dependency Inversion).
type InputSystem struct {
	input domain.IInputManager
}

// NewInputSystem tạo InputSystem với IInputManager từ Engine.
func NewInputSystem(input domain.IInputManager) *InputSystem {
	return &InputSystem{input: input}
}

// Update duyệt objectList, xử lý keyboard bindings (theo EventType) và mouse bindings.
func (s *InputSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		// ── Keyboard ──────────────────────────────────────────────────────────
		data := enginetype.GetComponent(obj, enginetype.Input)
		if data != nil {
			for _, binding := range data.Bindings {
				for _, key := range binding.Keys {
					if s.checkKey(key, binding.EventType) {
						keyName := domain.KeyReverseMap[key]
						binding.Handler(keyName) // OR logic: gọi 1 lần với phím đầu tiên khớp
						break
					}
				}
			}
		}

		// ── Mouse ─────────────────────────────────────────────────────────────
		// Bất kỳ object nào nhúng MouseComponent đều implement IMouseUpdater,
		// InputSystem gọi UpdateMouseBindings() để xử lý bindings đã đăng ký.
		if mu, ok := obj.(iMouseUpdater); ok {
			mu.UpdateMouseBindings()
		}
	}
}

// checkKey kiểm tra trạng thái phím theo EventType.
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

// iMouseUpdater là interface nội bộ để InputSystem phát hiện object có MouseComponent.
// Không export ra ngoài để tránh lộ chi tiết triển khai.
type iMouseUpdater interface {
	UpdateMouseBindings()
}

// Đảm bảo MouseComponent implement iMouseUpdater tại compile time.
var _ iMouseUpdater = (*components.MouseComponent)(nil)
