package domain

import "github.com/yohamta/donburi"

// IScene đại diện cho một phân cảnh độc lập trong trò chơi.
// Mỗi Scene chứa một Physical Map (thế giới game), tùy chọn GUI Map (HUD),
// và một Camera (viewport + render). Scene điều phối Update và Draw toàn bộ.
type IScene interface {
	// Update được gọi mỗi frame để cập nhật logic của toàn Scene.
	// Thứ tự: Map.Update() → Camera.Update(). Trả về error để dừng game.
	Update() error

	// Draw render Scene lên màn hình qua Camera.
	// Camera sẽ vẽ Physical Map (có offset) rồi GUI Map (không offset).
	Draw()

	// Destroy được gọi khi Scene bị xóa khỏi SceneManager.
	// Dùng để giải phóng tài nguyên nếu cần.
	Destroy()

	// AddObject đăng ký IObject vào Physical Map của Scene.
	AddObject(obj IObject)

	// AddGUIObject đăng ký IObject vào GUI Map của Scene (screen-space).
	// GUI Map được tạo tự động nếu chưa tồn tại khi gọi hàm này.
	AddGUIObject(obj IObject)

	// World trả về donburi.World của Physical Map.
	// Dùng bởi napi khi tạo ECS entity. Tương đương GetMap().World().
	World() donburi.World

	// GetMap trả về Physical Map của Scene.
	GetMap() IMap

	// GetGUIMap trả về GUI Map của Scene (có thể nil nếu chưa tạo).
	GetGUIMap() IMap

	// GetCamera trả về Camera của Scene.
	GetCamera() ICamera
}

// ISceneManager đại diện cho trình quản lý scene — điều phối stack các Scene.
type ISceneManager interface {
	// Update được gọi mỗi frame — chạy Global Scene rồi Current Scene.
	Update() error

	// Draw render Current Scene lên màn hình.
	Draw() error

	// Layout thiết lập kích thước màn hình logic cho Ebitengine.
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)

	// ChangeSceneFromList chuyển sang Scene đã có trong danh sách theo id.
	// Scene hiện tại bị pause (không Destroy) — có thể quay lại sau.
	ChangeSceneFromList(id string) error

	// ChangeSceneForce ép thay thế Scene hiện tại bằng Scene mới.
	// Scene cũ bị gọi Destroy() — dùng khi không cần quay lại Scene cũ.
	ChangeSceneForce(next IScene) error

	// GetCurrentScene trả về Scene đang hoạt động. Trả về nil nếu chưa có.
	GetCurrentScene() IScene

	// GetSceneFromList trả về IScene theo id. Trả về nil nếu không tìm thấy.
	GetSceneFromList(id string) IScene

	// GetGlobalScene trả về Global Hidden Scene — luôn Update mọi frame,
	// không Draw. Dùng để chứa Object persistent (data, audio không thuộc scene).
	GetGlobalScene() IScene

	// AddScene thêm Scene vào danh sách quản lý theo id.
	// Trả về error nếu id đã tồn tại.
	AddScene(id string, scene IScene) error

	// RemoveScene xóa Scene khỏi danh sách và gọi Destroy.
	RemoveScene(id string) error

	// RemoveAllScene xóa toàn bộ Scene và đặt currentScene = nil.
	RemoveAllScene() error
}
