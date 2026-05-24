package core

import (
	"autoworld/domain"

	"github.com/yohamta/donburi"
)

// Scene quản lý một màn chơi độc lập trong Engine.
// Mỗi Scene chứa:
//   - Physical Map: ECS World + Logic/Audio/Input system + objectList (tọa độ map space)
//   - GUI Map: screen-space overlay cho HUD (optional, tự tạo khi cần)
//   - Camera: viewport, follow-target, DrawSystem
type Scene struct {
	map_   *Map // physical game world
	guiMap *Map // screen-space GUI/HUD (optional, nil nếu không dùng)
	camera *Camera
	input  domain.IInputManager // lưu để tạo guiMap lazy
}

// NewScene khởi tạo Scene mới với Physical Map và Camera.
// mapW, mapH là kích thước bản đồ (pixel). Truyền 0, 0 nếu không giới hạn.
// viewW, viewH là kích thước viewport của Camera (thường bằng screen size).
func NewScene(input domain.IInputManager, mapW, mapH, viewW, viewH int) *Scene {
	return &Scene{
		map_:   NewMap(input, mapW, mapH),
		camera: NewCamera(viewW, viewH),
		input:  input,
	}
}

// ─── IScene interface ──────────────────────────────────────────────────────────

// Update cập nhật toàn bộ Scene mỗi frame.
// Thứ tự: Physical Map → GUI Map (nếu có) → Camera follow.
func (s *Scene) Update() error {
	if err := s.map_.Update(); err != nil {
		return err
	}
	if s.guiMap != nil {
		if err := s.guiMap.Update(); err != nil {
			return err
		}
	}
	s.camera.Update(s.map_.Width(), s.map_.Height())
	return nil
}

// Draw render Scene lên màn hình qua Camera.
// Camera vẽ Physical Map (có camera offset) rồi GUI Map (không offset, đè lên trên).
func (s *Scene) Draw() {
	if s.guiMap != nil {
		s.camera.Draw(s.map_.World(), s.guiMap.World())
	} else {
		s.camera.Draw(s.map_.World(), nil)
	}
}

// Destroy được gọi khi Scene bị xóa khỏi SceneManager.
func (s *Scene) Destroy() {}

// AddObject đăng ký IObject vào Physical Map.
func (s *Scene) AddObject(obj IObject) {
	s.map_.AddObject(obj)
}

// AddGUIObject đăng ký IObject vào GUI Map (screen-space, không camera offset).
// GUI Map được tạo tự động nếu chưa tồn tại.
func (s *Scene) AddGUIObject(obj IObject) {
	if s.guiMap == nil {
		s.guiMap = NewGUIMap(s.input, s.camera.Width(), s.camera.Height())
	}
	s.guiMap.AddObject(obj)
}

// World trả về donburi.World của Physical Map. Tương đương GetMap().World().
func (s *Scene) World() donburi.World {
	return s.map_.World()
}

// ─── Getters ──────────────────────────────────────────────────────────────────

// GetMap trả về Physical Map của Scene.
func (s *Scene) GetMap() IMap { return s.map_ }

// GetGUIMap trả về GUI Map của Scene. Trả về nil nếu chưa khởi tạo.
func (s *Scene) GetGUIMap() IMap {
	if s.guiMap == nil {
		return nil
	}
	return s.guiMap
}

// GetCamera trả về Camera của Scene.
func (s *Scene) GetCamera() ICamera { return s.camera }

// ─── Setters ──────────────────────────────────────────────────────────────────

// SetGuiMap
func (s *Scene) SetGUIMap(input domain.IInputManager, mapW, mapH int) {
	s.guiMap = NewMap(input, mapW, mapH)
}
