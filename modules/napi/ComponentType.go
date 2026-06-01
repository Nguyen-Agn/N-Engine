package napi

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// ─── ECS Entry Helper ─────────────────────────────────────────────────────────

// NewEntry tạo một ECS entity (donburi.Entry) trong World của scene với các component types cho sẵn.
// Dùng khi cần tạo entity thủ công thay vì qua NewObject.
//
// Ví dụ:
//
//	entry := napi.NewEntry(scene, napi.Position, napi.Sprite)
func newEntry(scene IScene, components ...donburi.IComponentType) *donburi.Entry {
	return scene.World().Entry(scene.World().Create(components...))
}

// ─── Custom Component Type ────────────────────────────────────────────────────

// NewComponentType tạo một donburi component type mới cho game dev tùy biến.
// Nếu token != "", component sẽ được đăng ký vào registry và có thể dùng
// trong componentCode của NewObject (ví dụ: "pos spr sta").
//
// Ví dụ:
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//	napi.NewObject(hero, "hero", "pos spr sta")
func newComponentType[T any](token string) *donburi.ComponentType[T] {
	return enginetype.NewComponentType[T](token)
}

// GetComponentType tra cứu component type đã đăng ký theo token.
// Trả về nil nếu token chưa được đăng ký.
func getComponentType(token string) donburi.IComponentType {
	return enginetype.GetComponentType(token)
}
