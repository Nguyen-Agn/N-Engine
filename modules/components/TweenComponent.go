package components

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Tween = enginetype.Tween

func init() {
	enginetype.RegisterComponentInitializer("twn", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Tween, domain.TweenData{
			Tweens: make([]domain.Tween, 0),
		})
	})
}

// TweenComponent là Mixin để nhúng vào Custom Object.
type TweenComponent struct {
	IObject
	data *domain.TweenData
}

func (t *TweenComponent) BindComponent(base IObject) {
	t.IObject = base
	t.data = enginetype.GetComponent(base, Tween)
}

func (t TweenComponent) addTween(tw domain.Tween) {
	if t.data == nil {
		return
	}
	// Ghi đè nếu đã có cùng loại target (VD: đang di chuyển thì hủy lệnh di chuyển cũ)
	for i := range t.data.Tweens {
		if t.data.Tweens[i].TargetType == tw.TargetType && t.data.Tweens[i].IsActive {
			t.data.Tweens[i] = tw
			return
		}
	}
	// Hoặc thêm mới nếu chưa có
	// Tìm chỗ trống
	for i := range t.data.Tweens {
		if !t.data.Tweens[i].IsActive {
			t.data.Tweens[i] = tw
			return
		}
	}
	// Không có chỗ trống thì append
	t.data.Tweens = append(t.data.Tweens, tw)
}

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

func (t TweenComponent) TweenAlpha(targetAlpha uint8, duration int) {
	if duration <= 0 {
		return
	}
	t.addTween(domain.Tween{
		TargetType: "alpha",
		IsActive:   true,
		Duration:   duration,
		Elapsed:    0,
		StartX:     0,             // Sẽ được set bởi TweenSystem khi bắt đầu
		EndX:       float32(targetAlpha),
	})
}
