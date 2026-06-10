package core

import (
	"log"

	"github.com/Nguyen-Agn/N-Engine/modules/components"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"
	globalconfig "github.com/Nguyen-Agn/N-Engine/modules/globaLConfig"
	"github.com/Nguyen-Agn/N-Engine/modules/nsave"
	"github.com/Nguyen-Agn/N-Engine/modules/nsys"

	"github.com/hajimehoshi/ebiten/v2"
	ebitenAudio "github.com/hajimehoshi/ebiten/v2/audio"
)

// Engine is center struct (Root) of entire game system.
//
// Contain all manager to run system
type Engine struct {
	Scene    ISceneManager
	Config   IGlobalConfig
	Input    IInputManager
	Store    IGlobal              // Global resource store (Singleton từ nglobal)
	AudioCtx *ebitenAudio.Context // Ebitengine audio context dùng chung
	Save     ISaveManager         // Save manager system
}

// GameConfig contain all game' status to begin Game
type GameConfig struct {
	// title
	Title string
	// width of scene|windows (pixel)
	Width int
	// height of scene|windows (pixel)
	Height int
	// Sample rate for audio (default 44100)
	SampleRate int
	// SaveDir is default path for file save (mặc định "./saves")
	SaveDir string
	// AutoSaveVars : is? auto-save all values from nsys/Global?
	AutoSaveVars bool
}

// NewGame initializes the Engine using the provided GameConfig.
//
// Purpose: Automates the setup of core engine components including the SceneManager, GlobalConfig, InputManager, AudioContext, and GlobalStore.
// After initialization, you must call Start() to open the window and begin the main loop.
//
// Inputs:
// - cfg (GameConfig): A struct containing the initial settings for the game (e.g., window size, title, sample rate).
//
// Outputs:
// - *Engine: A pointer to the newly created Engine instance.
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
	enginetype.LogError("[Game] Input-system installed")

	return &Engine{
		Scene:    _scene,
		Config:   _config,
		Input:    _input,
		Store:    _store,
		AudioCtx: ebitenAudio.NewContext(sampleRate),
		Save:     nsave.NewSaveManager(_scene, _store, cfg.SaveDir, cfg.AutoSaveVars),
	}
}

// Start opens the game window and begins the main game loop.
//
// Purpose: Sets window properties based on the global configuration and hands control over to Ebitengine's RunGame.
//
// Special Requirements:
// - This must be the final function called in your main routine, as it blocks execution until the game is closed.
func (e *Engine) Start() {
	config := e.Config

	// set windows
	ebiten.SetWindowSize(config.GetValue("game-width").(int), config.GetValue("game-height").(int))
	ebiten.SetWindowTitle(config.GetValue("game-title").(string))

	// open
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// main loop
	if err := ebiten.RunGame(NewEbitenGame(e)); err != nil {
		log.Fatal(err)
	}
}
