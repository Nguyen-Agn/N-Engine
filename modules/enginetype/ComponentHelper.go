package enginetype

import (
	"hash/fnv"

	"github.com/yohamta/donburi"
)

// entryProvider là interface nội bộ để lấy *donburi.Entry từ bất kỳ IObject nào.
// nobject.Object implement interface này qua method Entry().
type entryProvider interface {
	Entry() *donburi.Entry
}

// HashString băm chuỗi thành uint64 dùng FNV-1a.
func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// SetComponent gán giá trị cho một custom component trên object.
// comp phải là *donburi.ComponentType[T] đã được tạo bằng napi.NewComponentType[T]().
//
// Ví dụ:
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//	napi.SetComponent(player, StatsComp, StatsData{Health: 100, Mana: 50})
func SetComponent[T any](obj IObject, comp *donburi.ComponentType[T], value T) {
	ep, ok := obj.(entryProvider)
	if !ok {
		return
	}
	donburi.SetValue(ep.Entry(), comp, value)
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
func GetComponent[T any](obj IObject, comp *donburi.ComponentType[T]) *T {
	ep, ok := obj.(entryProvider)
	if !ok {
		return nil
	}
	entry := ep.Entry()
	if !entry.HasComponent(comp) {
		return nil
	}
	return donburi.Get[T](entry, comp)
}

// AddComponentType thêm component type vào ECS entry của object.
// Dùng khi muốn thêm custom component sau khi object đã được tạo.
//
// Ví dụ:
//
//	napi.AddComponentType(player, StatsComp)
//	napi.SetComponent(player, StatsComp, StatsData{Health: 100})
func AddComponentType[T any](obj IObject, comp *donburi.ComponentType[T]) {
	ep, ok := obj.(entryProvider)
	if !ok {
		return
	}
	ep.Entry().AddComponent(comp)
}
