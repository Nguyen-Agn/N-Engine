package nmath

import (
	"testing"
)

func TestVec2Dot(t *testing.T) {
	v1 := Vec2{X: 1, Y: 0}
	v2 := Vec2{X: 0, Y: 1}
	
	if got := v1.Dot(v2); got != 0 {
		t.Errorf("Vec2.Dot() orthogonal = %v; want 0", got)
	}
	
	v3 := Vec2{X: 2, Y: 0}
	if got := v1.Dot(v3); got != 2 {
		t.Errorf("Vec2.Dot() parallel = %v; want 2", got)
	}
}

func TestVec2Norm(t *testing.T) {
	v := Vec2{X: 3, Y: 4}
	n := v.Norm()
	
	if !floatEq(n.X, 0.6) || !floatEq(n.Y, 0.8) {
		t.Errorf("Vec2.Norm() = %v; want {0.6, 0.8}", n)
	}
}

func TestVec2AddSub(t *testing.T) {
	v1 := Vec2{1, 2}
	v2 := Vec2{3, 4}
	
	sum := v1.Add(v2)
	if sum.X != 4 || sum.Y != 6 {
		t.Errorf("Vec2.Add() = %v; want {4, 6}", sum)
	}
	
	sub := v1.Sub(v2)
	if sub.X != -2 || sub.Y != -2 {
		t.Errorf("Vec2.Sub() = %v; want {-2, -2}", sub)
	}
}
