//go:build ignore

// Demo tính năng Object Pooling (GamePooling)
// Run: go run .\tests\simulation\PoolDemo.go
package main

import (
	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
	"fmt"
	"image/color"
)

// --- 1. Bullet (Sử dụng Pool) ---
type Bullet struct {
	napi.IObject
	ncom.Pos
	ncom.Drw
	ncom.Info
	ncom.Box
	ncom.Deb

	VX, VY float32
	Life   int
}

// Khai báo Pool cho đạn. MaxSize = 50.
var BulletPool *napi.ObjectPool[*Bullet]

func init() {
	BulletPool = napi.NewObjectPool(napi.PoolConfig[*Bullet]{
		New: func() *Bullet {
			b := &Bullet{}
			// Nạp Component
			napi.Obj.NewObject(b, "bullet", "pos drw inf box deb")
			b.SetBoxW(10)
			b.SetBoxH(10)
			b.SetShape(napi.BSCircle)
			return b
		},
		Reset: func(b *Bullet) {
			// Hàm này được gọi mỗi khi Get() đạn từ kho
			b.Life = 120 // sống 120 frames (2s)
			b.Debug("")  // mặc định không bật debug
		},
		MaxSize: 50,
	})
}

func (b *Bullet) OnStep() {
	b.SetX(b.X() + b.VX)
	b.SetY(b.Y() + b.VY)
	b.Life--

	// Khi đạn hết máu, gọi Remove.
	// Nhờ hệ thống Auto-Routing, Engine sẽ tự đưa đạn vào lại BulletPool.
	if b.Life <= 0 {
		napi.Obj.Remove(b)
	}
}

func (b *Bullet) Draw() {
	b.Circle(b.X(), b.Y(), 5, color.RGBA{255, 255, 0, 255})
}

// --- 2. Boss (KHÔNG sử dụng Pool) ---
type Boss struct {
	napi.IObject
	ncom.Pos
	ncom.Drw
	ncom.Info
	ncom.Deb

	Life int
}

func NewBoss() *Boss {
	b := &Boss{Life: 300}
	napi.Obj.NewObject(b, "boss", "pos drw inf deb")
	b.SetX(320)
	b.SetY(240)
	napi.Obj.Register(b, "")
	b.Debug("info") // Bật debug để in tên
	fmt.Println("new boss")
	return b
}

func (b *Boss) OnStep() {
	b.Life--

	// Khi Boss chết, gọi Remove.
	// Vì Boss KHÔNG CÓ POOL, Engine sẽ xóa sạch Boss khỏi ECS RAM.
	if b.Life <= 0 {
		napi.Obj.Remove(b)
	}
}

func (b *Boss) Draw() {
	b.Rect(b.X()-20, b.Y()-20, 40, 40, color.RGBA{255, 0, 0, 255})
	b.Text(fmt.Sprintf("Boss HP: %d", b.Life), b.X()-25, b.Y()-30, color.RGBA{255, 255, 255, 255})
}

// --- 3. Spawner System ---
type Spawner struct {
	napi.IObject
	ncom.Drw
	ncom.Info
	ncom.Inp
	timer int
}

func NewSpawner() *Spawner {
	s := &Spawner{}
	napi.Obj.NewObject(s, "spawner", "drw inf inp")
	napi.Obj.Register(s, "")
	s.ListenOn("space", "p", func(key string) { NewBoss() })
	return s
}

func (s *Spawner) OnStep() {
	s.timer++

	// Cứ mỗi 10 frame, bắn ra 2 viên đạn (dùng Pool)
	if s.timer%10 == 0 {
		b1 := BulletPool.Get("") // Get() tự lấy ra hoặc tạo mới
		b1.SetX(320)
		b1.SetY(400)
		b1.VX = -2
		b1.VY = -5

		b2 := BulletPool.Get("")
		b2.SetX(320)
		b2.SetY(400)
		b2.VX = 2
		b2.VY = -5
	}

	// Bấm Space để Spawn Boss

}

func (s *Spawner) Draw() {
	// Lấy số lượng Object trong Engine hiện tại
	totalObjects := len(napi.Scene.GetCurrentScene().GetMap().GetObjects())

	msg := fmt.Sprintf("In Scene: %d - Space: Spawn Boss (True)", totalObjects)
	s.Text(msg, 10, 10, color.RGBA{255, 255, 255, 255})
}

// --- Main ---
func main() {
	napi.Game.Init(napi.GameConfig{
		Title:  "Object Pooling Demo",
		Width:  640,
		Height: 480,
	})

	napi.Scene.NewSceneAndGo("main", "map-640-480")
	//	napi.Game.LoadFromFile("./SharedObject/Config.toml")

	NewSpawner()

	napi.Game.GameStart()
}
