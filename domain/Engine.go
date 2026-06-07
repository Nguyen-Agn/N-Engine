package domain

import "github.com/hajimehoshi/ebiten/v2/audio"

type Engine struct {
	Scene    ISceneManager
	Config   IGlobalConfig
	Input    IInputManager
	Store    IGlobal        // Global resource store (Singleton từ nglobal)
	AudioCtx *audio.Context // Ebitengine audio context dùng chung

	IEngine
}

type IEngine interface {
	// Start starts the game loop. Window title and size are retrieved from global config.
	Start()
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
