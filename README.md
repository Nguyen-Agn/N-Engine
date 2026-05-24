# N-Engine (AutoWorld)

> A lightweight 2D Game Engine built in Go, powered by **Ebitengine** for rendering/audio/input and **Donburi** for ECS (Entity Component System).

N-Engine focuses on developer experience: game code stays clean and OOP-style, while the engine handles ECS data management, rendering, physics, and audio behind the scenes.

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **ECS + OOP Hybrid** | Write game objects as plain Go structs. Engine manages ECS internals via reflection + mixins |
| **Component Mixins** | Embed `napi.IPosition`, `napi.ISprite`, `napi.IVelocity`, etc. to gain capabilities |
| **Custom Components** | Define your own components with `napi.NewComponentType[T]` and `GenericComponent[T]` |
| **Scene Management** | Multi-scene system with Physical Map (world space) + GUI Map (screen space) |
| **Camera** | Viewport with follow-target and map-bounds clamping |
| **Tilemap** | 2D grid-based tile rendering with culling |
| **Background** | Scrolling/repeating/stretching backgrounds |
| **Input System** | Event-driven keyboard input with named key groups |
| **Alarm System** | Trigger callbacks after N frames (like GameMaker's Alarms) |
| **Velocity/Physics** | Per-object velocity, friction, and max-speed, auto-applied to position |
| **Tween System** | Smooth lerp transitions for position, scale, and alpha |
| **Asset Manifest** | Load all sprites & audio from a single `.toml` file |

---

## 🏗 Architecture

The engine is organized in strict dependency layers. Higher layers use lower layers **only through interfaces** defined in `domain`.

```
┌──────────────────────────────────────┐
│          [ Game Code ]               │  ← Developer writes here
└─────────────────┬────────────────────┘
                  │ import (only napi)
┌─────────────────▼────────────────────┐
│              [ napi ]                │  ← Public API layer (only entry point)
└──┬──────────┬──────────┬────────────┘
   │          │          │
┌──▼──┐  ┌───▼──┐  ┌────▼──────┐
│core │  │nasset│  │ nsys/nsys │   ← Engine subsystems
└──┬──┘  └──────┘  └───────────┘
   │
┌──▼──────────────────────────────────┐
│  domain / enginetype / nsystem      │  ← Interfaces, ECS types, systems
└─────────────────────────────────────┘
```

### Modules

| Module | Role |
|--------|------|
| `domain` | Pure interfaces & data structs. Zero logic. Shared contract between all modules |
| `enginetype` | Registers all ECS `ComponentType` tokens globally |
| `components` | Built-in Component Mixins (Position, Sprite, Box, Audio, Alarm, Velocity, Tween...) |
| `nsystem` | All ECS Systems (LogicSystem, InputSystem, DrawSystem, AlarmSystem, PhysicsSystem, TweenSystem, AudioSystem) |
| `core` | Engine heart: Scene, Map, Camera, SceneManager, EbitenGame loop |
| `nsprite` | Sprite/image data structures and rendering helpers |
| `naudio` | Audio lightweight wrappers |
| `nasset` | Asset loading (images, audio) from manifest files |
| `nsys` | Global resource store (sprites, audio, objects, vars, constants) |
| `nlayout` | Flexbox-style UI layout system |
| `nobject` | Lightweight object wrapper |
| `napi` | **The only module game developers import.** Re-exports everything needed |

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+
- Dependencies are vendored in `.libs/` — no additional `go get` needed

### Running the demo

```powershell
cd tests/simulation
go run .\TilemapDemo.go
```

---

## 📦 Core Concepts

### 1. Creating a Game Object

A game object is a plain Go struct embedding the capabilities it needs:

```go
import "autoworld/modules/napi"

type Player struct {
    napi.IObject    // Required: lifecycle (Create, StepUpdate, Destroy)
    napi.IPosition  // Gives: X(), Y(), SetX(), SetY()
    napi.ISprite    // Gives: SetSprite(), SetCurrentSprite(), ScaleX()...
    napi.IInput     // Gives: ListenOn(key, handler)
    napi.IVelocity  // Gives: SetVelocity(), VelocityX(), Friction()...

    // Your own fields
    Health int
}

func NewPlayer(x, y float32) *Player {
    p := &Player{Health: 100}

    // 1. Initialize ECS entity and inject component mixins
    //    Component tokens: "pos" "spr" "inp" "vel" → binds Position, Sprite, Input, Velocity
    napi.NewObject(p, "player", "pos spr inp vel")

    // 2. Configure object
    p.SetX(x)
    p.SetY(y)
    p.ListenOn("space", p.OnJump)

    // 3. Register into the current scene's map
    napi.Register(p, "main")
    return p
}

// Called once when the object enters the scene
func (p *Player) Create() {
    p.SetSprite("idle", napi.GetSprite("hero_idle"))
    p.SetCurrentSprite("idle")
    p.SetVelocityX(2.0)
    p.SetFriction(0.1)
}

// Called every frame
func (p *Player) StepUpdate() {
    // Velocity is automatically applied to Position by PhysicsSystem
    // Just write game logic here
    if p.X() > 640 {
        p.SetVelocityX(-p.VelocityX())
    }
}

// Called when the object is removed from the scene
func (p *Player) Destroy() {}

// Custom handler — called by InputSystem when Space is held
func (p *Player) OnJump() {
    p.SetVelocityY(-5.0)
}
```

---

### 2. Component Tokens

When calling `napi.NewObject(obj, name, "tokens...")`, use space-separated tokens:

| Token | Component | Capabilities |
|-------|-----------|-------------|
| `pos` | PositionData | `X()`, `Y()`, `SetX()`, `SetY()` |
| `spr` | SpriteData | `SetSprite()`, `SetCurrentSprite()`, `ScaleX/Y()`, `Rotation()`, `ImageColor()`, `IsVisible()` |
| `box` | BoxData | `BoxW/H()`, `SetBoxW/H()`, `IsCollidable()`, `Shape()` |
| `aud` | AudioData | `Play()`, `PlayDefault()`, `StopAudio()`, `Volume()`, `Pitch()` |
| `dir` | DirectionData | `Direction()`, `SetDirection()`, `Rotate()` |
| `inp` | InputData | `ListenOn(key, handler)` |
| `bg` | BackgroundData | `SetSprite()`, `SetColor()`, `RepeatX/Y()`, `ScrollSpeed()` |
| `til` | TilemapData | `SetSprite()`, `SetGrid()`, `SetTile()`, `TileWidth/Height()` |
| `alr` | AlarmData | `SetAlarm(id, frames, callback)`, `GetAlarm(id)`, `StopAlarm(id)` |
| `vel` | VelocityData | `SetVelocity()`, `VelocityX/Y()`, `Friction()`, `MaxSpeed()` |
| `twn` | TweenData | `TweenMove()`, `TweenScale()`, `TweenAlpha()` |
| `inf` | InforData | `GetName()`, `GetId()` — **added automatically to every object** |

---

### 3. Game Initialization

```go
package main

import (
    "log"
    "autoworld/modules/napi"
)

func main() {
    // 1. Configure the engine
    napi.Init(napi.GameConfig{
        Title:  "My Game",
        Width:  640,
        Height: 480,
    })

    // 2. Load assets from a manifest file
    napi.LoadFromFile("./assets/manifest.toml")

    // 3. Create a scene
    //    Format: "map-WxH" creates a physical map of W×H pixels
    //    Other options: "gui", "map-gui-WxH"
    _, err := napi.NewSceneAndGo("main", "map-640-480")
    if err != nil {
        log.Fatal(err)
    }

    // 4. Spawn your objects
    NewPlayer(100, 200)
    NewBackground()

    // 5. Start the game loop (blocks until window is closed)
    napi.GameStart()
}
```

---

### 4. Asset Manifest (TOML)

Create a `.toml` file to list all assets:

```toml
[[sprites]]
key  = "hero_idle"
path = "assets/hero_idle.png"
cols = 4       # How many columns (frames) in the sprite sheet
rows = 1

[[sprites]]
key  = "tileset"
path = "assets/tileset.png"
cols = 8
rows = 4

[[audios]]
key  = "bgm_main"
path = "assets/bgm.wav"

[[audios]]
key  = "sfx_jump"
path = "assets/jump.wav"
```

Then load it:
```go
napi.LoadFromFile("./assets/manifest.toml")

sprite := napi.GetSprite("hero_idle") // Returns ISpriteLW
audio  := napi.GetAudio("bgm_main")   // Returns IAudioLW
```

---

### 5. Built-in Systems (Frame Order)

Each frame, the `Map.Update()` runs systems in this exact order:

```
LogicSystem     → Runs Create() for new objects, then StepUpdate() for all, then Destroy()
InputSystem     → Checks key state and calls ListenOn handlers
AlarmSystem     → Decrements frame counters, fires callbacks at 0
TweenSystem     → Advances lerp transitions (move, scale, alpha)
PhysicsSystem   → Applies friction, clamps MaxSpeed, adds velocity to position
AudioSystem     → Processes Play/Stop audio commands
```

---

### 6. Scene System

```go
// Create and immediately activate a scene
scene, _ := napi.NewSceneAndGo("game", "map-1280-720")

// Navigate to a pre-existing scene (pauses current, doesn't destroy)
napi.GoToScene("menu")

// Replace current scene (destroys current)
napi.ReplaceScene(newScene)

// Access current scene
scene := napi.GetCurrentScene()
```

#### Scene Layout Strings

| String | Creates |
|--------|---------|
| `"map-WxH"` | Physical map W×H pixels |
| `"gui"` | GUI-only map (screen space, no camera) |
| `"map-640-480"` | Physical map 640×480 |

---

### 7. Tilemap

```go
type TilemapObject struct {
    napi.IObject
    napi.ITilemap
    napi.IPosition
}

func NewTilemapObject() *TilemapObject {
    tm := &TilemapObject{}
    napi.NewObject(tm, "tilemap", "pos til")

    tm.SetSprite(napi.GetSprite("tileset"))
    tm.SetTileWidth(32)
    tm.SetTileHeight(32)

    // 2D grid: -1 = empty, 0..N = tile index from sprite sheet
    tm.SetGrid([][]int{
        {-1, -1, -1, -1, -1},
        { 0,  1,  2,  3,  4},
        { 5,  6,  7,  8,  9},
    })

    tm.SetX(0)
    tm.SetY(200)

    napi.Register(tm, "main")
    return tm
}
```

---

### 8. Alarm Component

```go
// Set alarm "fire-bullet" to trigger after 60 frames (~1 second at 60fps)
p.SetAlarm("fire-bullet", 60, func() {
    // This runs exactly once after 60 frames
    SpawnBullet(p.X(), p.Y())

    // Re-arm the alarm for the next shot
    p.SetAlarm("fire-bullet", 60, ...)
})

// Check remaining frames
framesLeft := p.GetAlarm("fire-bullet")

// Cancel it
p.StopAlarm("fire-bullet")
```

---

### 9. Tween Component

```go
// Smoothly move to (400, 300) over 60 frames
p.TweenMove(400, 300, 60)

// Smoothly scale to (2x, 2x) over 30 frames
p.TweenScale(2.0, 2.0, 30)

// Fade out (alpha → 0) over 90 frames
p.TweenAlpha(0, 90)
```

---

### 10. Custom Components

Define your own data and component type:

```go
// 1. Define your data struct
type StatsData struct {
    Health int
    Mana   int
    Level  int
}

// 2. Register a component type (use a unique 3-letter token)
var StatsComp = napi.NewComponentType[StatsData]("sta")

// 3. Create a mixin with custom methods by wrapping GenericComponent
type StatsComponent struct {
    napi.GenericComponent[StatsData]
}

func (s *StatsComponent) TakeDamage(amount int) {
    if s.Get() != nil {
        s.Get().Health -= amount
    }
}

func (s *StatsComponent) IsAlive() bool {
    return s.Get() != nil && s.Get().Health > 0
}

// 4. Embed in your object
type Hero struct {
    napi.IObject
    napi.IPosition
    StatsComponent
}

func NewHero() *Hero {
    h := &Hero{
        StatsComponent: StatsComponent{
            GenericComponent: napi.NewGenericComponent(StatsComp),
        },
    }
    napi.NewObject(h, "hero", "pos sta")
    napi.Register(h, "main")
    return h
}

func (h *Hero) StepUpdate() {
    h.TakeDamage(1)          // method promoted from StatsComponent
    if !h.IsAlive() {
        h.Destroy()
    }
}
```

---

## 📁 Project Structure

```
AutoWorld/
├── domain/                  # Interfaces and data structs (no logic)
│   ├── ObjectInterface.go   # IObject, IPosition, ISprite, IBox, IAudio, IAlarm, IVelocity, ITween...
│   ├── ObjectData.go        # PositionData, SpriteData, AlarmData, VelocityData, TweenData...
│   ├── System.go            # ILogicSystem, IDrawSystem, IAudioSystem, IUpdateSystem...
│   └── ...
│
├── modules/
│   ├── napi/                # ← Game developers ONLY import this
│   │   ├── Game.go          # Init, GameStart, LoadFromFile
│   │   ├── Scene.go         # NewSceneAndGo, GoToScene, ReplaceScene
│   │   ├── ObjectHelper.go  # NewObject, Register
│   │   ├── StoreHelper.go   # GetSprite, GetAudio
│   │   ├── Config.go        # VarInt, VarString, VarBool
│   │   └── Register.go      # Type aliases: IPosition, ISprite, AlarmComponent...
│   │
│   ├── components/          # Component Mixin implementations
│   │   ├── PositionComponent.go
│   │   ├── SpriteComponent.go
│   │   ├── AlarmComponent.go
│   │   ├── VelocityComponent.go
│   │   ├── TweenComponent.go
│   │   └── GenericComponent.go  # For custom components
│   │
│   ├── nsystem/             # ECS Systems
│   │   ├── LogicSystem.go   # Create/StepUpdate/Destroy lifecycle
│   │   ├── InputSystem.go   # Keyboard input dispatch
│   │   ├── DrawSystem.go    # Sprite/Tilemap/Background rendering
│   │   ├── AlarmSystem.go   # Frame countdown and callback firing
│   │   ├── PhysicsSystem.go # Velocity, friction, position integration
│   │   ├── TweenSystem.go   # Lerp transitions
│   │   └── AudioSystem.go   # Audio play/stop command processing
│   │
│   ├── core/                # Engine internals
│   │   ├── EbitenGame.go    # Ebitengine integration (Update/Draw loop)
│   │   ├── SceneManager.go  # Scene registry and navigation
│   │   ├── Scene.go         # Scene with Physical Map + GUI Map + Camera
│   │   ├── Map.go           # ECS World + objectList + all systems
│   │   ├── Camera.go        # Viewport, follow-target, bounds clamping
│   │   └── InputManager.go  # Raw key state from Ebiten
│   │
│   ├── enginetype/          # ECS ComponentType registration
│   ├── nasset/              # Asset loading (image, audio, manifest)
│   ├── nsys/                # Global resource store (singleton)
│   ├── naudio/              # Audio lightweight wrapper
│   ├── nsprite/             # Sprite data and rendering helpers
│   └── nlayout/             # Flexbox-style UI layout
│
├── tests/
│   └── simulation/
│       └── TilemapDemo.go   # Working demo: background + tilemap + bouncing player
│
└── domain/explanation.md    # Architecture notes (read before modifying)
```

---

## 🔧 Dependencies

| Library | Purpose | Location |
|---------|---------|----------|
| [Ebitengine v2](https://ebitengine.org/) | Rendering, audio, input | `.libs/github.com/hajimehoshi/ebiten/v2` |
| [Donburi](https://github.com/yohamta/donburi) | ECS (Entity Component System) | `.libs/github.com/yohamta/donburi` |
| [BurntSushi/toml](https://github.com/BurntSushi/toml) | Asset manifest parsing | via `go.sum` |

Dependencies are vendored locally in `.libs/` — the project does **not** require network access to build.

---

## 🧩 Design Principles

1. **Interface First** — Modules never talk directly to each other. All communication goes through interfaces defined in `domain`.
2. **Open/Closed Principle** — Extend via new files/types; avoid modifying existing stable modules.
3. **Single Entry Point** — Game code only ever imports `napi`. All engine complexity is hidden behind it.
4. **ECS + OOP Hybrid** — Data lives in ECS (Donburi); behavior is expressed through Go struct embedding and method promotion.
5. **`explanation.md`-first** — Every module has a `explanation.md`. Read it before touching source code.

---

## 📄 License

This project is developed by **Nguyen-Agn**. See repository at [github.com/Nguyen-Agn/N-Engine](https://github.com/Nguyen-Agn/N-Engine).
