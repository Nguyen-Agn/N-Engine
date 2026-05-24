package components

import (
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
// Cung cấp hàm ListenOn để đăng ký lắng nghe phím và gán hàm xử lý tương ứng.
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
//	func (p *Player) Create() {
//	    p.ListenOn("w", p.MoveUp)
//	    p.ListenOn("number", p.OnNumberKey)  // nhóm phím
//	}
type InputComponent struct {
	IObject
	data *InputData
}

func (c *InputComponent) BindComponent(base IObject) {
	c.IObject = base
	c.data = enginetype.GetComponent(base, enginetype.Input)
}

// ListenOn đăng ký một handler sẽ được gọi mỗi frame khi phím (hoặc nhóm phím) được giữ.
// key là tên phím đơn ("w", "space", "enter"...) hoặc tên nhóm đặc biệt:
//   - "alpha"   — bất kỳ phím chữ nào (a-z)
//   - "number"  — bất kỳ phím số nào (0-9)
//   - "arrows"  — bất kỳ phím mũi tên nào
//   - "wasd"    — bất kỳ phím W/A/S/D nào
//   - "all"     — bất kỳ phím nào (dùng cho nhập liệu)
func (c *InputComponent) ListenOn(key string, handler func()) {
	if c.data == nil {
		return
	}

	// Thử tra nhóm đặc biệt trước
	if group, ok := domain.KeyGroupMap[key]; ok {
		c.data.Bindings = append(c.data.Bindings, domain.KeyBinding{
			Keys:    group,
			Handler: handler,
		})
		return
	}

	// Tra phím đơn
	if k, ok := domain.KeyNameMap[key]; ok {
		c.data.Bindings = append(c.data.Bindings, domain.KeyBinding{
			Keys:    []domain.Key{k},
			Handler: handler,
		})
	}
}
