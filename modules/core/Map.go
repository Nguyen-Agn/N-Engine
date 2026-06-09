package core

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/nsystem"

	"github.com/yohamta/donburi"
)

// Map manages the ECS World and all per-frame logic for one game map.
// Map does not know about rendering — Camera owns the DrawSystem.
//
// Two map types exist per Scene:
//   - Physical Map: game world in map space, supports camera scrolling.
//   - GUI Map: screen-space HUD overlay, no camera offset.
type Map struct {
	logicSystem     ILogicSystem
	audioSystem     IAudioSystem
	inputSystem     IInputSystem
	alarmSystem     IAlarmSystem
	tweenSystem     ITweenSystem
	velocitySystem  IVelocitySystem
	collisionSystem domain.IUpdateSystem
	// drawRegistry is an optional reference to DrawSystem.
	// When set, AddObject auto-registers IDraw objects for rendering.
	drawRegistry  domain.IDrawObjectRegistry
	world         donburi.World
	objectList    []IObject
	pendingRemove []IObject
	width         int // map width in pixels; 0 = unbounded
	height        int // map height in pixels; 0 = unbounded
}

// NewMap creates a Map instance with the specified dimensions.
//
// Purpose: Initializes a physical map to manage the ECS World and per-frame logic, with a given boundary.
//
// Inputs:
// - input (domain.IInputManager): The input manager to pass to the InputSystem.
// - width (int): The width of the map in pixels (0 for unbounded).
// - height (int): The height of the map in pixels (0 for unbounded).
//
// Outputs:
// - *Map: A new map instance configured for physical world space.
func NewMap(input domain.IInputManager, width, height int) *Map {
	return &Map{
		logicSystem:     nsystem.NewLogicSystem(),
		audioSystem:     nsystem.NewAudioSystem(),
		alarmSystem:     nsystem.NewAlarmSystem(),
		tweenSystem:     nsystem.NewTweenSystem(),
		velocitySystem:  nsystem.NewVelocitySystem(),
		collisionSystem: nsystem.NewCollisionSystem(),
		inputSystem:     nsystem.NewInputSystem(input),
		world:           donburi.NewWorld(),
		width:           width,
		height:          height,
	}
}

// NewGUIMap creates a Map instance specifically for screen-space UI elements.
//
// Purpose: Initializes a GUI map that operates in screen coordinates, independent of the camera's physical offset. Audio systems are skipped or distinct.
//
// Inputs:
// - input (domain.IInputManager): The input manager to pass to the InputSystem.
// - viewW (int): The width of the viewport.
// - viewH (int): The height of the viewport.
//
// Outputs:
// - *Map: A new map instance configured for GUI world space.
func NewGUIMap(input domain.IInputManager, viewW, viewH int) *Map {
	return &Map{
		logicSystem:     nsystem.NewLogicSystem(),
		audioSystem:     nsystem.NewAudioSystem(),
		alarmSystem:     nsystem.NewAlarmSystem(),
		tweenSystem:     nsystem.NewTweenSystem(),
		velocitySystem:  nsystem.NewVelocitySystem(),
		collisionSystem: nsystem.NewCollisionSystem(),
		inputSystem:     nsystem.NewInputSystem(input),
		world:           donburi.NewWorld(),
		width:           viewW,
		height:          viewH,
	}
}

// SetDrawRegistry assigns the drawing registry interface.
//
// Purpose: Connects the DrawSystem with the map so that when IDraw objects are added, they are automatically registered for rendering.
//
// Inputs:
// - r (domain.IDrawObjectRegistry): The drawing registry to use for rendering registration.
func (m *Map) SetDrawRegistry(r domain.IDrawObjectRegistry) {
	m.drawRegistry = r
}

// Update runs all systems sequentially for the current frame.
//
// Purpose: Executes the logic loop in order: Logic -> Input -> Alarm -> Tween -> Velocity -> Collision -> Audio, and finally flushes queued removals.
//
// Outputs:
// - error: Always returns nil; satisfies standard update signatures.
func (m *Map) Update() error {
	// LogicSystem runs first: Create → StepUpdate → Destroy
	m.logicSystem.Update(m.objectList)

	// Event-driven systems (may fire callbacks activated by logic)
	m.inputSystem.Update(m.objectList)
	m.alarmSystem.Update(m.objectList)

	// Support systems — run after primary logic is settled
	m.tweenSystem.Update(m.objectList)
	m.velocitySystem.Update(m.objectList)
	m.collisionSystem.Update(m.objectList)

	// Audio — no ordering dependency
	m.audioSystem.Update(m.world)

	// Deferred remove: clean up objects queued for removal this frame
	m.flushRemove()
	return nil
}

