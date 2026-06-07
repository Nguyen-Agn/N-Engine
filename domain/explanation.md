# Module domain — Interface và Data Contract của Engine

## Mục tiêu

Module `domain` là **tầng hợp đồng (Contract Layer)** của toàn bộ Engine.
Nó định nghĩa **chỉ các Interface và Data Struct** — không chứa bất kỳ logic triển khai nào. Mọi module khác phụ thuộc vào `domain`, nhưng `domain` không phụ thuộc vào bất kỳ module nào khác (ngoại trừ `ebiten` cho một số kiểu dữ liệu đồ họa).

---

## Các file chính

### `ObjectInterface.go` — Các Interface của Object và Component

| Interface | Mô tả |
|-----------|-------|
| `IObject` | Interface gốc của mọi đối tượng: `OnCreate()`, `OnStep()`, `OnDestroy()`, `OnSave()`, `OnLoad()`, `Entry()` |
| `IPosition` | Tọa độ: `X()`, `Y()`, `SetX()`, `SetY()` |
| `ISprite` | Hình ảnh, animation, scale, rotation |
| `IBox` | Hitbox, va chạm |
| `IAudio` | Phát, dừng âm thanh |
| `IDirection` | Góc hướng di chuyển |
| `IInput` | Đăng ký lắng nghe phím: `ListenOn(key, eventType string, handler)` |
| `IMouse` | Chuột: `MouseX()`, `MouseY()`, `WheelX()`, `WheelY()`, `ListenMouseOn(btn, eventType string, handler)` |
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
| `InputData` | Danh sách `KeyBinding` (phím + EventType + callback nhận tên phím) |
| `KeyBinding` | Liên kết keys + EventType + `func(key string)` handler |
| `MouseBinding` | Liên kết MouseButton + EventType + `func(button string)` handler |
| `InforData` | ID, Name, SaveTag |
| `DrawData` | Thông tin màn hình và camera nội bộ dành cho `DrawComponent` |
| `DebugData` | Cấu hình bật/tắt vẽ hitbox, tâm tọa độ, info và string log |

### `ISave.go` — Interface của hệ thống Save/Load

| Interface / Struct | Mô tả |
|--------------------|-------|
| `ISaveManager` | Interface quản lý Save/Load: `SaveGame`, `LoadGame`, `ReadSnapshot`... |
| `SaveSnapshot` | Cấu trúc dữ liệu để lưu vào file JSON |

### `IInput.go` — Định nghĩa Input System

- **`Key`**: Type và toàn bộ hằng số phím bấm (độc lập với ebiten).
- **`MouseButton`**: Các nút chuột (Left, Right, Middle).
- **`EventType`**: Loại sự kiện cần bắt — `EventPressed` (giữ), `EventJustPressed` (vừa nhấn), `EventJustReleased` (vừa thả).
- **`EventTypeNameMap`**: Map từ chuỗi `""`, `"pressed"`, `"released"` sang hằng số `EventType`.
- **`IInputManager`**: Interface polling input (Pressed, JustPressed, JustReleased), chuột, wheel và Virtual Action.
- **`KeyNameMap`**: Map từ chuỗi dễ nhớ (`"w"`, `"space"`) sang hằng số `Key`.
- **`KeyReverseMap`**: Map ngược từ `Key` hằng số → tên chuỗi. Dùng để truyền tên phím vào handler.
- **`KeyGroupMap`**: Map nhóm phím đặc biệt (`"alpha"`, `"number"`, `"arrows"`, `"wasd"`, `"all"`).
- **`MouseButtonNameMap`**: Map tên chuỗi → `MouseButton` hằng số (`"left"`, `"right"`, `"middle"`).
- **`MouseButtonReverseMap`**: Map ngược `MouseButton` → tên chuỗi.

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
