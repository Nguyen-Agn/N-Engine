package nsystem

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"
)

type IObject = domain.IObject
type PositionData = domain.PositionData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type AudioData = domain.AudioData
type AudioTrackState = domain.AudioTrackState
type InputData = domain.InputData

type BackgroundData = domain.BackgroundData
type DrawData = domain.DrawData
type TilemapData = domain.TilemapData

var (
	Draw       = enginetype.Draw
	Position   = enginetype.Position
	Sprite     = enginetype.Sprite
	Box        = enginetype.Box
	Background = enginetype.Background
	Tilemap    = enginetype.Tilemap
	Audio      = enginetype.Audio
	Debug      = enginetype.Debug
	Infor      = enginetype.Infor
)
