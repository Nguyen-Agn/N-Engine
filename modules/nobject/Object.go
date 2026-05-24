package nobject

import (
	"github.com/yohamta/donburi"
)

type Object struct {
	entry *donburi.Entry
}

func NewObject(entry *donburi.Entry) *Object {
	return &Object{entry: entry}
}

// Entry trả về *donburi.Entry bên trong object.
// Dùng bởi napi.SetComponent và napi.GetComponent để gán custom component data.
func (this *Object) Entry() *donburi.Entry {
	return this.entry
}

// #region Event

// step update (call every frame)
func (this *Object) Create()     {}
func (this *Object) StepUpdate() {}
func (this *Object) Destroy()    {}

// #endregion
