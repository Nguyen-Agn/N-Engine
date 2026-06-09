package core

import (
	"log"

	"autoworld/modules/components"
	globalconfig "autoworld/modules/globaLConfig"
	"autoworld/modules/nsave"
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
	Save     ISaveManager         // Save manager system
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
	// SaveDir là thư mục lưu trữ file save (mặc định "./saves")
	SaveDir string
	// AutoSaveVars cho biết có tự động lưu toàn bộ biến từ nsys/Global không
	AutoSaveVars bool
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
	_scene := NewSceneManager(cfg.Width, cfg.Height, _input)
	_store := nsys.GetInstance()

	// Inject IInputManager vào components package để MouseComponent có thể truy cập.
	components.SetGlobalInputManager(_input)

	return &Engine{
		Scene:    _scene,
		Config:   _config,
		Input:    _input,
		Store:    _store,
		AudioCtx: ebitenAudio.NewContext(sampleRate),
		Save:     nsave.NewSaveManager(_scene, _store, cfg.SaveDir, cfg.AutoSaveVars),
	}
}

// Start cấu hình cửa sổ và chạy vòng lặp Ebitengine.
// Đây là lời gọi cuối cùng trong hàm main — hàm này sẽ block cho đến khi game đóng.
func (e *Engine) Start() {
	config := e.Config

	ebiten.SetWindowSize(config.GetValue("game-width").(int), config.GetValue("game-height").(int))
	ebiten.SetWindowTitle(config.GetValue("game-title").(string))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewEbitenGame(e)); err != nil {
		log.Fatal(err)
	}
}
