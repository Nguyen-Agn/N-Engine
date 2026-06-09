package napi

import "autoworld/modules/nmath"

// Math cung cấp các hàm hỗ trợ tính toán 2D đặc thù cho game (Vector, Góc, Lưới, Random, v.v.).
// Được viết chuyên biệt cho float32 để không phải ép kiểu từ float64 của Go.
var Math = nmath.NewMathHelper()

// =============================================================================
// Wrapper cho các hàm Generic (Go không hỗ trợ Generic Method)
// =============================================================================

// Choose selects one element randomly from the provided arguments.
//
// Purpose: Provides a generic, zero-allocation way to randomly pick a value.
//
// Outputs:
// - T: The randomly chosen element.
func Choose[T any](args ...T) T {
	return nmath.Choose(&Math, args...)
}

// ChooseWeighted selects one element randomly based on a slice of weights.
//
// Purpose: Allows probability-based random selection.
//
// Inputs:
// - weights ([]float32): The relative probabilities for each corresponding argument.
// - args (...T): The values to choose from.
//
// Outputs:
// - T: The randomly chosen element.
func ChooseWeighted[T any](weights []float32, args ...T) T {
	return nmath.ChooseWeighted(&Math, weights, args...)
}

// ChooseFromArray randomly selects one element from a slice.
//
// Outputs:
// - T: The chosen element.
// - bool: True if an element was chosen, false if the slice was empty.
func ChooseFromArray[T any](arr []T) (T, bool) {
	return nmath.ChooseFromArray(&Math, arr)
}

// ChooseFromArrayN randomly selects 'n' distinct elements from a slice.
//
// Outputs:
// - []T: A slice of the chosen elements.
func ChooseFromArrayN[T any](arr []T, n int) []T {
	return nmath.ChooseFromArrayN(&Math, arr, n)
}

// Shuffle randomly reorders the elements of a slice in place.
func Shuffle[T any](arr []T) {
	nmath.Shuffle(&Math, arr)
}
