package enginetype

import (
	"autoworld/domain"
)

type IScene = domain.IScene

type IObject = domain.IObject

var (
	Position   = NewComponentType[domain.PositionData]("pos")
	Sprite     = NewComponentType[domain.SpriteData]("spr")
	Box        = NewComponentType[domain.BoxData]("box")
	Audio      = NewComponentType[domain.AudioData]("aud")
	Infor      = NewComponentType[domain.InforData]("inf")
	Direction  = NewComponentType[domain.DirectionData]("dir")
	Input      = NewComponentType[domain.InputData]("inp")
	Background = NewComponentType[domain.BackgroundData]("bg")
	Tilemap    = NewComponentType[domain.TilemapData]("til")
	Alarm      = NewComponentType[domain.AlarmData]("alr")
	Velocity   = NewComponentType[domain.VelocityData]("vel")
	Tween      = NewComponentType[domain.TweenData]("twn")
)
