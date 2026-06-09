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
// Purpose: Hashes a string into a uint64 value using the FNV-1a algorithm.
// Inputs: s string - The string to hash.
// Outputs: uint64 - The resulting hash value.
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
// Purpose: Assigns a value to a custom donburi component on an object.
// Inputs:
//   - obj IObject: The target game object.
//   - comp *donburi.ComponentType[T]: The component type definition.
//   - value T: The data to assign to the component.
// Outputs: None.
// Special requirements: The object must implement entryProvider to access its ECS entry.
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
// Purpose: Retrieves a pointer to a custom component's data on an object.
// Inputs:
//   - obj IObject: The target game object.
//   - comp *donburi.ComponentType[T]: The component type definition.
// Outputs: *T - Pointer to the component data, or nil if the object lacks this component.
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
// Purpose: Dynamically adds a new component type to an existing object's ECS entry.
// Inputs:
//   - obj IObject: The target game object.
//   - comp *donburi.ComponentType[T]: The component type to add.
// Outputs: None.
func AddComponentType[T any](obj IObject, comp *donburi.ComponentType[T]) {
	ep, ok := obj.(entryProvider)
	if !ok {
		return
	}
	ep.Entry().AddComponent(comp)
}
