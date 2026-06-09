package core

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"
)

// ─── Domain Interface Aliases ─────────────────────────────────────────────────
// Type alias center: to keep 'import' at other files shorter
// Help centrelize on dependency at domain -> usable & change 1 for all

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

// Scene support
type IMap = domain.IMap
type ICamera = domain.ICamera

// Object interface
type IObject = domain.IObject

// ─── Domain Data Struct Aliases ───────────────────────────────────────────────
// Alias for Component Data structs in ECS.

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
type ISaveManager = domain.ISaveManager

// ─── Misc Type & Constant Aliases ────────────────────────────────────────────

type BoxShape = domain.BoxShape

const (
	BSRectangle = domain.BSRectangle
	BSCircle    = domain.BSCircle
)

// ─── ECS Component Type Variables ────────────────────────────────────────────
// Pointer directly -> static vars in enginetype to ensure start' order.
// No new ComponentType — just re-export from enginetype.

var (
	Position   = enginetype.Position
	Sprite     = enginetype.Sprite
	Box        = enginetype.Box
	Audio      = enginetype.Audio
	Input      = enginetype.Input
	Background = enginetype.Background
	Tilemap    = enginetype.Tilemap
)
