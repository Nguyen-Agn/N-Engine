# Module domain — Interface và Data Contract của Engine

## Mục tiêu

Module `domain` là **tầng hợp đồng (Contract Layer)** của toàn bộ Engine.
Nó định nghĩa **chỉ các Interface và Data Struct** — không chứa bất kỳ logic triển khai nào. Mọi module khác phụ thuộc vào `domain`, nhưng `domain` không phụ thuộc vào bất kỳ module nào khác (ngoại trừ `ebiten` cho một số kiểu dữ liệu đồ họa).

---

## Các file chính

### `ObjectInterface.go` — Các Interface của Object và Component

| Interface | Mô tả |
|-----------|-------|
| `IObject` | Interface gốc của mọi đối tượng: `Create()`, `StepUpdate()`, `Destroy()`, `Entry()` |
| `IPosition` | Tọa độ: `X()`, `Y()`, `SetX()`, `SetY()` |
| `ISprite` | Hình ảnh, animation, scale, rotation |
| `IBox` | Hitbox, va chạm |
| `IAudio` | Phát, dừng âm thanh |
| `IDirection` | Góc hướng di chuyển |
| `IInput` | Đăng ký lắng nghe phím: `ListenOn(key, handler)` |
| `IInfor` | ID và tên: `GetId()`, `GetName()` |

### `ObjectData.go` — Cấu trúc dữ liệu ECS

Định nghĩa các Data Struct được lưu trữ trong ECS (donburi):

| Struct | Mô tả |
|--------|-------|
| `PositionData` | X, Y |
| `SpriteData` | Sprite map, frame hiện tại, scale, rotation |
| `BoxData` | Kích thước và vị trí hitbox |
| `AudioData` | Âm thanh, volume, pitch, flag phát/dừng |
| `DirectionData` | Góc hướng |
| `InputData` | Danh sách `KeyBinding` (phím + callback) |
| `InforData` | ID, Name |

### `IInput.go` — Định nghĩa Input System

- **`Key`**: Type và toàn bộ hằng số phím bấm (độc lập với ebiten).
- **`MouseButton`**: Các nút chuột.
- **`IInputManager`**: Interface polling input (Pressed, JustPressed, JustReleased), chuột, wheel và Virtual Action.
- **`KeyNameMap`**: Map từ chuỗi dễ nhớ (`"w"`, `"space"`) sang hằng số `Key`.
- **`KeyGroupMap`**: Map nhóm phím đặc biệt (`"alpha"`, `"number"`, `"arrows"`, `"wasd"`, `"all"`).

### `System.go` — Interface của các ECS System

| Interface | Mô tả |
|-----------|-------|
| `ILogicSystem` | Điều phối logic: `Update(objectList)`, `AddObjectCreated`, `AddObjectDestroy` |
| `IDrawSystem` | Vẽ: `Draw(world)`, `SetScreen(screen)` |
| `IAudioSystem` | Âm thanh: `Update(world)` |
| `IInputSystem` | Input: `Update(objectList)` |

### `Engine.go` / `IScene.go` / `IGlobalConfig.go` / `nglobal.go`

Các interface điều phối cấp Engine và Scene.

---

## Quy tắc quan trọng cho Agent

1. **Chỉ Interface và Struct** — tuyệt đối không thêm logic vào đây.
2. **Không import module nội bộ** — `domain` chỉ được import `ebiten` và `donburi`.
3. Khi thêm Component mới, cần: (1) thêm `XxxData` vào `ObjectData.go`, (2) thêm `IXxx` vào `ObjectInterface.go`, (3) đăng ký type tại `enginetype/Register.go`.
