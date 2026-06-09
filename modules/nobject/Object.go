package nobject

import (
	"autoworld/domain"

	"github.com/yohamta/donburi"
)

type Object struct {
	entry *donburi.Entry
	pool  domain.IPool
}

func NewObject(entry *donburi.Entry) *Object {
	return &Object{entry: entry}
}

// Entry trả về *donburi.Entry bên trong object.
// Dùng bởi napi.SetComponent và napi.GetComponent để gán custom component data.
func (this *Object) Entry() *donburi.Entry {
	return this.entry
}

func (this *Object) GetPool() domain.IPool {
	return this.pool
}

func (this *Object) SetPool(pool domain.IPool) {
	this.pool = pool
}

// #region Event

// step update (call every frame)
func (this *Object) OnCreate()                     {}
func (this *Object) OnStep()                       {}
func (this *Object) OnDestroy()                    {}
func (this *Object) OnSave(data map[string]any)    {}
func (this *Object) OnLoad(data map[string]any)    {}
func (this *Object) SetTokens(tokenClasses string) {}
func (this *Object) Remove()                       {}

// #endregion
