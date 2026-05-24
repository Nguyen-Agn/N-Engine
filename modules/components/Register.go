package components

import (
	"autoworld/domain"
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
type BackgroundData = domain.BackgroundData
type TilemapData = domain.TilemapData
