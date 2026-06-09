# Module components - Component Mixins

## Purpose

The `components` module provides **ECS Component Mixins** — plain Go structs that
Custom Objects embed to gain engine capabilities (position, sprite, audio, drawing, etc.).

Each mixin wraps a corresponding ECS data struct from `domain` and exposes a clean
getter/setter interface. The engine binds mixins automatically via reflection in `napi.NewObject`.

> **Technical note**: Component types (Position, Sprite, Box, etc.) are defined and
> initialized centrally in `enginetype` to prevent race conditions at Go init time.

---

## Interface Segregation Principle (ISP)

Interfaces are kept small and focused. Each mixin implements exactly the interface
defined for it in `domain/ObjectInterface.go` — no more, no less.

---

## Available Component Mixins

### 0. InforComponent (always auto-added)
Provides identity and label for the entity. Auto-added by `napi.NewObject`.

| Method | Description |
|--------|-------------|
| `GetName() string` | Returns the display name |
| `GetId() int` | Returns the auto-incremented unique ID |

### 1. PositionComponent (token: `"pos"`)
Provides position tracking and movement.

| Method | Description |
|--------|-------------|
| `X() float32` | Current X coordinate |
| `Y() float32` | Current Y coordinate |
| `SetX(x float32)` | Set X coordinate |
| `SetY(y float32)` | Set Y coordinate |

### 2. SpriteComponent (token: `"spr"`)
Provides image display, animation, scaling and rotation.

### 3. BoxComponent (token: `"box"`)
Provides hitbox for physics and collision.

### 4. AudioComponent (token: `"aud"`)
Provides audio playback and volume/pitch control.

### 5. DirectionComponent (token: `"dir"`)
Provides movement direction angle.

### 6. InputComponent (token: `"inp"`)
Cung cấp lắng nghe sự kiện bàn phím với đầy đủ loại sự kiện, tự được poll bởi InputSystem.

| Phương thức | Mô tả |
|------------|-------|
| `ListenOn(key, eventType string, handler)` | Đăng ký lắng nghe phím/nhóm phím theo loại sự kiện |

**eventType string:**
- `"pressed"` — mỗi frame khi phím đang được GIỮ
- `"just_pressed"` — duy nhất 1 lần khi vừa NHẤN XUỐNG
- `"just_released"` — duy nhất 1 lần khi vừa THẢ RA

**handler** nhận `key string` là tên phím đã trigger (ví dụ: `"w"`, `"space"`, `"a"`).
Kéo dài tớt cho nhóm phím như `"alpha"`, `"wasd"` — handler biết được phím nào cụ thể đã trigger.

### 7. MouseComponent (token: không cần)
Cung cấp lắng nghe sự kiện chuột và truy cập trạng thái chuột. Không cần ECS token.
Injected `IInputManager` từ singleton của package.

| Phương thức | Mô tả |
|------------|-------|
| `MouseX() int` | Tọa độ ngang con trỏ chuột (pixel) |
| `MouseY() int` | Tọa độ dọc con trỏ chuột (pixel) |
| `WheelX() float64` | Độ cuộn trục X frame hiện tại |
| `WheelY() float64` | Độ cuộn trục Y frame hiện tại (dương=xuống, âm=lên) |
| `ListenMouseOn(button, eventType string, handler)` | Đăng ký lắng nghe nút chuột |

**button**: `"left"`, `"right"`, `"middle"`
**eventType string**: `"pressed"`, `"just_pressed"`, `"just_released"`
**handler** nhận `button string` là tên nút đã trigger.

### 7. BackgroundComponent (token: `"back"`)
Scene background image or color.

### 8. TilemapComponent (token: `"tile"`)
Tilemap grid rendering.

### 9. AlarmComponent (token: `"alrm"`)
Timer-based callbacks (countdown in frames).

### 10. VelocityComponent (token: `"velo"`)
Basic physics: velocity, friction, max speed.

### 11. TweenComponent (token: `"twn"`)
Smooth interpolation over time (move, scale, alpha).

### 12. CollisionComponent (token: `"col"`)
Tag-based collision callbacks.

### 13. GenericComponent[T] — Custom Mixin
Allows game devs to define their own ECS component with full getter/setter
and custom methods, auto-bound by the engine just like built-in mixins.

### 14. DrawComponent (token: `"drw"`)
Provides primitive drawing methods called inside the object's own `Draw()` method.
Requires token `"pos"` (added automatically via constraint).
The object must also implement `domain.IDraw` (i.e., have a `Draw()` method).

| Method | Description |
|--------|-------------|
| `Rect(x, y, w, h float32, c color.RGBA)` | Filled rectangle |
| `RectStroke(x, y, w, h float32, c color.RGBA, sw float32)` | Rectangle outline |
| `Circle(x, y, r float32, c color.RGBA)` | Filled circle |
| `CircleStroke(x, y, r float32, c color.RGBA, sw float32)` | Circle outline |
| `Text(text string, x, y float32, c color.RGBA)` | Text with default font |
| `TextEx(text string, x, y float32, c color.RGBA, scale float64)` | Text with scale |
| `SetTextAlign(hAlign, vAlign string)` | Căn lề chữ ("left"/"l", "center"/"c", "right"/"r", "justify"/"j") |
| `SetTextOverflow(maxWidth, maxHeight float32, mode string)` | Xử lý tràn chữ: mode = "visible"/"v", "hidden"/"h", "scale"/"s" |
| `Image(sprite ISpriteLW, idx int, x, y float32)` | Manual sprite frame draw |
| `SetFont(f font.Face)` | Override per-instance font |

All coordinates are in **map space** — camera offset is applied automatically.

**Usage:**
```go
type HUD struct {
    napi.IObject
    napi.Pos
    napi.Drw   // token "drw"
    score int
}

func (h *HUD) Draw() {
    h.Rect(10, 10, 200, 40, color.RGBA{0, 0, 0, 150})
    h.Text(fmt.Sprintf("Score: %d", h.score), 20, 35, color.RGBA{255, 255, 0, 255})
}
```

---

## How to use

Freely combine mixins by embedding them in a Custom Object:

```go
type Monster struct {
    napi.IObject
    napi.Pos
    napi.Spr
}

func NewMonster() *Monster {
    m := &Monster{}
    napi.Obj.NewObject(m, "monster", "pos spr sce-cur")
    return m
}
```
