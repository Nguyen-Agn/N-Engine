package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// GenericComponent là mixin generic cho phép game dev tạo Component Mixin riêng
// mà vẫn được engine tự động bind, không cần thao tác ECS thủ công.
// Thường được re-export qua napi — dev chỉ cần import napi để dùng.
//
// # Cách dùng cơ bản (chỉ Get/Set) — qua napi
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//
//	type Hero struct {
//	    napi.IObject
//	    napi.IPosition
//	    napi.GenericComponent[StatsData]
//	}
//
//	func NewHero() *Hero {
//	    h := &Hero{
//	        GenericComponent: napi.NewGenericComponent(StatsComp),
//	    }
//	    napi.NewObject(h, "hero", "pos sta")
//	    napi.Register(h, "")
//	    return h
//	}
//
//	func (h *Hero) OnStep() {
//	    h.Get().Health -= 1
//	}
//
// # Cách dùng nâng cao (thêm method riêng) — qua napi
//
//	type StatsComponent struct {
//	    napi.GenericComponent[StatsData]
//	}
//
//	func (s *StatsComponent) TakeDamage(amount int) {
//	    if s.Get() != nil {
//	        s.Get().Health -= amount
//	    }
//	}
//
//	type Hero struct {
//	    napi.IObject
//	    napi.IPosition
//	    StatsComponent   // method TakeDamage được promoted lên Hero
//	}
//
//	func NewHero() *Hero {
//	    h := &Hero{
//	        StatsComponent: StatsComponent{
//	            GenericComponent: napi.NewGenericComponent(StatsComp),
//	        },
//	    }
//	    napi.NewObject(h, "hero", "pos sta")
//	    napi.Register(h, "")
//	    return h
//	}
type GenericComponent[T any] struct {
	IObject
	comp *donburi.ComponentType[T]
	data *T
}

// NewGenericComponent tạo GenericComponent đã gắn sẵn ComponentType.
// Phải gọi trước napi.NewObject để BindComponent hoạt động đúng.
func NewGenericComponent[T any](comp *donburi.ComponentType[T]) GenericComponent[T] {
	return GenericComponent[T]{comp: comp}
}

// BindComponent được gọi tự động bởi napi.bind() bên trong napi.NewObject.
// Gán IObject và lấy con trỏ data từ ECS entry của object.
// Game dev không cần gọi hàm này trực tiếp.
func (c *GenericComponent[T]) BindComponent(base IObject) {
	c.IObject = base
	if c.comp != nil {
		c.data = enginetype.GetComponent(base, c.comp)
	}
}

// Get trả về con trỏ trực tiếp đến data của component trong ECS.
// Thay đổi qua pointer sẽ tác động ngay lên ECS data.
// Trả về nil nếu component chưa được bind hoặc object không có component này.
//
// Ví dụ:
//
//	s.Get().Health -= 10
func (c *GenericComponent[T]) Data() *T {
	return c.data
}

// Set gán toàn bộ giá trị mới cho component data.
// Không làm gì nếu component chưa được bind.
//
// Ví dụ:
//
//	s.Set(StatsData{Health: 100, Speed: 3.0})
func (c *GenericComponent[T]) SetData(val T) {
	if c.data != nil {
		*c.data = val
	}
}
