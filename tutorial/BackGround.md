# Hình nền / Background

> **Vision**: Một object, vài method — background đơn giản, hiệu quả.
> **Vision**: One object, a few methods — simple, effective background.

---

## 1. Giải thích / Explanation

Background object là một thực thể nhúng mixin `ncom.Back` (`BackgroundComponent`) để hiển thị màu sắc hoặc sprite làm hình nền phía dưới tất cả các Object khác.
Background object embeds `ncom.Back` to display a color or sprite background.

**Tính năng / Features:**
- Màu nền đơn sắc / Solid color background
- Sprite nền tĩnh hoặc cuộn (Parallax) / Static or scrolling sprite
- Lặp theo trục X, Y / Repeat along X, Y axes
- Kéo giãn toàn màn hình / Full-screen stretch

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"image/color"

	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

// Background - hình nền game / game background
type Background struct {
	ncom.Object // Lifecycle
	ncom.Back    // Background component (Back)
}

// NewBackground - tạo background / create background
func NewBackground() *Background {
	bg := &Background{}
	
	// "bg" là token tương ứng của ncom.Back
	// "sce-main" để tự động thêm vào scene main
	napi.Obj.NewObject(bg, "MainBackground", "bg sce-main")
	return bg
}

// Create - thiết lập khi khởi tạo / setup on creation
func (bg *Background) Create() {
	// --- Dùng màu nền / Use solid color ---
	// Thiết lập màu sắc trực tiếp (Tùy chọn)
	bg.SetColor(color.RGBA{R: 30, G: 30, B: 50, A: 255})

	// --- Hoặc dùng sprite / Or use sprite ---
	// bg.SetSprite(napi.Assert.GetSprite("bg_sky"))

	// Lặp theo cả hai trục / Repeat on both axes
	// bg.SetRepeatX(true)
	// bg.SetRepeatY(true)

	// Cuộn nền (parallax) / Scrolling background (parallax)
	// Trục X cuộn 50 pixels mỗi giây, Trục Y không cuộn
	// bg.SetScrollSpeedX(50) 
	// bg.SetScrollSpeedY(0)

	// Kéo giãn toàn màn hình / Stretch to fill screen
	// bg.SetStretch(true)
}
```

Background được xử lý bởi `DrawSystem` để luôn vẽ đầu tiên, phía dưới Map và Sprite. Bạn có thể thay đổi thuộc tính cuộn (ScrollSpeed) ngay trong quá trình chơi (vd. `StepUpdate()`) để tạo hiệu ứng Parallax chân thực khi nhân vật chạy.
