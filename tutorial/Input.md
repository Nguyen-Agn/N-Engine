# Bàn phím & Chuột / Keyboard & Mouse (ncom.Inp & ncom.Mouse)

> **Vision**: Gọi callback khi người dùng bấm phím/chuột thay vì check trạng thái liên tục trong update loop.
> **Vision**: Call callbacks on input events instead of continuously checking state in the update loop.

---

## 1. Giải thích / Explanation

- `ncom.Inp` (Keyboard): Lắng nghe sự kiện từ bàn phím.
- `ncom.Mouse` (Mouse): Cung cấp tọa độ chuột và lắng nghe sự kiện nút bấm chuột. (Lưu ý: `ncom.Mouse` không cần chuỗi token lúc tạo `NewObject` vì nó không lưu ECS data).

Các loại sự kiện / Event Types (từ package `domain`):
- `domain.EventJustPressed`: Vừa mới bấm xuống (Kích hoạt 1 lần). / Just pressed down (Fires once).
- `domain.EventPressed`: Đang giữ (Kích hoạt mỗi frame). / Holding down (Fires every frame).
- `domain.EventJustReleased`: Vừa nhả ra (Kích hoạt 1 lần). / Just released (Fires once).

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"fmt"

	"autoworld/domain"
	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
)

type InputController struct {
	ncom.Object
	ncom.Inp   // Lắng nghe Keyboard
	ncom.Mouse // Xử lý Chuột
}

func NewInputController() *InputController {
	i := &InputController{}
	// Token "inp" dành cho bàn phím, Mouse không cần token
	napi.Obj.NewObject(i, "GlobalInput", "inp sce-main")
	return i
}

func (i *InputController) Create() {
	// --- Bàn Phím / Keyboard ---
	// Lắng nghe phím Space vừa mới bấm xuống
	// Listen to Space key just pressed
	i.ListenOn("space", domain.EventJustPressed, func(key string) {
		fmt.Println("Vừa bấm phím / Just pressed:", key)
	})

	// Lắng nghe nhiều phím cùng lúc (Dùng mảng)
	// Listen to multiple keys (Using array)
	i.ListenList([]string{"w", "a", "s", "d"}, domain.EventPressed, i.onMove)

	// --- Chuột / Mouse ---
	// Lắng nghe chuột trái vừa nhả ra
	// Listen to left mouse button just released
	i.ListenMouseOn("left", domain.EventJustReleased, func(btn string) {
		// Lấy tọa độ chuột hiện tại
		// Get current mouse coordinates
		fmt.Printf("Nhả chuột trái tại / Released left at: X=%d, Y=%d\n", i.MouseX(), i.MouseY())
	})
}

// Handler cho di chuyển (Gọi mỗi frame nếu phím WASD đang bị giữ)
// Movement handler (Called every frame if WASD is held)
func (i *InputController) onMove(key string) {
	switch key {
	case "w":
		// Di chuyển lên
	case "s":
		// Di chuyển xuống
	}
}

func (i *InputController) StepUpdate() {
	// Lấy tốc độ cuộn chuột (Scroll wheel)
	// Get scroll wheel velocity
	if wheelY := i.WheelY(); wheelY != 0 {
		fmt.Println("Cuộn chuột / Scrolled:", wheelY)
	}
}
```
