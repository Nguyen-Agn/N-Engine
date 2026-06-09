# N-Engine (AutoWorld)

> A lightweight 2D Game Engine built in Go, powered by **Ebitengine** for rendering/audio/input and **Donburi** for ECS (Entity Component System).

N-Engine focuses on developer experience: game code stays clean and OOP-style, while the engine handles ECS data management, rendering, physics, and audio behind the scenes.

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **ECS + OOP Hybrid** | Write game objects as plain Go structs. Engine manages ECS internals via reflection + mixins |
| **Component Mixins** | Embed `ncom.Pos`, `ncom.Spr`, `ncom.Velo`, etc. to gain capabilities |
| **Custom Components** | Define your own components with `napi.NewComponentType[T]` and `ncom.Generic[T]` |
| **Scene Management** | Multi-scene system with Physical Map (world space) + GUI Map (screen space) |
| **Camera** | Viewport with follow-target and map-bounds clamping |
| **Tilemap** | 2D grid-based tile rendering with culling |
| **Background** | Scrolling/repeating/stretching backgrounds |
| **Input System** | Event-driven keyboard & mouse input with named callbacks |
| **Physics & Alarms** | Auto-applied velocity/friction and frame-based alarm triggers |
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
│core │  │nasset│  │ nsys      │   ← Engine subsystems
└──┬──┘  └──────┘  └───────────┘
   │
┌──▼──────────────────────────────────┐
│  domain / enginetype / nsystem      │  ← Interfaces, ECS types, systems
└─────────────────────────────────────┘
```

### Modules

| Module | Role |
|--------|------|
| `domain` | Pure interfaces & data structs. Zero logic. Shared contract |
| `enginetype` | Registers all ECS `ComponentType` tokens globally |
| `components` | Built-in Component Mixins (Position, Sprite, Box, Audio...) |
| `nsystem` | All ECS Systems (Logic, Input, Draw, Alarm, Physics, Tween, Audio) |
| `core` | Engine heart: Scene, Map, Camera, SceneManager, EbitenGame loop |
| `nasset` | Asset loading (images, audio) from manifest files |
| `napi` | **The only module game developers import.** Re-exports API methods and `ncom` |
| `...` | Various other modules for layout, object pooling, physics math, etc. |

---

## 🚀 Quick Start & Tutorials

### 1. Installation

If you are creating a new game project, initialize your module and get N-Engine:

```powershell
go mod init mygame
go get github.com/Nguyen-Agn/N-Engine@latest
```

*(Note: N-Engine requires Go 1.21+)*

### 2. Basic Setup (`main.go`)

Here is the absolute minimum code to open a window and start the engine:

```go
package main

import (
	"log"
	"github.com/Nguyen-Agn/N-Engine/modules/napi"
)

func main() {
	// 1. Initialize engine
	napi.Game.Init(napi.GameConfig{
		Title:  "My First Game",
		Width:  800,
		Height: 600,
	})

	// 2. Create the first scene and transition to it
	_, err := napi.Scene.NewSceneAndGo("main", "map-800-600")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Start the game loop
	napi.Game.GameStart()
}
```

### 3. Running the built-in Demo

If you cloned the repository locally, you can run the demo directly:

```powershell
cd tests/simulation
go run .\TilemapDemo.go
```

### 4. Learning N-Engine

To master N-Engine, please refer to the detailed bilingual tutorials in the `tutorial/` folder:
- **[createGame.md](tutorial/createGame.md)**: Initialize engine, load assets, and start the game loop.
- **[CreateObject.md](tutorial/CreateObject.md)**: Create ECS objects using OOP structs and `ncom` mixins.
- **[Global.md](tutorial/Global.md)**: Manage global variables, constants, and the `napi.Store`.
- **[Layout.md](tutorial/Layout.md)**: Use the flexbox-like UI layout system (`nlayout`).
- **[TileSet.md](tutorial/TileSet.md)**: Render optimized tilemaps.
- **[BackGround.md](tutorial/BackGround.md)**: Set up scrolling backgrounds.
- **[Audio.md](tutorial/Audio.md)**: Play music and sound effects.
- **[Tween.md](tutorial/Tween.md)**: Smooth animations (Move, Scale, Alpha).
- **[Alarm.md](tutorial/Alarm.md)**: Schedule events and callbacks.
- **[Collision.md](tutorial/Collision.md)**: Configure hitboxes and collision tags.
- **[Input.md](tutorial/Input.md)**: Listen for keyboard and mouse events.
- **[Physics.md](tutorial/Physics.md)**: Control velocity, friction, and direction.
- **[Draw.md](tutorial/Draw.md)**: Draw custom shapes and text directly to the screen.
- **[NewComponent.md](tutorial/NewComponent.md)**: Define and inject your own custom data components.

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
4. **ECS + OOP Hybrid** — Data lives in ECS (Donburi); behavior is expressed through Go struct embedding.
5. **`explanation.md`-first** — Every module has an `explanation.md`. Read it before touching source code.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
