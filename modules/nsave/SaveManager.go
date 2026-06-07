package nsave

import (
	"autoworld/domain"
	"errors"
	"time"
)

const CURRENT_VERSION = 1

type saveManager struct {
	file         iFileAdapter
	sceneManager domain.ISceneManager
	store        domain.IGlobal
	collector    *collector
	autoSaveVars bool
}

// NewSaveManager creates a new instance of ISaveManager
func NewSaveManager(sceneManager domain.ISceneManager, store domain.IGlobal, saveDir string, autoSaveVars bool) domain.ISaveManager {
	return &saveManager{
		file:         newJsonFileAdapter(saveDir),
		sceneManager: sceneManager,
		store:        store,
		collector:    &collector{},
		autoSaveVars: autoSaveVars,
	}
}

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

func (s *saveManager) HasSave(path string) bool {
	if path == "" {
		path = "default"
	}
	return s.file.Exists(path)
}

func (s *saveManager) DeleteSave(path string) error {
	if path == "" {
		path = "default"
	}
	return s.file.Delete(path)
}

func (s *saveManager) ListSlots() []string {
	return s.file.ListAll()
}

func (s *saveManager) ReadSnapshot(path string) (domain.SaveSnapshot, error) {
	if path == "" {
		path = "default"
	}
	return s.file.Read(path)
}
