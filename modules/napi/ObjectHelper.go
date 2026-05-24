package napi

import (
	"reflect"
	"strings"
	"sync/atomic"

	"autoworld/modules/components"
	"autoworld/modules/enginetype"
	"autoworld/modules/nobject"

	"github.com/yohamta/donburi"
)

// ─── Object Interface ─────────────────────────────────────────────────────────

// Object là interface tối thiểu mà Custom Object phải implement để dùng với NewObject.
// Game object chỉ cần có Create() — các lifecycle khác (StepUpdate, Destroy) là optional override.
type Object interface {
	Create()
}

// ─── Object ID Counter ────────────────────────────────────────────────────────

// objectIDCounter là bộ đếm ID tự tăng thread-safe cho mọi Object được tạo.
var objectIDCounter int64 = 0

// getNextObjectID trả về ID tiếp theo (bắt đầu từ 1), thread-safe dùng atomic.
func getNextObjectID() int {
	return int(atomic.AddInt64(&objectIDCounter, 1))
}

// ─── NewObject ────────────────────────────────────────────────────────────────

// NewObject tạo game object theo component code, inject vào target và tự động
// khởi tạo giá trị mặc định cho từng component.
// Sau khi gọi NewObject, hãy gọi Register/RegisterGlobal để đưa object vào update loop.
//
// target phải là con trỏ tới struct có nhúng napi.IObject và các Component Mixin.
// name là tên định danh — object sẽ được lưu vào global store để tra cứu sau.
//
// componentCode là chuỗi token cách nhau bởi khoảng trắng:
//
//		pos  — Position component (X, Y)
//		spr  — Sprite component   (hình ảnh, animation)
//		box  — Box component      (hitbox, collision)
//		aud  — Audio component    (âm thanh)
//		dir  — Direction component (góc hướng)
//		glo  — modifier: tạo object trong Global Scene (persistent)
//	 inp  — Input component (lắng nghe sự kiện từ bàn phím, chuột)
//	 bg   — BackGround component (ảnh nền)
//	 til  — TileMap component (vẽ nền tilemap)
//	 alr  — Alarm component (chạy hàm theo hẹn giờ)
//	 vel  — Velocỉy component (vận tốc)
//	 twn  — Tween component (gia ảnh)
//
// inf (Infor) được thêm tự động vào mọi object.
//
// Ví dụ:
//
//	p := &Player{}
//	napi.NewObject(p, "player", "pos spr")
//	napi.Register(p, false)
func NewObject(target Object, name string, componentCode string) {
	tokens := strings.Fields(componentCode)
	tokenSet, isGlobal := filter(tokens)

	// Lấy map từ scene phù hợp để tạo entry
	targetMap := getScene(isGlobal).GetMap()

	// Xây danh sách component types từ token
	componentsList := []donburi.IComponentType{}
	for token := range tokenSet {
		comp := GetComponentType(token)
		if comp != nil {
			componentsList = append(componentsList, comp)
		}
	}

	// Tạo ECS entry trong world của map
	entry := targetMap.World().Entry(targetMap.World().Create(componentsList...))

	// Khởi tạo data mặc định từ registry cho từng component
	for token := range tokenSet {
		enginetype.InitializeComponent(token, entry)
	}

	// Luôn ghi InforData (bắt buộc): ID tự tăng + name
	nextID := getNextObjectID()
	donburi.SetValue(entry, enginetype.Infor, components.InforData{
		Id:   nextID,
		Name: name,
	})

	obj := nobject.NewObject(entry)

	// Lưu vào global store nếu có tên (để tra cứu sau bằng napi.GetObject)
	if name != "" {
		engine().Store.AddObject(name, obj)
	}

	// Inject IObject và tất cả Component Mixin vào target
	bind(target, obj)
}

// ─── Register Helpers ─────────────────────────────────────────────────────────

// Register đăng ký một IObject vào scene's update loop.
//
// scene: name of scene, if nil -> global var
func Register(obj IObject, scene string) {
	var _scene interface{ AddObject(IObject) }
	if scene != "" {
		_scene = GetSceneByID(scene)
	} else {
		_scene = GetGlobalScene()
	}

	if _scene != nil {
		_scene.AddObject(obj)
	}
}

// RegisterIn đăng ký một IObject vào scene cụ thể do caller chỉ định.
func RegisterIn(scene IScene, obj IObject) {
	if scene != nil {
		scene.AddObject(obj)
	}
}

// ─── Internal Helpers ─────────────────────────────────────────────────────────

// getScene trả về scene phù hợp: Global Scene nếu global=true, Current Scene nếu false.
func getScene(global bool) IScene {
	if global {
		return GetGlobalScene()
	}
	return GetCurrentScene()
}

// filter phân tích danh sách token, tách modifier "glo" và đảm bảo "inf" luôn có mặt.
// Trả về: set token duy nhất và cờ isGlobal.
func filter(tokens []string) (map[string]bool, bool) {
	tokenSet := make(map[string]bool, len(tokens))
	var global bool = false
	for _, t := range tokens {
		if t == "glo" {
			global = true
		}
		tokenSet[t] = true
	}

	// Infor là component bắt buộc đối với mọi Object
	tokenSet["inf"] = true
	return tokenSet, global
}

// bind inject IObject vào target struct và tất cả Component Mixin nhúng trong nó.
// Dùng reflection để tự động gán — game dev không cần khởi tạo thủ công từng mixin.
//
// Quy tắc inject:
//  1. Field tên "IObject" → gán trực tiếp IObject base.
//  2. Field là struct có sub-field "IObject" → gán vào sub-field đó.
//  3. Field implement interface BindComponent(IObject) → gọi BindComponent.
func bind(target any, base IObject) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		panic("[napi.Bind] target phải là con trỏ tới struct (ví dụ: *Player)")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		// 1. Gán IObject chính của Custom Object
		if fieldType.Name == "IObject" {
			if fieldVal.CanSet() {
				fieldVal.Set(reflect.ValueOf(base))
			}
			continue
		}

		// 2. Gán IObject cho các Component Mixin nhúng (PositionComponent, SpriteComponent…)
		if fieldVal.Kind() == reflect.Struct {
			subField := fieldVal.FieldByName("IObject")
			if subField.IsValid() && subField.CanSet() {
				subField.Set(reflect.ValueOf(base))
			}
		}

		// 3. Gọi BindComponent nếu mixin implement interface đó
		if fieldType.PkgPath == "" && fieldVal.CanAddr() {
			if binder, ok := fieldVal.Addr().Interface().(interface{ BindComponent(IObject) }); ok {
				binder.BindComponent(base)
				continue
			}
		}
	}
}
