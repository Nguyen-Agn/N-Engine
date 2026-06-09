package nobject

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
)

type IObject = domain.IObject
type IAudioLW = domain.IAudioLW
type ISpriteLW = domain.ISpriteLW
type PositionData = domain.PositionData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type AudioData = domain.AudioData
type BoxShape = domain.BoxShape

const (
	BSRectangle = domain.BSRectangle
	BSCircle    = domain.BSCircle
)

// var (
// 	Position = napi.Position
// 	Sprite   = napi.Sprite
// 	Box      = napi.Box
// 	Audio    = napi.Audio
// )
