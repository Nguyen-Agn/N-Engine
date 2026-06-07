package napi

import (
	"autoworld/domain"
	"autoworld/modules/components"
	"autoworld/modules/core"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// â”€â”€â”€ Core Concrete Types â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// Re-export cÃ¡c struct cá»¥ thá»ƒ tá»« core â€” dÃ¹ng trá»±c tiáº¿p trong game code.

// Engine lÃ  *core.Engine â€” struct trung tÃ¢m chá»©a Scene, Store, AudioCtx.
type Engine = core.Engine

// GameConfig lÃ  cáº¥u hÃ¬nh khá»Ÿi Ä‘á»™ng: Title, Width, Height, SampleRate.
type GameConfig = core.GameConfig

// SceneType lÃ  *core.SceneType â€” má»™t mÃ n chÆ¡i Ä‘á»™c láº­p vá»›i world donburi riÃªng.
type SceneType = core.Scene

// NewGame khá»Ÿi táº¡o Engine tá»« GameConfig. ThÆ°á»ng Ä‘Æ°á»£c gá»i giÃ¡n tiáº¿p qua napi.Init().
func NewGame(cfg GameConfig) *Engine {
	return core.NewGame(cfg)
}

// â”€â”€â”€ Domain Interface Aliases â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// Re-export cÃ¡c interface tá»« domain â€” game code khÃ´ng cáº§n import domain trá»±c tiáº¿p.

type IEngine = domain.IEngine
type IScene = domain.IScene
type IMap = domain.IMap
type ICamera = domain.ICamera
type IObject = domain.IObject
type ISpriteLW = domain.ISpriteLW
type IAudioLW = domain.IAudioLW
type IGlobal = domain.IGlobal

// Re-export cÃ¡c Component Data struct dÃ¹ng trong ECS.

type PositionData = domain.PositionData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type AudioData = domain.AudioData
type BackgroundData = domain.BackgroundData
type TilemapData = domain.TilemapData
type BoxShape = domain.BoxShape
type DrawData = domain.DrawData

const (
	BSRectangle = domain.BSRectangle
	BSCircle    = domain.BSCircle
)

// IDraw is the interface Objects implement so the engine calls Draw() each frame.
// Combine with the Drw mixin (token "drw") to get drawing methods (Rect, Circle, Text...).
type IDraw = domain.IDraw

// IDrawComponent defines the drawing primitives exposed by DrawComponent.
type IDrawComponent = domain.IDrawComponent

var (
	Position   = enginetype.Position
	Sprite     = enginetype.Sprite
	Box_       = enginetype.Box
	Audio      = enginetype.Audio
	Infor      = enginetype.Infor
	Direction  = enginetype.Direction
	Background = enginetype.Background
	Tilemap    = enginetype.Tilemap
	Alarm      = enginetype.Alarm
	Velocity   = enginetype.Velocity
	Tween      = enginetype.Tween
)

// â”€â”€â”€ Component Mixin Aliases â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
// Re-export cÃ¡c Component Mixin struct tá»« modules/components.
// NhÃºng (embed) cÃ¡c struct nÃ y vÃ o Custom Object Ä‘á»ƒ nháº­n Ä‘áº§y Ä‘á»§ getter/setter.
// VÃ­ dá»¥: type Player struct { napi.IObject; napi.IPosition; napi.ISprite }

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

func NewGenericComponent[T any](comp *donburi.ComponentType[T]) components.GenericComponent[T] {
	return components.NewGenericComponent(comp)
}
