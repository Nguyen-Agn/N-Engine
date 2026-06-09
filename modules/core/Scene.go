package core

import (
	"autoworld/domain"

	"github.com/yohamta/donburi"
)

// Scene is the space for object active in Engine.
// Per-Scene contains:
//   - Physical Map: ECS World + Logic/Audio/Input system + objectList
//   - GUI Map: screen-space overlay cho HUD (optional)
//   - Camera: viewport, follow-target, DrawSystem
type Scene struct {
	map_   *Map // physical game world (not nil)
	guiMap *Map // screen-space GUI/HUD (optional)
	camera *Camera
	input  domain.IInputManager
	id     string // ID of the scene
}

// NewScene creates a new Scene instance.
//
// Purpose: Initializes a complete scene environment, including a physical map and a camera.
// It also auto-registers the map with the camera's draw system.
//
// Inputs:
// - input (domain.IInputManager): The global input manager to be used by the physical map.
// - mapW (int): The width of the physical map in pixels (0 for unbounded).
// - mapH (int): The height of the physical map in pixels (0 for unbounded).
// - viewW (int): The width of the camera viewport in pixels.
// - viewH (int): The height of the camera viewport in pixels.
//
// Outputs:
// - *Scene: A newly created Scene.
func NewScene(input domain.IInputManager, mapW, mapH, viewW, viewH int) *Scene {
	scene := &Scene{
		map_:   NewMap(input, mapW, mapH),
		camera: NewCamera(viewW, viewH),
		input:  input,
	}

	// Inject DrawSystem into Map so IDraw objects are auto-registered on AddObject
	if dr := scene.camera.GetDrawSystem(); dr != nil {
		scene.map_.SetDrawRegistry(dr)
	}
	return scene
}

// Update processes logic for the scene and its maps.
//
// Purpose: Called once per frame to update the physical map, the GUI map (if it exists), and the camera's position.
//
// Outputs:
// - error: Returns an error if map updates fail, nil otherwise.
func (s *Scene) Update() error {
	if err := s.map_.Update(); err != nil {
		return err
	}
	if s.guiMap != nil {
		if err := s.guiMap.Update(); err != nil {
			return err
		}
	}
	s.camera.Update(s.map_.Width(), s.map_.Height())
	return nil
}

// Draw renders the scene to the screen.
//
// Purpose: Instructs the camera to draw the physical map's ECS world, and optionally the GUI map's world.
func (s *Scene) Draw() {
	if s.guiMap != nil {
		s.camera.Draw(s.map_.World(), s.guiMap.World())
	} else {
		s.camera.Draw(s.map_.World(), nil)
	}
}

// Destroy cleans up the scene resources.
//
// Purpose: Acts as a hook for destruction logic when the SceneManager removes the scene.
func (s *Scene) Destroy() {}

// AddObject inserts a new game object into the physical map.
//
// Purpose: Simplifies the process of registering entities in the active game world.
//
// Inputs:
// - obj (IObject): The object to add.
func (s *Scene) AddObject(obj IObject) {
	s.map_.AddObject(obj)
}

// AddGUIObject inserts a new object into the GUI map.
//
// Purpose: Registers a UI entity. If the GUI map does not exist, it is created automatically with dimensions matching the camera's viewport.
//
// Inputs:
// - obj (IObject): The GUI object to add.
func (s *Scene) AddGUIObject(obj IObject) {
	if s.guiMap == nil {
		s.guiMap = NewGUIMap(s.input, s.camera.Width(), s.camera.Height())
	}
	s.guiMap.AddObject(obj)
}

// World returns the ECS world belonging to the physical map.
//
// Purpose: Exposes the primary ECS instance for querying game entities.
//
// Outputs:
// - donburi.World: The physical map's world.
func (s *Scene) World() donburi.World {
	return s.map_.World()
}

// GetMap retrieves the physical map instance.
//
// Outputs:
// - IMap: The interface to the physical map.
func (s *Scene) GetMap() IMap { return s.map_ }

// GetGUIMap retrieves the GUI map instance, if it exists.
//
// Outputs:
// - IMap: The interface to the GUI map, or nil if uninitialized.
func (s *Scene) GetGUIMap() IMap {
	if s.guiMap == nil {
		return nil
	}
	return s.guiMap
}

// GetCamera retrieves the scene's camera.
//
// Outputs:
// - ICamera: The camera assigned to this scene.
func (s *Scene) GetCamera() ICamera { return s.camera }

// SetGUIMap manually initializes a new GUI map with specified dimensions.
//
// Purpose: Allows explicit configuration of a GUI map rather than letting AddGUIObject auto-create it with camera bounds.
//
// Inputs:
// - input (domain.IInputManager): The input manager to pass to the GUI map.
// - mapW (int): The width of the GUI map in pixels.
// - mapH (int): The height of the GUI map in pixels.
func (s *Scene) SetGUIMap(input domain.IInputManager, mapW, mapH int) {
	s.guiMap = NewMap(input, mapW, mapH)
}

// setID assigns a unique identifier to the scene.
//
// Purpose: Internally stores the ID used by the SceneManager to track this scene in the waitlist.
//
// Inputs:
// - id (string): The identifier string.
func (s *Scene) setID(id string) {
	s.id = id
}

// GetID retrieves the scene's identifier.
//
// Outputs:
// - string: The ID assigned to this scene.
func (s *Scene) GetID() string {
	return s.id
}
