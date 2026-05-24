package core

import (
	"autoworld/domain"
	"autoworld/modules/nsystem"

	"github.com/yohamta/donburi"
)

// Map quản lý ECS World và toàn bộ logic update của một bản đồ game.
// Map không biết gì về render — Camera chịu trách nhiệm đó.
// Hai loại Map trong một Scene:
//   - Physical Map: thế giới game, tọa độ map space, có camera scroll.
//   - GUI Map: HUD/overlay, tọa độ screen space, không camera offset.
type Map struct {
	logicSystem    ILogicSystem
	audioSystem    IAudioSystem
	inputSystem    IInputSystem
	alarmSystem    IAlarmSystem
	tweenSystem    ITweenSystem
	velocitySystem IVelocitySystem
	world          donburi.World
	objectList     []IObject
	width          int // chiều rộng bản đồ (pixel), 0 = không giới hạn
	height         int // chiều cao bản đồ (pixel), 0 = không giới hạn
}

// NewMap khởi tạo Map mới với kích thước cho trước.
// Truyền width=0, height=0 nếu bản đồ không có giới hạn (cuộn vô hạn).
// input được dùng bởi InputSystem để kiểm tra phím mỗi frame.
func NewMap(input domain.IInputManager, width, height int) *Map {
	return &Map{
		logicSystem:    nsystem.NewLogicSystem(),
		audioSystem:    nsystem.NewAudioSystem(),
		alarmSystem:    nsystem.NewAlarmSystem(),
		tweenSystem:    nsystem.NewTweenSystem(),
		velocitySystem: nsystem.NewVelocitySystem(),
		inputSystem:    nsystem.NewInputSystem(input),
		world:          donburi.NewWorld(),
		width:          width,
		height:         height,
	}
}

// NewGUIMap khởi tạo GUI Map (screen-space overlay).
// GUI Map không cần AudioSystem — âm thanh thuộc Physical Map.
// Kích thước = kích thước màn hình (viewW, viewH).
func NewGUIMap(input domain.IInputManager, viewW, viewH int) *Map {
	return &Map{
		logicSystem:    nsystem.NewLogicSystem(),
		audioSystem:    nsystem.NewAudioSystem(),
		alarmSystem:    nsystem.NewAlarmSystem(),
		tweenSystem:    nsystem.NewTweenSystem(),
		velocitySystem: nsystem.NewVelocitySystem(),
		inputSystem:    nsystem.NewInputSystem(input),
		world:          donburi.NewWorld(),
		width:          viewW,
		height:         viewH,
	}
}

// Update chạy toàn bộ logic mỗi frame theo thứ tự: Input → Logic → Audio.
func (m *Map) Update() error {
	// hệ thống cốt lõi
	// Ưu tiên trước -> có kích hoạt các sýtem khác trong cùng 1 frames
	// chạy hàm Create -> phải được chạy trước tiên
	m.logicSystem.Update(m.objectList)

	// 2 system có thể kích hoạt function
	// có sự kiện input -> thé tận dụng cơ chế ghi đè để được ưu tiên trước StepUpdate trong logicSystem
	m.inputSystem.Update(m.objectList)
	m.alarmSystem.Update(m.objectList)

	// 2 system hỗ trợ tính năng
	// Sau khi logic đã chính đã hoàn thành, triển khai các hỗ trợ
	// Luôn chạy sau các sự kiện cập nhật, Dev có thể setup trước trong các hàm chính
	m.tweenSystem.Update(m.objectList)
	m.velocitySystem.Update(m.objectList)

	// Cần được active từ các system chính
	// Không bị phụ thuộc nên không cần ưu tiên
	m.audioSystem.Update(m.world)
	return nil
}

// AddObject đăng ký IObject vào Map để được xử lý mỗi frame.
// LogicSystem sẽ gọi Create() vào frame tiếp theo.
func (m *Map) AddObject(obj IObject) {
	m.logicSystem.AddObjectCreated(obj)
	m.objectList = append(m.objectList, obj)
}

// World trả về donburi.World của Map.
func (m *Map) World() donburi.World {
	return m.world
}

// Width trả về chiều rộng bản đồ (pixel). 0 = không giới hạn.
func (m *Map) Width() int {
	return m.width
}

// Height trả về chiều cao bản đồ (pixel). 0 = không giới hạn.
func (m *Map) Height() int {
	return m.height
}