// AddObject inserts an object into the Map for per-frame updates.
//
// Purpose: Adds an object to the logic system, the collision system, the main object list, and automatically registers it for rendering if it supports drawing.
//
// Inputs:
// - obj (IObject): The game object to add to the map.
func (m *Map) AddObject(obj IObject) {
	// Absoulutly added to logicSystem
	m.logicSystem.AddObjectCreated(obj)

	// Auto-register objects that join collisionSystem
	if sys, ok := m.collisionSystem.(interface{ AddObject(IObject) }); ok {
		sys.AddObject(obj)
	}
	m.objectList = append(m.objectList, obj)

	// Auto-register objects that implement IDraw into DrawSystem
	if m.drawRegistry != nil {
		if _, ok := obj.(domain.IDraw); ok {
			m.drawRegistry.AddDrawObject(obj)
		}
	}
}

// RemoveObject schedules an object for deferred removal at the end of the current frame.
//
// Purpose: Safely queues an object for deletion and immediately marks it dead so it won't be processed further in the current step.
//
// Inputs:
// - obj (IObject): The game object to remove.
func (m *Map) RemoveObject(obj IObject) {
	if dead, ok := obj.(interface {
		IsDead() bool
		MarkDead()
	}); ok {
		if dead.IsDead() {
			return // already queued
		}
		dead.MarkDead()
	}
	m.pendingRemove = append(m.pendingRemove, obj)
}

// flushRemove processes the queue of objects scheduled for removal.
//
// Purpose: Internally cleans up dead objects at the end of a frame, either moving them to a pool or destroying them entirely from the ECS world and DrawRegistry.
func (m *Map) flushRemove() {
	if len(m.pendingRemove) == 0 {
		return
	}

	// Build a set for O(1) lookup
	removeSet := make(map[IObject]bool, len(m.pendingRemove))
	for _, obj := range m.pendingRemove {
		removeSet[obj] = true

		pool := obj.GetPool()
		isPooled := false
		if pool != nil {
			isPooled = pool.Put(obj) // Try to pool. Returns false if pool is full.
		}

		if isPooled {
			// Auto-Routing: Object is hibernated into Pool
			// Remove from drawRegistry but skip OnDestroy and DO NOT remove from ECS world.
			if m.drawRegistry != nil {
				if _, ok := obj.(domain.IDraw); ok {
					m.drawRegistry.RemoveDrawObject(obj)
				}
			}
		} else {
			// True Deletion: No pool, or pool is full
			m.logicSystem.AddObjectDestroy(obj)
			if m.drawRegistry != nil {
				if _, ok := obj.(domain.IDraw); ok {
					m.drawRegistry.RemoveDrawObject(obj)
				}
			}
			// Clean up Donburi ECS Entity to prevent memory leak
			if obj.Entry() != nil {
				m.world.Remove(obj.Entry().Entity())
			}
		}
	}

	// Filter objectList in-place
	filtered := m.objectList[:0]
	for _, obj := range m.objectList {
		if !removeSet[obj] {
			filtered = append(filtered, obj)
		}
	}
	m.objectList = filtered
	m.pendingRemove = m.pendingRemove[:0]
}

// World returns the ECS world managed by this map.
//
// Purpose: Provides access to the underlying donburi.World for querying components or entities.
//
// Outputs:
// - donburi.World: The ECS world instance.
func (m *Map) World() donburi.World {
	return m.world
}

// Width returns the horizontal boundary of the map.
//
// Purpose: Gets the map's width.
//
// Outputs:
// - int: Width in pixels. 0 if unbounded.
func (m *Map) Width() int {
	return m.width
}

// Height returns the vertical boundary of the map.
//
// Purpose: Gets the map's height.
//
// Outputs:
// - int: Height in pixels. 0 if unbounded.
func (m *Map) Height() int {
	return m.height
}

// GetObjects returns all objects currently active in the map.
//
// Purpose: Retrieves the list of all registered objects for processing or inspection.
//
// Outputs:
// - []IObject: A slice of all game objects in the map.
func (m *Map) GetObjects() []IObject {
	return m.objectList
}
