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

// Purpose: Multiplies two Vec2 instances treating them as complex numbers, equivalent to rotation and scaling.
// Inputs: b (Vec2) - The vector to multiply with.
// Outputs: (Vec2) - The result of the multiplication.
// If b is a unit vector, this results in pure rotation.
func (a Vec2) Mul(b Vec2) Vec2 {
	return Vec2{
		X: a.X*b.X - a.Y*b.Y,
		Y: a.X*b.Y + a.Y*b.X,
	}
}

// Purpose: Returns the complex conjugate (negating the imaginary part) of the Vec2.
// Outputs: (Vec2) - The conjugated vector.
func (a Vec2) Conj() Vec2 { return Vec2{a.X, -a.Y} }

// Purpose: Normalizes the Vec2 to a unit vector (length 1).
// Outputs: (Vec2) - The normalized unit vector.
func (a Vec2) Norm() Vec2 {
	l := float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y)))
	return Vec2{a.X / l, a.Y / l}
}

// Purpose: Calculates the squared length of the vector.
// Outputs: (float32) - The squared length. Useful for fast distance comparisons without Sqrt.
func (a Vec2) LenSq() float32 { return a.X*a.X + a.Y*a.Y }

// Purpose: Calculates the actual length of the Vec2.
// Outputs: (float32) - The length.
func (a Vec2) Len() float32 { return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y))) }

// Purpose: Multiplies the vector by a scalar value.
// Inputs: s (float32) - The scalar to multiply by.
// Outputs: (Vec2) - The scaled vector.
func (a Vec2) Scale(s float32) Vec2 { return Vec2{a.X * s, a.Y * s} }

// Purpose: Adds another Vec2 to this one.
// Inputs: b (Vec2) - The vector to add.
// Outputs: (Vec2) - The sum vector.
func (a Vec2) Add(b Vec2) Vec2 { return Vec2{a.X + b.X, a.Y + b.Y} }

// Purpose: Subtracts another Vec2 from this one.
// Inputs: b (Vec2) - The vector to subtract.
// Outputs: (Vec2) - The difference vector.
func (a Vec2) Sub(b Vec2) Vec2 { return Vec2{a.X - b.X, a.Y - b.Y} }

// Purpose: Calculates the dot product of this vector and another.
// Inputs: b (Vec2) - The other vector.
// Outputs: (float32) - The dot product value.
func (a Vec2) Dot(b Vec2) float32 { return a.X*b.X + a.Y*b.Y }

// Purpose: Calculates the 2D cross product (scalar) of this vector and another.
// Inputs: b (Vec2) - The other vector.
// Outputs: (float32) - The cross product value. Positive if b is left of a, negative if right.
func (a Vec2) Cross(b Vec2) float32 { return a.X*b.Y - a.Y*b.X }
