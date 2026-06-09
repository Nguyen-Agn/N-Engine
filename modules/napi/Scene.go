package napi

import (
	"strconv"
	"strings"

	"github.com/Nguyen-Agn/N-Engine/modules/core"
)

// ─────────────────────────────────────────────
//
// Group: Scene — quản lý Scene
//
// Group: Scene — Scene management
//
// ─────────────────────────────────────────────
type sceneGroup struct{}

// Scene là nhóm hàm quản lý Scene và Camera.
// Scene is the function group for Scene and Camera management.
var Scene = &sceneGroup{}

// ─── Scene Management ──────────────────────────────────────────────────────────

// AddScene inserts a new scene into the engine's management list without activating it.
//
// Inputs:
// - id (string): The unique identifier for the scene.
// - scene (IScene): The scene instance to add.
//
// Outputs:
// - error: An error if the ID already exists.
func (s *sceneGroup) AddScene(id string, scene IScene) error {
	return engine().Scene.AddScene(id, scene)
}

// GoToScene switches to an existing scene from the management list by its ID.
//
// Purpose: Pauses the currently active scene (without destroying it) and resumes the target scene. Useful for pausing menus or reusable states.
//
// Inputs:
// - id (string): The target scene's ID.
//
// Outputs:
// - error: An error if the ID does not exist.
func (s *sceneGroup) GoToScene(id string) error {
	return engine().Scene.ChangeSceneFromList(id)
}

// ReplaceScene forcefully transitions to a new scene and permanently destroys the current one.
//
// Purpose: Used when the current scene is no longer needed and its resources should be freed.
//
// Inputs:
// - next (IScene): The new scene instance to transition to.
//
// Outputs:
// - error: An error if the transition fails.
func (s *sceneGroup) ReplaceScene(next IScene) error {
	return engine().Scene.ChangeSceneForce(next)
}

// RemoveScene permanently removes and destroys a paused scene from the management list.
//
// Inputs:
// - id (string): The ID of the scene to remove.
//
// Outputs:
// - error: An error if the scene cannot be found or is currently active.
func (s *sceneGroup) RemoveScene(id string) error {
	return engine().Scene.RemoveScene(id)
}

// RemoveAllScenes destroys all managed scenes and clears the current scene.
//
// Outputs:
// - error: An error if the cleanup fails, nil otherwise.
func (s *sceneGroup) RemoveAllScenes() error {
	return engine().Scene.RemoveAllScene()
}

// GetCurrentScene retrieves the active scene that is currently running.
//
// Outputs:
// - IScene: The active scene, or nil if none are running.
func (s *sceneGroup) GetCurrentScene() IScene {
	_s := engine().Scene.GetCurrentScene()
	if _s == nil {
		return nil
	}
	return _s
}

// GetSceneByID retrieves a scene from the waitlist by its ID without activating it.
//
// Inputs:
// - id (string): The ID of the scene.
//
// Outputs:
// - IScene: The found scene, or nil if not found.
func (s *sceneGroup) GetSceneByID(id string) IScene {
	_s := engine().Scene.GetSceneFromList(id)
	if _s == nil {
		return nil
	}
	return _s
}

// NewScene creates and registers a new scene based on the given configuration tokens.
//
// Purpose: Simplifies the setup of physical and GUI maps for a new scene using string tokens (e.g., "gui-640-480 map-1280-1280").
//
// Inputs:
// - id (string): The unique ID for the new scene.
// - component (string): Configuration tokens for map sizes.
//
// Outputs:
// - IScene: The created scene.
// - error: An error if registration fails (e.g., duplicate ID).
func (s *sceneGroup) NewScene(id, component string) (IScene, error) {
	e := engine()
	viewW := e.Config.GetValue("game-width").(int)
	viewH := e.Config.GetValue("game-height").(int)

	var mapW, mapH int
	var guiW, guiH int
	var hasGui bool

	// Parse component string
	tokens := strings.FieldsSeq(component)
	for token := range tokens {
		parts := strings.Split(token, "-")
		if len(parts) >= 3 {
			switch parts[0] {
			case "map":
				mapW, _ = strconv.Atoi(parts[1])
				mapH, _ = strconv.Atoi(parts[2])
			case "gui":
				hasGui = true
				guiW, _ = strconv.Atoi(parts[1])
				guiH, _ = strconv.Atoi(parts[2])
			}
		} else if len(parts) >= 1 {
			if parts[0] == "gui" {
				hasGui = true
				guiW = viewW
				guiH = viewH
			}
		}
	}

	scene := core.NewScene(e.Input, mapW, mapH, viewW, viewH)

	if hasGui {
		scene.SetGUIMap(e.Input, guiW, guiH)
	}

	if err := e.Scene.AddScene(id, scene); err != nil {
		return nil, err
	}
	return scene, nil
}

