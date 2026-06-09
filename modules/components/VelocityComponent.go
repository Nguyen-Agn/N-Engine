package components

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Velocity = enginetype.Velocity

// init registers the Velocity Component initializer with default values.
func init() {
	enginetype.RegisterComponentInitializer("vel", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Velocity, domain.VelocityData{
			Vx:       0,
			Vy:       0,
			Friction: 0,
			MaxSpeed: 0,
		})
	})
}

// VelocityComponent is a mixin to embed into Custom Objects.
type VelocityComponent struct {
	IObject
	data *domain.VelocityData
}

// BindComponent binds the base object and retrieves the VelocityData from the ECS.
// Inputs: base - the base IObject to bind.
func (v *VelocityComponent) BindComponent(base IObject) {
	v.IObject = base
	v.data = enginetype.GetComponent(base, Velocity)
}

// VelocityX returns the current X velocity.
// Outputs: float32 representing the X velocity.
func (v VelocityComponent) VelocityX() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.Vx
}

// VelocityY returns the current Y velocity.
// Outputs: float32 representing the Y velocity.
func (v VelocityComponent) VelocityY() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.Vy
}

// SetVelocityX sets the X velocity.
// Inputs: vx - the new X velocity.
func (v VelocityComponent) SetVelocityX(vx float32) {
	if v.data != nil {
		v.data.Vx = vx
	}
}

// SetVelocityY sets the Y velocity.
// Inputs: vy - the new Y velocity.
func (v VelocityComponent) SetVelocityY(vy float32) {
	if v.data != nil {
		v.data.Vy = vy
	}
}

// SetVelocity sets both X and Y velocities simultaneously.
// Inputs: vx, vy - the new X and Y velocities.
func (v VelocityComponent) SetVelocity(vx, vy float32) {
	if v.data != nil {
		v.data.Vx = vx
		v.data.Vy = vy
	}
}

// AddVelocity adds the specified amounts to the current X and Y velocities.
// Inputs: vx, vy - the amounts to add to the X and Y velocities.
func (v VelocityComponent) AddVelocity(vx, vy float32) {
	if v.data != nil {
		v.data.Vx += vx
		v.data.Vy += vy
	}
}

// Friction returns the current friction value applied to velocity.
// Outputs: float32 representing the friction.
func (v VelocityComponent) Friction() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.Friction
}

// SetFriction sets the friction value applied to velocity over time.
// Inputs: f - the new friction value.
func (v VelocityComponent) SetFriction(f float32) {
	if v.data != nil {
		v.data.Friction = f
	}
}

// MaxSpeed returns the maximum speed limit for the object.
// Outputs: float32 representing the maximum speed.
func (v VelocityComponent) MaxSpeed() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.MaxSpeed
}

// SetMaxSpeed sets the maximum speed limit for the object.
// Inputs: speed - the new maximum speed.
func (v VelocityComponent) SetMaxSpeed(speed float32) {
	if v.data != nil {
		v.data.MaxSpeed = speed
	}
}
