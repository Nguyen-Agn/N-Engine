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

type objGroup struct{}

// Obj lÃ  nhÃ³m hÃ m táº¡o vÃ  Ä‘Äƒng kÃ½ Game Object.
// Obj is the function group for creating and registering Game Objects.
var Obj = &objGroup{}

// Object lÃ  interface tá»‘i thiá»ƒu mÃ  Custom Object pháº£i implement Ä‘á»ƒ dÃ¹ng vá»›i NewObject.
// Game object chá»‰ cáº§n cÃ³ OnCreate() — cÃ¡c lifecycle khÃ¡c (OnStep, OnDestroy) lÃ  optional override.
// type Object interface {
// 	OnCreate()
// }

// â”€â”€â”€ Object ID Counter â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// objectIDCounter lÃ  bá»™ Ä‘áº¿m ID tá»± tÄƒng thread-safe cho má» i Object Ä‘Æ°á»£c táº¡o.
var objectIDCounter int64 = 0

// getNextObjectID tráº£ vá» ID tiáº¿p theo (báº¯t Ä‘áº§u tá»« 1), thread-safe dÃ¹ng atomic.
func getNextObjectID() int {
	return int(atomic.AddInt64(&objectIDCounter, 1))
}

func (o *objGroup) NewObject(target Object, name string, componentCode string) {
	tokens := strings.Fields(componentCode)
	tokenSet, sceneName, hasSceneToken := filter(tokens)

	// Lấy map từ scene phù hợp để tạo entry
	var targetScene = getScene(sceneName)

	targetMap := targetScene.GetMap()

	// XÃ¢y danh sÃ¡ch component types tá»« token
	componentsList := []donburi.IComponentType{}
	for token := range tokenSet {
		comp := getComponentType(token)
		if comp != nil {
			componentsList = append(componentsList, comp)
		}
	}

	// Táº¡o ECS entry trong world cá»§a map
	entry := targetMap.World().Entry(targetMap.World().Create(componentsList...))

	// Khá»Ÿi táº¡o data máº·c Ä‘á»‹nh tá»« registry cho tá»«ng component
	for token := range tokenSet {
		enginetype.InitializeComponent(token, entry)
	}

	// LuÃ´n ghi InforData (báº¯t buá»™c): ID tá»± tÄƒng + name
	nextID := getNextObjectID()
	donburi.SetValue(entry, enginetype.Infor, components.InforData{
		Id:   nextID,
		Name: name,
	})

	obj := nobject.NewObject(entry)
	bind(target, obj)
	specilalCase(obj, tokenSet)

	// LÆ°u vÃ o global store náº¿u cÃ³ tÃªn (Ä‘á»ƒ tra cá»©u sau báº±ng napi.GetObject)
	if target != nil && hasSceneToken {
		o.RegisterIn(target, targetScene)
	}
}

func specilalCase(target Object, tokenSet map[string]bool) {}

// â”€â”€â”€ Register Helpers â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// Register Ä‘Äƒng kÃ½ má»™t IObject vÃ o scene's update loop.
//
// scene: name of scene, if nil -> global var
func (o *objGroup) Register(obj IObject, scene string) {
	var _scene interface{ AddObject(IObject) }
	_scene = getScene(scene)

	if _scene != nil {
		_scene.AddObject(obj)
	}
}

// RegisterIn Ä‘Äƒng kÃ½ má»™t IObject vÃ o scene cá»¥ thá»ƒ do caller chá»‰ Ä‘á»‹nh.
func (o *objGroup) RegisterIn(obj IObject, scene IScene) {
	if scene != nil {
		scene.AddObject(obj)
	}
}

// Remove xóa IObject khỏi scene hiện tại theo cơ chế deferred (cuối frame).
// MarkDead() được gọi ngay; OnDestroy() sẽ được gọi ở frame tiếp theo.
// Nếu object đã bị remove trước đó, hàm không làm gì thêm.
func (o *objGroup) Remove(obj IObject) {
	currentScene := Scene.GetCurrentScene()
	if currentScene == nil {
		return
	}
	m := currentScene.GetMap()
	if m == nil {
		return
	}
	m.RemoveObject(obj)
}

