package enginetype

import (
	"github.com/yohamta/donburi"
)

func NewEntry(scene IScene, components ...donburi.IComponentType) *donburi.Entry {
	return scene.World().Entry(scene.World().Create(components...))
}

// Registry lưu trữ các custom component type để dùng với componentCode string
var customComponentRegistry = make(map[string]donburi.IComponentType)

// NewComponentType tạo một donburi component type mới cho game dev dùng.
// Nếu token != "", component sẽ được đăng ký để có thể dùng trong chuỗi componentCode
// của NewObject hoặc NewBaseObject (ví dụ: "sta").
//
// Ví dụ:
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//	obj := napi.NewBaseObject("hero", "pos spr sta")
func NewComponentType[T any](token string) *donburi.ComponentType[T] {
	comp := donburi.NewComponentType[T]()
	if token != "" {
		customComponentRegistry[token] = comp
	}
	return comp
}

// GetComponentType lấy component đã đăng ký dựa trên token.
func GetComponentType(token string) donburi.IComponentType {
	return customComponentRegistry[token]
}

// ComponentInitializer đại diện cho hàm khởi tạo giá trị mặc định của một component trên Entry.
type ComponentInitializer func(entry *donburi.Entry)

var componentInitializerRegistry = make(map[string]ComponentInitializer)

// RegisterComponentInitializer đăng ký hàm khởi tạo mặc định cho một token component.
func RegisterComponentInitializer(token string, initFn ComponentInitializer) {
	componentInitializerRegistry[token] = initFn
}

// InitializeComponent chạy hàm khởi tạo mặc định cho token nếu có.
func InitializeComponent(token string, entry *donburi.Entry) {
	if initFn, ok := componentInitializerRegistry[token]; ok {
		initFn(entry)
	}
}
