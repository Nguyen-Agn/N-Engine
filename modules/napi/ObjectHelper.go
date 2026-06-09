package napi

import (
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/Nguyen-Agn/N-Engine/modules/components"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"
	"github.com/Nguyen-Agn/N-Engine/modules/nobject"

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

// ——————————————————————————————————————————————————————————————————————————————
// Object ID Counter
// ——————————————————————————————————————————————————————————————————————————————

// objectIDCounter is a thread-safe auto-incrementing ID counter for every created Object.
var objectIDCounter int64 = 0

// getNextObjectID safely increments and returns the next available unique object ID.
//
// Purpose: Provides a thread-safe auto-incrementing ID for every newly instantiated game object.
//
// Outputs:
// - int: The next unique integer ID.
func getNextObjectID() int {
	return int(atomic.AddInt64(&objectIDCounter, 1))
}

// NewObject creates a new game object, binds its components, and registers it to a scene.
//
// Purpose: Parses a string of component tokens to dynamically construct an ECS entity, initializes its default data, and binds it to the provided custom object struct.
//
// Inputs:
// - target (Object): The custom object struct pointer to bind to.
// - name (string): The human-readable name of the object.
// - componentCode (string): Space-separated tokens (e.g., "pos spr") indicating which components to attach.
func (o *objGroup) NewObject(target Object, name string, componentCode string) {
	tokens := strings.Fields(componentCode)
	tokenSet, sceneName, hasSceneToken := filter(tokens)

	// Lấy map từ scene phù hợp để tạo entry
	var targetScene = getScene(sceneName)

	targetMap := targetScene.GetMap()

	// Xây danh sách component types từ token
	componentsList := []donburi.IComponentType{}
	for token := range tokenSet {
		comp := getComponentType(token)
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
	bind(target, obj)
	specilalCase(obj, tokenSet)

	// Lưu vào global store nếu có tên (để tra cứu sau bằng napi.GetObject)
	if target != nil && hasSceneToken {
		o.RegisterIn(target, targetScene)
	}
}

// specilalCase handles any special initialization logic for specific tokens.
//
// Purpose: A hook for custom post-creation logic based on the attached component tokens.
//
// Inputs:
// - target (Object): The newly created object.
// - tokenSet (map[string]bool): The set of component tokens attached to the object.
func specilalCase(target Object, tokenSet map[string]bool) {}

// ——————————————————————————————————————————————————————————————————————————————
// Register Helpers
// ——————————————————————————————————————————————————————————————————————————————

// Register adds an object to a scene's update loop by name.
//
// Purpose: Injects an initialized object into a scene so its lifecycle methods (like Update) are called.
//
// Inputs:
// - obj (IObject): The object to register.
// - scene (string): The identifier of the target scene. If "glo", registers to the global scene. If empty, registers to the current scene.
func (o *objGroup) Register(obj IObject, scene string) {
	var _scene interface{ AddObject(IObject) }
	_scene = getScene(scene)

	if _scene != nil {
		_scene.AddObject(obj)
	}
}

// RegisterIn adds an object directly to a specific scene instance.
//
// Purpose: Bypasses string lookup to insert an object directly into a given scene.
//
// Inputs:
// - obj (IObject): The object to register.
// - scene (IScene): The target scene instance.
func (o *objGroup) RegisterIn(obj IObject, scene IScene) {
	if scene != nil {
		scene.AddObject(obj)
	}
}

// Remove schedules an object for deletion at the end of the current frame.
//
// Purpose: Safely removes an object from the active scene's map. Immediately marks the object as dead so it is skipped during the rest of the frame's processing.
//
// Inputs:
// - obj (IObject): The object to remove.
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

// ——————————————————————————————————————————————————————————————————————————————
// Internal Helpers
// ——————————————————————————————————————————————————————————————————————————————

// getScene resolves a scene string identifier to an actual scene instance.
//
// Purpose: Translates "glo", "cur", empty strings, or specific IDs into their respective IScene instances.
//
// Inputs:
// - name (string): The identifier string.
//
// Outputs:
// - IScene: The resolved scene.
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

// filter extracts component tokens and any scene modifier from a raw list of tokens.
//
// Purpose: Parses the component code string, identifying "sce-*" modifiers to determine the target scene and ensuring required components like "info" are implicitly added.
//
// Inputs:
// - tokens ([]string): The raw list of tokens.
//
// Outputs:
// - map[string]bool: The final set of unique component tokens.
// - string: The extracted scene name, if any.
// - bool: True if a scene modifier was found.
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

	// Infor là component bắt buộc đối với mọi Object
	tokenSet["info"] = true
	return constraint(tokenSet), sceneName, hasSceneToken
}

// bind automatically links the base object and its mixins into the target struct using reflection.
//
// Purpose: Iterates through the target struct's fields to auto-populate the embedded IObject base and trigger BindComponent for any mixins that implement it.
//
// Inputs:
// - target (any): A pointer to the custom struct.
// - base (IObject): The base object to inject.
//
// Special Requirements:
// - Panics if the target is not a pointer to a struct.
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

// constraint applies implicit dependencies between component tokens.
//
// Purpose: Ensures that if an object requires certain components (like collision needing a box, or drawing needing position), the dependencies are automatically fulfilled.
//
// Inputs:
// - tokenSet (map[string]bool): The initial set of tokens.
//
// Outputs:
// - map[string]bool: The modified set including any implicitly required tokens.
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
