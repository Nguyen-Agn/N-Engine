package domain

import "github.com/yohamta/donburi"

// IScene đại diện cho một phân cảnh độc lập trong trò chơi.
// Mỗi Scene chứa một Physical Map (thế giới game), tùy chọn GUI Map (HUD),
// và một Camera (viewport + render). Scene điều phối Update và Draw toàn bộ.
type IScene interface {
	// Update được gọi mỗi frame để cập nhật logic của toàn Scene.
	// Thứ tự: Map.Update() → Camera.Update(). Trả về error để dừng game.
	// Purpose: Updates the logic of the entire scene for the current frame.
	// Inputs: None.
	// Outputs: error - Returns an error to signal the game loop to stop.
	Update() error

	// Draw render Scene lên màn hình qua Camera.
	// Camera sẽ vẽ Physical Map (có offset) rồi GUI Map (không offset).
	// Purpose: Renders the scene to the screen.
	// Inputs: None.
	// Outputs: None.
	Draw()

	// Destroy được gọi khi Scene bị xóa khỏi SceneManager.
	// Dùng để giải phóng tài nguyên nếu cần.
	// Purpose: Cleans up resources when the scene is removed or destroyed.
	// Inputs: None.
	// Outputs: None.
	Destroy()

	// AddObject đăng ký IObject vào Physical Map của Scene.
	// Purpose: Adds a game object to the scene's physical map.
	// Inputs: obj IObject - The object to register.
	// Outputs: None.
	AddObject(obj IObject)

	// AddGUIObject đăng ký IObject vào GUI Map của Scene (screen-space).
	// GUI Map được tạo tự động nếu chưa tồn tại khi gọi hàm này.
	// Purpose: Adds a game object to the scene's GUI map (screen-space rendering).
	// Inputs: obj IObject - The UI object to register.
	// Outputs: None.
	AddGUIObject(obj IObject)

	// World trả về donburi.World của Physical Map.
	// Dùng bởi napi khi tạo ECS entity. Tương đương GetMap().World().
	// Purpose: Retrieves the ECS World of the physical map.
	// Inputs: None.
	// Outputs: donburi.World - The donburi world instance.
	World() donburi.World

	// GetMap trả về Physical Map của Scene.
	// Purpose: Retrieves the physical map of the scene.
	// Inputs: None.
	// Outputs: IMap - The scene's physical map.
	GetMap() IMap

	// GetGUIMap trả về GUI Map của Scene (có thể nil nếu chưa tạo).
	// Purpose: Retrieves the GUI map of the scene.
	// Inputs: None.
	// Outputs: IMap - The scene's GUI map, or nil if none exists.
	GetGUIMap() IMap

	// GetCamera trả về Camera của Scene.
	// Purpose: Retrieves the camera associated with this scene.
	// Inputs: None.
	// Outputs: ICamera - The camera instance.
	GetCamera() ICamera

	// GetID trả về ID của Scene hiện tại.
	// Purpose: Retrieves the unique identifier of the scene.
	// Inputs: None.
	// Outputs: string - The scene ID.
	GetID() string
}

// ISceneManager đại diện cho trình quản lý scene — điều phối stack các Scene.
type ISceneManager interface {
	// Update được gọi mỗi frame — chạy Global Scene rồi Current Scene.
	// Purpose: Updates the active scenes (global and current) for the frame.
	// Inputs: None.
	// Outputs: error - Returns an error to signal game termination.
	Update() error

	// Draw render Current Scene lên màn hình.
	// Purpose: Renders the currently active scene.
	// Inputs: None.
	// Outputs: error - Returns an error if rendering fails.
	Draw() error

	// Layout thiết lập kích thước màn hình logic cho Ebitengine.
	// Purpose: Computes the logical screen size based on the window's outside dimensions.
	// Inputs:
	//   - outsideWidth int: The width of the game window.
	//   - outsideHeight int: The height of the game window.
	// Outputs:
	//   1. screenWidth int - The logical screen width.
	//   2. screenHeight int - The logical screen height.
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)

	// ChangeSceneFromList chuyển sang Scene đã có trong danh sách theo id.
	// Scene hiện tại bị pause (không Destroy) — có thể quay lại sau.
	// Purpose: Switches to an existing scene from the registered list by its ID, pausing the current one.
	// Inputs: id string - The ID of the target scene.
	// Outputs: error - Returns an error if the scene ID is not found.
	ChangeSceneFromList(id string) error

	// ChangeSceneForce ép thay thế Scene hiện tại bằng Scene mới.
	// Scene cũ bị gọi Destroy() — dùng khi không cần quay lại Scene cũ.
	// Purpose: Forcefully replaces the current scene with a new one, destroying the old scene.
	// Inputs: next IScene - The new scene instance to transition to.
	// Outputs: error - Returns an error if the transition fails.
	ChangeSceneForce(next IScene) error

	// GetCurrentScene trả về Scene đang hoạt động. Trả về nil nếu chưa có.
	// Purpose: Retrieves the currently active scene.
	// Inputs: None.
	// Outputs: IScene - The active scene, or nil if none is set.
	GetCurrentScene() IScene

	// GetSceneFromList trả về IScene theo id. Trả về nil nếu không tìm thấy.
	// Purpose: Retrieves a registered scene by its ID without switching to it.
	// Inputs: id string - The scene ID to look up.
	// Outputs: IScene - The requested scene, or nil if not found.
	GetSceneFromList(id string) IScene

	// GetGlobalScene trả về Global Hidden Scene — luôn Update mọi frame,
	// không Draw. Dùng để chứa Object persistent (data, audio không thuộc scene).
	// Purpose: Retrieves the global hidden scene that persists across scene transitions.
	// Inputs: None.
	// Outputs: IScene - The global scene instance.
	GetGlobalScene() IScene

	// AddScene thêm Scene vào danh sách quản lý theo id.
	// Trả về error nếu id đã tồn tại.
	// Purpose: Registers a new scene to the manager under a specific ID.
	// Inputs:
	//   - id string: The unique identifier for the scene.
	//   - scene IScene: The scene instance to add.
	// Outputs: error - Returns an error if a scene with this ID already exists.
	AddScene(id string, scene IScene) error

	// RemoveScene xóa Scene khỏi danh sách và gọi Destroy.
	// Purpose: Removes a registered scene from the manager and destroys it.
	// Inputs: id string - The ID of the scene to remove.
	// Outputs: error - Returns an error if the scene is not found.
	RemoveScene(id string) error

	// RemoveAllScene xóa toàn bộ Scene và đặt currentScene = nil.
	// Purpose: Destroys all registered scenes and clears the current scene state.
	// Inputs: None.
	// Outputs: error - Returns an error if deletion fails.
	RemoveAllScene() error
}
