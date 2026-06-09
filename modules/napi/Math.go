package napi

import "autoworld/modules/nmath"

// Math cung cấp các hàm hỗ trợ tính toán 2D đặc thù cho game (Vector, Góc, Lưới, Random, v.v.).
// Được viết chuyên biệt cho float32 để không phải ép kiểu từ float64 của Go.
var Math = nmath.NewMathHelper()

// =============================================================================
// Wrapper cho các hàm Generic (Go không hỗ trợ Generic Method)
// =============================================================================

// Choose chọn ngẫu nhiên một phần tử. Hoạt động với mọi kiểu dữ liệu (zero alloc).
func Choose[T any](args ...T) T {
	return nmath.Choose(&Math, args...)
}

// ChooseWeighted chọn ngẫu nhiên có trọng số.
func ChooseWeighted[T any](weights []float32, args ...T) T {
	return nmath.ChooseWeighted(&Math, weights, args...)
}

// ChooseFromArray chọn ngẫu nhiên một phần tử từ slice.
func ChooseFromArray[T any](arr []T) (T, bool) {
	return nmath.ChooseFromArray(&Math, arr)
}

// ChooseFromArrayN chọn n phần tử ngẫu nhiên không lặp lại từ slice.
func ChooseFromArrayN[T any](arr []T, n int) []T {
	return nmath.ChooseFromArrayN(&Math, arr, n)
}

// Shuffle xáo trộn slice tại chỗ.
func Shuffle[T any](arr []T) {
	nmath.Shuffle(&Math, arr)
}
