package napi

import (
	"strconv"
	"strings"

	"autoworld/modules/core"
)

// ─── Scene Management ──────────────────────────────────────────────────────────

// AddScene thêm một Scene vào danh sách quản lý của engine.
func AddScene(id string, scene IScene) error {
	return engine().Scene.AddScene(id, scene)
}

// GoToScene chuyển sang Scene đã có trong danh sách theo id.
// Scene hiện tại bị pause (không Destroy), có thể quay lại sau.
func GoToScene(id string) error {
	return engine().Scene.ChangeSceneFromList(id)
}

// ReplaceScene ép thay thế Scene hiện tại bằng Scene mới.
// Scene hiện tại bị gọi Destroy() — dùng khi không cần quay lại Scene cũ.
func ReplaceScene(next IScene) error {
	return engine().Scene.ChangeSceneForce(next)
}

// RemoveScene xóa Scene khỏi danh sách theo id và gọi Destroy trên nó.
func RemoveScene(id string) error {
	return engine().Scene.RemoveScene(id)
}

// RemoveAllScenes xóa toàn bộ Scene và đặt currentScene = nil.
func RemoveAllScenes() error {
	return engine().Scene.RemoveAllScene()
}

// GetCurrentScene trả về Scene đang hoạt động.
// Trả về nil nếu không có scene nào đang chạy.
func GetCurrentScene() IScene {
	s := engine().Scene.GetCurrentScene()
	if s == nil {
		return nil
	}
	return s
}

// GetSceneByID lấy Scene từ danh sách theo id.
// Trả về nil nếu không tìm thấy.
func GetSceneByID(id string) IScene {
	s := engine().Scene.GetSceneFromList(id)
	if s == nil {
		return nil
	}
	return s
}

// NewScene tạo Scene mới và đăng ký vào engine với id cho trước.
// component: "gui-640-480 map-1280-1280" hoặc "map-0-0"
// Trả về IScene và error nếu id đã tồn tại.
func NewScene(id, component string) (IScene, error) {
	e := engine()
	viewW := e.Config.GetInt("game-width")
	viewH := e.Config.GetInt("game-height")

	var mapW, mapH int
	var guiW, guiH int
	var hasGui bool

	// Parse component string
	tokens := strings.FieldsSeq(component)
	for token := range tokens {
		parts := strings.Split(token, "-")
		if len(parts) >= 3 {
			switch parts[0] {
			case "map":
				mapW, _ = strconv.Atoi(parts[1])
				mapH, _ = strconv.Atoi(parts[2])
			case "gui":
				hasGui = true
				guiW, _ = strconv.Atoi(parts[1])
				guiH, _ = strconv.Atoi(parts[2])
			}
		} else if len(parts) >= 1 {
			if parts[0] == "gui" {
				hasGui = true
				guiW = viewW
				guiH = viewH
			}
		}
	}

	scene := core.NewScene(e.Input, mapW, mapH, viewW, viewH)

	if hasGui {
		scene.SetGUIMap(e.Input, guiW, guiH)
	}

	if err := e.Scene.AddScene(id, scene); err != nil {
		return nil, err
	}
	return scene, nil
}

// NewSceneAndGo tạo Scene mới, đăng ký và chuyển ngay sang nó.
// Shortcut cho khởi động: tạo scene đầu tiên và activate liền.
func NewSceneAndGo(id, component string) (IScene, error) {
	scene, err := NewScene(id, component)
	if err != nil {
		return nil, err
	}
	if err := engine().Scene.ChangeSceneFromList(id); err != nil {
		return nil, err
	}
	return scene, nil
}

// ─── Global Hidden Scene ──────────────────────────────────────────────────────

// GetGlobalScene trả về Global Hidden Scene.
// Scene này luôn chạy Update mọi frame bất kể scene nào đang active.
// Không Draw. Dùng để chứa Object persistent (data, audio xuyên scene).
func GetGlobalScene() IScene {
	s := engine().Scene.GetGlobalScene()
	if s == nil {
		return nil
	}
	return s
}

// ─── Map & Camera Helpers ──────────────────────────────────────────────────────

// GetMap trả về Physical Map của scene chỉ định.
// Trả về nil nếu scene nil hoặc không có map.
func GetMap(scene IScene) IMap {
	if scene == nil {
		return nil
	}
	return scene.GetMap()
}

// GetGUIMap trả về GUI Map của scene chỉ định.
// Trả về nil nếu chưa có GUI object nào được thêm vào scene.
func GetGUIMap(scene IScene) IMap {
	if scene == nil {
		return nil
	}
	return scene.GetGUIMap()
}

// GetCamera trả về Camera của scene chỉ định.
func GetCamera(scene IScene) ICamera {
	if scene == nil {
		return nil
	}
	return scene.GetCamera()
}

// SetCameraTarget đặt IObject làm mục tiêu để camera tự động theo dõi mỗi frame.
// Truyền nil để tắt follow.
func SetCameraTarget(scene IScene, target IObject) {
	if scene == nil {
		return
	}
	cam := scene.GetCamera()
	if cam != nil {
		cam.SetTarget(target)
	}
}

// MoveCamera dịch chuyển camera đến vị trí (x, y) trong map space.
func MoveCamera(scene IScene, x, y float32) {
	if scene == nil {
		return
	}
	cam := scene.GetCamera()
	if cam != nil {
		cam.SetX(x)
		cam.SetY(y)
	}
}
