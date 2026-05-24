package domain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ─── ECS System Interfaces ────────────────────────────────────────────────────

// ILogicSystem điều phối logic game mỗi frame trên tập hợp Object.
type ILogicSystem interface {
	// Update gọi Create() cho object mới, StepUpdate() cho tất cả, Destroy() cho object bị xóa.
	Update(objectList []IObject)

	// AddObjectCreated đăng ký Object vừa tạo — sẽ gọi Create() vào frame tiếp theo.
	AddObjectCreated(o IObject)

	// AddObjectDestroy đăng ký Object sẽ xóa — sẽ gọi Destroy() cuối frame.
	AddObjectDestroy(o IObject)
}

// IDrawSystem xử lý render đồ họa mỗi frame.
// Truy vấn ECS World để tìm entity cần vẽ và áp dụng camera offset.
type IDrawSystem interface {
	// Draw render tất cả entity có Sprite + Position lên screen.
	// camX, camY là tọa độ camera trong map space — được trừ khỏi tọa độ entity
	// để chuyển sang screen space. Truyền 0, 0 cho GUI layer (không có camera offset).
	Draw(w donburi.World, camX, camY float32)

	// SetScreen thiết lập canvas đích (ebiten.Image) để DrawSystem vẽ lên.
	// Phải được gọi trước Draw() mỗi frame.
	SetScreen(s *ebiten.Image)
}

// IAudioSystem xử lý phát âm thanh mỗi frame.
type IAudioSystem interface {
	// Update đọc AudioData của từng entity và thực thi lệnh phát/dừng âm thanh.
	Update(w donburi.World)
}

// IUpdateSystem là interface chung cho các system thực thi logic cập nhật mỗi frame.
type IUpdateSystem interface {
	Update(objectList []IObject)
}

type IAlarmSystem = IUpdateSystem
type ITweenSystem = IUpdateSystem
type IVelocitySystem = IUpdateSystem
