package nsystem

import (
	"autoworld/modules/enginetype"
)

// PhysicsSystem calculates basic physics for entities: friction, max speed limiting, and applying velocity to position.
type PhysicsSystem struct {
}

// NewVelocitySystem creates and returns a new instance of PhysicsSystem.
// Outputs: Returns a pointer to a newly initialized PhysicsSystem.
func NewVelocitySystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

// Update applies physics logic to all objects with a Velocity component.
// Inputs: objectList ([]IObject) - The active list of objects in the scene.
// Purpose: For each object with Velocity, it applies friction to slow down movement, limits velocity components to MaxSpeed, and updates the object's Position based on the final velocity.
func (s *PhysicsSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		velData := enginetype.GetComponent(obj, enginetype.Velocity)
		if velData == nil {
			continue
		}

		// Apply Friction
		if velData.Friction > 0 {
			// Reduce speed based on friction
			if velData.Vx > 0 {
				velData.Vx -= velData.Friction
				if velData.Vx < 0 {
					velData.Vx = 0
				}
			} else if velData.Vx < 0 {
				velData.Vx += velData.Friction
				if velData.Vx > 0 {
					velData.Vx = 0
				}
			}

			if velData.Vy > 0 {
				velData.Vy -= velData.Friction
				if velData.Vy < 0 {
					velData.Vy = 0
				}
			} else if velData.Vy < 0 {
				velData.Vy += velData.Friction
				if velData.Vy > 0 {
					velData.Vy = 0
				}
			}
		}

		// Apply MaxSpeed
		if velData.MaxSpeed > 0 {
			// Simplified: limit each axis independently (could use vector magnitude for accuracy)
			if velData.Vx > velData.MaxSpeed {
				velData.Vx = velData.MaxSpeed
			}
			if velData.Vx < -velData.MaxSpeed {
				velData.Vx = -velData.MaxSpeed
			}
			if velData.Vy > velData.MaxSpeed {
				velData.Vy = velData.MaxSpeed
			}
			if velData.Vy < -velData.MaxSpeed {
				velData.Vy = -velData.MaxSpeed
			}
		}

		// Apply Velocity to Position
		posData := enginetype.GetComponent(obj, enginetype.Position)
		if posData != nil {
			posData.X += velData.Vx
			posData.Y += velData.Vy
		}
	}
}
