package nsystem

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"
)

type IObject = domain.IObject
type PositionData = domain.PositionData
type SpriteData = domain.SpriteData
type BoxData = domain.BoxData
type AudioData = domain.AudioData

type BackgroundData = domain.BackgroundData
type TilemapData = domain.TilemapData

var (
	Position   = enginetype.Position
	Sprite     = enginetype.Sprite
	Box        = enginetype.Box
	Background = enginetype.Background
	Tilemap    = enginetype.Tilemap
	Audio      = enginetype.Audio
)
