# Module napi — Game Logic API Layer

## Mục tiêu

Module `napi` là **API layer nằm trên cùng** của engine stack.
Nó che giấu sự phức tạp của ECS (donburi), asset loading và audio, giúp game developer viết game logic mà không cần biết chi tiết bên trong engine.

> **Nguyên tắc**: `napi` chỉ gọi các module bên dưới qua interface. Không module engine nào được import ngược lại vào `napi`.

---

## Cấu trúc file

| File | Mô tả |
|------|-------|
| `Game.go` | Singleton Engine — `Init`, `GetEngine`, `LoadFromFile`, `GameStart` |
| `Register.go` | Type alias gateway + `NewGame` helper |
| `Scene.go` | Scene management — `NewScene`, `GoToScene`, `ReplaceScene`... |
| `ObjectHelper.go` | Object factory — `NewObject`, `Register`, `bind` (internal) |
| `StoreHelper.go` | Asset & audio helpers — `LoadManifest`, `GetSprite`, `Play`... |
| `Config.go` | Global config accessors — `VarInt`, `VarString`, `VarBool`... |
| `Component.go` | Generic component helpers — `SetComponent`, `GetComponent` |
| `ComponentType.go` | Custom component creation — `NewComponentType`, `NewEntry` |

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
| `glo` | *(modifier)* | Tạo / đăng ký vào Global Scene |
| `inf` | InforData | ID + Name — **tự động thêm vào mọi object** |

#### Custom Object (có type riêng)

```go
// player/player.go
type Player struct {
    napi.IObject           // lifecycle: Create, StepUpdate, Destroy, Entry
    napi.IPosition         // getter/setter: X, Y, SetX, SetY
    napi.ISprite           // getter/setter: Sprite, Scale, Color...

    Health int
    Speed  float32
}

func NewPlayer(x, y float32) *Player {
    p := &Player{Health: 100, Speed: 3.0}

    // Tạo ECS entry + inject component mixin vào p
    napi.NewObject(p, "player", "pos spr")

    // Dùng getter/setter ngay sau NewObject
    p.SetX(x)
    p.SetY(y)
    p.SetSprite("idle", napi.GetSprite("hero_idle"))

    // Đăng ký vào current scene
    napi.Register(p, false)
    return p
}

func (p *Player) StepUpdate() {
    p.SetX(p.X() + p.Speed)
}
```

#### Register helpers

| Hàm | Mô tả |
|-----|-------|
| `Register(obj, global)` | Đăng ký vào current scene (false) hoặc Global Scene (true) |
| `RegisterIn(scene, obj)` | Đăng ký vào scene cụ thể |
| `RegisterGlobal(obj)` | Shortcut: đăng ký vào Global Scene |

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
