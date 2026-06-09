# Module napi — Game Logic API Layer

## Mục tiêu

Module `napi` là **API layer nằm trên cùng** của engine stack.
Nó che giấu sự phức tạp của ECS (donburi), asset loading và audio, giúp game developer viết game logic mà không cần biết chi tiết bên trong engine.

> **Nguyên tắc**: `napi` chỉ gọi các module bên dưới qua interface. Không module engine nào được import ngược lại vào `napi`.

---

## Cấu trúc file

| File | Mô tả |
|------|-------|
| `Game.go` | Singleton Engine — `Init`, `GetEngine`, `LoadFromFile`, `GameStart`, **`SaveGame`, `LoadGame`, ...** |
| `Register.go` | Type alias gateway + `NewGame` helper (các aliases legacy giữ lại cho tham chiếu nội bộ) |
| `Scene.go` | Scene management — `NewScene`, `GoToScene`, `ReplaceScene`... |
| `ObjectHelper.go` | Object factory — `NewObject`, `Register`, `Remove`, `bind` (internal) |
| `StoreHelper.go` | Asset & audio helpers — `LoadManifest`, `GetSprite`, `Play`... |
| `Config.go` | Global config accessors — `VarInt`, `VarString`, `VarBool`... |
| `Component.go` | Generic component helpers — `SetComponent`, `GetComponent` |
| `ComponentType.go` | Custom component creation — `NewComponentType`, `NewEntry` |
| `ncom/Components.go` | **Sub-package** chứa toàn bộ Component Mixin type aliases — `Pos`, `Spr`, `Drw`... |

---

## Khởi động (main.go)

```go
napi.Init(napi.GameConfig{Title: "MyGame", Width: 640, Height: 480})
// hoặc dùng LoadFromFile thay vì LoadManifest nếu cần config audio riêng:
napi.LoadManifest("assets/manifest.json")

scene, _ := napi.NewSceneAndGo("main")
// ... setup scene, tạo object ...

napi.GameStart() // block đến khi đóng cửa sổ
```

---

## API tham khảo

### Game Singleton — `Game.go`

| Hàm | Mô tả |
|-----|-------|
| `Init(cfg)` | Khởi tạo Engine và đăng ký singleton — PHẢI gọi đầu tiên |
| `GetEngine()` | Lấy *Engine singleton hiện tại |
| `LoadFromFile(path)` | Load manifest dùng AudioCtx của engine (SampleRate từ config) |
| `GameStart()` | Cấu hình cửa sổ và chạy vòng lặp game (block) |

---

### Scene Management — `Scene.go`

| Hàm | Mô tả |
|-----|-------|
| `NewScene(id)` | Tạo Scene mới, đăng ký với id (chưa active) |
| `NewSceneAndGo(id)` | Tạo Scene mới và activate ngay |
| `AddScene(id, scene)` | Đăng ký IScene đã tạo sẵn |
| `GoToScene(id)` | Chuyển sang Scene đã có (scene cũ pause, không destroy) |
| `ReplaceScene(scene)` | Ép thay thế Scene hiện tại (scene cũ bị Destroy) |
| `RemoveScene(id)` | Xóa Scene khỏi danh sách, gọi Destroy |
| `RemoveAllScenes()` | Xóa toàn bộ Scene, currentScene = nil |
| `GetCurrentScene()` | Lấy Scene đang chạy |
| `GetSceneByID(id)` | Lấy Scene theo id từ danh sách |
| `GetGlobalScene()` | Lấy Global Hidden Scene (luôn update mọi frame) |

---

### Object System — `ObjectHelper.go`

#### Component tokens

