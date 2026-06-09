package nsystem

import (
	"autoworld/modules/enginetype"
)

// TweenSystem calculates linear interpolation (lerp) for values such as position, scale, and color/alpha.
type TweenSystem struct {
}

// NewTweenSystem creates and returns a new instance of TweenSystem.
// Outputs: Returns a pointer to a newly initialized TweenSystem.
func NewTweenSystem() *TweenSystem {
	return &TweenSystem{}
}

// Update processes all active tweens on entities with a Tween component.
// Inputs: objectList ([]IObject) - The list of active objects in the scene.
// Purpose: It iterates through tweens. If a tween is newly started, it records the initial component values. It then advances the elapsed time, calculates progress, and interpolates target properties (move, scale, alpha). Deactivates tweens upon completion.
func (s *TweenSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		tweenData := enginetype.GetComponent(obj, enginetype.Tween)
		if tweenData == nil {
			continue
		}

		for i := range tweenData.Tweens {
			tw := &tweenData.Tweens[i]
			if !tw.IsActive {
				continue
			}

			// Initialize StartValue for the first time (if Elapsed == 0)
			if tw.Elapsed == 0 {
				switch tw.TargetType {
				case "move":
					posData := enginetype.GetComponent(obj, enginetype.Position)
					if posData != nil {
						tw.StartX = posData.X
						tw.StartY = posData.Y
					}
				case "scale":
					sprData := enginetype.GetComponent(obj, enginetype.Sprite)
					if sprData != nil {
						tw.StartX = sprData.ScaleX
						tw.StartY = sprData.ScaleY
					}
				case "alpha":
					sprData := enginetype.GetComponent(obj, enginetype.Sprite)
					if sprData != nil {
						tw.StartX = float32(sprData.ImageColor.A)
					}
				}
			}

			// Increment elapsed time
			tw.Elapsed++
			if tw.Elapsed > tw.Duration {
				tw.Elapsed = tw.Duration
			}

			// Calculate completion percentage (progress)
			progress := float32(tw.Elapsed) / float32(tw.Duration)

			// Update target values
			switch tw.TargetType {
			case "move":
				posData := enginetype.GetComponent(obj, enginetype.Position)
				if posData != nil {
					posData.X = lerp(tw.StartX, tw.EndX, progress)
					posData.Y = lerp(tw.StartY, tw.EndY, progress)
				}
			case "scale":
				sprData := enginetype.GetComponent(obj, enginetype.Sprite)
				if sprData != nil {
					sprData.ScaleX = lerp(tw.StartX, tw.EndX, progress)
					sprData.ScaleY = lerp(tw.StartY, tw.EndY, progress)
				}
			case "alpha":
				sprData := enginetype.GetComponent(obj, enginetype.Sprite)
				if sprData != nil {
					sprData.ImageColor.A = uint8(lerp(tw.StartX, tw.EndX, progress))
				}
			}

			// Complete tween
			if tw.Elapsed == tw.Duration {
				tw.IsActive = false
			}
		}
	}
}

// lerp performs linear interpolation between a and b using ratio t (0.0 - 1.0).
// Inputs: 
//   a (float32) - Start value.
//   b (float32) - End value.
//   t (float32) - Progress ratio from 0.0 to 1.0.
// Outputs: Returns the interpolated float32 value.
func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
