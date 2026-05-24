package components

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Velocity = enginetype.Velocity

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

// VelocityComponent là Mixin để nhúng vào Custom Object.
type VelocityComponent struct {
	IObject
	data *domain.VelocityData
}

func (v *VelocityComponent) BindComponent(base IObject) {
	v.IObject = base
	v.data = enginetype.GetComponent(base, Velocity)
}

func (v VelocityComponent) VelocityX() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.Vx
}

func (v VelocityComponent) VelocityY() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.Vy
}

func (v VelocityComponent) SetVelocityX(vx float32) {
	if v.data != nil {
		v.data.Vx = vx
	}
}

func (v VelocityComponent) SetVelocityY(vy float32) {
	if v.data != nil {
		v.data.Vy = vy
	}
}

func (v VelocityComponent) SetVelocity(vx, vy float32) {
	if v.data != nil {
		v.data.Vx = vx
		v.data.Vy = vy
	}
}

func (v VelocityComponent) AddVelocity(vx, vy float32) {
	if v.data != nil {
		v.data.Vx += vx
		v.data.Vy += vy
	}
}

func (v VelocityComponent) Friction() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.Friction
}

func (v VelocityComponent) SetFriction(f float32) {
	if v.data != nil {
		v.data.Friction = f
	}
}

func (v VelocityComponent) MaxSpeed() float32 {
	if v.data == nil {
		return 0
	}
	return v.data.MaxSpeed
}

func (v VelocityComponent) SetMaxSpeed(speed float32) {
	if v.data != nil {
		v.data.MaxSpeed = speed
	}
}
