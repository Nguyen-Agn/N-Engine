package napi

import (
	"testing"
)

func TestNapiMath(t *testing.T) {
	// Test Math.Clamp (direct wrapper)
	if got := Math.Clamp(10, 0, 5); got != 5 {
		t.Errorf("Math.Clamp failed")
	}

	// Test Choose
	Math.SetSeed(1) // Fixed seed for reproducibility
	val := Choose(1, 2, 3)
	if val != 1 && val != 2 && val != 3 {
		t.Errorf("Choose returned unexpected value: %v", val)
	}

	// Test ChooseFromArray
	arr := []string{"A", "B", "C"}
	v, ok := ChooseFromArray(arr)
	if !ok || (v != "A" && v != "B" && v != "C") {
		t.Errorf("ChooseFromArray failed")
	}

	// Test ChooseFromArrayN
	arrN := ChooseFromArrayN(arr, 2)
	if len(arrN) != 2 {
		t.Errorf("ChooseFromArrayN returned wrong size")
	}

	// Test Shuffle
	Shuffle(arr)
	if len(arr) != 3 {
		t.Errorf("Shuffle changed array size")
	}
}
