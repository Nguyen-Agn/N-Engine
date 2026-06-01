package napi

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// SetComponent gán giá trị cho một custom component trên object.
// comp phải là *donburi.ComponentType[T] đã được tạo bằng napi.NewComponentType[T]().
//
// Ví dụ:
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//	napi.SetComponent(player, StatsComp, StatsData{Health: 100, Mana: 50})
func setComponent[T any](obj IObject, comp *donburi.ComponentType[T], value T) {
	enginetype.SetComponent(obj, comp, value)
}

// GetComponent lấy con trỏ đến data của một custom component trên object.
// Trả về nil nếu object không có component này.
//
// Ví dụ:
//
//	stats := napi.GetComponent(player, StatsComp)
//	if stats != nil {
//	    stats.Health -= 10
//	}
func getComponent[T any](obj IObject, comp *donburi.ComponentType[T]) *T {
	return enginetype.GetComponent(obj, comp)
}

// AddComponentType thêm component type vào ECS entry của object.
// Dùng khi muốn thêm custom component sau khi object đã được tạo.
//
// Ví dụ:
//
//	napi.AddComponentType(player, StatsComp)
//	napi.SetComponent(player, StatsComp, StatsData{Health: 100})
func addComponentType[T any](obj IObject, comp *donburi.ComponentType[T]) {
	enginetype.AddComponentType(obj, comp)
}
