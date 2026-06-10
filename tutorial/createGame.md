# Khởi động Game / Start a Game

> **Vision**: Code ngắn, ít hàm gọi — chỉ cần vài dòng là chạy được game.
> **Vision**: Short code, few calls — just a few lines to run a game.

---

## 1. Giải thích / Explanation

Để khởi động một game với N-Engine, bạn cần:
To start a game with N-Engine, you need:

1. **Khởi tạo engine** với file config — `napi.Game.Init(cfg)`
2. **Nạp tài nguyên** từ manifest TOML — `napi.Game.LoadFromFile(path)`
3. **Tạo scene đầu tiên** và chuyển đến đó — `napi.Scene.NewSceneAndGo(id, component)`
4. **Chạy game** — `napi.Game.GameStart()`

---

## 2. Ví dụ / Code Example

```go
package main

import (
	"log"

	"github.com/Nguyen-Agn/N-Engine/modules/napi"
)

func main() {
	// Cấu hình cửa sổ / Window config
	cfg := napi.GameConfig{
		Title:      "My Game",
		Width:      800,
		Height:     600,
	}

	// Khởi tạo engine / Initialize engine
	napi.Game.Init(cfg)

	// Nạp tài nguyên từ file TOML / Load resources from TOML [Optional]
	napi.Game.LoadFromFile("assets/manifest.toml")

	// Tạo scene "main" với một physical map và chuyển đến ngay
	// Create "main" scene with a physical map and go to it immediately
	// Format: "map-WxH" (Physical Map) hoặc "gui-WxH" (GUI Map)
	_, err := napi.Scene.NewSceneAndGo("main", "map-800-600")
	if err != nil {
		log.Fatal(err)
	}

	// YOUR GAME LOGIN
	// Đăng ký objects vào scene / Register objects into scene
	// NewPlayer() -> Sẽ gọi napi.Obj.NewObject() và napi.Obj.Register() bên trong.



	// Final;
	// Chạy game / Start game loop
	napi.Game.Start()
}
```

---

## 3. Quản lý Scene / Scene Management

Engine cung cấp các API để tạo và chuyển Scene rất dễ dàng:
The engine provides easy APIs to create and switch Scenes:

```go
// Tạo scene nhưng không chuyển đến / Create scene without navigating
napi.Scene.NewScene("menu", "map-800-600")
napi.Scene.NewScene("game", "map-1600-1200")

// Chuyển scene / Navigate to scene (Scene cũ sẽ bị pause, có thể quay lại)
napi.Scene.GoToScene("game")

// Thay thế scene / Replace scene (Scene cũ sẽ bị destroy)
// napi.Scene.ReplaceScene(nextScene)

// Lấy scene hiện tại / Get current scene
current := napi.Scene.GetCurrentScene()

// Lấy scene theo ID / Get scene by ID
menuScene := napi.Scene.GetSceneByID("menu")
```

## 4. Map & Camera

Bản đồ vật lý (Physical Map) và Bản đồ giao diện (GUI Map) được gắn với Scene. Camera cũng là một phần của Scene:

```go
// Đặt camera bám theo một Object
napi.Scene.SetCameraTarget(napi.Scene.GetCurrentScene(), playerObj)

// Dịch chuyển camera tức thì
napi.Scene.MoveCamera(napi.Scene.GetCurrentScene(), 100, 200)
```
