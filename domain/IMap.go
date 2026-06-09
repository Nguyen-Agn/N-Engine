package domain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ─── IMap ─────────────────────────────────────────────────────────────────────

// IMap quản lý ECS World và toàn bộ logic update của một bản đồ game.
// Mỗi Scene có một Physical Map (thế giới game) và tùy chọn một GUI Map (HUD, screen-space).
// IMap không biết gì về render — đó là trách nhiệm của ICamera.
type IMap interface {
	// World trả về donburi.World của Map — dùng để tạo/query ECS entity.
	// Purpose: Retrieves the underlying ECS World instance used by this map.
	// Inputs: None.
	// Outputs: donburi.World - The donburi ECS world managing entities.
	World() donburi.World

	// AddObject đăng ký IObject vào Map để được xử lý mỗi frame.
	// LogicSystem sẽ gọi Create() vào frame tiếp theo, StepUpdate() mỗi frame sau đó.
	// Purpose: Registers a new game object into the map for per-frame processing.
	// Inputs: obj IObject - The game object to add.
	// Outputs: None.
	// Special requirements: The object's Create() method will be called on the next frame.
	AddObject(obj IObject)

	// RemoveObject yêu cầu xóa IObject khỏi Map theo cơ chế deferred (cuối frame).
	// MarkDead() được gọi ngay lập tức; object sẽ bị cắt khỏi objectList và
	// DrawRegistry sau khi tất cả System trong frame đó chạy xong.
	// OnDestroy() sẽ được gọi ở đầu frame tiếp theo.
	// Purpose: Requests the removal of a game object from the map (deferred until the end of the frame).
	// Inputs: obj IObject - The game object to remove.
	// Outputs: None.
	// Special requirements: Marks the object as dead immediately. OnDestroy() is called at the start of the next frame.
	RemoveObject(obj IObject)

	// Update chạy toàn bộ logic mỗi frame theo thứ tự: Input → Logic → Audio.
	// Trả về error nếu cần dừng game.
	// Purpose: Executes the main logic loop for the map for a single frame.
	// Inputs: None.
	// Outputs: error - Returns an error if the update loop encounters a critical failure or game stop condition.
	Update() error

	// Width trả về chiều rộng bản đồ (pixel). 0 = không giới hạn.
	// Purpose: Retrieves the total physical width of the map.
	// Inputs: None.
	// Outputs: int - The width in pixels. 0 means unbounded.
	Width() int

	// Height trả về chiều cao bản đồ (pixel). 0 = không giới hạn.
	// Purpose: Retrieves the total physical height of the map.
	// Inputs: None.
	// Outputs: int - The height in pixels. 0 means unbounded.
	Height() int

	// GetObjects trả về danh sách các IObject hiện có trong Map.
	// Purpose: Retrieves a list of all currently active game objects registered in the map.
	// Inputs: None.
	// Outputs: []IObject - A slice containing all active objects.
	GetObjects() []IObject
}

// ─── ICamera ──────────────────────────────────────────────────────────────────

// ICamera định nghĩa một viewport nhìn vào một IMap.
// Camera chịu trách nhiệm: xác định vùng nhìn, follow target, và render qua DrawSystem.
//
// Tọa độ camera (X, Y) là vị trí góc trên-trái của viewport trong **map space**.
// DrawSystem sẽ trừ (camX, camY) khỏi tọa độ entity để chuyển sang screen space.
type ICamera interface {
	// X trả về tọa độ ngang của camera trong map space (pixel).
	// Purpose: Retrieves the X coordinate of the camera's top-left corner in map space.
	// Inputs: None.
	// Outputs: float32 - The X position in pixels.
	X() float32
	
	// Y trả về tọa độ dọc của camera trong map space (pixel).
	// Purpose: Retrieves the Y coordinate of the camera's top-left corner in map space.
	// Inputs: None.
	// Outputs: float32 - The Y position in pixels.
	Y() float32
	
	// SetX dịch chuyển camera theo trục ngang.
	// Purpose: Sets the X coordinate of the camera in map space.
	// Inputs: x float32 - The new X position.
	// Outputs: None.
	SetX(x float32)
	
	// SetY dịch chuyển camera theo trục dọc.
	// Purpose: Sets the Y coordinate of the camera in map space.
	// Inputs: y float32 - The new Y position.
	// Outputs: None.
	SetY(y float32)

	// Width trả về chiều rộng viewport (pixel) — thường bằng chiều rộng cửa sổ.
	// Purpose: Retrieves the width of the camera's viewport.
	// Inputs: None.
	// Outputs: int - The viewport width in pixels.
	Width() int
	
	// Height trả về chiều cao viewport (pixel) — thường bằng chiều cao cửa sổ.
	// Purpose: Retrieves the height of the camera's viewport.
	// Inputs: None.
	// Outputs: int - The viewport height in pixels.
	Height() int

	// SetTarget đặt IObject làm mục tiêu để camera tự động theo dõi mỗi frame.
	// Truyền nil để tắt follow.
	// Purpose: Sets a target game object for the camera to follow automatically.
	// Inputs: obj IObject - The target object, or nil to disable following.
	// Outputs: None.
	SetTarget(obj IObject)

	// Update cập nhật vị trí camera mỗi frame (follow target + clamp trong map bounds).
	// mapW, mapH là kích thước bản đồ để giới hạn camera không ra ngoài biên.
	// Truyền 0, 0 nếu bản đồ không có giới hạn.
	// Purpose: Updates the camera position, handling target following and boundary clamping.
	// Inputs:
	//   - mapW int: Total map width for boundary clamping.
	//   - mapH int: Total map height for boundary clamping.
	// Outputs: None.
	// Special requirements: Pass (0, 0) if the map has no boundaries.
	Update(mapW, mapH int)

	// SetScreen thiết lập canvas đích (ebiten.Image) để render lên.
	// Phải gọi mỗi frame trước Draw() — thực hiện bởi EbitenGame.
	// Purpose: Sets the destination canvas where the camera will render the frame.
	// Inputs: screen *ebiten.Image - The destination surface.
	// Outputs: None.
	// Special requirements: Must be called each frame prior to calling Draw().
	SetScreen(screen *ebiten.Image)

	// Draw render Physical Map (có camera offset) rồi GUI Map (không có offset, đè lên trên).
	// guiWorld có thể nil nếu Scene không có GUI Map.
	// Purpose: Renders the physical world through the camera view, followed by the GUI world overlay.
	// Inputs:
	//   - physicalWorld donburi.World: The ECS world for in-game objects.
	//   - guiWorld donburi.World: The ECS world for UI elements (can be nil).
	// Outputs: None.
	Draw(physicalWorld donburi.World, guiWorld donburi.World)
}
