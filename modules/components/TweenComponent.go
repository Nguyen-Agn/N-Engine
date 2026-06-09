package components

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Tween = enginetype.Tween

// init registers the Tween Component initializer with an empty slice of tweens.
func init() {
	enginetype.RegisterComponentInitializer("twn", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Tween, domain.TweenData{
			Tweens: make([]domain.Tween, 0),
		})
	})
}

// TweenComponent is a mixin to embed into Custom Objects to provide tweening capabilities.
type TweenComponent struct {
	IObject
	data *domain.TweenData
}

// BindComponent binds the base object and retrieves the TweenData from the ECS.
// Inputs: base - the base IObject to bind.
func (t *TweenComponent) BindComponent(base IObject) {
	t.IObject = base
	t.data = enginetype.GetComponent(base, Tween)
}

// addTween adds a new tween to the component or overwrites an existing one of the same type.
// Inputs: tw - the Tween struct to add.
func (t TweenComponent) addTween(tw domain.Tween) {
	if t.data == nil {
		return
	}
	// Overwrite if there is an active tween of the same target type (e.g., replace old move command)
	for i := range t.data.Tweens {
		if t.data.Tweens[i].TargetType == tw.TargetType && t.data.Tweens[i].IsActive {
			t.data.Tweens[i] = tw
			return
		}
	}
	// Otherwise add a new one
	// Find an empty slot
	for i := range t.data.Tweens {
		if !t.data.Tweens[i].IsActive {
			t.data.Tweens[i] = tw
			return
		}
	}
	// If no empty slot is available, append it
	t.data.Tweens = append(t.data.Tweens, tw)
}

// TweenMove initiates a move tween to the specified target coordinates.
// Inputs: targetX, targetY - the destination coordinates, duration - time in frames/ticks to complete the move.
func (t TweenComponent) TweenMove(targetX, targetY float32, duration int) {
	if duration <= 0 {
		return
	}
	t.addTween(domain.Tween{
		TargetType: "move",
		IsActive:   true,
		Duration:   duration,
		Elapsed:    0,
		EndX:       targetX,
		EndY:       targetY,
	})
}

// TweenScale initiates a scaling tween to the specified target scales.
// Inputs: targetScaleX, targetScaleY - the destination scale values, duration - time in frames/ticks to complete.
func (t TweenComponent) TweenScale(targetScaleX, targetScaleY float32, duration int) {
	if duration <= 0 {
		return
	}
	t.addTween(domain.Tween{
		TargetType: "scale",
		IsActive:   true,
		Duration:   duration,
		Elapsed:    0,
		EndX:       targetScaleX,
		EndY:       targetScaleY,
	})
}

// TweenAlpha initiates an alpha (transparency) tween to the specified target alpha.
// Inputs: targetAlpha - the destination alpha value (0-255), duration - time in frames/ticks to complete.
func (t TweenComponent) TweenAlpha(targetAlpha uint8, duration int) {
	if duration <= 0 {
		return
	}
	t.addTween(domain.Tween{
		TargetType: "alpha",
		IsActive:   true,
		Duration:   duration,
		Elapsed:    0,
		StartX:     0, // Will be set by TweenSystem when it starts
		EndX:       float32(targetAlpha),
	})
}