| Token | Component | Nội dung |
|-------|-----------|---------|
| `pos` | PositionData | X, Y |
| `spr` | SpriteData | Sprite map, Scale, Color, Visibility |
| `box` | BoxData | Hitbox W/H, IsCollidable, Shape |
| `aud` | AudioData | Audio map, Volume, Pitch |
| `dir` | DirectionData | Góc hướng di chuyển |
| `"twn"` | `ncom.Twn` / `TweenComponent` | Hoạt ảnh thay đổi thông số |
| `"col"` | `ncom.Col` / `CollisionComponent` | Logic va chạm |
| `"drw"` | `ncom.Drw` / `DrawComponent` | Tự định nghĩa hàm Draw() |
| `"deb"` | `ncom.Deb` / `DebugComponent` | Hiển thị hitbox, log, origin |
| `inp` | InputData | Keyboard bindings (ListenOn + EventType) |
| `glo` | *(modifier)* | Tạo / đăng ký vào Global Scene |
| `inf` | InforData | ID + Name — **tự động thêm vào mọi object** |

#### Custom Object (có type riêng) — dùng `ncom`

Dev import `autoworld/modules/napi/ncom` để nhúng các Component Mixin:

```go
import (
    "autoworld/modules/napi"
    "autoworld/modules/napi/ncom"
)

type Player struct {
    napi.IObject   // lifecycle: OnCreate, OnStep, OnDestroy, OnSave, OnLoad
    ncom.Pos       // X(), Y(), SetX(), SetY()
    ncom.Spr       // AddSprite(), SetCurrentSprite()...
    ncom.Info      // GetId(), SaveTag(), SetSaveTag(), IsDead()...
    ncom.Inp       // ListenOn(key, eventType, handler)
    ncom.Mouse     // MouseX/Y, WheelX/Y, ListenMouseOn(btn, eventType, handler)

    Health int
    Speed  float32
}

func NewPlayer(x, y float32) *Player {
    p := &Player{Health: 100, Speed: 3.0}
    napi.Obj.NewObject(p, "player", "pos spr inf inp sce-cur")
    // Đăng ký input (token "inp" cần thiết cho keyboard)
    p.ListenOn("wasd", "pressed", func(key string) { /* di chuyển */ })
    p.ListenOn("space", "just_pressed", func(key string) { /* nhảy */ })
    // Mouse không cần token ECS
    p.ListenMouseOn("left", "just_pressed", func(btn string) { /* bắn */ })
    return p
}
```

#### Register helpers

| Hàm | Mô tả |
|-----|-------|
| `Obj.Register(obj, scene)` | Đăng ký vào scene theo tên ("" = current) |
| `Obj.RegisterIn(scene, obj)` | Đăng ký vào scene cụ thể |
| `Obj.Remove(obj)` | Xóa object khỏi scene hiện tại (deferred, cuối frame). Gọi `MarkDead()` ngay; `OnDestroy()` ở frame tiếp. |

> **Lưu ý**: `Remove` chỉ hoạt động với Current Scene. Với object trong Global Scene, gọi trực tiếp trên `IMap` của Global Scene.

---

### Asset & Audio — `StoreHelper.go`

```go
napi.LoadManifest("assets/manifest.json") // load JSON, đưa sprite+audio vào store

sprite := napi.GetSprite("hero_idle")  // lấy ISpriteLW
audio  := napi.GetAudio("bgm_main")    // lấy IAudioLW
obj    := napi.GetObject("player")     // lấy IObject đã đặt tên
const_ := napi.GetConst("gravity")    // lấy hằng số

napi.Play("bgm_main")                 // phát với vol=1.0, pitch=1.0
napi.PlayAt("sfx_jump", 0.8, 1.0)    // phát với tùy chỉnh
napi.Stop("bgm_main")                 // dừng âm thanh
```

### Save/Load System — `napi.Game`

Save/Load được gộm vào `napi.Game` cùng với các hàm vòng đời Engine:

