package components

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
)

type IObject = domain.IObject
type IScene = domain.IScene

type IAudioLW = domain.IAudioLW
type ISpriteLW = domain.ISpriteLW

type BoxShape = domain.BoxShape

const (
	BSRectangle BoxShape = "rectangle"
	BSCircle    BoxShape = "circle"
)

type AudioData = domain.AudioData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type PositionData = domain.PositionData
type InforData = domain.InforData
type DirectionData = domain.DirectionData
type InputData = domain.InputData
type KeyBinding = domain.KeyBinding
type MouseBinding = domain.MouseBinding
type EventType = domain.EventType
type BackgroundData = domain.BackgroundData
type TilemapData = domain.TilemapData
type CollisionData = domain.CollisionData

type DrawData = domain.DrawData
type DebugData = domain.DebugData
