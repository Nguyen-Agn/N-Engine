package components

import (
	"strings"

	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

func init() {
	enginetype.RegisterComponentInitializer("inp", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Input, domain.InputData{
			Bindings: []domain.KeyBinding{},
		})
	})
}

// InputComponent là Mixin để nhúng vào Custom Object.
// Cung cấp hàm ListenOn để đăng ký lắng nghe phím bàn phím với đầy đủ loại sự kiện.
// Yêu cầu Object phải khai báo token "inp" khi gọi napi.NewObject.
//
// Cách dùng:
//
//	type Player struct {
//	    napi.IObject
//	    napi.IPosition
//	    components.InputComponent
//	}
//
//	func (p *Player) OnCreate() {
//	    // Giữ phím W mỗi frame
//	    p.ListenOn("w", "pressed", func(key string) { p.MoveUp() })
//	    // Nhấn Space 1 lần
//	    p.ListenOn("space", "just_pressed", func(key string) { p.Jump() })
//	    // Thả Shift
//	    p.ListenOn("shift", "just_released", func(key string) { p.StopRun() })
//	    // Nhóm phím: handler nhận tên phím cụ thể đã trigger
//	    p.ListenOn("alpha", "just_pressed", func(key string) {
//	        if key == "a" { ... }
//	    })
//	}
type InputComponent struct {
	IObject
	data *InputData
}

// BindComponent binds the base object and its corresponding ECS data to the component.
// Purpose: Automatically invoked by the engine to fetch the input data pointer.
// Inputs:
//   - base: The IObject representing the base entity to bind.
// Outputs: None.
func (c *InputComponent) BindComponent(base IObject) {
	c.IObject = base
	c.data = enginetype.GetComponent(base, enginetype.Input)
}

// ListenOn registers a handler to be called when a specific key or key group triggers an event.
// Purpose: Subscribes to keyboard input events.
// Inputs:
//   - key: The string name of the key or group (e.g., "w a s d alpha", "all").
//   - eventType: The type of event ("pressed", "just_pressed", "just_released").
//   - handler: The function to execute when the event occurs. It receives the triggered key name.
// Outputs: None.
func (c *InputComponent) ListenOn(key string, eventType string, handler func(key string)) {
	if c.data == nil {
		return
	}

	evt, ok := domain.EventTypeNameMap[eventType]
	if !ok {
		evt = domain.EventPressed // fallback
	}

	tokens := strings.FieldsSeq(key)
	for token := range tokens {
		// Thử tra nhóm đặc biệt trước
		if group, ok := domain.KeyGroupMap[token]; ok {
			c.data.Bindings = append(c.data.Bindings, domain.KeyBinding{
				Keys:      group,
				EventType: evt,
				Handler:   handler,
			})
			continue
		}

		// Tra phím đơn
		if k, ok := domain.KeyNameMap[token]; ok {
			c.data.Bindings = append(c.data.Bindings, domain.KeyBinding{
				Keys:      []domain.Key{k},
				EventType: evt,
				Handler:   handler,
			})
		}
	}
}
