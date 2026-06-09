package nmath

// =============================================================================
// Vec2 — kiểu vector 2D dựa trên số phức
//
// Vec2{X, Y} tương đương số phức X + Yi.
// Nhân hai Vec2 = nhân số phức = xoay + scale trong một bước.
// Đây là nền tảng cho toàn bộ nhóm hàm góc và hướng.
// =============================================================================

import "math"

// Vec2 biểu diễn vector 2D hoặc số phức X + Yi.
type Vec2 struct{ X, Y float32 }

const (
	deg2rad = math.Pi / 180.0
	rad2deg = 180.0 / math.Pi
)

// --- Các phép toán nội bộ của Vec2 -------------------------------------------

// Mul nhân hai Vec2 như số phức: tương đương xoay + scale.
// Nếu b là unit vector thì kết quả là xoay thuần túy.
//
//	(a.X + a.Y·i) × (b.X + b.Y·i) = (a.X·b.X − a.Y·b.Y) + (a.X·b.Y + a.Y·b.X)·i
//
// Chi phí: 4 phép nhân + 2 phép cộng. Không dùng sin/cos.
func (a Vec2) Mul(b Vec2) Vec2 {
	return Vec2{
		X: a.X*b.X - a.Y*b.Y,
		Y: a.X*b.Y + a.Y*b.X,
	}
}

// Conj trả về conjugate (đảo chiều phần ảo) của Vec2.
// Dùng để tính góc tương đối: conj(a).Mul(b) cho góc từ a đến b.
func (a Vec2) Conj() Vec2 { return Vec2{a.X, -a.Y} }

// Norm chuẩn hóa Vec2 thành unit vector (độ dài = 1).
// Tránh gọi khi |v| ≈ 0.
func (a Vec2) Norm() Vec2 {
	l := float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y)))
	return Vec2{a.X / l, a.Y / l}
}

// LenSq trả về bình phương độ dài. Dùng để so sánh khoảng cách mà không cần Sqrt.
func (a Vec2) LenSq() float32 { return a.X*a.X + a.Y*a.Y }

// Len trả về độ dài thực của Vec2.
func (a Vec2) Len() float32 { return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y))) }

// Scale nhân Vec2 với một scalar.
func (a Vec2) Scale(s float32) Vec2 { return Vec2{a.X * s, a.Y * s} }

// Add cộng hai Vec2.
func (a Vec2) Add(b Vec2) Vec2 { return Vec2{a.X + b.X, a.Y + b.Y} }

// Sub trừ hai Vec2.
func (a Vec2) Sub(b Vec2) Vec2 { return Vec2{a.X - b.X, a.Y - b.Y} }

// Dot trả về tích vô hướng. cos(θ) = Dot(a.Norm(), b.Norm()).
func (a Vec2) Dot(b Vec2) float32 { return a.X*b.X + a.Y*b.Y }

// Cross trả về tích có hướng (scalar 2D). sin(θ) = Cross(a.Norm(), b.Norm()).
// Dương = b ở phía trái a; âm = b ở phía phải a.
func (a Vec2) Cross(b Vec2) float32 { return a.X*b.Y - a.Y*b.X }
