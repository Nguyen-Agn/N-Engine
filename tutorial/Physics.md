# Vật lý cơ bản / Velocity & Direction (ncom.Velo & ncom.Dir)

> **Vision**: Di chuyển Object theo các tham số vật lý mượt mà thay vì cộng trừ tọa độ thủ công.
> **Vision**: Move Objects using smooth physics parameters instead of manually manipulating coordinates.

---

## 1. Giải thích / Explanation

- `ncom.Velo` (Velocity): Quản lý vận tốc, tốc độ tối đa, và lực ma sát. Object sẽ tự động cập nhật vị trí (`ncom.Pos`) dựa vào vận tốc mỗi frame.
- `ncom.Dir` (Direction): Quản lý góc quay (Tính bằng Độ / Degrees), giúp Object tự động di chuyển tới một hướng nhất định và tự động xoay góc của chính nó.

Lưu ý: Khi dùng `ncom.Velo`, bạn luôn phải gắn kèm token `"pos"` để hệ thống có thể cập nhật vị trí.
Note: When using `ncom.Velo`, you must always include the `"pos"` token so the system can update coordinates.

---

## 2. Ví dụ Vận tốc cơ bản / Velocity Example

```go
package objects

import (
	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

type Car struct {
	ncom.Object
	ncom.Pos
	ncom.Spr
	ncom.Velo // Kích hoạt Vật lý vận tốc
}

func NewCar() *Car {
	c := &Car{}
	napi.Obj.NewObject(c, "Car", "pos spr vel sce-main")
	return c
}

func (c *Car) OnCreate() {
	// Giới hạn tốc độ tối đa (Max Speed)
	// Limit maximum speed
	c.SetMaxSpeed(10.0)

	// Đặt ma sát (Friction) để xe tự động giảm tốc dần về 0 khi không có lực tác dụng
	// Set friction so the car slows down automatically
	c.SetFriction(0.1)
}

func (c *Car) Accelerate() {
	// Thêm vận tốc vào trục X (1.0 mỗi lần gọi)
	// Add velocity to X axis
	c.AddVelocityX(1.0)
	
	// Hoặc Set cứng: c.SetVelocityX(5.0)
}
```

---

## 3. Ví dụ Phương hướng / Direction Example

Sử dụng `ncom.Dir` kết hợp `ncom.Velo` để di chuyển theo tọa độ cực (Góc và Vận tốc).
Use `ncom.Dir` combined with `ncom.Velo` to move using polar coordinates (Angle and Velocity).

```go
type Missile struct {
	ncom.Object
	ncom.Pos
	ncom.Spr
	ncom.Velo
	ncom.Dir // Kích hoạt Hướng
}

func NewMissile() *Missile {
	m := &Missile{}
	// Sử dụng token "dir"
	napi.Obj.NewObject(m, "Missile", "pos spr vel dir sce-main")
	return m
}

func (m *Missile) Create() {
	// Cài đặt hướng bay ban đầu là 90 độ (Đi xuống dưới theo hệ tọa độ Game 2D)
	// Set initial direction to 90 degrees (Downwards in 2D Game coords)
	m.SetDirection(90)

	// Xoay Object từ từ về hướng 180 độ, mỗi frame xoay 5 độ
	// Rotate object slowly towards 180 degrees by 5 degrees per frame
	// (Hàm này có thể được gọi trong StepUpdate để bám theo mục tiêu)
	m.RotateTowards(180, 5.0)

	// Đặt tốc độ bay là 8.0 theo góc hiện tại
	// Set flying speed to 8.0 in the current angle
	m.SetVelocityInDirection(8.0)
}
```
