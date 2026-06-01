package main

import (
	"autoworld/modules/core"
	"autoworld/modules/napi"
	"log"
)

func main() {
	// 1. Khởi tạo Engine
	engine := core.GameConfig{
		Title:      "AutoWorld",
		Width:      640,
		Height:     480,
		SampleRate: 44100,
	}

	// 2. Đăng ký Engine vào napi singleton — sau bước này mọi hàm napi đều dùng được
	napi.Game.Init(engine)

	// 3. Tạo và activate Scene đầu tiên
	_, err := napi.Scene.NewSceneAndGo("main", "map-640-480")
	if err != nil {
		log.Fatal(err)
	}

	// 4. Khởi chạy game
	napi.Game.GameStart()
}
