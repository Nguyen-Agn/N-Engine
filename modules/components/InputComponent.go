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

func (c *InputComponent) BindComponent(base IObject) {
	c.IObject = base
	c.data = enginetype.GetComponent(base, enginetype.Input)
}

// ListenOn đăng ký một handler được gọi khi phím (hoặc nhóm phím) kích hoạt theo loại sự kiện.
//
// Tham số:
//   - key: tên phím hoặc nhóm phím, hỗ trợ nhận nhiều phím cách nhau bằng dấu cách (VD: "w a s d alpha")
//     "alpha" — bất kỳ phím chữ nào (a-z)
//     "number" — bất kỳ phím số nào (0-9)
//     "arrows" — bất kỳ phím mũi tên nào
//     "wasd" — W/A/S/D
//     "all" — bất kỳ phím nào
//   - eventType: loại sự kiện ("pressed", "just_pressed", "just_released")
//   - handler: hàm nhận tên phím đã trigger (ví dụ: "w", "space", "a"...)
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
