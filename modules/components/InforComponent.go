package components

import "autoworld/modules/enginetype"

type InforComponent struct {
	IObject
}

var Infor = enginetype.Infor

func (p InforComponent) GetId() int {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return 0
	}
	return data.Id
}

func (p InforComponent) GetName() string {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return ""
	}
	return data.Name
}

func (p InforComponent) AddTag(tag string) {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.Tags = append(data.Tags, enginetype.HashString(tag))
	}
}

func (p InforComponent) HasTag(tag string) bool {
	return p.HasTagHash(enginetype.HashString(tag))
}

func (p InforComponent) HasTagHash(hash uint64) bool {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return false
	}
	for _, t := range data.Tags {
		if t == hash {
			return true
		}
	}
	return false
}

func (p InforComponent) IsDead() bool {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return false
	}
	return data.IsDead
}

func (p InforComponent) MarkDead() {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.IsDead = true
	}
}

func (p InforComponent) SetIsDead(dead bool) {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.IsDead = dead
	}
}

func (p InforComponent) SaveTag() string {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return ""
	}
	return data.SaveTag
}

func (p InforComponent) SetSaveTag(tag string) {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.SaveTag = tag
	}
}
