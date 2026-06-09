package ncom

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/components"

	"github.com/yohamta/donburi"
)

// Object represents the base interface that game objects must implement.
// Token: No
// Provide: (you can override them)
// - OnCreate(): called once after object is created
// - OnDestroy(): called once after object is destroyed
// - OnStep(): called every logic frame
// - OnSave(map[string]any): called when Game.SaveGame is executed
// - OnLoad(map[string]any): called when Game.LoadGame is executed
// - SetTokens(string): Convenient function to set ECS properties via string classes
type Object = domain.IObject

// Pos is the position component mixin.
// Token: "pos"
// Provide: X(), Y(), SetX(), SetY()
type Pos = components.PositionComponent

// Spr is the sprite component mixin.
// Token: "spr"
// Provide: AddSprite(), SetCurrentSprite(), SetScaleX/Y(), SetImageSpeed()...
type Spr = components.SpriteComponent

// Box is the collision hitbox component mixin.
// Token: "box" (automatically added when using "col")
// Provide: BoxW(), BoxH(), SetBoxW/H(), IsCollidable()...
type Box = components.BoxComponent

// Aud is the audio component mixin.
// Token: "aud"
// Provide: Play(), Stop(), SetVolume(), SetPitch()...
type Aud = components.AudioComponent

// Info is the identity component mixin.
// Token: "info" (automatically added to all objects)
// Provide: GetName(), GetId(), SaveTag(), SetSaveTag(), IsDead(), MarkDead()...
type Info = components.InforComponent

// Dir is the direction component mixin.
// Token: "dir"
// Provide: Direction(), SetDirection(), Rotate()
type Dir = components.DirectionComponent

// Inp is the keyboard input component mixin.
// Token: "inp"
// Provide: ListenOn(key, eventType, handler)
//   - handler receives the triggered key name (e.g. "w", "space", "a")
//   - eventType: domain.EventPressed / EventJustPressed / EventJustReleased
type Inp = components.InputComponent

// Mouse is the mouse input component mixin.
// Token: none (does not store ECS data)
// Provide:
//   - MouseX() int, MouseY() int         — cursor coordinates (pixel)
//   - WheelX() float64, WheelY() float64 — scroll wheel speed X/Y
//   - ListenMouseOn(button, eventType, handler)
//     button: "left", "right", "middle"
//     eventType: domain.EventPressed / EventJustPressed / EventJustReleased
//     handler receives the triggered mouse button name (e.g. "left")
type Mouse = components.MouseComponent

// Back is the background component mixin.
// Token: "back"
// Provide: SetBackground(), SetRepeat(), SetScrollSpeed()...
type Back = components.BackgroundComponent

// Tile is the tilemap component mixin.
// Token: "til"
// Provide: SetTilemap(), SetGrid()...
type Tile = components.TilemapComponent

// Alrm is the alarm component mixin.
// Token: "alr"
// Provide: SetAlarm(id, frames, callback), GetAlarm(id), StopAlarm(id)
type Alrm = components.AlarmComponent

// Velo is the velocity component mixin.
// Token: "vel"
// Provide: VelocityX/Y(), SetVelocity(), Friction(), MaxSpeed()...
type Velo = components.VelocityComponent

// Twn is the tweening component mixin.
// Token: "twn" (automatically adds "spr")
// Provide: TweenMove(), TweenScale(), TweenAlpha()
type Twn = components.TweenComponent

// Col is the collision component mixin.
// Token: "col" (automatically adds "box")
// Provide: AddTag(), OnCollisionTag(tag, handler)
type Col = components.CollisionComponent

// Drw is the draw component mixin.
// Token: "drw" (automatically adds "pos")
// Provide: Rect(), Circle(), Line(), Text(), Image()...
// Note: The object must also implement the IDraw interface (Draw() method).
type Drw = components.DrawComponent

// Deb is the debug component mixin.
// Token: "deb"
// Provide: Log(), DrawDebug()...
type Deb = components.DebugComponent

// Generic is a base component for creating custom components with user-defined data.
// Intended to be combined with napi.NewComponentType[T]().
//
// Example:
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//
//	type StatsComponent struct {
//	    ncom.Generic[StatsData]
//	}
//
//	func (s *StatsComponent) HP() int { return s.Get().Health }
type Generic[T any] = components.GenericComponent[T]

// NewGeneric creates a Generic component bounded to the specified ComponentType.
//
// Purpose: Must be called when instantiating the struct, before calling napi.Obj.NewObject().
//
// Inputs:
// - comp (*donburi.ComponentType[T]): The component type to bind.
//
// Outputs:
// - components.GenericComponent[T]: The bound generic component.
func NewGeneric[T any](comp *donburi.ComponentType[T]) components.GenericComponent[T] {
	return components.NewGenericComponent(comp)
}
