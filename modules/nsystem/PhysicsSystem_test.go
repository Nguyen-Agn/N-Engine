package nsystem

import (
	"autoworld/modules/enginetype"
	"testing"

	"github.com/yohamta/donburi"
)

func TestPhysicsSystem(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(enginetype.Velocity, enginetype.Position)
	entry := world.Entry(entity)

	enginetype.InitializeComponent("vel", entry)
	enginetype.InitializeComponent("pos", entry)

	obj := &mockObject{entry: entry}

	velData := enginetype.GetComponent(obj, enginetype.Velocity)
	posData := enginetype.GetComponent(obj, enginetype.Position)

	velData.Vx = 10
	velData.Vy = -5
	velData.MaxSpeed = 8
	velData.Friction = 1

	sys := NewVelocitySystem()

	// Apply Physics: 
	// Friction: Vx 10 -> 9, Vy -5 -> -4
	// MaxSpeed: Vx 9 -> 8, Vy -4 -> -4
	// Pos: X 0 -> 8, Y 0 -> -4
	sys.Update([]IObject{obj})

	if velData.Vx != 8 || velData.Vy != -4 {
		t.Errorf("Physics System Velocity mismatch: Vx=%v, Vy=%v", velData.Vx, velData.Vy)
	}

	if posData.X != 8 || posData.Y != -4 {
		t.Errorf("Physics System Position mismatch: X=%v, Y=%v", posData.X, posData.Y)
	}
}
