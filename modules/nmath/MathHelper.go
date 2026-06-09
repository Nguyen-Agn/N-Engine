package nmath

import (
	"math"
	"math/rand/v2"
	"sync"
)

// MathHelper là struct rỗng dùng để gom nhóm các hàm toán học hỗ trợ Game 2D.
// Việc dùng struct thay vì func tự do giúp napi có thể export thành 1 object duy nhất (napi.Math).
// =============================================================================
// MathHelper — nhóm hàm góc và vector
//
// Có hai lớp API cho mỗi hàm:
//   - Lớp degree: nhận/trả float32 degree — quen thuộc, dễ dùng
//   - Lớp Vec2 (suffix V): nhận/trả Vec2 — zero sin/cos, dùng cho hot path
//
// Quy tắc chọn:
//   - Logic game thông thường → dùng lớp degree
//   - Vòng lặp cập nhật hàng trăm entity mỗi frame → dùng lớp Vec2
// =============================================================================

type MathHelper struct {
	rng *rand.Rand
}

var (
	instance MathHelper
	Once     sync.Once
)

// Purpose: Returns a singleton instance of MathHelper.
func GetInstance() MathHelper {
	Once.Do(func() {
		instance = NewMathHelper()
	})
	return instance
}

// Purpose: Clamps a value between a minimum and maximum.
// Inputs: val (float32) - Value to clamp, min (float32) - Minimum value, max (float32) - Maximum value.
// Outputs: (float32) - Clamped value.
func (m MathHelper) Clamp(val, min, max float32) float32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// Purpose: Performs linear interpolation between two values.
// Inputs: a (float32) - Start value, b (float32) - End value, t (float32) - Interpolation factor [0-1].
// Outputs: (float32) - Interpolated value.
func (m MathHelper) Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

// Purpose: Approaches a target value by a given amount.
// Inputs: current (float32) - Current value, target (float32) - Target value, amount (float32) - Maximum change step.
// Outputs: (float32) - The new value closer to the target.
func (m MathHelper) Approach(current, target, amount float32) float32 {
	diff := target - current
	if diff > amount {
		return current + amount
	}
	if diff < -amount {
		return current - amount
	}
	return target
}

// Purpose: Wraps a value within a specified range [min, max).
// Inputs: val (float32) - Value to wrap, min (float32) - Minimum bound, max (float32) - Maximum bound.
// Outputs: (float32) - Wrapped value.
func (m MathHelper) Wrap(val, min, max float32) float32 {
	span := max - min
	if span == 0 {
		return min
	}
	val -= min
	val -= span * float32(int32(val/span))
	if val < 0 {
		val += span
	}
	return val + min
}

// Purpose: Returns the sign of a value.
// Outputs: (int) - 1 if positive, -1 if negative, 0 if zero.
func (m MathHelper) Sign(val float64) int {
	if float32(val) < 0 {
		return -1
	}
	if float32(val) > 0 {
		return 1
	}
	return 0
}

// Purpose: Snaps a value to the nearest grid step.
// Inputs: val (float32) - Value to snap, step (float32) - Grid step size.
// Outputs: (float32) - Snapped value.
func (m MathHelper) Snap(val, step float32) float32 {
	return float32(int32(val/step+0.5)) * step
}

// --- Helper convert ----------------------------------------------------------

// Purpose: Creates a unit vector from an angle in degrees.
// Inputs: deg (float32) - Angle in degrees.
// Outputs: (Vec2) - The resulting unit direction vector.
func (m MathHelper) DirFromDeg(deg float32) Vec2 {
	r := float64(deg) * deg2rad
	return Vec2{float32(math.Cos(r)), float32(math.Sin(r))}
}

// Purpose: Converts a unit vector to an angle in degrees [-180, 180].
// Inputs: v (Vec2) - The direction vector.
// Outputs: (float32) - The angle in degrees.
func (m MathHelper) ToDeg(v Vec2) float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X))) * rad2deg
}

// =============================================================================
// LengthDir — tính toạ độ từ độ dài và hướng
// =============================================================================

// Purpose: Returns the X component of a vector given length and direction in degrees.
// Inputs: length (float32) - The length, direction (float32) - The angle in degrees.
// Outputs: (float32) - The X component.
func (m MathHelper) LengthDirX(length, direction float32) float32 {
	return length * float32(math.Cos(float64(direction)*deg2rad))
}

// Purpose: Returns the Y component of a vector given length and direction in degrees.
// Inputs: length (float32) - The length, direction (float32) - The angle in degrees.
// Outputs: (float32) - The Y component. Note: Y axis points downwards in 2D.
func (m MathHelper) LengthDirY(length, direction float32) float32 {
	return length * float32(math.Sin(float64(direction)*deg2rad))
}

