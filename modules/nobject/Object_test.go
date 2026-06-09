package nobject

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"

	"github.com/yohamta/donburi"
)

// Mock IPool
type mockPool struct {
	putCalled bool
}

func (m *mockPool) Put(obj domain.IObject) bool {
	m.putCalled = true
	return true
}

func TestObject(t *testing.T) {
	world := donburi.NewWorld()
	dummyComp := donburi.NewComponentType[int]()
	entity := world.Create(dummyComp)
	entry := world.Entry(entity)

	obj := NewObject(entry)

	if obj.Entry() != entry {
		t.Errorf("Entry() mismatch")
	}

	pool := &mockPool{}
	obj.SetPool(pool)

	if obj.GetPool() != pool {
		t.Errorf("GetPool() mismatch")
	}

	// Just calling these to ensure no panic (default empty implementation)
	obj.OnCreate()
	obj.OnStep()
	obj.OnDestroy()
	obj.OnSave(nil)
	obj.OnLoad(nil)
	obj.SetTokens("pos spr")
	obj.Remove()
}
