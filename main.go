package main

import (
	"log"

	"github.com/user/autoworld/modules/core"
	"github.com/user/autoworld/modules/napi"
)

func main() {
	// 1. Khởi tạo Engine
	engine := core.NewGame(core.GameConfig{
		Title:      "AutoWorld",
		Width:      640,
		Height:     480,
		SampleRate: 44100,
	})

	// 2. Đăng ký Engine vào napi singleton — sau bước này mọi hàm napi đều dùng được
	napi.Init(engine)

	// 3. Tạo và activate Scene đầu tiên
	_, err := napi.NewSceneAndGo("main")
	if err != nil {
		log.Fatal(err)
	}

	// 4. Khởi chạy game
	engine.Start()
}
