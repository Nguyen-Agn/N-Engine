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
