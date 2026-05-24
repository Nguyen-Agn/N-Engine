package core

import (
	"log"

	globalconfig "autoworld/modules/globaLConfig"
	"autoworld/modules/nsys"

	"github.com/hajimehoshi/ebiten/v2"
	ebitenAudio "github.com/hajimehoshi/ebiten/v2/audio"
)

// Engine là struct trung tâm (Root) của toàn bộ hệ thống game.
// Chứa tất cả manager và resource cần thiết để vận hành game.
type Engine struct {
	Scene    ISceneManager
	Config   IGlobalConfig
	Input    IInputManager
	Store    IGlobal              // Global resource store (Singleton từ nglobal)
	AudioCtx *ebitenAudio.Context // Ebitengine audio context dùng chung
}

// GameConfig chứa các tham số cấu hình ban đầu khi khởi động game.
type GameConfig struct {
	// Tiêu đề cửa sổ game
	Title string
	// Chiều rộng cửa sổ (pixel)
	Width int
	// Chiều cao cửa sổ (pixel)
	Height int
	// Sample rate cho âm thanh (thường là 44100)
	SampleRate int
}

// NewGame khởi tạo Engine với GameConfig cho sẵn.
// Tự động setup: SceneManager, GlobalConfig, InputManager, AudioContext, GlobalStore.
// Gọi engine.Start() sau khi đã thêm scene để chạy game.
func NewGame(cfg GameConfig) *Engine {
	sampleRate := cfg.SampleRate
	if sampleRate == 0 {
		sampleRate = 44100
	}

	_config := globalconfig.NewGlobalConfig()
	_config.SetValue("game-title", cfg.Title)
	_config.SetValue("game-width", cfg.Width)
	_config.SetValue("game-height", cfg.Height)
	_config.SetValue("game-rate", cfg.SampleRate)

	_input := NewInputManager()

	return &Engine{
		Scene:    NewSceneManager(cfg.Width, cfg.Height, _input),
		Config:   _config,
		Input:    _input,
		Store:    nsys.GetInstance(),
		AudioCtx: ebitenAudio.NewContext(sampleRate),
	}
}

// Start cấu hình cửa sổ và chạy vòng lặp Ebitengine.
// Đây là lời gọi cuối cùng trong hàm main — hàm này sẽ block cho đến khi game đóng.
func (e *Engine) Start() {
	config := e.Config

	ebiten.SetWindowSize(config.GetInt("game-width"), config.GetInt("game-height"))
	ebiten.SetWindowTitle(config.GetString("game-title"))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewEbitenGame(e)); err != nil {
		log.Fatal(err)
	}
}
