# Vẽ tùy biến / Custom Drawing (ncom.Drw)

> **Vision**: Vẽ hình dạng cơ bản hoặc chữ viết trực tiếp lên màn hình mà không cần sprite.
> **Vision**: Draw basic shapes or text directly to the screen without needing sprites.

---

## 1. Giải thích / Explanation

Mặc định Engine vẽ `ncom.Spr` tự động. Tuy nhiên, nếu bạn muốn vẽ các thanh máu (HP bar), vẽ chữ (Text) hoặc các hình khối cơ bản (Hình tròn, Hình chữ nhật), bạn có thể dùng `ncom.Drw`.
By default, the engine draws `ncom.Spr` automatically. However, if you want to draw HP bars, Text, or basic shapes (Circle, Rectangle), you can use `ncom.Drw`.

**Điều kiện bắt buộc / Requirements:**
1. Thêm `ncom.Drw` vào Object của bạn. / Add `ncom.Drw` to your Object.
2. Gắn token `"drw"` trong `NewObject` (Token này tự động kéo theo `"pos"`). / Add `"drw"` token (which automatically pulls `"pos"`).
3. Ghi đè hàm `Draw()` — Đây là interface bắt buộc của Engine dành cho Custom Drawing. / Override the `Draw()` function.

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"image/color"

	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

// Khai báo màu sắc
var Red = color.RGBA{R: 255, G: 0, B: 0, A: 255}
var White = color.RGBA{R: 255, G: 255, B: 255, A: 255}

type HUD struct {
	ncom.Object
	ncom.Pos
	ncom.Drw // Bật tính năng vẽ tùy biến
}

func NewHUD() *HUD {
	h := &HUD{}
	// Gắn token "drw"
	napi.Obj.NewObject(h, "PlayerHUD", "drw sce-main")
	return h
}

// Draw - Hàm này sẽ được Engine tự động gọi mỗi frame (thuộc tính bắt buộc của IDraw)
// Draw - This function is automatically called by the Engine every frame
func (h *HUD) Draw() {
	// Lấy tọa độ hiện tại (nếu HUD có Pos)
	x := h.X()
	y := h.Y()

	// 1. Vẽ một hình chữ nhật (Ví dụ: Thanh máu)
	// Draw a rectangle (Example: HP Bar)
	// Tham số: Tọa độ X, Y, Chiều rộng, Chiều cao, Độ dày (0 = Fill nguyên khối), Màu sắc, Z-Index
	h.Rect(x+10, y+10, 100, 20, 0, Red, 100)

	// 2. Vẽ chữ (Text)
	// Tham số: Tọa độ X, Y, Nội dung, Cỡ chữ, Tên Font, Màu sắc, Z-Index
	h.Text(x+10, y+40, "Player: Hero", 16, "default", White, 100)

	// 3. Vẽ hình tròn (Ví dụ: Tâm ngắm)
	// Tham số: Tâm X, Tâm Y, Bán kính, Độ dày (0 = Fill), Màu sắc, Z-Index
	h.Circle(x+150, y+20, 10, 2, White, 100)
}
```

> **Lưu ý quan trọng:** Z-Index (tham số cuối cùng) trong các hàm vẽ `Drw` cho phép bạn quyết định thứ tự vẽ (Số càng lớn vẽ càng sau, đè lên số nhỏ). Mặc định các Sprite sẽ vẽ ở Z-Index = 0.
> **Important Note:** Z-Index (the last parameter) in `Drw` functions allows you to decide the drawing order (Larger numbers draw later, overlapping smaller numbers). Default sprites draw at Z-Index = 0.
