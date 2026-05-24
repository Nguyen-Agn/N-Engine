package nsystem

import (
	"autoworld/modules/enginetype"
)

// PhysicsSystem tính toán vật lý cơ bản cho entity: ma sát, giới hạn vận tốc, và cộng vận tốc vào vị trí.
type PhysicsSystem struct {
}

func NewVelocitySystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (s *PhysicsSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		velData := enginetype.GetComponent(obj, enginetype.Velocity)
		if velData == nil {
			continue
		}

		// Apply Friction
		if velData.Friction > 0 {
			// Giảm tốc độ dựa trên ma sát
			if velData.Vx > 0 {
				velData.Vx -= velData.Friction
				if velData.Vx < 0 {
					velData.Vx = 0
				}
			} else if velData.Vx < 0 {
				velData.Vx += velData.Friction
				if velData.Vx > 0 {
					velData.Vx = 0
				}
			}

			if velData.Vy > 0 {
				velData.Vy -= velData.Friction
				if velData.Vy < 0 {
					velData.Vy = 0
				}
			} else if velData.Vy < 0 {
				velData.Vy += velData.Friction
				if velData.Vy > 0 {
					velData.Vy = 0
				}
			}
		}

		// Apply MaxSpeed
		if velData.MaxSpeed > 0 {
			// Đơn giản hóa: chỉ limit từng trục (có thể dùng vector magnitude nếu cần chính xác)
			if velData.Vx > velData.MaxSpeed {
				velData.Vx = velData.MaxSpeed
			}
			if velData.Vx < -velData.MaxSpeed {
				velData.Vx = -velData.MaxSpeed
			}
			if velData.Vy > velData.MaxSpeed {
				velData.Vy = velData.MaxSpeed
			}
			if velData.Vy < -velData.MaxSpeed {
				velData.Vy = -velData.MaxSpeed
			}
		}

		// Apply Velocity to Position
		posData := enginetype.GetComponent(obj, enginetype.Position)
		if posData != nil {
			posData.X += velData.Vx
			posData.Y += velData.Vy
		}
	}
}
