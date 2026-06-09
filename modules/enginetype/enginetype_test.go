package enginetype

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/yohamta/donburi"
)

func TestHashString(t *testing.T) {
	h1 := HashString("player")
	h2 := HashString("player")
	h3 := HashString("enemy")

	if h1 != h2 {
		t.Errorf("HashString should be deterministic. Got %v and %v", h1, h2)
	}
	if h1 == h3 {
		t.Errorf("HashString collision for different strings. Got %v", h1)
	}
}

type DummyData struct {
	HP int
}

func TestComponentRegistry(t *testing.T) {
	comp := NewComponentType[DummyData]("dmy")

	if comp == nil {
		t.Fatalf("NewComponentType returned nil")
	}

	fetched := GetComponentType("dmy")
	if fetched != comp {
		t.Errorf("GetComponentType did not return registered component")
	}

	// Test Initializer Registry
	initCalled := false
	RegisterComponentInitializer("dmy", func(entry *donburi.Entry) {
		initCalled = true
	})

	// We need a dummy entry to test InitializeComponent
	world := donburi.NewWorld()
	entry := world.Entry(world.Create(comp))

	InitializeComponent("dmy", entry)

	if !initCalled {
		t.Errorf("InitializeComponent did not call the registered function")
	}
}

// Mock IObject that implements entryProvider
type mockObject struct {
	entry *donburi.Entry
}

// Implement minimal IObject methods
func (m *mockObject) OnCreate()                     {}
func (m *mockObject) OnStep()                       {}
func (m *mockObject) OnDestroy()                    {}
func (m *mockObject) OnSave(data map[string]any)    {}
func (m *mockObject) OnLoad(data map[string]any)    {}
func (m *mockObject) GetPool() domain.IPool         { return nil }
func (m *mockObject) SetPool(pool domain.IPool)     {}
func (m *mockObject) SetTokens(tokenClasses string) {}
func (m *mockObject) Remove()                       {}
func (m *mockObject) Entry() *donburi.Entry         { return m.entry }

func TestComponentHelper(t *testing.T) {
	world := donburi.NewWorld()
	comp := donburi.NewComponentType[DummyData]()

	entity := world.Create(comp)
	entry := world.Entry(entity)

	obj := &mockObject{entry: entry}

	// Test SetComponent
	SetComponent(obj, comp, DummyData{HP: 100})

	// Test GetComponent
	data := GetComponent(obj, comp)
	if data == nil {
		t.Fatalf("GetComponent returned nil")
	}
	if data.HP != 100 {
		t.Errorf("GetComponent returned wrong data: %v", data.HP)
	}

	// Test AddComponentType on fly
	comp2 := donburi.NewComponentType[int]()
	AddComponentType(obj, comp2)

	if !entry.HasComponent(comp2) {
		t.Errorf("AddComponentType failed to add component")
	}
}
