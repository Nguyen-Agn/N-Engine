package nsystem

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

func TestTweenSystem(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(enginetype.Tween, enginetype.Position)
	entry := world.Entry(entity)

	enginetype.InitializeComponent("twn", entry)
	enginetype.InitializeComponent("pos", entry)

	obj := &mockObject{entry: entry}

	tweenData := enginetype.GetComponent(obj, enginetype.Tween)
	posData := enginetype.GetComponent(obj, enginetype.Position)

	posData.X = 0
	posData.Y = 0

	tweenData.Tweens = append(tweenData.Tweens, domain.Tween{
		TargetType: "move",
		EndX:       100,
		EndY:       50,
		Duration:   2,
		IsActive:   true,
	})

	sys := NewTweenSystem()

	// Step 1: Elapsed 1, Duration 2 -> Progress 0.5 -> X=50, Y=25
	sys.Update([]IObject{obj})

	if posData.X != 50 || posData.Y != 25 {
		t.Errorf("Tween progress 50%% mismatch: X=%v, Y=%v", posData.X, posData.Y)
	}

	// Step 2: Elapsed 2, Duration 2 -> Progress 1.0 -> X=100, Y=50
	sys.Update([]IObject{obj})

	if posData.X != 100 || posData.Y != 50 {
		t.Errorf("Tween progress 100%% mismatch: X=%v, Y=%v", posData.X, posData.Y)
	}

	if tweenData.Tweens[0].IsActive {
		t.Errorf("Tween should be inactive after finishing")
	}
}
