package core

import (
	"autoworld/modules/nsystem"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Camera định nghĩa một viewport nhìn vào Physical Map.
// Camera chịu trách nhiệm: xác định vùng nhìn, follow target, và render qua DrawSystem.
//
// Tọa độ (x, y) là góc trên-trái của viewport trong **map space**.
// DrawSystem sẽ trừ (x, y) khỏi tọa độ entity để chuyển sang screen space.
type Camera struct {
	x, y          float32 // vị trí camera trong map space
	width, height int     // kích thước viewport (pixel)
	target        IObject // optional: follow target (nil = không follow)
	drawSystem    IDrawSystem
	screen        *ebiten.Image
}

// NewCamera khởi tạo Camera mới với kích thước viewport cho trước.
// viewW, viewH thường bằng kích thước cửa sổ game.
func NewCamera(viewW, viewH int) *Camera {
	return &Camera{
		x:          0,
		y:          0,
		width:      viewW,
		height:     viewH,
		drawSystem: nsystem.NewDrawSystem(),
	}
}

// ─── Getters / Setters ────────────────────────────────────────────────────────

// X trả về tọa độ ngang của camera trong map space.
func (c *Camera) X() float32 { return c.x }

// Y trả về tọa độ dọc của camera trong map space.
func (c *Camera) Y() float32 { return c.y }

// SetX dịch chuyển camera đến tọa độ ngang mới trong map space.
func (c *Camera) SetX(x float32) { c.x = x }

// SetY dịch chuyển camera đến tọa độ dọc mới trong map space.
func (c *Camera) SetY(y float32) { c.y = y }

// Width trả về chiều rộng viewport (pixel).
func (c *Camera) Width() int { return c.width }

// Height trả về chiều cao viewport (pixel).
func (c *Camera) Height() int { return c.height }

// SetTarget đặt IObject làm mục tiêu để camera tự động theo dõi mỗi frame.
// Truyền nil để tắt follow.
func (c *Camera) SetTarget(obj IObject) { c.target = obj }

// ─── Update ───────────────────────────────────────────────────────────────────

// Update cập nhật vị trí camera mỗi frame.
// Nếu có target, camera được đặt sao cho target nằm ở giữa viewport.
// mapW, mapH dùng để clamp camera trong biên bản đồ (0 = không giới hạn).
func (c *Camera) Update(mapW, mapH int) {
	if c.target == nil {
		return
	}

	// Lấy vị trí target từ ECS (cần PositionData)
	entry := c.target.Entry()
	if entry == nil || !entry.HasComponent(Position) {
		return
	}
	posData := donburi.Get[PositionData](entry, Position)

	// Đặt camera sao cho target ở giữa viewport
	c.x = posData.X - float32(c.width)/2
	c.y = posData.Y - float32(c.height)/2

	// Clamp trong biên bản đồ
	if mapW > 0 {
		if c.x < 0 {
			c.x = 0
		}
		if c.x+float32(c.width) > float32(mapW) {
			c.x = float32(mapW - c.width)
		}
	}
	if mapH > 0 {
		if c.y < 0 {
			c.y = 0
		}
		if c.y+float32(c.height) > float32(mapH) {
			c.y = float32(mapH - c.height)
		}
	}
}

// ─── Render ───────────────────────────────────────────────────────────────────

// SetScreen thiết lập canvas đích (ebiten.Image) để render lên.
// Phải gọi mỗi frame trước Draw() — thực hiện bởi EbitenGame.
func (c *Camera) SetScreen(screen *ebiten.Image) {
	c.screen = screen
	c.drawSystem.SetScreen(screen)
}

// Draw render scene lên màn hình theo thứ tự:
//  1. Physical Map (physWorld) với camera offset (camX, camY) — entity ngoài viewport bị skip.
//  2. GUI Map (guiWorld) với offset 0,0 — luôn đè lên trên, không bị culling camera.
//
// guiWorld có thể nil nếu Scene không có GUI Map.
func (c *Camera) Draw(physWorld donburi.World, guiWorld donburi.World) {
	if c.screen == nil {
		return
	}
	// Vẽ Physical Map với camera offset
	c.drawSystem.Draw(physWorld, c.x, c.y)

	// Vẽ GUI Map không có camera offset (screen space)
	if guiWorld != nil {
		c.drawSystem.Draw(guiWorld, 0, 0)
	}
}
