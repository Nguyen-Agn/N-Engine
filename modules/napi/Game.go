package napi

import (
	"log"

	"autoworld/modules/nasset"
	"autoworld/modules/nsys"
)

type gameGroup struct{}

// Game là nhóm hàm quản lý vòng đời engine.
// Game is the function group for engine lifecycle management.
var Game = &gameGroup{}

// _engine là singleton Engine của toàn bộ game.
// Được khởi tạo một lần qua Init(), sau đó mọi hàm napi sử dụng trực tiếp.
var _engine *Engine

// Init khởi tạo Engine từ GameConfig và đăng ký làm singleton toàn cục.
// PHẢI gọi trước khi sử dụng bất kỳ hàm napi nào.
// Thứ tự khuyến nghị trong main.go: Init → LoadManifest → NewSceneAndGo → GameStart.
func (g *gameGroup) Init(cfg GameConfig) *gameGroup {
	_engine = NewGame(cfg)
	return g
}

// GetEngine trả về singleton *Engine hiện tại.
// Trả về nil nếu Init() chưa được gọi.
func (g *gameGroup) GetEngine() *Engine {
	return _engine
}

// engine là hàm nội bộ trả về singleton engine.
// Panic (Fatal) nếu Init() chưa được gọi — giúp phát hiện lỗi khởi tạo sớm.
func engine() *Engine {
	if _engine == nil {
		log.Fatal("[napi] chưa gọi napi.Init() — hãy gọi trước khi dùng bất kỳ hàm napi nào")
	}
	return _engine
}

// LoadFromFile đọc file manifest tại path và load toàn bộ sprite/audio vào global store.
// Khác LoadManifest ở chỗ: dùng SampleRate từ engine config thay vì giá trị cố định.
// Trả về lỗi qua log nếu file không đọc được hoặc có resource nào load thất bại.
func (g *gameGroup) LoadFromFile(path string) *gameGroup {
	spriteLoader := nasset.NewSpriteLoader()
	audioLoader := nasset.NewAudioLoader(_engine.AudioCtx, 60)
	manifestLoader := nasset.NewManifestLoader(spriteLoader, audioLoader)

	manifestLoader.LoadFromFile(path, nsys.GetInstance())

	return g
}

// GameStart cấu hình cửa sổ và chạy vòng lặp Ebitengine.
// Đây là lời gọi cuối cùng trong main.go — hàm block cho đến khi game đóng.
func (g *gameGroup) GameStart() {
	_engine.Start()
}
