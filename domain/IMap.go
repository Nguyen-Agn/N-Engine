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
	World() donburi.World

	// AddObject đăng ký IObject vào Map để được xử lý mỗi frame.
	// LogicSystem sẽ gọi Create() vào frame tiếp theo, StepUpdate() mỗi frame sau đó.
	AddObject(obj IObject)

	// RemoveObject yêu cầu xóa IObject khỏi Map theo cơ chế deferred (cuối frame).
	// MarkDead() được gọi ngay lập tức; object sẽ bị cắt khỏi objectList và
	// DrawRegistry sau khi tất cả System trong frame đó chạy xong.
	// OnDestroy() sẽ được gọi ở đầu frame tiếp theo.
	RemoveObject(obj IObject)

	// Update chạy toàn bộ logic mỗi frame theo thứ tự: Input → Logic → Audio.
	// Trả về error nếu cần dừng game.
	Update() error

	// Width trả về chiều rộng bản đồ (pixel). 0 = không giới hạn.
	Width() int

	// Height trả về chiều cao bản đồ (pixel). 0 = không giới hạn.
	Height() int

	// GetObjects trả về danh sách các IObject hiện có trong Map.
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
	X() float32
	// Y trả về tọa độ dọc của camera trong map space (pixel).
	Y() float32
	// SetX dịch chuyển camera theo trục ngang.
	SetX(x float32)
	// SetY dịch chuyển camera theo trục dọc.
	SetY(y float32)

	// Width trả về chiều rộng viewport (pixel) — thường bằng chiều rộng cửa sổ.
	Width() int
	// Height trả về chiều cao viewport (pixel) — thường bằng chiều cao cửa sổ.
	Height() int

	// SetTarget đặt IObject làm mục tiêu để camera tự động theo dõi mỗi frame.
	// Truyền nil để tắt follow.
	SetTarget(obj IObject)

	// Update cập nhật vị trí camera mỗi frame (follow target + clamp trong map bounds).
	// mapW, mapH là kích thước bản đồ để giới hạn camera không ra ngoài biên.
	// Truyền 0, 0 nếu bản đồ không có giới hạn.
	Update(mapW, mapH int)

	// SetScreen thiết lập canvas đích (ebiten.Image) để render lên.
	// Phải gọi mỗi frame trước Draw() — thực hiện bởi EbitenGame.
	SetScreen(screen *ebiten.Image)

	// Draw render Physical Map (có camera offset) rồi GUI Map (không có offset, đè lên trên).
	// guiWorld có thể nil nếu Scene không có GUI Map.
	Draw(physicalWorld donburi.World, guiWorld donburi.World)
}
