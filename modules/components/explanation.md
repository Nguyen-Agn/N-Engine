# Module components — Engine Support Components (Mixins)

## Mục tiêu

Module `components` cung cấp các **Component Mixins** cho game developer.
Từ khi `domain.IObject` được rút gọn về các tác vụ vòng đời cốt lõi (`StepUpdate`, `Destroy`), việc tương tác với các thuộc tính cụ thể (tọa độ, hình ảnh, âm thanh, hitbox) được tách biệt vào các mixin chuyên biệt này để đảm bảo:
1. **An toàn kiểu dữ liệu lúc biên dịch (Compile-time Safety)**: Nếu Object không nhúng `AudioComponent`, IDE sẽ báo đỏ ngay khi gõ `Play` hoặc `SetVolume`.
2. **Nguyên lý Interface Segregation (ISP)**: Giao diện sạch sẽ, linh hoạt, tối ưu cho lập trình hướng đối tượng trong Go.

> **Lưu ý kỹ thuật**: Các component types tương ứng (`Position`, `Sprite`, `Box`, `Audio`, `Infor`, `Direction`) được định nghĩa và khởi tạo tập trung tại module `enginetype` để đảm bảo an toàn về thứ tự khởi tạo trong Go, loại bỏ hoàn toàn các lỗi race condition lúc runtime.


---

## Các Component Mixins có sẵn

### 0. `InforComponent` (Bắt buộc - Tự động đính kèm)
Cung cấp định danh và nhãn hiển thị cho thực thể. Mọi Object khi được khởi tạo qua `napi.NewObject` đều tự động đính kèm component này.

| Phương thức | Mô tả |
|-------------|-------|
| `GetName() string` | Trả về tên hiển thị của thực thể |
| `GetId() int` | Trả về ID số nguyên duy nhất tự động tăng |

### 1. `PositionComponent` (Token: `"pos"`)
Cung cấp khả năng di chuyển và quản lý tọa độ của thực thể.

| Phương thức | Mô tả |
|-------------|-------|
| `X() float32` | Trả về tọa độ X hiện tại |
| `Y() float32` | Trả về tọa độ Y hiện tại |
| `SetX(x float32)` | Đặt tọa độ X mới |
| `SetY(y float32)` | Đặt tọa độ Y mới |

### 2. `SpriteComponent` (Token: `"spr"`)
Cung cấp khả năng hiển thị, lật ảnh, xoay và chỉnh tỉ lệ (Scale).

| Phương thức | Mô tả |
|-------------|-------|
| `SpriteIdx() int` / `ImageIndex() int` | Lấy frame hình hiện tại |
| `SetSpriteIdx(idx int)` / `SetImageIndex(idx int)` | Đặt frame hình cụ thể |
| `ImageSpeed() float32` / `SetImageSpeed(speed float32)` | Tốc độ chuyển frame hình |
| `Rotation() float32` / `SetRotation(r float32)` | Góc xoay của ảnh |
| `ScaleX() float32` / `SetScaleX(sx float32)` | Tỉ lệ co giãn ngang |
| `ScaleY() float32` / `SetScaleY(sy float32)` | Tỉ lệ co giãn dọc |
| `Sprite(name string) ISpriteLW` | Lấy Sprite theo tên |
| `SetSprite(name string, sprite ISpriteLW)` | Gán Sprite |
| `NextImage()` | Chuyển sang frame hình kế tiếp |

### 3. `BoxComponent` (Token: `"box"`)
Cung cấp hitbox cho vật lý và va chạm.

| Phương thức | Mô tả |
|-------------|-------|
| `BoxW() float32` / `SetBoxW(w float32)` | Chiều rộng hitbox |
| `BoxH() float32` / `SetBoxH(h float32)` | Chiều cao hitbox |
| `BoxX() float32` / `SetBoxX(x float32)` | Điểm neo X của hitbox |
| `BoxY() float32` / `SetBoxY(y float32)` | Điểm neo Y của hitbox |
| `IsCollidable() bool` / `SetIsCollidable(c bool)` | Trạng thái có thể va chạm hay không |

### 4. `AudioComponent` (Token: `"aud"`)
Cung cấp khả năng phát âm thanh và quản lý âm lượng, cao độ.

| Phương thức | Mô tả |
|-------------|-------|
| `Play(name string, volume, pitch float32)` | Phát âm thanh bất kỳ |
| `PlayDefault(name string)` | Phát âm thanh với Volume=1.0, Pitch=1.0 |
| `StopAudio()` | Dừng âm thanh đang phát |
| `SetVolume(vol float32)` / `Volume() float32` | Quản lý âm lượng |
| `SetPitch(pitch float32)` / `Pitch() float32` | Quản lý cao độ |

