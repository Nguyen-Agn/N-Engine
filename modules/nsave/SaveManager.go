package nsave

import (
	"errors"
	"time"

	"github.com/Nguyen-Agn/N-Engine/domain"
)

const CURRENT_VERSION = 1

type saveManager struct {
	file         iFileAdapter
	sceneManager domain.ISceneManager
	store        domain.IGlobal
	collector    *collector
	autoSaveVars bool
}

// NewSaveManager creates and returns a new instance of ISaveManager.
// Inputs: sceneManager - the system controlling scenes, store - global variable registry, saveDir - folder path, autoSaveVars - boolean flag to include variables.
// Outputs: configured domain.ISaveManager implementation.
func NewSaveManager(sceneManager domain.ISceneManager, store domain.IGlobal, saveDir string, autoSaveVars bool) domain.ISaveManager {
	return &saveManager{
		file:         newJsonFileAdapter(saveDir),
		sceneManager: sceneManager,
		store:        store,
		collector:    &collector{},
		autoSaveVars: autoSaveVars,
	}
}

// SaveGame aggregates data from the current scene and global variables, storing them to the specified path.
// Inputs: path - exact file path or save name. Uses "default" if empty.
// Outputs: an error if the save process fails.
func (s *saveManager) SaveGame(path string) error {
	if path == "" {
		path = "default"
	}

	if s.sceneManager == nil {
		return errors.New("cannot save: SceneManager is nil")
	}

	currentScene := s.sceneManager.GetCurrentScene()
	if currentScene == nil {
		return errors.New("cannot save: CurrentScene is nil")
	}

	sceneID := currentScene.GetID()

	// Collect data
	objects := s.collector.collectObjects(currentScene)
	var vars map[string]any
	if s.autoSaveVars {
		vars = s.collector.collectVariables(s.store)
	}

	snap := domain.SaveSnapshot{
		Version:        CURRENT_VERSION,
		Path:           path,
		SavedAt:        time.Now().Unix(),
		CurrentSceneID: sceneID,
		Variables:      vars,
		Objects:        objects,
	}

	return s.file.Write(path, snap)
}

// LoadGame reads a save from disk and injects its data back into the variables and the current scene objects.
// Inputs: path - exact file path or save name.
// Outputs: an error if loading, parsing, or applying fails.
func (s *saveManager) LoadGame(path string) error {
	if path == "" {
		path = "default"
	}

	snap, err := s.file.Read(path)
	if err != nil {
		return err
	}

	if s.sceneManager == nil {
		return errors.New("cannot load: SceneManager is nil")
	}

	currentScene := s.sceneManager.GetCurrentScene()
	if currentScene == nil {
		return errors.New("cannot load: CurrentScene is nil")
	}

	app := &applier{}

	// Load variables if flag is true
	if s.autoSaveVars && snap.Variables != nil {
		app.applyVariables(s.store, snap.Variables)
	}

	// Load objects
	if snap.Objects != nil {
		app.applyObjects(currentScene, snap.Objects)
	}

	return nil
}

// HasSave verifies whether a specific save file slot exists.
// Inputs: path - name or path of the save.
// Outputs: boolean indicating presence.
func (s *saveManager) HasSave(path string) bool {
	if path == "" {
		path = "default"
	}
	return s.file.Exists(path)
}

// DeleteSave removes a specific save file slot permanently.
// Inputs: path - name or path of the save.
// Outputs: error if file deletion encounters an issue.
func (s *saveManager) DeleteSave(path string) error {
	if path == "" {
		path = "default"
	}
	return s.file.Delete(path)
}

// ListSlots enumerates all readable save slots discovered in the designated save directory.
// Outputs: array of save slot name strings.
func (s *saveManager) ListSlots() []string {
	return s.file.ListAll()
}

// ReadSnapshot retrieves the snapshot structure strictly as data without applying it to the game.
// Inputs: path - name or path of the save.
// Outputs: the snapshot structure domain.SaveSnapshot, and any error that arose.
func (s *saveManager) ReadSnapshot(path string) (domain.SaveSnapshot, error) {
	if path == "" {
		path = "default"
	}
	return s.file.Read(path)
}
