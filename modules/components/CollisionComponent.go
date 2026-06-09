package components

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

type CollisionComponent struct {
	domain.IObject
	data *CollisionData
}

var Collision = enginetype.Collision

// BindComponent binds the base object and its ECS data to this component.
// Inputs:
//   - base: The base IObject to bind to.
func (p *CollisionComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Collision)
}

// init initializes the default data for the collision component.
// It registers the "col" component token and sets the object to be collidable by default.
func init() {
	enginetype.RegisterComponentInitializer("col", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Collision, domain.CollisionData{
			Handlers:     make(map[uint64]func(other IObject)),
			IsCollidable: true,
		})
	})
}

// OnCollisionTag registers a callback that will be invoked every frame when colliding with an object that has the specified tag.
// Inputs:
//   - tag: The tag string to detect collisions with.
//   - handler: The function to call when a collision occurs, receiving the other object as a parameter.
func (p CollisionComponent) OnCollisionTag(tag string, handler func(other IObject)) {

	if p.data != nil {
		hash := enginetype.HashString(tag)
		p.data.Handlers[hash] = handler
	}
}

// IsCollidable checks whether the object is currently active for collision detection.
// Outputs: Returns true if it can collide, false otherwise.
func (p CollisionComponent) IsCollidable() bool {
	if p.data == nil {
		return false
	}
	return p.data.IsCollidable
}

// SetIsCollidable enables or disables collision detection for the object.
// Inputs:
//   - isCollidable: Boolean indicating whether collisions should be processed.
func (p CollisionComponent) SetIsCollidable(isCollidable bool) {
	if p.data != nil {
		p.data.IsCollidable = isCollidable
	}
}
