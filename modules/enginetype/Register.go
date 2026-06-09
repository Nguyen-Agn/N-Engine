package enginetype

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
)

type IScene = domain.IScene

type IObject = domain.IObject

var (
	Position   = NewComponentType[domain.PositionData]("pos")
	Sprite     = NewComponentType[domain.SpriteData]("spr")
	Box        = NewComponentType[domain.BoxData]("box")
	Audio      = NewComponentType[domain.AudioData]("aud")
	Infor      = NewComponentType[domain.InforData]("info")
	Direction  = NewComponentType[domain.DirectionData]("dir")
	Input      = NewComponentType[domain.InputData]("inp")
	Background = NewComponentType[domain.BackgroundData]("back")
	Tilemap    = NewComponentType[domain.TilemapData]("tile")
	Alarm      = NewComponentType[domain.AlarmData]("alrm")
	Velocity   = NewComponentType[domain.VelocityData]("velo")
	Tween      = NewComponentType[domain.TweenData]("twn")
	Collision  = NewComponentType[domain.CollisionData]("col")
	Draw       = NewComponentType[domain.DrawData]("drw")
	Debug      = NewComponentType[domain.DebugData]("deb")
)
