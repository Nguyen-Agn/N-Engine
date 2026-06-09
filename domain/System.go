package domain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ─── ECS System Interfaces ────────────────────────────────────────────────────

// ILogicSystem điều phối logic game mỗi frame trên tập hợp Object.
type ILogicSystem interface {
	// Update gọi Create() cho object mới, StepUpdate() cho tất cả, Destroy() cho object bị xóa.
	// Purpose: Orchestrates the logic lifecycle (Create, Step, Destroy) for all managed objects in a frame.
	// Inputs: objectList []IObject - The list of objects to update.
	// Outputs: None.
	Update(objectList []IObject)

	// AddObjectCreated đăng ký Object vừa tạo — sẽ gọi Create() vào frame tiếp theo.
	// Purpose: Registers a newly spawned object to have its Create() method called on the next frame.
	// Inputs: o IObject - The newly created object.
	// Outputs: None.
	AddObjectCreated(o IObject)

	// AddObjectDestroy đăng ký Object sẽ xóa — sẽ gọi Destroy() cuối frame.
	// Purpose: Registers an object scheduled for deletion to have its Destroy() method called.
	// Inputs: o IObject - The object marked for destruction.
	// Outputs: None.
	AddObjectDestroy(o IObject)
}

// IDrawSystem xử lý render đồ họa mỗi frame.
// Truy vấn ECS World để tìm entity cần vẽ và áp dụng camera offset.
type IDrawSystem interface {
	// Draw render tất cả entity có Sprite + Position lên screen.
	// camX, camY là tọa độ camera trong map space — được trừ khỏi tọa độ entity
	// để chuyển sang screen space. Truyền 0, 0 cho GUI layer (không có camera offset).
	// Purpose: Renders all drawable entities to the screen, applying the camera offset.
	// Inputs:
	//   - w donburi.World: The ECS world containing the entities.
	//   - camX float32: The camera's X position in the world.
	//   - camY float32: The camera's Y position in the world.
	// Outputs: None.
	Draw(w donburi.World, camX, camY float32)

	// SetScreen thiết lập canvas đích (ebiten.Image) để DrawSystem vẽ lên.
	// Phải được gọi trước Draw() mỗi frame.
	// Purpose: Sets the destination screen buffer for rendering operations.
	// Inputs: s *ebiten.Image - The destination canvas.
	// Outputs: None.
	SetScreen(s *ebiten.Image)
}

// IAudioSystem xử lý phát âm thanh mỗi frame.
type IAudioSystem interface {
	// Update đọc AudioData của từng entity và thực thi lệnh phát/dừng âm thanh.
	// Purpose: Processes queued audio commands (play, stop, pause) for all audio-enabled entities.
	// Inputs: w donburi.World - The ECS world containing the entities.
	// Outputs: None.
	Update(w donburi.World)
}

// IUpdateSystem là interface chung cho các system thực thi logic cập nhật mỗi frame.
type IUpdateSystem interface {
	// Purpose: Executes a per-frame update loop for a specific subsystem logic.
	// Inputs: objectList []IObject - The list of active objects to process.
	// Outputs: None.
	Update(objectList []IObject)
}

type IAlarmSystem = IUpdateSystem
type ITweenSystem = IUpdateSystem
type IVelocitySystem = IUpdateSystem

// IDrawObjectRegistry allows Map to register IDraw objects into DrawSystem
// without depending directly on the nsystem package.
// DrawSystem implements this interface.
type IDrawObjectRegistry interface {
	// AddDrawObject registers an object whose Draw() will be called each frame.
	// Purpose: Registers a custom object into the draw loop for manual rendering.
	// Inputs: obj IObject - The object implementing IDraw.
	// Outputs: None.
	AddDrawObject(obj IObject)
	
	// RemoveDrawObject unregisters an object from the draw loop.
	// Safe to call even if the object was never registered.
	// Purpose: Unregisters an object from the draw loop.
	// Inputs: obj IObject - The object to remove.
	// Outputs: None.
	RemoveDrawObject(obj IObject)
}
