package enginetype

import (
	"github.com/yohamta/donburi"
)

// Purpose: Creates a new ECS entity in the scene and returns its donburi entry.
// Inputs:
//   - scene IScene: The scene where the entity will be created.
//   - components ...donburi.IComponentType: Variadic list of component types to attach.
// Outputs: *donburi.Entry - The ECS entry of the newly created entity.
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
// Purpose: Creates a new generic donburi component type and optionally registers it with a token.
// Inputs: token string - A short string identifier for the component (empty string skips registration).
// Outputs: *donburi.ComponentType[T] - The newly created component type.
func NewComponentType[T any](token string) *donburi.ComponentType[T] {
	comp := donburi.NewComponentType[T]()
	if token != "" {
		customComponentRegistry[token] = comp
	}
	return comp
}

// GetComponentType lấy component đã đăng ký dựa trên token.
// Purpose: Retrieves a registered component type by its string token.
// Inputs: token string - The short string identifier.
// Outputs: donburi.IComponentType - The component type, or nil if not found.
func GetComponentType(token string) donburi.IComponentType {
	return customComponentRegistry[token]
}

// ComponentInitializer đại diện cho hàm khởi tạo giá trị mặc định của một component trên Entry.
type ComponentInitializer func(entry *donburi.Entry)

var componentInitializerRegistry = make(map[string]ComponentInitializer)

// RegisterComponentInitializer đăng ký hàm khởi tạo mặc định cho một token component.
// Purpose: Registers a default initialization function for a specific component token.
// Inputs:
//   - token string: The component identifier.
//   - initFn ComponentInitializer: The initialization function to call when the component is added.
// Outputs: None.
func RegisterComponentInitializer(token string, initFn ComponentInitializer) {
	componentInitializerRegistry[token] = initFn
}

// InitializeComponent chạy hàm khởi tạo mặc định cho token nếu có.
// Purpose: Executes the registered default initialization function for a component on an entity entry.
// Inputs:
//   - token string: The component identifier.
//   - entry *donburi.Entry: The ECS entry being initialized.
// Outputs: None.
func InitializeComponent(token string, entry *donburi.Entry) {
	if initFn, ok := componentInitializerRegistry[token]; ok {
		initFn(entry)
	}
}
