package nsystem

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"
)

// InputSystem duyệt tất cả Object có InputData và kích hoạt Handler
// tương ứng mỗi frame khi phím được giữ.
//
// OR logic: nếu bất kỳ phím nào trong KeyBinding.Keys đang được nhấn,
// Handler sẽ được gọi đúng một lần.
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

// Update duyệt objectList, kiểm tra phím và kích hoạt callback tương ứng.
func (s *InputSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		data := enginetype.GetComponent(obj, enginetype.Input)
		if data == nil {
			continue
		}
		for _, binding := range data.Bindings {
			for _, key := range binding.Keys {
				if s.input.IsKeyPressed(key) {
					binding.Handler()
					break // OR logic: 1 phím khớp là đủ
				}
			}
		}
	}
}
