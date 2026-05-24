package nsystem

import (
	"autoworld/modules/enginetype"
)

// TweenSystem tính toán lerp (nội suy) cho các giá trị như vị trí, scale, màu sắc.
type TweenSystem struct {
}

func NewTweenSystem() *TweenSystem {
	return &TweenSystem{}
}

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

			// Khởi tạo StartValue lần đầu tiên (nếu Elapsed == 0)
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

			// Tăng thời gian
			tw.Elapsed++
			if tw.Elapsed > tw.Duration {
				tw.Elapsed = tw.Duration
			}

			// Tính % hoàn thành (progress)
			progress := float32(tw.Elapsed) / float32(tw.Duration)

			// Cập nhật giá trị mục tiêu
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

			// Kết thúc
			if tw.Elapsed == tw.Duration {
				tw.IsActive = false
			}
		}
	}
}

// lerp nội suy tuyến tính giữa a và b với tỷ lệ t (0.0 - 1.0)
func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
