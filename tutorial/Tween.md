# Chuyển động mượt / Tweening (ncom.Twn)

> **Vision**: Hiệu ứng chuyển động mượt mà chỉ với một hàm duy nhất.
> **Vision**: Smooth transition animations with a single function call.

---

## 1. Giải thích / Explanation

Tween Component (`ncom.Twn`) cung cấp các hàm để tạo hiệu ứng biến đổi (di chuyển, phóng to, làm mờ) một cách mượt mà theo thời gian.
The Tween Component (`ncom.Twn`) provides functions to create smooth transition effects (movement, scaling, fading) over time.

Lưu ý: Nhúng `ncom.Twn` (sử dụng token `"twn"`) sẽ tự động yêu cầu và cấp phát các component `pos` và `spr`.
Note: Embedding `ncom.Twn` (using the `"twn"` token) will automatically require and allocate `pos` and `spr` components.

**Các loại Tween / Tween Types:**
- `domain.TLinear`: Chuyển động đều / Linear
- `domain.TEaseInQuad`, `domain.TEaseOutQuad`: Nhanh dần, chậm dần / Ease In, Ease Out
- `domain.TEaseInOutQuad`: Chậm-Nhanh-Chậm / Ease In Out
- `domain.TBounceOut`: Nảy / Bounce

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"autoworld/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

type Popup struct {
	ncom.Object
	ncom.Pos
	ncom.Spr
	ncom.Twn // Thêm Tween component
}

func NewPopup(startX, startY float32) *Popup {
	p := &Popup{}
	
	// "twn" tự động kéo theo "pos" và "spr"
	napi.Obj.NewObject(p, "UIPopup", "twn sce-main")
	
	p.SetX(startX)
	p.SetY(startY)
	return p
}

func (p *Popup) Create() {
	p.SetSprite("panel", napi.Assert.GetSprite("ui_panel"))
	p.SetCurrentSprite("panel")

	// 1. Di chuyển mượt (TweenMove)
	// Di chuyển từ vị trí hiện tại đến (400, 300) trong 60 frames, dùng hiệu ứng EaseOutQuad
	p.TweenMove(400, 300, 60, domain.TEaseOutQuad)

	// 2. Phóng to mượt (TweenScale)
	// Bắt đầu từ 0.0 (ẩn), phóng to lên 1.0 trong 30 frames, dùng hiệu ứng nảy (BounceOut)
	p.SetScaleX(0.0)
	p.SetScaleY(0.0)
	p.TweenScale(1.0, 1.0, 30, domain.TBounceOut)

	// 3. Làm mờ (TweenAlpha)
	// (Ví dụ: mờ dần đi về 0.0 trong vòng 120 frames)
	// p.TweenAlpha(0.0, 120, domain.TLinear)
}
```

Các Tween có thể chạy đồng thời trên cùng một object. Hệ thống sẽ tự động cập nhật và kết thúc chúng khi thời gian (`frames`) chạy hết.
Tweens can run simultaneously on the same object. The system automatically updates and ends them when the time (`frames`) expires.
