package core

import (
	"fmt"
)

// ─── SceneManager ────────────────────────────────────────────────────────────

// sceneEntry là internal record to save IScene and pause-state.
type sceneEntry struct {
	scene    IScene
	isPaused bool // true = no Update/Draw
}

// SceneManager manger all Scene's lifecycle.
type SceneManager struct {
	currentScene IScene

	// sceneList is list of scene is being not using
	sceneList map[string]*sceneEntry

	// globalScene is hidden scene, always run Update, but no Draw
	// Should save consitency object (not killed by changing Scene).
	globalScene *Scene

	screenW, screenH int

	input IInputManager
}

// NewSceneManager creates a new SceneManager instance.
//
// Purpose: Initializes the manager responsible for tracking active and paused scenes, as well as maintaining the global background scene.
//
// Inputs:
// - screenW (int): The width of the game screen in pixels.
// - screenH (int): The height of the game screen in pixels.
// - input (IInputManager): The global input manager to distribute to scenes.
//
// Outputs:
// - *SceneManager: The configured manager.
func NewSceneManager(screenW, screenH int, input IInputManager) *SceneManager {
	return &SceneManager{
		sceneList: make(map[string]*sceneEntry),
		screenW:   screenW,
		screenH:   screenH,
		input:     input,
		// Global Scene: unlimit size (0,0), viewport = screen
		globalScene: NewScene(input, 0, 0, screenW, screenH),
	}
}

// ─── ISceneManager interface — Lifecycle ─────────────────────────────────────

// Update advances the logic of all active scenes for one frame.
//
// Purpose: Updates the global scene first, followed by the active current scene.
//
// Outputs:
// - error: An error if any scene fails to update, otherwise nil.
func (r *SceneManager) Update() error {
	// Global scene allways run
	if err := r.globalScene.Update(); err != nil {
		return err
	}
	if r.currentScene == nil {
		return nil
	}

	return r.currentScene.Update()
}

// Draw renders the active current scene.
//
// Purpose: Called by Ebitengine's render loop. Only the currently active scene is drawn.
//
// Outputs:
// - error: Always returns nil.
func (r *SceneManager) Draw() error {
	if r.currentScene == nil {
		return nil
	}
	r.currentScene.Draw()
	return nil
}

// Layout dictates the logical screen dimensions for Ebitengine.
//
// Outputs:
// - (int, int): The logical width and height defined at creation.
func (r *SceneManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return r.screenW, r.screenH
}

// ─── ISceneManager interface — Scene Management ─────────────────────────────

// AddScene registers a new scene into the waiting list without activating it.
//
// Purpose: Stores a scene for later use, pausing it initially.
//
// Inputs:
// - id (string): The unique identifier for the scene.
// - scene (IScene): The scene instance to add.
//
// Outputs:
// - error: Returns an error if the ID already exists in the list.
func (r *SceneManager) AddScene(id string, scene IScene) error {
	if _, exists := r.sceneList[id]; exists {
		return fmt.Errorf("scene '%s' already exists in list", id)
	}
	if s, ok := scene.(interface{ setID(string) }); ok {
		s.setID(id)
	}
	r.sceneList[id] = &sceneEntry{
		scene:    scene,
		isPaused: true, // Mặc định pause khi mới thêm vào
	}
	return nil
}

// ChangeSceneFromList switches the active scene to one stored in the waitlist.
//
// Purpose: Pauses the current scene (without destroying it) and resumes the requested scene. Allows for reusable scenes.
//
// Inputs:
// - id (string): The ID of the scene to activate.
//
// Outputs:
// - error: Returns an error if the ID does not exist in the list.
func (r *SceneManager) ChangeSceneFromList(id string) error {
	entry, exists := r.sceneList[id]
	if !exists {
		return fmt.Errorf("scene '%s' not found in list", id)
	}

	// Pause current Scene  (not delete)
	if r.currentScene != nil {
		// Tìm và đánh dấu pause Scene hiện tại trong list
		for _, e := range r.sceneList {
			if e.scene == r.currentScene {
				e.isPaused = true
				break
			}
		}
	}

	// active new scene
	entry.isPaused = false
	r.currentScene = entry.scene
	return nil
}

// ChangeSceneForce forcefully switches to a new scene, permanently destroying the current one.
//
// Purpose: Transitions to a fresh scene instance and cleans up the active scene.
//
// Inputs:
// - next (IScene): The new scene instance to activate.
//
// Outputs:
// - error: Returns an error if the provided scene is nil.
func (r *SceneManager) ChangeSceneForce(next IScene) error {
	if next == nil {
		return fmt.Errorf("cannot force change to nil scene")
	}

	// Destroy Scene hiện tại nếu có
	if r.currentScene != nil {
		// Xóa Scene hiện tại khỏi list nếu nó đang được lưu trong list
		for id, entry := range r.sceneList {
			if entry.scene == r.currentScene {
				delete(r.sceneList, id)
				break
			}
		}
		r.currentScene.Destroy()
	}

	r.currentScene = next
	return nil
}

// RemoveScene permanently removes and destroys a paused scene from the waitlist.
//
// Purpose: Frees resources associated with a scene that is no longer needed.
//
// Inputs:
// - id (string): The ID of the scene to remove.
//
// Outputs:
// - error: Returns an error if the scene is not found or is currently active.
func (r *SceneManager) RemoveScene(id string) error {
	entry, exists := r.sceneList[id]
	if !exists {
		return fmt.Errorf("scene '%s' not found in list", id)
	}

	// Không cho xóa Scene đang chạy qua method này
	if entry.scene == r.currentScene {
		return fmt.Errorf("cannot remove active scene '%s', use ChangeSceneForce first", id)
	}

	entry.scene.Destroy()
	delete(r.sceneList, id)
	return nil
}

// RemoveAllScene destroys every scene in the waitlist and clears the current scene.
//
// Purpose: Performs a complete cleanup of all managed scenes, usually during a hard reset or exit.
//
// Outputs:
// - error: Always returns nil.
func (r *SceneManager) RemoveAllScene() error {
	for id, entry := range r.sceneList {
		entry.scene.Destroy()
		delete(r.sceneList, id)
	}
	r.currentScene = nil
	return nil
}

// GetCurrentScene retrieves the currently running active scene.
//
// Outputs:
// - IScene: The active scene interface.
func (r *SceneManager) GetCurrentScene() IScene {
	return r.currentScene
}

// GetSceneFromList retrieves a scene from the waitlist without activating it.
//
// Inputs:
// - id (string): The scene's identifier.
//
// Outputs:
// - IScene: The found scene, or nil if it doesn't exist.
func (r *SceneManager) GetSceneFromList(id string) IScene {
	entry, exists := r.sceneList[id]
	if !exists {
		return nil
	}
	return entry.scene
}

// GetGlobalScene retrieves the hidden background scene.
//
// Purpose: Used to add persistent objects that should not be destroyed across regular scene transitions.
//
// Outputs:
// - IScene: The global scene interface.
func (r *SceneManager) GetGlobalScene() IScene {
	return r.globalScene
}
