package components

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

// Mock IObject
type mockObject struct {
	entry *donburi.Entry
}

// OnCreate is a mock implementation.

func (m *mockObject) OnCreate() {}

// OnStep is a mock implementation.

func (m *mockObject) OnStep() {}

// OnDestroy is a mock implementation.

func (m *mockObject) OnDestroy() {}

// OnSave is a mock implementation.

func (m *mockObject) OnSave(data map[string]any) {}

// OnLoad is a mock implementation.

func (m *mockObject) OnLoad(data map[string]any) {}

// GetPool is a mock implementation.

func (m *mockObject) GetPool() domain.IPool { return nil }

// SetPool is a mock implementation.

func (m *mockObject) SetPool(pool domain.IPool) {}

// SetTokens is a mock implementation.

func (m *mockObject) SetTokens(tokenClasses string) {}

// Remove is a mock implementation.
func (m *mockObject) Remove() {}

// Entry is a mock implementation.
func (m *mockObject) Entry() *donburi.Entry { return m.entry }

// TestPositionComponent verifies the creation and modification of PositionComponent.
func TestPositionComponent(t *testing.T) {
	world := donburi.NewWorld()

	// Create entity with Position component
	entity := world.Create(enginetype.Position)
	entry := world.Entry(entity)

	// Initialize default values using initializer
	enginetype.InitializeComponent("pos", entry)

	obj := &mockObject{entry: entry}

	// Create and bind mixin
	pos := &PositionComponent{}
	pos.BindComponent(obj)

	// Test Getters
	if pos.X() != 0 || pos.Y() != 0 {
		t.Errorf("Initial position should be 0,0, got %v,%v", pos.X(), pos.Y())
	}

	// Test Setters
	pos.SetX(100)
	pos.SetY(200)

	if pos.X() != 100 || pos.Y() != 200 {
		t.Errorf("Position not updated. Got %v,%v", pos.X(), pos.Y())
	}

	// Verify data is actually stored in donburi
	data := enginetype.GetComponent(obj, enginetype.Position)
	if data.X != 100 || data.Y != 200 {
		t.Errorf("Donburi data mismatch")
	}
}

// TestBoxComponent verifies the creation and modification of BoxComponent.
func TestBoxComponent(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(enginetype.Box)
	entry := world.Entry(entity)
	enginetype.InitializeComponent("box", entry)

	obj := &mockObject{entry: entry}

	box := &BoxComponent{}
	box.BindComponent(obj)

	box.SetBoxW(50)
	box.SetBoxH(60)
	box.SetBoxX(-25)
	box.SetBoxY(-30)
	box.SetShape(BSRectangle)

	if box.BoxW() != 50 || box.BoxH() != 60 || box.BoxX() != -25 || box.BoxY() != -30 || box.Shape() != BSRectangle {
		t.Errorf("Box Component data mismatch")
	}
}

// TestCollisionComponent verifies the creation and modification of CollisionComponent.
func TestCollisionComponent(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(enginetype.Collision)
	entry := world.Entry(entity)
	enginetype.InitializeComponent("col", entry)

	obj := &mockObject{entry: entry}
	col := &CollisionComponent{}
	col.BindComponent(obj)

	col.SetIsCollidable(false)
	if col.IsCollidable() {
		t.Errorf("IsCollidable should be false")
	}

	col.OnCollisionTag("enemy", func(other IObject) {})
	// Just ensuring no panic
}

// TestVelocityComponent verifies the creation and modification of VelocityComponent.
func TestVelocityComponent(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(enginetype.Velocity)
	entry := world.Entry(entity)
	enginetype.InitializeComponent("vel", entry)

	obj := &mockObject{entry: entry}
	vel := &VelocityComponent{}
	vel.BindComponent(obj)

	vel.SetVelocity(10, -5)
	if vel.VelocityX() != 10 || vel.VelocityY() != -5 {
		t.Errorf("Velocity mismatch")
	}

	vel.AddVelocity(2, 3)
	if vel.VelocityX() != 12 || vel.VelocityY() != -2 {
		t.Errorf("AddVelocity failed")
	}
}