### 5. `DirectionComponent` (Token: `"dir"`)
Cung cấp góc hướng di chuyển cho thực thể.

| Phương thức | Mô tả |
|-------------|-------|
| `Direction() float32` | Trả về góc xoay hiện tại (0 - 360) |
| `SetDirection(dir float32)` | Thiết lập góc xoay mới |
| `Rotate(dir float32)` | Xoay thêm một góc tương ứng |

### 6. `InputComponent` (Token: `"inp"`)
Cung cấp khả năng lắng nghe phím và kích hoạt callback tương ứng mỗi frame.
ECS System tự động kiểm tra phím và gọi hàm — Object không cần viết code polling thủ công.

| Phương thức | Mô tả |
|-------------|-------|
| `ListenOn(key string, handler func())` | Đăng ký phím/nhóm phím và hàm xử lý tương ứng |

**Key Groups đặc biệt**:
| Tên nhóm | Phím được lắng nghe |
|----------|--------------------|
| `"alpha"` | a–z (dùng cho nhập văn bản) |
| `"number"` | 0–9 (dùng cho nhập số) |
| `"arrows"` | ↑ ↓ ← → |
| `"wasd"` | W A S D |
| `"all"` | Toàn bộ phím |

### 7. `GenericComponent[T]` — Custom Mixin cho Game Dev

Cho phép game dev tự tạo Component Mixin riêng với đầy đủ getter/setter và method tùy biến,
mà vẫn được engine **auto-bind** giống như các mixin built-in.

| Hàm / Phương thức | Mô tả |
|-------------------|-------|
| `NewGenericComponent(comp)` | Constructor — phải gọi trước `napi.NewObject` |
| `Get() *T` | Trả về con trỏ đến data trong ECS (thay đổi qua pointer có hiệu lực ngay) |
| `Set(val T)` | Gán toàn bộ giá trị mới cho component data |
| `BindComponent(base IObject)` | Tự động gọi bởi engine — game dev không cần gọi trực tiếp |

**Cách 1 — Dùng trực tiếp (chỉ cần Get/Set)**:

```go
var StatsComp = napi.NewComponentType[StatsData]("sta")

type Hero struct {
    napi.IObject
    components.PositionComponent
    components.GenericComponent[StatsData]
}

func NewHero() *Hero {
    h := &Hero{
        GenericComponent: components.NewGenericComponent(StatsComp),
    }
    napi.NewObject(h, "hero", "pos sta")
    napi.Register(h, "")
    return h
}

func (h *Hero) StepUpdate() {
    h.Get().Health -= 1  // truy cập trực tiếp qua Get()
}
```

**Cách 2 — Bọc trong struct riêng (thêm method tùy biến)**:

```go
var StatsComp = napi.NewComponentType[StatsData]("sta")

// Bọc GenericComponent để thêm method
type StatsComponent struct {
    components.GenericComponent[StatsData]
}

func (s *StatsComponent) TakeDamage(amount int) {
    if s.Get() != nil {
        s.Get().Health -= amount
    }
}

func (s *StatsComponent) IsAlive() bool {
    return s.Get() != nil && s.Get().Health > 0
}

// Object nhúng StatsComponent — method được promoted tự nhiên
type Hero struct {
    napi.IObject
    components.PositionComponent
    StatsComponent
}

func NewHero() *Hero {
    h := &Hero{
        StatsComponent: StatsComponent{
            GenericComponent: components.NewGenericComponent(StatsComp),
        },
    }
    napi.NewObject(h, "hero", "pos sta")
    napi.Register(h, "")
    return h
}

func (h *Hero) StepUpdate() {
    h.TakeDamage(1)       // method promoted từ StatsComponent lên Hero
    if !h.IsAlive() {
        h.Destroy()
    }
}
```

---

## Cách sử dụng

Tự do phối hợp nhúng các component cần thiết vào Custom Object của bạn:

```go
package game

import (
	"github.com/user/autoworld/modules/napi"
	"github.com/user/autoworld/modules/components"
)

type Monster struct {
	napi.IObject
	components.InforComponent    // Luôn nhúng vì InforComponent là bắt buộc
	components.PositionComponent // Cần tọa độ di chuyển
	components.SpriteComponent   // Cần hình vẽ hiển thị
}

func NewMonster() *Monster {
	base := napi.NewObject("monster", "pos spr") // Không truyền "aud"
	m := &Monster{
		IObject:           base,
		InforComponent:    components.InforComponent{IObject: base},
		PositionComponent: components.PositionComponent{IObject: base},
		SpriteComponent:   components.SpriteComponent{IObject: base},
	}
	napi.Register(m, false)
	return m
}
```