| Hàm | Mô tả |
|-----|-------|
| `napi.Game.SaveGame(path)` | Lưu toàn bộ state game vào file (path đầy đủ hoặc "" = "default"). |
| `napi.Game.LoadGame(path)` | Nạp lại state từ file. |
| `napi.Game.HasSave(path)` | Trả về `true` nếu file tồn tại. |
| `napi.Game.DeleteSave(path)` | Xóa file save theo path. |
| `napi.Game.ListSaveSlots()` | Trả về mảng `[]string` chứa tất cả file save trong `SaveDir`. |
| `napi.Game.ReadSaveSnapshot(path)` | Trả về data (kèm `CurrentSceneID`) mà không tự động load. |

---

### Global Config — `Config.go`

Đọc giá trị từ `IGlobalConfig` (key-value store có kiểu):

```go
w := napi.VarInt("game-width")        // trả về 0 nếu không có
t := napi.VarString("game-title")     // trả về "" nếu không có
```

---

### Custom Component — `Component.go` & `ComponentType.go`

```go
// Khai báo custom component (1 lần, ở cấp package)
var StatsComp = napi.NewComponentType[StatsData]("sta")

// Tạo object với custom component
napi.NewObject(hero, "hero", "pos spr sta")

// Đọc/ghi data
napi.SetComponent(hero, StatsComp, StatsData{Health: 100})
stats := napi.GetComponent(hero, StatsComp) // *StatsData hoặc nil

// Thêm component sau khi object đã tạo
napi.AddComponentType(hero, StatsComp)
```

---

## Dependency graph

```
     [ game code ]
           ↓ import
       [ napi ]
           ↓ import
[ core ] [ nasset ] [ nobject ] [ nsys ]
    ↓         ↓          ↓
  [ domain / enginetype / naudio / nsprite ]
```

`napi` **KHÔNG được import** bởi bất kỳ module engine nào (tránh circular dependency).

---

## Quy tắc quan trọng cho Agent

1. **Interface First** — Mọi tương tác qua interface từ domain. Không dùng struct cụ thể của module khác.
2. **Singleton Guard** — Luôn check `engine()` trước mọi thao tác (panic nếu chưa Init).
3. **Không thêm logic vào Register.go** — Chỉ type alias và re-export.
4. **NewObject → Register** — Tạo object và đăng ký là 2 bước riêng biệt có chủ đích.

---

## 5. Object Pooling & Quản lý bộ nhớ

N-Engine trang bị sẵn một hệ thống `ObjectPool` **Zero-Allocation** dùng để tái sử dụng object và chống Memory Leak ECS.

### Nguyên lý Xóa Object (`napi.Obj.Remove`)
Hàm `napi.Obj.Remove(obj)` hoạt động cực kỳ thông minh:
- **Nếu Object KHÔNG có Pool (Ví dụ: Boss, Player):** Gọi `OnDestroy()` và xóa sổ hoàn toàn Entity khỏi bộ nhớ Donburi (Chống rò rỉ RAM ECS).
- **Nếu Object CÓ Pool (Ví dụ: Bullet):** Bỏ qua `OnDestroy()`, giữ nguyên Entity, tắt logic và chuyển về cất vào `ObjectPool` cấu hình sẵn.

### Cách sử dụng ObjectPool
```go
// 1. Khai báo Pool toàn cục
var BulletPool *napi.ObjectPool[*Bullet]

func init() {
    BulletPool = napi.NewObjectPool(napi.PoolConfig[*Bullet]{
        New: func() *Bullet {
            b := &Bullet{}
            napi.Obj.NewObject(b, "bullet", "pos drw inf")
            return b
        },
        Reset: func(b *Bullet) {
            b.Life = 100 // Phục hồi máu khi rút từ Pool ra
        },
        MaxSize: 100, // Cất tối đa 100 viên. Chết viên thứ 101 sẽ bị Engine xóa thẳng tay.
    })
}

// 2. Tạo hoặc Lấy đạn từ Pool
b := BulletPool.Get("tên_scene_nếu_có_hoặc_chuỗi_rỗng")
b.SetX(10)

// 3. Xóa đạn (Engine tự động Cất vào Pool)
napi.Obj.Remove(b)
```
