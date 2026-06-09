package core

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/yohamta/donburi"
)

// Minimal mock for IScene
type mockScene struct {
	id           string
	updateCount  int
	drawCount    int
	destroyCount int
}

func (m *mockScene) setID(id string)                 { m.id = id }
func (m *mockScene) GetID() string                   { m.updateCount++; return "" }
func (m *mockScene) Draw()                           { m.drawCount++ }
func (m *mockScene) Update() error                   { m.drawCount++; return nil }
func (m *mockScene) Destroy()                        { m.destroyCount++ }
func (m *mockScene) AddObject(obj domain.IObject)    {}
func (m *mockScene) AddGUIObject(obj domain.IObject) {}
func (m *mockScene) World() donburi.World            { return donburi.NewWorld() } // Need import donburi
func (m *mockScene) GetMap() domain.IMap             { return nil }
func (m *mockScene) GetGUIMap() domain.IMap          { return nil }
func (m *mockScene) GetCamera() domain.ICamera       { return nil }

func TestSceneManager(t *testing.T) {
	// The real NewSceneManager creates a real globalScene which requires real dependencies.
	// But we can test it if NewScene allows nil or minimal dependencies.
	// In core, NewSceneManager might try to create a real Scene.
	// Let's create an empty SceneManager manually to avoid complex real NewScene creation.
	sm := &SceneManager{
		sceneList:   make(map[string]*sceneEntry),
		screenW:     800,
		screenH:     600,
		globalScene: nil, // Leave nil for this test if possible, or we skip Update() if it crashes.
	}

	scene1 := &mockScene{}
	scene2 := &mockScene{}

	err := sm.AddScene("main", scene1)
	if err != nil {
		t.Fatalf("Failed to add scene: %v", err)
	}

	err = sm.ChangeSceneFromList("main")
	if err != nil {
		t.Fatalf("Failed to change scene: %v", err)
	}

	if sm.GetCurrentScene() != domain.IScene(scene1) {
		t.Errorf("Current scene is not scene1")
	}

	sm.Draw()
	if scene1.drawCount != 1 {
		t.Errorf("Expected scene1.Draw to be called once")
	}

	err = sm.ChangeSceneForce(scene2)
	if err != nil {
		t.Fatalf("Failed to force change scene: %v", err)
	}

	if sm.GetCurrentScene() != domain.IScene(scene2) {
		t.Errorf("Current scene is not scene2")
	}
	// scene1 should have been destroyed because it was in list but forced out
	if scene1.destroyCount != 1 {
		t.Errorf("Expected scene1 to be destroyed, got %v", scene1.destroyCount)
	}

	// Add another scene to test remove
	scene3 := &mockScene{}
	sm.AddScene("menu", scene3)
	err = sm.RemoveScene("menu")
	if err != nil {
		t.Fatalf("Failed to remove scene: %v", err)
	}
	if scene3.destroyCount != 1 {
		t.Errorf("Expected scene3 to be destroyed")
	}
}
