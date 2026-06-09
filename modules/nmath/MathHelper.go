package nmath

import (
	"math"
	"math/rand/v2"
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

// Hàm khống chế giá trị
func (m MathHelper) Clamp(val, min, max float32) float32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (m MathHelper) Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

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

// Trả về 1.0 hoặc -1.0, bỏ qua zero — dùng cho normalize hướng nhanh
func (m MathHelper) Sign(val float64) int {
	if float32(val) < 0 {
		return -1
	}
	if float32(val) > 0 {
		return 1
	}
	return 0
}

// Hàm lưới tọa độ
func (m MathHelper) Snap(val, step float32) float32 {
	return float32(int32(val/step+0.5)) * step
}

// --- Helper convert ----------------------------------------------------------

// DirFromDeg tạo unit vector từ góc degree.
// Chỉ gọi một lần khi cần tạo hướng từ số degree thô.
// Sau đó dùng Vec2.Mul để xoay — không cần gọi lại sin/cos.
func (m MathHelper) DirFromDeg(deg float32) Vec2 {
	r := float64(deg) * deg2rad
	return Vec2{float32(math.Cos(r)), float32(math.Sin(r))}
}

// ToDeg chuyển unit vector thành góc degree [-180, 180].
func (m MathHelper) ToDeg(v Vec2) float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X))) * rad2deg
}

// =============================================================================
// LengthDir — tính toạ độ từ độ dài và hướng
// =============================================================================

// LengthDirX trả về thành phần X khi đi theo hướng direction (degree) một đoạn length.
// Tương đương: length × cos(direction).
func (m MathHelper) LengthDirX(length, direction float32) float32 {
	return length * float32(math.Cos(float64(direction)*deg2rad))
}

// LengthDirY trả về thành phần Y khi đi theo hướng direction (degree) một đoạn length.
// Tương đương: length × sin(direction).
// Ghi chú: trục Y hướng xuống trong game 2D thông thường — điều chỉnh dấu nếu cần.
func (m MathHelper) LengthDirY(length, direction float32) float32 {
	return length * float32(math.Sin(float64(direction)*deg2rad))
}

// LengthDirV trả về Vec2 từ độ dài và unit vector hướng.
// Không cần sin/cos — chỉ scale trực tiếp.
//
// Ví dụ dùng:
//
//	dir := nmath.DirFromDeg(45)               // tạo một lần
//	vel := m.LengthDirV(speed, dir)           // dùng mỗi frame, không sin/cos
func (m MathHelper) LengthDirV(length float32, dir Vec2) Vec2 {
	return Vec2{dir.X * length, dir.Y * length}
}

// =============================================================================
// Distance — khoảng cách giữa hai điểm
// =============================================================================

