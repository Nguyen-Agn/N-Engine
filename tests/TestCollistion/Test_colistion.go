package main

import (
	"autoworld/modules/napi"
	sharedobject "autoworld/tests/SharedObject"
)

func main() {
	napi.Game.Init(napi.GameConfig{
		Title:  "test",
		Width:  640,
		Height: 480,
	})

	napi.Game.LoadFromFile("../SharedObject/Config.toml")

	napi.Scene.NewSceneAndGo("main", "map-640-480")

	sharedobject.NewMovingBox(320, 240)
	sharedobject.NewMovingBox(480, 240)

	napi.Game.GameStart()
}
