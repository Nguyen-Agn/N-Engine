package core

import (
	"autoworld/domain"

	"github.com/yohamta/donburi"
)

// Scene quáº£n lÃ½ má»™t mÃ n chÆ¡i Ä‘á»™c láº­p trong Engine.
// Má»—i Scene chá»©a:
//   - Physical Map: ECS World + Logic/Audio/Input system + objectList (tá» a Ä‘á»™ map space)
//   - GUI Map: screen-space overlay cho HUD (optional, tá»± táº¡o khi cáº§n)
//   - Camera: viewport, follow-target, DrawSystem
type Scene struct {
	map_   *Map // physical game world
	guiMap *Map // screen-space GUI/HUD (optional, nil náº¿u khÃ´ng dÃ¹ng)
	camera *Camera
	input  domain.IInputManager // lÆ°u Ä‘á»ƒ táº¡o guiMap lazy
	id     string               // ID of the scene
}

// NewScene khá»Ÿi táº¡o Scene má»›i vá»›i Physical Map vÃ  Camera.
// mapW, mapH lÃ  kÃ­ch thÆ°á»›c báº£n Ä‘á»“ (pixel). Truyá» n 0, 0 náº¿u khÃ´ng giá»›i háº¡n.
// viewW, viewH lÃ  kÃ­ch thÆ°á»›c viewport cá»§a Camera (thÆ°á» ng báº±ng screen size).
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

// â”€â”€â”€ IScene interface â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// Update cáº­p nháº­t toÃ n bá»™ Scene má»—i frame.
// Thá»© tá»±: Physical Map â†’ GUI Map (náº¿u cÃ³) â†’ Camera follow.
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

// Draw render Scene lÃªn mÃ n hÃ¬nh qua Camera.
// Camera váº½ Physical Map (cÃ³ camera offset) rá»“i GUI Map (khÃ´ng offset, Ä‘Ã¨ lÃªn trÃªn).
func (s *Scene) Draw() {
	if s.guiMap != nil {
		s.camera.Draw(s.map_.World(), s.guiMap.World())
	} else {
		s.camera.Draw(s.map_.World(), nil)
	}
}

// Destroy Ä‘Æ°á»£c gá»i khi Scene bá»‹ xÃ³a khá»i SceneManager.
func (s *Scene) Destroy() {}

// AddObject Ä‘Äƒng kÃ½ IObject vÃ o Physical Map.
func (s *Scene) AddObject(obj IObject) {
	s.map_.AddObject(obj)
}

// AddGUIObject Ä‘Äƒng kÃ½ IObject vÃ o GUI Map (screen-space, khÃ´ng camera offset).
// GUI Map Ä‘Æ°á»£c táº¡o tá»± Ä‘á»™ng náº¿u chÆ°a tá»“n táº¡i.
func (s *Scene) AddGUIObject(obj IObject) {
	if s.guiMap == nil {
		s.guiMap = NewGUIMap(s.input, s.camera.Width(), s.camera.Height())
	}
	s.guiMap.AddObject(obj)
}

// World tráº£ vá» donburi.World cá»§a Physical Map. TÆ°Æ¡ng Ä‘Æ°Æ¡ng GetMap().World().
func (s *Scene) World() donburi.World {
	return s.map_.World()
}

// â”€â”€â”€ Getters â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// GetMap tráº£ vá» Physical Map cá»§a Scene.
func (s *Scene) GetMap() IMap { return s.map_ }

// GetGUIMap tráº£ vá» GUI Map cá»§a Scene. Tráº£ vá» nil náº¿u chÆ°a khá»Ÿi táº¡o.
func (s *Scene) GetGUIMap() IMap {
	if s.guiMap == nil {
		return nil
	}
	return s.guiMap
}

// GetCamera tráº£ vá» Camera cá»§a Scene.
func (s *Scene) GetCamera() ICamera { return s.camera }

// â”€â”€â”€ Setters â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// SetGuiMap
func (s *Scene) SetGUIMap(input domain.IInputManager, mapW, mapH int) {
	s.guiMap = NewMap(input, mapW, mapH)
}

// setID sets the scene's ID internally
func (s *Scene) setID(id string) {
	s.id = id
}

// GetID returns the ID of the scene
func (s *Scene) GetID() string {
	return s.id
}

