package components

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

type CollisionComponent struct {
	domain.IObject
	data *CollisionData
}

var Collision = enginetype.Collision

func (p *CollisionComponent) BindComponent(base IObject) {
	p.IObject = base
	p.data = enginetype.GetComponent(base, Collision)
}

func init() {
	enginetype.RegisterComponentInitializer("col", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Collision, domain.CollisionData{
			Handlers:     make(map[uint64]func(other IObject)),
			IsCollidable: true,
		})
	})
}

// OnCollisionTag dang ky m?t callback s? du?c g?i m?i frame khi va ch?m v?i Object c ch?a tag ch? d?nh.

func (p CollisionComponent) OnCollisionTag(tag string, handler func(other IObject)) {

	if p.data != nil {
		hash := enginetype.HashString(tag)
		p.data.Handlers[hash] = handler
	}
}

func (p CollisionComponent) IsCollidable() bool {
	if p.data == nil {
		return false
	}
	return p.data.IsCollidable
}

func (p CollisionComponent) SetIsCollidable(isCollidable bool) {
	if p.data != nil {
		p.data.IsCollidable = isCollidable
	}
}
