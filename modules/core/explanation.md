# AutoWorld — Core Module Reference

**Ý nghĩa**:
    Module `core` đóng vai trò là **"Trái tim" của Engine**, chịu trách nhiệm khởi tạo, kết nối và điều phối tất cả các module con (nsprite, naudio, nobject, nsys…).
    Đây là module **duy nhất** mà user-code (main.go, scene scripts) tương tác trực tiếp — thay vì import từng sub-module riêng lẻ.

---

## Các file chính

### `Register.go` — Type Alias Gateway
Cửa ngõ tập trung: re-export toàn bộ interface và type từ `domain` và `enginetype`.
Nhờ đó, user-code chỉ cần `import "autoworld/modules/core"` là đủ.

| Nhóm | Nội dung |
|------|----------|
| Manager interfaces | ISceneManager, IGlobalConfig, IInputManager, IObserver |
| System interfaces | IScene, ILogicSystem, IDrawSystem, IAudioSystem, IInputSystem |
| Object | IObject |
| Data structs | PositionData, SpriteData, BoxData, AudioData, BackgroundData, TilemapData |
| Asset interfaces | ISpriteLW, IAudioLW, IGlobal, ISpriteLoader, IAudioLoader, IManifestLoader |
| Misc | BoxShape, BSRectangle, BSCircle |
| ECS Component Types | Position, Sprite, Box, Audio, Input, Background, Tilemap (từ enginetype) |

### `Root.go` — Engine (Root)
Struct trung tâm `Engine` chứa toàn bộ manager:

| Field | Kiểu | Mô tả |
|-------|------|-------|
| Scene | ISceneManager | Quản lý vòng đời các Scene |
| Config | IGlobalConfig | Cấu hình toàn cục (Singleton + Observer) |
| Input | IInputManager | Xử lý input bàn phím/chuột |
| Store | IGlobal | Global resource store |
| AudioCtx | *audio.Context | Ebiten audio context dùng chung |

- `NewGame(cfg GameConfig)` — khởi tạo toàn bộ Engine từ 1 config struct.
- `Start()` — set cửa sổ và chạy `ebiten.RunGame()` (block đến khi game đóng).

### `EbitenGame.go` — Adapter (Ebiten ↔ Engine)
Implement `ebiten.Game`, chuyển vòng đời Ebitengine vào Engine:
- `Update()` — gọi `Input.Update()` rồi `Scene.Update()`
- `Draw(screen)` — type-assert sang `*Scene` để set screen cho DrawSystem
- `Layout()` — ủy quyền cho SceneManager

> **Lý do type-assert**: `IScene` không expose `SetScreen` (tách biệt concern vẽ khỏi logic). Chỉ `EbitenGame` biết về `*ebiten.Image`, nên chỉ ở đây mới cần cast.

### `Scene.go` — Scene
Quản lý một màn chơi độc lập, sở hữu `donburi.World` và 4 System riêng:

| System | Nguồn | Trách nhiệm |
|--------|-------|-------------|
| logicSystem | nobject | Gọi OnStep mỗi frame |
| drawSystem | nsprite | Render entity lên screen |
| audioSystem | naudio | Phát/dừng âm thanh |
| inputSystem | nobject | Kích hoạt KeyBinding callback |

Thứ tự Update mỗi frame: **Input → Logic → Audio**.

### `Map.go` — Map (Physical/GUI)
Quản lý ECS World và vòng lặp update cho một bản đồ:

- `AddObject(obj)` — đăng ký object vào update loop; auto-register IDraw vào DrawRegistry.
- `RemoveObject(obj)` — xóa object theo cơ chế **deferred** (cuối frame):
  - Gọi `MarkDead()` ngay lập tức để Collector/Applier bỏ qua object này.
  - Đưa vào `pendingRemove` queue.
  - Cuối frame (`flushRemove`): cắt khỏi `objectList`, xóa khỏi `DrawRegistry`, gọi `AddObjectDestroy` để `OnDestroy()` chạy ở frame tiếp theo.
- **Lý do deferred**: tránh race condition khi đang duyệt objectList giữa frame.

### `SceneManager.go` — SceneManager
Quản lý stack/list Scene, xử lý chuyển cảnh:
- `AddScene` / `RemoveScene` / `RemoveAllScene` — quản lý danh sách
- `ChangeSceneFromList` — chuyển scene, giữ scene cũ (có thể quay lại)
- `ChangeSceneForce` — destroy scene cũ, thay bằng scene mới
- `GetGlobalScene` — trả về hidden scene luôn chạy (persistent objects)

### `InputManager.go` — InputManager
Bridge duy nhất giữa `domain.IInputManager` (trừu tượng) và Ebitengine input API:
- **ebitenKeyMap** — ánh xạ 100+ domain.Key → ebiten.Key
- **ebitenMouseMap** — ánh xạ 3 nút chuột
- Hỗ trợ: Keyboard, Mouse, Wheel, Virtual Action Mapping

---

## Quy trình hoạt động

```
Update Flow:
  ebiten → EbitenGame.Update()
              ├─ InputManager.Update()
              └─ SceneManager.Update()
                    ├─ globalScene.Update() [luôn chạy]
                    └─ currentScene.Update()
                           ├─ InputSystem
                           ├─ LogicSystem (OnCreate → OnStep → OnDestroy)
                           ├─ AudioSystem
                           └─ Map.flushRemove() ← Deferred Remove

Draw Flow:
  ebiten → EbitenGame.Draw(screen)
              ├─ scene.drawSystem.SetScreen(screen)
              └─ currentScene.Draw()
                    └─ DrawSystem.Draw(world) [nsprite]
```

---

## Quy tắc quan trọng cho Agent

1. **Interface First** — Luôn dùng interface từ `domain/` hoặc alias trong `Register.go`.
2. **Engine Access** — Đi qua `Engine` struct để truy cập các Manager toàn cục.
3. **ECS Pattern** — Core chỉ điều phối System; dữ liệu thực tế nằm trong Component (donburi).
4. **Không thêm logic vào Register.go** — Chỉ type alias và var re-export.