// â”€â”€â”€ Internal Helpers â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// getScene tráº£ vá»  scene phÃ¹ há»£p: Global Scene náº¿u global=true, Current Scene náº¿u false.
func getScene(name string) IScene {
	switch name {
	case "glo":
		return Scene.GetGlobalScene()
	case "cur", "":
		return Scene.GetCurrentScene()
	default:
		return Scene.GetSceneByID(name)
	}
}

// filter phÃ¢n tÃ­ch danh sÃ¡ch token, tÃ¡ch modifier "sce-*" (hoáº·c "glo").
// Tráº£ vá» : set token duy nháº¥t vÃ  tÃªn scene Ä‘á»ƒ auto-register.
func filter(tokens []string) (map[string]bool, string, bool) {
	tokenSet := make(map[string]bool, len(tokens))
	var sceneName string = ""
	hasSceneToken := false
	for _, t := range tokens {
		if after, ok := strings.CutPrefix(t, "sce-"); ok {
			sceneName = after
			hasSceneToken = true
			continue
		}
		tokenSet[t] = true
	}

	// Infor lÃ  component báº¯t buá»™c Ä‘á»‘i vá»›i má» i Object
	tokenSet["info"] = true
	return constraint(tokenSet), sceneName, hasSceneToken
}

// bind inject IObject vÃ o target struct vÃ  táº¥t cáº£ Component Mixin nhÃºng trong nÃ³.
// DÃ¹ng reflection Ä‘á»ƒ tá»± Ä‘á»™ng gÃ¡n â€” game dev khÃ´ng cáº§n khá»Ÿi táº¡o thá»§ cÃ´ng tá»«ng mixin.
//
// Quy táº¯c inject:
//  1. Field tÃªn "IObject" â†’ gÃ¡n trá»±c tiáº¿p IObject base.
//  2. Field lÃ  struct cÃ³ sub-field "IObject" â†’ gÃ¡n vÃ o sub-field Ä‘Ã³.
//  3. Field implement interface BindComponent(IObject) â†’ gá»i BindComponent.
func bind(target any, base IObject) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		panic("[napi.Bind] target pháº£i lÃ  con trá» tá»›i struct (vÃ­ dá»¥: *Player)")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		// 1. GÃ¡n IObject chÃ­nh cá»§a Custom Object
		if fieldType.Name == "IObject" || fieldType.Name == "Object" {
			if fieldVal.CanSet() {
				fieldVal.Set(reflect.ValueOf(base))
			}
			continue
		}

		// 2. GÃ¡n IObject cho cÃ¡c Component Mixin nhÃºng (PositionComponent, SpriteComponentâ€¦)
		if fieldVal.Kind() == reflect.Struct {
			subField := fieldVal.FieldByName("IObject")
			if subField.IsValid() && subField.CanSet() {
				subField.Set(reflect.ValueOf(base))
			}
		}

		// 3. Gá»i BindComponent náº¿u mixin implement interface Ä‘Ã³
		if fieldType.PkgPath == "" && fieldVal.CanAddr() {
			if binder, ok := fieldVal.Addr().Interface().(interface{ BindComponent(IObject) }); ok {
				binder.BindComponent(base)
				continue
			}
		}
	}
}

func constraint(tokenSet map[string]bool) map[string]bool {
	if tokenSet["col"] {
		tokenSet["box"] = true
	}
	if tokenSet["twn"] {
		tokenSet["spr"] = true
	}
	if tokenSet["back"] || tokenSet["velo"] || tokenSet["spr"] {
		tokenSet["pos"] = true
	}
	// DrawComponent requires a Position to anchor coordinates
	if tokenSet["drw"] {
		tokenSet["pos"] = true
	}
	return tokenSet
}
