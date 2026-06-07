package core

import (
	"fmt"
)

// ─── SceneManager ────────────────────────────────────────────────────────────

// sceneEntry là record nội bộ lưu một IScene kèm trạng thái pause trong danh sách chờ.
type sceneEntry struct {
	scene    IScene
	isPaused bool // true = đang bị pause, không Update/Draw
}

// SceneManager quản lý toàn bộ vòng đời của các Scene.
type SceneManager struct {
	// currentScene là Scene đang hoạt động
	currentScene IScene

	// sceneList là kho lưu các Scene theo ID để tái sử dụng
	sceneList map[string]*sceneEntry

	// globalScene là màn hình ẩn, luôn chạy Update mọi frame bất kể Scene nào đang active.
	// Dùng để chứa các Object tồn tại xuyên suốt (không bị xóa khi đổi Scene).
	// Không được Draw.
	globalScene *Scene

	// screenW, screenH là kích thước màn hình logic
	screenW, screenH int

	// input được truyền vào mỗi Scene mới để InputSystem có thể hoạt động
	input IInputManager
}

// NewSceneManager khởi tạo SceneManager với kích thước màn hình và IInputManager.
func NewSceneManager(screenW, screenH int, input IInputManager) *SceneManager {
	return &SceneManager{
		sceneList:   make(map[string]*sceneEntry),
		screenW:     screenW,
		screenH:     screenH,
		input:       input,
		// Global Scene: kích thước map không giới hạn (0,0), viewport = screen
		globalScene: NewScene(input, 0, 0, screenW, screenH),
	}
}

// ─── ISceneManager interface — Lifecycle ─────────────────────────────────────

// Update được Ebitengine gọi mỗi frame để xử lý logic.
// Luôn cập nhật globalScene trước, sau đó mới đến currentScene.
func (r *SceneManager) Update() error {
	// Global scene luôn chạy bất kể scene nào đang active
	if err := r.globalScene.Update(); err != nil {
		return err
	}
	if r.currentScene == nil {
		return nil
	}
	return r.currentScene.Update()
}

// Draw được Ebitengine gọi mỗi frame để vẽ.
// Chỉ vẽ currentScene. Trả về error để phù hợp với IRootManager.
func (r *SceneManager) Draw() error {
	if r.currentScene == nil {
		return nil
	}
	r.currentScene.Draw()
	return nil
}

// Layout trả về kích thước màn hình logic cho Ebitengine.
func (r *SceneManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return r.screenW, r.screenH
}

// ─── ISceneManager interface — Scene Management ─────────────────────────────

// AddScene thêm một Scene vào danh sách chờ (chưa kích hoạt).
// Trả về error nếu ID đã tồn tại.
func (r *SceneManager) AddScene(id string, scene IScene) error {
	if _, exists := r.sceneList[id]; exists {
		return fmt.Errorf("scene '%s' already exists in list", id)
	}
	if s, ok := scene.(interface{ setID(string) }); ok {
		s.setID(id)
	}
	r.sceneList[id] = &sceneEntry{
		scene:    scene,
		isPaused: true, // Mặc định pause khi mới thêm vào
	}
	return nil
}

// ChangeSceneFromList chuyển sang Scene đã có trong danh sách theo ID.
// Scene hiện tại sẽ bị đánh dấu pause và ẩn (không Destroy).
// Điều này cho phép tái sử dụng Scene (ví dụ: quay lại màn hình trước).
// Trả về error nếu ID không tồn tại.
func (r *SceneManager) ChangeSceneFromList(id string) error {
	entry, exists := r.sceneList[id]
	if !exists {
		return fmt.Errorf("scene '%s' not found in list", id)
	}

	// Pause Scene hiện tại (không xóa)
	if r.currentScene != nil {
		// Tìm và đánh dấu pause Scene hiện tại trong list
		for _, e := range r.sceneList {
			if e.scene == r.currentScene {
				e.isPaused = true
				break
			}
		}
	}

	// Kích hoạt Scene mới từ list
	entry.isPaused = false
	r.currentScene = entry.scene
	return nil
}

// ChangeSceneForce ép xóa Scene hiện tại (gọi Destroy) và thay bằng Scene mới.
// Scene mới KHÔNG được thêm vào list — đây là chế độ "dùng rồi xóa".
// Dùng khi bạn chắc chắn không cần quay lại Scene cũ.
func (r *SceneManager) ChangeSceneForce(next IScene) error {
	if next == nil {
		return fmt.Errorf("cannot force change to nil scene")
	}

	// Destroy Scene hiện tại nếu có
	if r.currentScene != nil {
		// Xóa Scene hiện tại khỏi list nếu nó đang được lưu trong list
		for id, entry := range r.sceneList {
			if entry.scene == r.currentScene {
				delete(r.sceneList, id)
				break
			}
		}
		r.currentScene.Destroy()
	}

	r.currentScene = next
	return nil
}

// RemoveScene xóa một Scene khỏi danh sách theo ID và gọi Destroy.
// Không làm gì nếu Scene đang là currentScene (cần ChangeSceneForce trước).
// Trả về error nếu ID không tồn tại.
func (r *SceneManager) RemoveScene(id string) error {
	entry, exists := r.sceneList[id]
	if !exists {
		return fmt.Errorf("scene '%s' not found in list", id)
	}

	// Không cho xóa Scene đang chạy qua method này
	if entry.scene == r.currentScene {
		return fmt.Errorf("cannot remove active scene '%s', use ChangeSceneForce first", id)
	}

	entry.scene.Destroy()
	delete(r.sceneList, id)
	return nil
}

// RemoveAllScene xóa toàn bộ danh sách Scene (gọi Destroy từng cái) và đặt currentScene = nil.
func (r *SceneManager) RemoveAllScene() error {
	for id, entry := range r.sceneList {
		entry.scene.Destroy()
		delete(r.sceneList, id)
	}
	r.currentScene = nil
	return nil
}

// GetCurrentScene trả về Scene đang hoạt động.
func (r *SceneManager) GetCurrentScene() IScene {
	return r.currentScene
}

// GetSceneFromList trả về IScene theo id từ danh sách.
// Trả về nil nếu không tìm thấy.
func (r *SceneManager) GetSceneFromList(id string) IScene {
	entry, exists := r.sceneList[id]
	if !exists {
		return nil
	}
	return entry.scene
}

// GetGlobalScene trả về Global Hidden Scene — scene ẩn luôn chạy Update mọi frame.
// Dùng để đặt các Object tồn tại xuyên suốt game (không bị xóa khi đổi Scene).
func (r *SceneManager) GetGlobalScene() IScene {
	return r.globalScene
}
