package core

import (
	"autoworld/domain"
	"autoworld/modules/nsystem"

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
	drawRegistry domain.IDrawObjectRegistry
	world        donburi.World
	objectList   []IObject
	pendingRemove []IObject
	width        int // map width in pixels; 0 = unbounded
	height       int // map height in pixels; 0 = unbounded
}

// NewMap creates a Map with the given bounds.
// Pass width=0, height=0 for an unbounded scrolling map.
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

// NewGUIMap creates a GUI Map for screen-space HUD/overlay.
// GUI Maps share no audio system — audio belongs to the Physical Map.
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

// SetDrawRegistry injects a DrawSystem reference so AddObject can automatically
// register IDraw objects for per-frame Draw() calls.
// Called by Scene after creating both Map and Camera.
func (m *Map) SetDrawRegistry(r domain.IDrawObjectRegistry) {
	m.drawRegistry = r
}

// Update runs all systems every frame in order:
// Logic → Input → Alarm → Tween → Velocity → Collision → Audio → flush removes.
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

// AddObject registers an IObject into the Map for per-frame updates.
// LogicSystem will call OnCreate() on the next frame.
// If the object implements IDraw and drawRegistry is set, it is also registered for rendering.
func (m *Map) AddObject(obj IObject) {
	m.logicSystem.AddObjectCreated(obj)
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

// RemoveObject schedules an IObject for deferred removal at the end of the current frame.
// MarkDead() is called immediately so Collector and Applier can skip this object.
// OnDestroy() will be invoked by LogicSystem at the start of the next frame.
func (m *Map) RemoveObject(obj IObject) {
	if dead, ok := obj.(interface{ IsDead() bool; MarkDead() }); ok {
		if dead.IsDead() {
			return // already queued
		}
		dead.MarkDead()
	}
	m.pendingRemove = append(m.pendingRemove, obj)
}

// flushRemove processes the pendingRemove queue:
//  1. Calls OnDestroy() via LogicSystem.
//  2. Removes from drawRegistry (if IDraw).
//  3. Filters objectList.
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

// World returns the donburi.World owned by this Map.
func (m *Map) World() donburi.World {
	return m.world
}

// Width returns the map width in pixels. 0 means unbounded.
func (m *Map) Width() int {
	return m.width
}

// Height returns the map height in pixels. 0 means unbounded.
func (m *Map) Height() int {
	return m.height
}

// GetObjects returns all registered objects.
func (m *Map) GetObjects() []IObject {
	return m.objectList
}
