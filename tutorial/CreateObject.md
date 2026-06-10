# Tạo Object / Create an Object

> **Vision**: Struct nhúng component, constructor ngắn gọn, lifecycle rõ ràng.
> **Vision**: Embed components in struct, short constructor, clear lifecycle.

---

## 1. Giải thích / Explanation

Mỗi object trong N-Engine là một **Go struct** nhúng các component token (Mixins).
Each object in N-Engine is a **Go struct** embedding component tokens (Mixins).

Các bước tạo object / Steps to create an object:
1. Định nghĩa struct nhúng `ncom.Object` và các Mixins (như `ncom.Pos`, `ncom.Spr`).
2. Viết constructor dùng `napi.Obj.NewObject` với chuỗi component.
3. Cấu hình Object.
4. Tùy chọn đăng ký vào scene (nếu không dùng modifier auto-register).
5. Ghi đè lifecycle methods (`Create`, `StepUpdate`, `Destroy`).

---

## 2. Component Tokens & Mixins

Khi tạo struct, bạn nhúng Mixin tương ứng. Khi gọi `NewObject`, bạn truyền token bằng chuỗi cách nhau bởi khoảng trắng.

| Token | Component Mixin | Chức năng / Purpose |
|---|---|---|
| (auto) | `ncom.Object` | Lifecycle cơ bản (Create, StepUpdate, Destroy) / Base lifecycle |
| (auto) | `ncom.Info` | Thông tin định danh / Info (ID, Name) |
| `pos` | `ncom.Pos` | Vị trí X, Y / Position X, Y |
| `spr` | `ncom.Spr` | Sprite/hình ảnh / Sprite/image |
| `inp` | `ncom.Inp` | Lắng nghe input / Input listener |
| `vel` | `ncom.Velo` | Vận tốc / Velocity |
| `box` | `ncom.Box` | Hitbox va chạm / Collision hitbox |
| `alr` | `ncom.Alrm` | Bộ đếm giờ / Alarm timer |
| `twn` | `ncom.Twn` | Tween animation |
| `aud` | `ncom.Aud` | Âm thanh / Audio |
| `bg`  | `ncom.Back`| Hình nền / Background |
| `til` | `ncom.Tile`| Tilemap |
| `dir` | `ncom.Dir` | Góc quay di chuyển / Direction |

---

## 3. Ví dụ / Code Example

```go
package objects

import (
	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

// Player - nhân vật người chơi / player character
type Player struct {
	ncom.Object // Bắt buộc / Required
	ncom.Pos     // Lấy getter/setter cho Tọa độ (SetX, SetY...)
	ncom.Spr     // Lấy getter/setter cho Ảnh (SetSprite...)
	ncom.Velo    // Lấy getter/setter cho Vận tốc (SetVelocity...)
	ncom.Inp     // Lấy khả năng lắng nghe Input (ListenOn)

	speed float32
	hp    int
}

// NewPlayer - constructor tạo player / constructor to create player
func NewPlayer(x, y float32) *Player {
	p := &Player{
		speed: 3.0,
		hp:    100,
	}

	// 1. Khởi tạo Object trong ECS với các tokens
	// "pos spr vel inp" tương ứng với Pos, Spr, Velo, Inp
	// "sce-main" là modifier để tự động đăng ký vào scene "main".
	napi.Obj.NewObject(p, "Player1", "pos spr vel inp sce-main")

	// 2. Cấu hình ban đầu
	p.SetX(x)
	p.SetY(y)

	// Lắng nghe phím
	p.ListenOn("space", p.OnJump)

	return p
}

// Create được gọi 1 lần khi object vào scene
func (p *Player) Create() {
	p.SetSprite("idle", napi.Assert.GetSprite("hero_idle"))
	p.SetCurrentSprite("idle")
	p.SetMaxSpeed(p.speed)
}

// StepUpdate được gọi mỗi khung hình (frame)
func (p *Player) StepUpdate() {
	if p.X() > 800 {
		p.SetX(0)
	}
}

// Hàm custom handler cho input
func (p *Player) OnJump() {
	p.SetVelocityY(-5.0)
}
```

---

## 4. Đăng ký tự động / Auto-Register

Trong chuỗi string khi gọi `napi.Obj.NewObject`, bạn có thể thêm các "modifier" đặc biệt để tự động đăng ký Object vào Scene:
- `sce-main` : Tự động thêm vào scene có ID "main".
- `sce-cur`  : Tự động thêm vào Scene hiện tại đang active.
- `sce-glo`  : Tự động thêm vào Global Hidden Scene (tồn tại xuyên suốt các scene).

Nếu không dùng modifier, bạn phải tự gọi `napi.Obj.Register(obj, "sceneId")`.
