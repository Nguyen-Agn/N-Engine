package napi

import (
	"autoworld/domain"
	"autoworld/modules/components"
	"autoworld/modules/core"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

type Engine = core.Engine

type GameConfig = core.GameConfig

type sceneType = core.Scene

// NewGame creates a new core.Engine instance based on the provided configuration.
//
// Purpose: Used indirectly via napi.Game.Init() to instantiate the game engine.
//
// Inputs:
// - cfg (GameConfig): Configuration parameters for the game engine.
//
// Outputs:
// - *Engine: The newly created core engine instance.
func NewGame(cfg GameConfig) *Engine {
	return core.NewGame(cfg)
}

// ————————————————————————————————————————————————————————————————————————————————
// Re-export cÃ¡c interface tá»« domain — game code khÃ´ng cáº§n import domain trá»±c tiáº¿p.

type IEngine = domain.Engine
type IScene = domain.IScene
type IMap = domain.IMap
type ICamera = domain.ICamera
type IObject = domain.IObject
type ISpriteLW = domain.ISpriteLW
type IAudioLW = domain.IAudioLW
type IGlobal = domain.IGlobal

// Re-export cÃ¡c Component Data struct dÃ¹ng trong ECS.

type positionData = domain.PositionData
type spriteData = domain.SpriteData
type boxData = domain.BoxData
type audioData = domain.AudioData
type backgroundData = domain.BackgroundData
type tilemapData = domain.TilemapData
type boxShape = domain.BoxShape
type drawData = domain.DrawData

const (
	BSRectangle = domain.BSRectangle
	BSCircle    = domain.BSCircle
)

// IDraw is the interface Objects implement so the engine calls Draw() each frame.
// Combine with the Drw mixin (token "drw") to get drawing methods (Rect, Circle, Text...).
type iDraw = domain.IDraw

// iDrawComponent defines the drawing primitives exposed by DrawComponent.
type iDrawComponent = domain.IDrawComponent

var (
	position   = enginetype.Position
	sprite     = enginetype.Sprite
	box_       = enginetype.Box
	audio      = enginetype.Audio
	infor      = enginetype.Infor
	direction  = enginetype.Direction
	background = enginetype.Background
	tilemap    = enginetype.Tilemap
	alarm      = enginetype.Alarm
	velocity   = enginetype.Velocity
	tween      = enginetype.Tween
)

type pos = components.PositionComponent
type box = components.BoxComponent
type aud = components.AudioComponent
type spr = components.SpriteComponent
type info = components.InforComponent
type dir = components.DirectionComponent
type inp = components.InputComponent
type back = components.BackgroundComponent
type tile = components.TilemapComponent
type alrm = components.AlarmComponent
type velo = components.VelocityComponent
type twn = components.TweenComponent
type Object = domain.IObject
type col = components.CollisionComponent
type drw = components.DrawComponent
type deb = components.DebugComponent

type GenericComponent[T any] = components.GenericComponent[T]

// NewGenericComponent creates a strongly typed wrapper around a Donburi component.
//
// Purpose: Allows component data to be retrieved or modified with type safety during gameplay.
//
// Inputs:
// - comp (*donburi.ComponentType[T]): The ECS component type.
//
// Outputs:
// - components.GenericComponent[T]: A strongly typed generic component struct.
func NewGenericComponent[T any](comp *donburi.ComponentType[T]) components.GenericComponent[T] {
	return components.NewGenericComponent(comp)
}
