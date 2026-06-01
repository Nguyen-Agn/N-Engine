# AutoWorld (N-Engine) - Project Context

AutoWorld (N-Engine) is a lightweight, high-performance 2D Game Engine built in Go, powered by **Ebitengine** (rendering, audio, input) and **Donburi** (ECS). It combines the efficiency of ECS with the developer-friendly patterns of Object-Oriented Programming (OOP) through component mixins and reflection-based binding.

## 🏗 Architecture Overview

The project follows a strict dependency layer architecture. Higher layers interact with lower layers **only through interfaces** defined in the `domain` module.

1.  **`domain` (Contract Layer)**: Pure interfaces and data structs. No logic. Zero dependencies on other internal modules.
2.  **`modules/enginetype`**: Global registry for ECS `ComponentType` tokens.
3.  **`modules/components`**: Built-in Component Mixins (Position, Sprite, Box, Audio, etc.) that game objects embed to gain capabilities.
4.  **`modules/nsystem`**: Implementation of ECS Systems (Logic, Input, Draw, Alarm, Physics, Tween, Audio).
5.  **`modules/core`**: The engine's heart, coordinating the game loop, scenes, maps, cameras, and managers.
6.  **`modules/napi` (Public API)**: The **only** module intended for use by game developers. It provides a simplified interface to all engine features.

## 🚀 Key Development Commands

| Task | Command |
| :--- | :--- |
| **Run Main Entry** | `go run main.go` |
| **Run Simulation Demo** | `go run .\tests\simulation\TilemapDemo.go` |
| **Run Tests** | `go test ./...` |
| **Clean Dependencies** | `go mod tidy` |

## 🛠 Development Conventions & Rules

### Core Mandates (from `AGENT_GUIDE.md`)
*   **Interface First**: Modules must communicate via interfaces defined in `domain`. Assume interfaces provide data exactly as commented.
*   **Explanation-First**: Always read the `explanation.md` file in a module folder before reading source code or making changes.
*   **Keep Docs Sync**: If a module's structure or logic changes, you **must** update its corresponding `explanation.md`.
*   **Open/Closed Principle**: Extend functionality by adding new types/files; minimize modifications to existing stable modules.
*   **Encapsulation**: Game logic should only import `napi`. Internal engine modules should not import `napi`.

### Coding Patterns
*   **Game Objects**: Plain Go structs embedding `napi.IObject` and various capability mixins (e.g., `napi.IPosition`, `napi.ISprite`).
*   **Lifecycle**: Objects implement `Create()`, `StepUpdate()`, and `Destroy()`.
*   **Initialization**: Use `napi.NewObject(obj, name, "tokens...")` to bind ECS data to the struct, then `napi.Register(obj, isGlobal)` to add it to a scene.
*   **Component Tokens**: Space-separated strings used in `NewObject` (e.g., `"pos spr vel"`).

## 📁 Directory Map

*   `domain/`: Interfaces (`ObjectInterface.go`) and ECS data structs (`ObjectData.go`).
*   `modules/napi/`: Simplified developer API.
*   `modules/core/`: Scene, Map, Camera, and Ebiten loop integration.
*   `modules/components/`: Component mixins implementation.
*   `modules/nsystem/`: ECS System implementations.
*   `modules/nasset/`: Asset loading from TOML/JSON manifests.
*   `.libs/`: Vendored dependencies (Ebiten, Donburi). Use `replace` in `go.mod` to point here.

## 🧩 ECS Integration (Donburi)
*   **World**: Managed within a `Map` (inside a `Scene`).
*   **Systems**: Executed in a specific order: `Input` -> `Logic` -> `Alarm` -> `Tween` -> `Physics` -> `Audio` -> `Draw`.
*   **Culling**: `DrawSystem` performs viewport culling, but `LogicSystem` updates all active entities.

---
*Note: This file is a living document. Update it as the engine evolves.*