// NewSceneAndGo creates a new scene, registers it, and immediately transitions to it.
//
// Purpose: A convenient shortcut for initial game setup where the first scene needs to be instantly activated.
//
// Inputs:
// - id (string): The unique ID for the new scene.
// - component (string): Configuration tokens.
//
// Outputs:
// - IScene: The created and activated scene.
// - error: An error if creation or transition fails.
func (s *sceneGroup) NewSceneAndGo(id, component string) (IScene, error) {
	scene, err := s.NewScene(id, component)
	if err != nil {
		return nil, err
	}
	if err := engine().Scene.ChangeSceneFromList(id); err != nil {
		return nil, err
	}
	return scene, nil
}

// ─── Global Hidden Scene ──────────────────────────────────────────────────────

// GetGlobalScene retrieves the persistent background scene.
//
// Purpose: Accesses the global scene that updates every frame regardless of which main scene is active. Useful for persistent objects like background audio or global managers.
//
// Outputs:
// - IScene: The global hidden scene.
func (s *sceneGroup) GetGlobalScene() IScene {
	_s := engine().Scene.GetGlobalScene()
	if _s == nil {
		return nil
	}
	return _s
}

// ─── Map & Camera Helpers ──────────────────────────────────────────────────────

// GetMap retrieves the physical map from the specified scene.
//
// Inputs:
// - scene (IScene): The scene to query.
//
// Outputs:
// - IMap: The physical map interface, or nil if the scene is nil.
func (s *sceneGroup) GetMap(scene IScene) IMap {
	if scene == nil {
		return nil
	}
	return scene.GetMap()
}

// GetGUIMap retrieves the GUI map from the specified scene.
//
// Inputs:
// - scene (IScene): The scene to query.
//
// Outputs:
// - IMap: The GUI map interface, or nil if not initialized.
func (s *sceneGroup) GetGUIMap(scene IScene) IMap {
	if scene == nil {
		return nil
	}
	return scene.GetGUIMap()
}

// GetCamera retrieves the camera associated with the specified scene.
//
// Inputs:
// - scene (IScene): The scene to query.
//
// Outputs:
// - ICamera: The camera interface, or nil if the scene is nil.
func (s *sceneGroup) GetCamera(scene IScene) ICamera {
	if scene == nil {
		return nil
	}
	return scene.GetCamera()
}

// SetCameraTarget sets a target object for the scene's camera to follow automatically.
//
// Purpose: Centers the camera viewport on the given object every frame. Pass nil to disable following.
//
// Inputs:
// - scene (IScene): The scene owning the camera.
// - target (IObject): The object to track.
func (s *sceneGroup) SetCameraTarget(scene IScene, target IObject) {
	if scene == nil {
		return
	}
	cam := scene.GetCamera()
	if cam != nil {
		cam.SetTarget(target)
	}
}

// MoveCamera instantly translates the camera to the specified map coordinates.
//
// Inputs:
// - scene (IScene): The scene owning the camera.
// - x (float32): The new X coordinate.
// - y (float32): The new Y coordinate.
func (s *sceneGroup) MoveCamera(scene IScene, x, y float32) {
	if scene == nil {
		return
	}
	cam := scene.GetCamera()
	if cam != nil {
		cam.SetX(x)
		cam.SetY(y)
	}
}
