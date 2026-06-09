package napi

import (
	"log"

	"autoworld/domain"
	"autoworld/modules/nasset"
	"autoworld/modules/nsys"
)

type gameGroup struct{}

// Game là nhóm hàm quản lý vòng đời engine và Save/Load game.
// Game is the function group for engine lifecycle management.
var Game = &gameGroup{}

// _engine là singleton Engine của toàn bộ game.
// Được khởi tạo một lần qua Init(), sau đó mọi hàm napi sử dụng trực tiếp.
var _engine *Engine

// Init initializes the global Engine instance using the provided configuration.
//
// Purpose: Sets up the core game engine components and registers it as a singleton. Must be called before any other napi functions.
//
// Inputs:
// - cfg (GameConfig): Configuration parameters for the engine.
//
// Outputs:
// - *gameGroup: Allows method chaining.
func (g *gameGroup) Init(cfg GameConfig) *gameGroup {
	_engine = NewGame(cfg)
	return g
}

// GetEngine returns the current singleton Engine instance.
//
// Outputs:
// - *Engine: The active engine instance, or nil if Init() has not been called.
func (g *gameGroup) GetEngine() *Engine {
	return _engine
}

// engine is an internal function that safely returns the singleton Engine.
//
// Purpose: Ensures that napi functions panic explicitly if the engine was not initialized via Init().
//
// Outputs:
// - *Engine: The active engine instance.
func engine() *Engine {
	if _engine == nil {
		log.Fatal("[napi] chưa gọi napi.Init() — hãy gọi trước khi dùng bất kỳ hàm napi nào")
	}
	return _engine
}

// LoadFromFile reads a manifest file to load sprites and audio into the global store.
//
// Purpose: Pre-loads game assets into memory based on a manifest specification. It uses the sample rate from the engine config.
//
// Inputs:
// - path (string): The path to the manifest file.
//
// Outputs:
// - *gameGroup: Allows method chaining.
func (g *gameGroup) LoadFromFile(path string) *gameGroup {
	spriteLoader := nasset.NewSpriteLoader()
	audioLoader := nasset.NewAudioLoader(_engine.AudioCtx, 60)
	manifestLoader := nasset.NewManifestLoader(spriteLoader, audioLoader)

	manifestLoader.LoadFromFile(path, nsys.GetInstance())

	return g
}

// GameStart executes the engine's main loop and handles window management.
//
// Purpose: Blocks execution and starts the Ebitengine event loop. This should be the last function called in your main routine.
func (g *gameGroup) GameStart() {
	_engine.Start()
}

// ─── Save / Load ──────────────────────────────────────────────────────────────

// SaveGame writes the current game state to a file.
//
// Purpose: Serializes and saves progress, settings, and other globally tracked values.
//
// Inputs:
// - path (string): The path for the save file. If empty, a default path is used.
//
// Outputs:
// - error: An error if saving fails.
func (g *gameGroup) SaveGame(path string) error {
	return engine().Save.SaveGame(path)
}

// LoadGame loads the game state from a specified file.
//
// Purpose: Deserializes the file and applies the stored state to current objects and global variables.
//
// Inputs:
// - path (string): The path to the save file.
//
// Outputs:
// - error: An error if loading fails.
func (g *gameGroup) LoadGame(path string) error {
	return engine().Save.LoadGame(path)
}

// HasSave checks if a save file exists at the specified path.
//
// Inputs:
// - path (string): The save file path to check.
//
// Outputs:
// - bool: True if the file exists, false otherwise.
func (g *gameGroup) HasSave(path string) bool {
	return engine().Save.HasSave(path)
}

// DeleteSave deletes a save file at the specified path.
//
// Inputs:
// - path (string): The save file path to delete.
//
// Outputs:
// - error: An error if deletion fails.
func (g *gameGroup) DeleteSave(path string) error {
	return engine().Save.DeleteSave(path)
}

// ListSaveSlots lists all existing save files in the default save directory.
//
// Outputs:
// - []string: A slice containing the names of all found save files.
func (g *gameGroup) ListSaveSlots() []string {
	return engine().Save.ListSlots()
}

// ReadSaveSnapshot reads a save file and extracts its metadata without applying the game state.
//
// Purpose: Useful for obtaining info (like the current scene) from a save before deciding to fully load it.
//
// Inputs:
// - path (string): The path to the save file.
//
// Outputs:
// - domain.SaveSnapshot: The metadata retrieved from the save file.
// - error: An error if reading or parsing fails.
func (g *gameGroup) ReadSaveSnapshot(path string) (domain.SaveSnapshot, error) {
	return engine().Save.ReadSnapshot(path)
}