// Purpose: Returns a scaled Vec2 from length and a unit direction vector.
// Inputs: length (float32) - The magnitude, dir (Vec2) - The unit direction vector.
// Outputs: (Vec2) - The resulting vector.
func (m MathHelper) LengthDirV(length float32, dir Vec2) Vec2 {
	return Vec2{dir.X * length, dir.Y * length}
}

// =============================================================================
// Distance — khoảng cách giữa hai điểm
// =============================================================================

// Purpose: Calculates the Euclidean distance between two points.
// Inputs: x1, y1 (float32) - First point, x2, y2 (float32) - Second point.
// Outputs: (float32) - The distance.
func (m MathHelper) Distance(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Purpose: Calculates the squared Euclidean distance between two points (faster than Distance as it skips Sqrt).
// Inputs: x1, y1 (float32) - First point, x2, y2 (float32) - Second point.
// Outputs: (float32) - The squared distance.
func (m MathHelper) DistanceSq(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

// Purpose: Compares the squared distance between (x1, y1) -> (x2, y2) and (x3, y3) -> (x4, y4).
// Outputs: (int) - Returns 1 if first distance is greater, -1 if smaller, 0 if equal.
func (m MathHelper) CompareDistance(x1, y1, x2, y2, x3, y3, x4, y4 float32) int {
	dx1, dy1 := x2-x1, y2-y1
	dx2, dy2 := x4-x3, y4-y3
	return m.Sign(float64(dx1*dx1 + dy1*dy1 - dx2*dx2 - dy2*dy2))
}

// =============================================================================
// Angle — góc hướng giữa hai điểm
// =============================================================================

// Purpose: Calculates the angle in degrees from point 1 to point 2.
// Inputs: x1, y1 (float32) - Start point, x2, y2 (float32) - End point.
// Outputs: (float32) - Angle in degrees [-180, 180].
func (m MathHelper) Angle(x1, y1, x2, y2 float32) float32 {
	return float32(math.Atan2(float64(y2-y1), float64(x2-x1))) * rad2deg
}

// Purpose: Calculates the unit direction vector from point 1 to point 2.
// Inputs: x1, y1 (float32) - Start point, x2, y2 (float32) - End point.
// Outputs: (Vec2) - Unit direction vector.
func (m MathHelper) AngleV(x1, y1, x2, y2 float32) Vec2 {
	return Vec2{x2 - x1, y2 - y1}.Norm()
}

// =============================================================================
// AngleDif — khoảng cách ngắn nhất giữa hai góc
// =============================================================================

// Purpose: Calculates the shortest angular difference between two angles in degrees.
// Inputs: angle1, angle2 (float32) - Angles in degrees.
// Outputs: (float32) - The shortest difference in degrees [-180, 180].
func (m MathHelper) AngleDif(angle1, angle2 float32) float32 {
	v1 := m.DirFromDeg(angle1)
	v2 := m.DirFromDeg(angle2)
	rel := v1.Conj().Mul(v2)
	// rel.X = cos(dif), rel.Y = sin(dif)
	// atan2 tự trả về [-π, π] — không cần Wrap thêm
	return float32(math.Atan2(float64(rel.Y), float64(rel.X))) * rad2deg
}

// Purpose: Calculates the angular difference between two unit vectors.
// Inputs: v1, v2 (Vec2) - The unit vectors.
// Outputs: (float32) - The shortest angular difference in degrees.
func (m MathHelper) AngleDifV(v1, v2 Vec2) float32 {
	rel := v1.Norm().Conj().Mul(v2.Norm())
	return float32(math.Atan2(float64(rel.Y), float64(rel.X))) * rad2deg
}

// =============================================================================
// AngleAdd — cộng góc và tự wrap
// =============================================================================

// Purpose: Adds an amount in degrees to an angle, handling wrapping automatically.
// Inputs: angle (float32) - Base angle, amount (float32) - Amount to add in degrees.
// Outputs: (float32) - The resulting wrapped angle.
func (m MathHelper) AngleAdd(angle, amount float32) float32 {
	dir := m.DirFromDeg(angle)
	rot := m.DirFromDeg(amount)
	return m.ToDeg(dir.Mul(rot))
	// ToDeg trả về [-180,180]; nếu cần [0,360] thêm:
	// if result < 0 { result += 360 }
}

// Purpose: Adds an angle in degrees to a direction vector.
// Inputs: dir (Vec2) - The initial direction vector, amountDeg (float32) - Angle to add in degrees.
// Outputs: (Vec2) - The rotated unit vector.
func (m MathHelper) AngleAddV(dir Vec2, amountDeg float32) Vec2 {
	rot := m.DirFromDeg(amountDeg)
	return dir.Mul(rot)
}

// =============================================================================
// RotateV — xoay vector quanh gốc toạ độ
// =============================================================================

// Purpose: Rotates a vector by a given direction vector (complex multiplication).
// Inputs: vec (Vec2) - The vector to rotate, dir (Vec2) - The rotation direction vector.
// Outputs: (Vec2) - The rotated vector.
func (m MathHelper) RotateV(vec, dir Vec2) Vec2 {
	return vec.Mul(dir)
}

// =============================================================================
// SlerpDir — nội suy hướng mượt trên vòng tròn đơn vị
// =============================================================================

// Purpose: Spherical linear interpolation between two direction vectors.
// Inputs: from, to (Vec2) - The unit vectors, t (float32) - Interpolation factor [0-1].
// Outputs: (Vec2) - The interpolated unit vector.
func (m MathHelper) SlerpDir(from, to Vec2, t float32) Vec2 {
	fn := from.Norm()
	tn := to.Norm()
	rel := fn.Conj().Mul(tn)
	relAng := float32(math.Atan2(float64(rel.Y), float64(rel.X)))
	partial := relAng * t
	rot := Vec2{
		float32(math.Cos(float64(partial))),
		float32(math.Sin(float64(partial))),
	}
	return fn.Mul(rot)
}

// Purpose: Wrapper around SlerpDir that operates with degree angles.
// Inputs: fromDeg, toDeg (float32) - Angles in degrees, t (float32) - Factor [0-1].
// Outputs: (float32) - The interpolated angle in degrees.
func (m MathHelper) SlerpDirDeg(fromDeg, toDeg, t float32) float32 {
	from := m.DirFromDeg(fromDeg)
	to := m.DirFromDeg(toDeg)
	return m.ToDeg(m.SlerpDir(from, to, t))
}

// =============================================================================
// IsInFOV — kiểm tra điểm có trong vùng tầm nhìn không
// =============================================================================

// Purpose: Checks if a target point is within the field of view of an observer.
// Inputs: ox, oy (float32) - Observer coordinates, forward (Vec2) - Observer's forward direction vector, tx, ty (float32) - Target coordinates, halfFOVDeg (float32) - Half the FOV angle in degrees.
// Outputs: (bool) - True if the target is within the FOV.
func (m MathHelper) IsInFOV(ox, oy float32, forward Vec2, tx, ty, halfFOVDeg float32) bool {
	toTarget := Vec2{tx - ox, ty - oy}.Norm()
	cosHalf := float32(math.Cos(float64(halfFOVDeg) * deg2rad))
	// Dot >= cos(halfFOV) ↔ góc <= halfFOV
	return forward.Dot(toTarget) >= cosHalf
}

// =============================================================================
// Reflect — phản xạ vector qua pháp tuyến
// =============================================================================

// Purpose: Reflects a vector across a surface normal.
// Inputs: v (Vec2) - The incoming vector, normal (Vec2) - The surface normal unit vector.
// Outputs: (Vec2) - The reflected vector.
func (m MathHelper) Reflect(v, normal Vec2) Vec2 {
	d := v.Dot(normal)
	return Vec2{
		X: v.X - 2*d*normal.X,
		Y: v.Y - 2*d*normal.Y,
	}
}

// =============================================================================
// Ghi chú thiết kế API
// =============================================================================
//
// Quy tắc chọn API:
//
//   Hàm degree (LengthDirX/Y, Angle, AngleDif, AngleAdd):
//   → Dùng khi nhận dữ liệu từ designer/editor/config (số người đọc được)
//   → Dùng cho logic không phải hot path
//
//   Hàm Vec2 suffix V (AngleV, AngleDifV, AngleAddV, RotateV, SlerpDir):
//   → Dùng trong game loop cập nhật hàng trăm entity
//   → Lưu hướng dạng Vec2 thay vì float32 angle để tránh convert mỗi frame
//
// Pattern hiệu quả cho entity có hướng:
//
//   type Enemy struct {
//       Pos Vec2
//       Dir Vec2    // lưu dạng unit vector, không phải degree
//       Speed float32
//   }
//
//   func (e *Enemy) Update(dt float32, playerPos Vec2) {
//       toPlayer := math.AngleV(e.Pos.X, e.Pos.Y, playerPos.X, playerPos.Y)
//       e.Dir = math.SlerpDir(e.Dir, toPlayer, 3*dt)   // xoay mượt
//       e.Pos = e.Pos.Add(e.Dir.Scale(e.Speed * dt))   // di chuyển
//   }
//
//   func (e *Enemy) AngleDeg() float32 {
//       return nmath.ToDeg(e.Dir)  // chỉ convert khi cần hiển thị
//   }
//
// Tổng kết chi phí (tất cả O(1), zero allocation):
//
//   Mul (xoay):     4 mul + 2 add    — không sin/cos
//   AngleDifV:      1 Conj + 1 Mul + 1 atan2
//   SlerpDir:       1 Conj + 2 Mul + 1 atan2 + 1 sin + 1 cos
//   IsInFOV:        1 Dot + 1 cos (precompute halfFOV nếu fixed)
//   Reflect:        1 Dot + 4 mul + 2 sub