// Distance trả về khoảng cách Euclidean giữa hai điểm.
// Dùng cho giá trị thực sự cần khoảng cách (UI, audio falloff...).
// Nếu chỉ cần so sánh "ai gần hơn" → dùng DistanceSq, tránh Sqrt.
func (m MathHelper) Distance(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// DistanceSq trả về bình phương khoảng cách — không dùng Sqrt.
// Nhanh hơn Distance ~3× trong vòng lặp so sánh.
//
// Ví dụ dùng:
//
//	if m.DistanceSq(px, py, ex, ey) < radius*radius {
//	    // trong vùng aggro
//	}
func (m MathHelper) DistanceSq(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func (m MathHelper) CompareDistance(x1, y1, x2, y2, x3, y3, x4, y4 float32) int {
	dx1, dy1 := x2-x1, y2-y1
	dx2, dy2 := x4-x3, y4-y3
	return m.Sign(float64(dx1*dx1 + dy1*dy1 - dx2*dx2 - dy2*dy2))
}

// =============================================================================
// Angle — góc hướng giữa hai điểm
// =============================================================================

// Angle trả về góc degree từ điểm (x1,y1) đến điểm (x2,y2).
// Kết quả trong [-180, 180]. Dùng atan2 — xử lý đúng mọi góc phần tư.
func (m MathHelper) Angle(x1, y1, x2, y2 float32) float32 {
	return float32(math.Atan2(float64(y2-y1), float64(x2-x1))) * rad2deg
}

// AngleV trả về unit vector hướng từ điểm 1 đến điểm 2.
// Dùng kết quả này trực tiếp với LengthDirV / RotateV để tránh convert qua lại.
func (m MathHelper) AngleV(x1, y1, x2, y2 float32) Vec2 {
	return Vec2{x2 - x1, y2 - y1}.Norm()
}

// =============================================================================
// AngleDif — khoảng cách ngắn nhất giữa hai góc
// =============================================================================

// AngleDif trả về khoảng cách góc ngắn nhất từ angle1 đến angle2, đơn vị degree.
// Kết quả trong [-180, 180]:
//   - Dương: angle2 ở phía ngược chiều kim đồng hồ so với angle1
//   - Âm: angle2 ở phía cùng chiều kim đồng hồ so với angle1
//
// Dùng số phức thay vì (angle2 - angle1 + 360) % 360:
// tránh cần Wrap thêm, không bao giờ sai góc phần tư.
func (m MathHelper) AngleDif(angle1, angle2 float32) float32 {
	v1 := m.DirFromDeg(angle1)
	v2 := m.DirFromDeg(angle2)
	rel := v1.Conj().Mul(v2)
	// rel.X = cos(dif), rel.Y = sin(dif)
	// atan2 tự trả về [-π, π] — không cần Wrap thêm
	return float32(math.Atan2(float64(rel.Y), float64(rel.X))) * rad2deg
}

// AngleDifV tính khoảng cách góc giữa hai unit vector.
// Không tạo Vec2 trung gian từ degree — dùng cho hot path.
//
// Ví dụ dùng:
//
//	forward := nmath.DirFromDeg(enemyAngle)
//	toPlayer := m.AngleV(ex, ey, px, py)
//	dif := m.AngleDifV(forward, toPlayer)  // enemy cần xoay bao nhiêu độ
func (m MathHelper) AngleDifV(v1, v2 Vec2) float32 {
	rel := v1.Norm().Conj().Mul(v2.Norm())
	return float32(math.Atan2(float64(rel.Y), float64(rel.X))) * rad2deg
}

// =============================================================================
// AngleAdd — cộng góc và tự wrap
// =============================================================================

// AngleAdd cộng amount (degree) vào angle và giữ kết quả trong [0, 360).
// Dùng số phức nội bộ — không cần gọi math.Mod thủ công.
func (m MathHelper) AngleAdd(angle, amount float32) float32 {
	dir := m.DirFromDeg(angle)
	rot := m.DirFromDeg(amount)
	return m.ToDeg(dir.Mul(rot))
	// ToDeg trả về [-180,180]; nếu cần [0,360] thêm:
	// if result < 0 { result += 360 }
}

// AngleAddV cộng góc degree vào unit vector hướng.
// Kết quả là unit vector mới — không cần Normalize lại.
//
// Ví dụ dùng:
//
//	bullet.Dir = m.AngleAddV(bullet.Dir, 5) // lệch 5 độ mỗi frame
func (m MathHelper) AngleAddV(dir Vec2, amountDeg float32) Vec2 {
	rot := m.DirFromDeg(amountDeg)
	return dir.Mul(rot)
}

// =============================================================================
// RotateV — xoay vector quanh gốc toạ độ
// =============================================================================

// RotateV xoay vec theo hướng dir (unit vector).
// Không cần sin/cos — chỉ nhân số phức.
//
// Ví dụ dùng:
//
//	forward  := nmath.DirFromDeg(90)          // hướng nhìn
//	rightward := m.RotateV(forward, nmath.DirFromDeg(-90)) // vector vuông phải
func (m MathHelper) RotateV(vec, dir Vec2) Vec2 {
	return vec.Mul(dir)
}

// =============================================================================
// SlerpDir — nội suy hướng mượt trên vòng tròn đơn vị
// =============================================================================

// SlerpDir nội suy hướng từ from đến to theo tham số t ∈ [0, 1].
// Kết quả luôn là unit vector — không cần Normalize lại.
// Mượt hơn Lerp vì đi theo cung tròn, không đi thẳng (không bị scale ngắn lại ở giữa).
//
// Dùng cho: enemy xoay về phía player, turret tracking, camera smoothing.
//
// Ví dụ dùng:
//
//	enemy.Dir = m.SlerpDir(enemy.Dir, toPlayer, 0.05) // xoay 5% mỗi frame
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

// SlerpDirDeg nội suy hướng dạng degree — wrapper tiện lợi cho SlerpDir.
//
// Ví dụ dùng:
//
//	enemy.Angle = m.SlerpDirDeg(enemy.Angle, targetAngle, 0.05)
func (m MathHelper) SlerpDirDeg(fromDeg, toDeg, t float32) float32 {
	from := m.DirFromDeg(fromDeg)
	to := m.DirFromDeg(toDeg)
	return m.ToDeg(m.SlerpDir(from, to, t))
}

// =============================================================================
// IsInFOV — kiểm tra điểm có trong vùng tầm nhìn không
// =============================================================================

// IsInFOV kiểm tra liệu điểm target có nằm trong góc nhìn (FOV) của observer không.
//   - forward: hướng nhìn (unit vector)
//   - halfFOVDeg: nửa góc tầm nhìn, ví dụ 45 = FOV 90°
//
// Dùng Dot thay vì tính góc thực sự — không cần atan2.
//
// Ví dụ dùng:
//
//	forward := nmath.DirFromDeg(enemy.Angle)
//	if m.IsInFOV(ex, ey, forward, px, py, 45) {
//	    // player trong tầm nhìn
//	}
func (m MathHelper) IsInFOV(ox, oy float32, forward Vec2, tx, ty, halfFOVDeg float32) bool {
	toTarget := Vec2{tx - ox, ty - oy}.Norm()
	cosHalf := float32(math.Cos(float64(halfFOVDeg) * deg2rad))
	// Dot >= cos(halfFOV) ↔ góc <= halfFOV
	return forward.Dot(toTarget) >= cosHalf
}

// =============================================================================
// Reflect — phản xạ vector qua pháp tuyến
// =============================================================================

// Reflect tính vector phản xạ của v qua normal n (unit vector).
// Dùng cho: bounce projectile, phản xạ ánh sáng, wall bouncing.
//
// Công thức: v - 2(v·n)n
//
// Ví dụ dùng:
//
//	bullet.Dir = m.Reflect(bullet.Dir, wallNormal)
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
