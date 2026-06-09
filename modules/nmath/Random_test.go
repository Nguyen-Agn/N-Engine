package nmath

import (
	"testing"
)

func TestRandomRange(t *testing.T) {
	m := NewMathHelper()
	
	for i := 0; i < 100; i++ {
		got := m.RandomRangeDouble(0, 10)
		if got < 0 || got >= 10 {
			t.Errorf("RandomRangeDouble(0, 10) = %v; out of bounds", got)
		}
		
		gotInt := m.RandomRangeInt(1, 5)
		if gotInt < 1 || gotInt >= 5 {
			t.Errorf("RandomRangeInt(1, 5) = %v; out of bounds", gotInt)
		}
	}
}

func TestChoose(t *testing.T) {
	m := NewMathHelper()
	
	// Test Choose
	res := Choose(&m, 1, 2, 3)
	if res != 1 && res != 2 && res != 3 {
		t.Errorf("Choose(1,2,3) = %v; unexpected result", res)
	}
	
	// Test ChooseFromArray
	arr := []string{"A", "B"}
	val, ok := ChooseFromArray(&m, arr)
	if !ok || (val != "A" && val != "B") {
		t.Errorf("ChooseFromArray() = %v, %v; unexpected", val, ok)
	}
}

func TestShuffle(t *testing.T) {
	m := NewMathHelper()
	arr := []int{1, 2, 3, 4, 5}
	Shuffle(&m, arr)
	
	// Just check if length is same and contents are somewhat there
	if len(arr) != 5 {
		t.Errorf("Shuffle changed length to %v", len(arr))
	}
	sum := 0
	for _, v := range arr {
		sum += v
	}
	if sum != 15 {
		t.Errorf("Shuffle altered elements, sum = %v", sum)
	}
}
