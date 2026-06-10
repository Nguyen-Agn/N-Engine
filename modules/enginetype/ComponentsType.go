package enginetype

import (
	"reflect"
	"strconv"
	"strings"

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

// componentSetterRegistry lưu trữ các hàm closure reflection để set field linh hoạt
var componentSetterRegistry = make(map[string]func(entry *donburi.Entry, fieldName string, valueStr string))

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
		
		componentSetterRegistry[token] = func(entry *donburi.Entry, fieldName string, valueStr string) {
			if !entry.HasComponent(comp) {
				return
			}
			
			dataPtr := donburi.Get[T](entry, comp)
			
			// Use reflect to set the field by name (case-insensitive)
			v := reflect.ValueOf(dataPtr)
			if v.Kind() != reflect.Ptr || v.IsNil() {
				return
			}
			v = v.Elem()
			if v.Kind() != reflect.Struct {
				return
			}

			var field reflect.Value
			for i := 0; i < v.NumField(); i++ {
				typeField := v.Type().Field(i)
				if strings.EqualFold(typeField.Name, fieldName) {
					field = v.Field(i)
					break
				}
			}

			if !field.IsValid() || !field.CanSet() {
				return
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(valueStr)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intVal, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
					field.SetInt(intVal)
				}
			case reflect.Float32, reflect.Float64:
				if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
					field.SetFloat(floatVal)
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(valueStr); err == nil {
					field.SetBool(boolVal)
				}
			}
		}
	}
	return comp
}

// SetComponentFieldByToken tìm component bằng token và gán field bằng reflection.
// Trả về true nếu gán thành công.
func SetComponentFieldByToken(entry *donburi.Entry, token string, fieldName string, valueStr string) bool {
	if setter, ok := componentSetterRegistry[token]; ok {
		setter(entry, fieldName, valueStr)
		return true
	}
	return false
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
