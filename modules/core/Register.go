package core

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"
)

// ─── Domain Interface Aliases ─────────────────────────────────────────────────
// Type alias tập trung: các module khác chỉ cần import "core" thay vì "domain" trực tiếp.
// Điều này giúp cô lập dependency vào domain và dễ thay thế nếu cần.

// Manager interfaces
type ISceneManager = domain.ISceneManager
type IGlobalConfig = domain.IGlobalConfig
type IInputManager = domain.IInputManager
type IObserver = domain.IObserver

// System interfaces
type IScene = domain.IScene
type ILogicSystem = domain.ILogicSystem
type IDrawSystem = domain.IDrawSystem
type IAudioSystem = domain.IAudioSystem
type IInputSystem = domain.IUpdateSystem
type IAlarmSystem = domain.IUpdateSystem
type ITweenSystem = domain.IUpdateSystem
type IVelocitySystem = domain.IUpdateSystem

type IMap = domain.IMap
type ICamera = domain.ICamera

// Object interface
type IObject = domain.IObject

// ─── Domain Data Struct Aliases ───────────────────────────────────────────────
// Alias cho các Component Data struct dùng trong ECS.

type PositionData = domain.PositionData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type AudioData = domain.AudioData
type BackgroundData = domain.BackgroundData
type TilemapData = domain.TilemapData

// ─── Asset & Resource Interface Aliases ──────────────────────────────────────

type ISpriteLW = domain.ISpriteLW
type IAudioLW = domain.IAudioLW
type IGlobal = domain.IGlobal
type IAudioLoader = domain.IAudioLoader
type ISpriteLoader = domain.ISpriteLoader
type IManifestLoader = domain.IManifestLoader

// ─── Misc Type & Constant Aliases ────────────────────────────────────────────

type BoxShape = domain.BoxShape

const (
	BSRectangle = domain.BSRectangle
	BSCircle    = domain.BSCircle
)

// ─── ECS Component Type Variables ────────────────────────────────────────────
// Trỏ thẳng đến biến static trong enginetype để đảm bảo an toàn thứ tự khởi tạo.
// Không tạo ComponentType mới ở đây — chỉ re-export từ enginetype.

var (
	Position   = enginetype.Position
	Sprite     = enginetype.Sprite
	Box        = enginetype.Box
	Audio      = enginetype.Audio
	Input      = enginetype.Input
	Background = enginetype.Background
	Tilemap    = enginetype.Tilemap
)
