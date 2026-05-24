package napi

import (
	"autoworld/domain"
	"autoworld/modules/components"
	"autoworld/modules/core"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// ─── Core Concrete Types ──────────────────────────────────────────────────────
// Re-export các struct cụ thể từ core — dùng trực tiếp trong game code.

// Engine là *core.Engine — struct trung tâm chứa Scene, Store, AudioCtx.
type Engine = core.Engine

// GameConfig là cấu hình khởi động: Title, Width, Height, SampleRate.
type GameConfig = core.GameConfig

// Scene là *core.Scene — một màn chơi độc lập với world donburi riêng.
type Scene = core.Scene

// NewGame khởi tạo Engine từ GameConfig. Thường được gọi gián tiếp qua napi.Init().
func NewGame(cfg GameConfig) *Engine {
	return core.NewGame(cfg)
}

// ─── Domain Interface Aliases ─────────────────────────────────────────────────
// Re-export các interface từ domain — game code không cần import domain trực tiếp.

type IEngine = domain.IEngine
type IScene = domain.IScene
type IMap = domain.IMap
type ICamera = domain.ICamera
type IObject = domain.IObject
type ISpriteLW = domain.ISpriteLW
type IAudioLW = domain.IAudioLW
type IGlobal = domain.IGlobal

// ─── ECS Data Struct Aliases ──────────────────────────────────────────────────
// Re-export các Component Data struct dùng trong ECS.

type PositionData = domain.PositionData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type AudioData = domain.AudioData
type BackgroundData = domain.BackgroundData
type TilemapData = domain.TilemapData
type BoxShape = domain.BoxShape

const (
	BSRectangle = domain.BSRectangle
	BSCircle    = domain.BSCircle
)

// ─── ECS Component Type Variables ────────────────────────────────────────────
// Trỏ thẳng đến biến static trong enginetype — đảm bảo chung ID với toàn hệ thống.

var (
	Position   = enginetype.Position
	Sprite     = enginetype.Sprite
	Box        = enginetype.Box
	Audio      = enginetype.Audio
	Infor      = enginetype.Infor
	Direction  = enginetype.Direction
	Background = enginetype.Background
	Tilemap    = enginetype.Tilemap
	Alarm      = enginetype.Alarm
	Velocity   = enginetype.Velocity
	Tween      = enginetype.Tween
)

// ─── Component Mixin Aliases ──────────────────────────────────────────────────
// Re-export các Component Mixin struct từ modules/components.
// Nhúng (embed) các struct này vào Custom Object để nhận đầy đủ getter/setter.
// Ví dụ: type Player struct { napi.IObject; napi.IPosition; napi.ISprite }

type IPosition = components.PositionComponent
type IBox = components.BoxComponent
type IAudio = components.AudioComponent
type ISprite = components.SpriteComponent
type IInfor = components.InforComponent
type IDirection = components.DirectionComponent
type IInput = components.InputComponent
type IBackground = components.BackgroundComponent
type ITilemap = components.TilemapComponent
type IAlarm = components.AlarmComponent
type IVelocity = components.VelocityComponent
type ITween = components.TweenComponent

// ─── Component Struct Aliases ─────────────────────────────────────────────────
// Cho phép game code dùng component struct trực tiếp mà không cần import modules/components.

type PositionComponent = components.PositionComponent
type BoxComponent = components.BoxComponent
type AudioComponent = components.AudioComponent
type SpriteComponent = components.SpriteComponent
type InforComponent = components.InforComponent
type DirectionComponent = components.DirectionComponent
type BackgroundComponent = components.BackgroundComponent
type TilemapComponent = components.TilemapComponent
type AlarmComponent = components.AlarmComponent
type VelocityComponent = components.VelocityComponent
type TweenComponent = components.TweenComponent

// ─── Custom Component (Generic Mixin) ────────────────────────────────────────
// GenericComponent[T] là mixin generic để game dev tự tạo Component với method riêng.
// Kết hợp với NewComponentType[T] để định nghĩa và sử dụng custom component
// mà không cần import bất kỳ package engine nào khác.
//
// Workflow:
//
//	// Bước 1 — Khai báo data và ComponentType (1 lần, cấp package):
//	type StatsData struct { Health int; Speed float32 }
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//
//	// Bước 2 — (Tùy chọn) Bọc GenericComponent để thêm method:
//	type StatsComponent struct {
//	    napi.GenericComponent[StatsData]
//	}
//	func (s *StatsComponent) TakeDamage(n int) { s.Get().Health -= n }
//	func (s *StatsComponent) IsAlive() bool    { return s.Get().Health > 0 }
//
//	// Bước 3 — Nhúng vào Custom Object:
//	type Hero struct {
//	    napi.IObject
//	    napi.IPosition
//	    StatsComponent
//	}
//
//	// Bước 4 — Khởi tạo:
//	func NewHero() *Hero {
//	    h := &Hero{
//	        StatsComponent: StatsComponent{
//	            GenericComponent: napi.NewGenericComponent(StatsComp),
//	        },
//	    }
//	    napi.NewObject(h, "hero", "pos sta")
//	    napi.Register(h, "")
//	    return h
//	}
type GenericComponent[T any] = components.GenericComponent[T]

// NewGenericComponent tạo GenericComponent đã gắn sẵn ComponentType.
// Phải được gọi khi khởi tạo struct, trước khi gọi napi.NewObject.
//
// Ví dụ:
//
//	h := &Hero{
//	    GenericComponent: napi.NewGenericComponent(StatsComp),
//	}
func NewGenericComponent[T any](comp *donburi.ComponentType[T]) components.GenericComponent[T] {
	return components.NewGenericComponent(comp)
}
