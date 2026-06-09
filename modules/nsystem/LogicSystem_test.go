package nsystem

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/yohamta/donburi"
)

type mockObject struct {
	entry        *donburi.Entry
	createCount  int
	stepCount    int
	destroyCount int
}

func (m *mockObject) OnCreate()                     { m.createCount++ }
func (m *mockObject) OnStep()                       { m.stepCount++ }
func (m *mockObject) OnDestroy()                    { m.destroyCount++ }
func (m *mockObject) OnSave(data map[string]any)    {}
func (m *mockObject) OnLoad(data map[string]any)    {}
func (m *mockObject) GetPool() domain.IPool         { return nil }
func (m *mockObject) SetPool(pool domain.IPool)     {}
func (m *mockObject) SetTokens(tokenClasses string) {}
func (m *mockObject) Remove()                       {}
func (m *mockObject) Entry() *donburi.Entry         { return m.entry }

func TestLogicSystem(t *testing.T) {
	ls := NewLogicSystem()
	obj := &mockObject{}

	ls.AddObjectCreated(obj)
	ls.AddObjectDestroy(obj)

	// Provide the object list for step
	ls.Update([]IObject{obj})

	if obj.createCount != 1 {
		t.Errorf("Expected OnCreate to be called once, got %v", obj.createCount)
	}
	if obj.stepCount != 1 {
		t.Errorf("Expected OnStep to be called once, got %v", obj.stepCount)
	}
	if obj.destroyCount != 1 {
		t.Errorf("Expected OnDestroy to be called once, got %v", obj.destroyCount)
	}

	// Verify slices are cleared
	if len(ls.CreateQuery) != 0 || len(ls.DestroyQuery) != 0 {
		t.Errorf("Expected CreateQuery and DestroyQuery to be cleared")
	}
}
