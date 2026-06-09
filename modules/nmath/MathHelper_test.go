package nmath

import (
	"math"
	"testing"
)

// Helper function để so sánh số thực float32 với sai số epsilon
func floatEq(a, b float32) bool {
	const epsilon = 1e-4
	return float32(math.Abs(float64(a-b))) < epsilon
}

func TestClamp(t *testing.T) {
	m := MathHelper{}
	if got := m.Clamp(150, 0, 100); got != 100 {
		t.Errorf("Clamp(150, 0, 100) = %v; want 100", got)
	}
	if got := m.Clamp(-50, 0, 100); got != 0 {
		t.Errorf("Clamp(-50, 0, 100) = %v; want 0", got)
	}
	if got := m.Clamp(50, 0, 100); got != 50 {
		t.Errorf("Clamp(50, 0, 100) = %v; want 50", got)
	}
}

func TestLerp(t *testing.T) {
	m := MathHelper{}
	if got := m.Lerp(0, 100, 0.5); !floatEq(got, 50) {
		t.Errorf("Lerp(0, 100, 0.5) = %v; want 50", got)
	}
}

func TestApproach(t *testing.T) {
	m := MathHelper{}
	if got := m.Approach(10, 100, 5); got != 15 {
		t.Errorf("Approach(10, 100, 5) = %v; want 15", got)
	}
	if got := m.Approach(98, 100, 5); got != 100 {
		t.Errorf("Approach(98, 100, 5) = %v; want 100", got)
	}
	if got := m.Approach(100, 10, 5); got != 95 {
		t.Errorf("Approach(100, 10, 5) = %v; want 95", got)
	}
}

func TestWrap(t *testing.T) {
	m := MathHelper{}
	if got := m.Wrap(370, 0, 360); !floatEq(got, 10) {
		t.Errorf("Wrap(370, 0, 360) = %v; want 10", got)
	}
	if got := m.Wrap(-10, 0, 360); !floatEq(got, 350) {
		t.Errorf("Wrap(-10, 0, 360) = %v; want 350", got)
	}
}



func TestSnap(t *testing.T) {
	m := MathHelper{}
	if got := m.Snap(22.5, 16); got != 16 {
		t.Errorf("Snap(22.5, 16) = %v; want 16", got)
	}
	if got := m.Snap(24.0, 16); got != 32 { // >= 24 rounds to 32? Wait, 24/16 = 1.5 -> int(2.0) = 2 * 16 = 32
		t.Errorf("Snap(24.0, 16) = %v; want 32", got)
	}
}

func TestGeometryAngles(t *testing.T) {
	m := MathHelper{}
	
	// Distance
	if got := m.Distance(0, 0, 3, 4); !floatEq(got, 5) {
		t.Errorf("Distance(0,0,3,4) = %v; want 5", got)
	}
	
	// Angle
	if got := m.Angle(0, 0, 1, 1); !floatEq(got, 45) {
		t.Errorf("Angle(0,0,1,1) = %v; want 45", got)
	}
	
	// AngleDif
	if got := m.AngleDif(350, 10); !floatEq(got, 20) {
		t.Errorf("AngleDif(350, 10) = %v; want 20", got)
	}
	
	// AngleAdd
	if got := m.AngleAdd(350, 30); !floatEq(got, 20) {
		t.Errorf("AngleAdd(350, 30) = %v; want 20", got)
	}
	
	// LengthDirX / Y
	if got := m.LengthDirX(100, 0); !floatEq(got, 100) {
		t.Errorf("LengthDirX(100, 0) = %v; want 100", got)
	}
	if got := m.LengthDirY(100, 90); !floatEq(got, 100) {
		t.Errorf("LengthDirY(100, 90) = %v; want 100", got)
	}
}

func TestAdvancedGeometry(t *testing.T) {
	m := MathHelper{}

	// DistanceSq
	if got := m.DistanceSq(0, 0, 3, 4); got != 25 {
		t.Errorf("DistanceSq(0,0,3,4) = %v; want 25", got)
	}

	// CompareDistance (0,0 to 3,4 is 5) vs (0,0 to 0,10 is 10) -> 5 < 10 -> -1
	if got := m.CompareDistance(0, 0, 3, 4, 0, 0, 0, 10); got >= 0 {
		t.Errorf("CompareDistance(5 vs 10) = %v; want negative", got)
	}

	// IsInFOV
	forward := m.DirFromDeg(0) // Pointing right (1, 0)
	if !m.IsInFOV(0, 0, forward, 10, 5, 45) { // 10,5 is within 45 deg? atan(5/10) = 26 deg. Yes.
		t.Errorf("IsInFOV expected true for (10,5) with 45 deg FOV")
	}
	if m.IsInFOV(0, 0, forward, 0, 10, 45) { // 0,10 is 90 deg. No.
		t.Errorf("IsInFOV expected false for (0,10) with 45 deg FOV")
	}

	// Reflect
	v := Vec2{X: 1, Y: -1} // Hitting floor moving down-right
	n := Vec2{X: 0, Y: -1} // Floor normal pointing up (y is down in game so -1 is up?)
	r := m.Reflect(v, n)
	if !floatEq(r.X, 1) || !floatEq(r.Y, 1) {
		t.Errorf("Reflect({1,-1}, {0,-1}) = %v; want {1,1}", r)
	}
	
	// SlerpDir
	dir1 := m.DirFromDeg(0)
	dir2 := m.DirFromDeg(90)
	slerp := m.SlerpDir(dir1, dir2, 0.5)
	if !floatEq(m.ToDeg(slerp), 45) {
		t.Errorf("SlerpDir 50%% between 0 and 90 = %v degree; want 45", m.ToDeg(slerp))
	}
}
