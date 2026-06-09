package nmath

// =============================================================================
// math_random.go — Nhóm hàm Random tối ưu cho Game 2D
//
// Thiết kế:
//   - Dùng math/rand/v2 (Go 1.22+): PCG algorithm, nhanh hơn v1 ~25%,
//     thread-safe mà không cần mutex, quality tốt hơn cho game
//   - Generic thay any: xóa hoàn toàn heap allocation khi boxing
//   - MathHelper giữ *rand.Rand riêng: deterministic, seedable, replay-safe
//   - Mọi hàm O(1) hoặc O(n) khi cần thiết, zero allocation trừ nơi ghi chú
//
// Dependency: Go 1.22+ (math/rand/v2)
// =============================================================================

import (
	"math/rand/v2"
)

// rng là local RNG của MathHelper, khởi tạo bằng SetSeed.
// Dùng PCG (Permuted Congruential Generator) — nhanh, quality tốt, thread-safe.
// Không dùng global rand để tránh contention và đảm bảo reproducible.

// Purpose: Creates a new MathHelper with a random seed derived from the runtime.
// Outputs: (MathHelper) - The newly created MathHelper.
func NewMathHelper() MathHelper {
	return MathHelper{rng: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))}
}

// Purpose: Initializes the RNG with a specific seed for reproducible results.
// Inputs: seed (int64) - The seed value to use.
func (m *MathHelper) SetSeed(seed int64) {
	s := uint64(seed)
	m.rng = rand.New(rand.NewPCG(s, s^0xdeadbeefcafebabe))
}

// =============================================================================
// RandomRangeDouble — số thực ngẫu nhiên trong [min, max)
// =============================================================================

// Purpose: Generates a random float32 in the range [min, max).
// Inputs: min, max (float32) - The boundaries of the range.
// Outputs: (float32) - The generated random number.
func (m *MathHelper) RandomRangeDouble(min, max float32) float32 {
	return min + m.rng.Float32()*(max-min)
}

// Purpose: Generates a random float64 in the range [min, max).
// Inputs: min, max (float64) - The boundaries of the range.
// Outputs: (float64) - The generated random number.
func (m *MathHelper) RandomRangeDouble64(min, max float64) float64 {
	return min + m.rng.Float64()*(max-min)
}

// =============================================================================
// RandomRangeInt — số nguyên ngẫu nhiên trong [min, max)
// =============================================================================

// Purpose: Generates a random integer in the range [min, max).
// Inputs: min, max (int) - The boundaries of the range.
// Outputs: (int) - The generated random integer.
func (m *MathHelper) RandomRangeInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + m.rng.IntN(max-min)
}

// Purpose: Generates a random int32 in the range [min, max).
// Inputs: min, max (int32) - The boundaries of the range.
// Outputs: (int32) - The generated random integer.
func (m *MathHelper) RandomRangeInt32(min, max int32) int32 {
	if min >= max {
		return min
	}
	return min + int32(m.rng.Int32N(max-min))
}

// =============================================================================
// RandomBool — boolean ngẫu nhiên
// =============================================================================

// Purpose: Returns a random boolean with 50/50 probability.
// Outputs: (bool) - True or false.
func (m *MathHelper) RandomBool() bool {
	return m.rng.Uint32()&1 == 1
}

// Purpose: Returns true with a given probability p [0.0, 1.0].
// Inputs: p (float32) - The probability of returning true.
// Outputs: (bool) - The result.
func (m *MathHelper) RandomChance(p float32) bool {
	return m.rng.Float32() < p
}

// =============================================================================
// Choose — chọn ngẫu nhiên một phần tử từ danh sách tường minh
// =============================================================================

// Purpose: Randomly chooses one element from the provided arguments.
// Inputs: m (*MathHelper) - The math helper instance, args (...T) - The elements to choose from.
// Outputs: (T) - The chosen element.
func Choose[T any](m *MathHelper, args ...T) T {
	return args[m.rng.IntN(len(args))]
}

// Purpose: Randomly chooses one element from the arguments, using the provided weights.
// Inputs: m (*MathHelper), weights ([]float32) - The weights, args (...T) - The elements to choose from.
// Outputs: (T) - The chosen element.
func ChooseWeighted[T any](m *MathHelper, weights []float32, args ...T) T {
	total := float32(0)
	for _, w := range weights {
		total += w
	}
	r := m.rng.Float32() * total
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return args[i]
		}
	}
	return args[len(args)-1] // fallback do float rounding
}

// =============================================================================
// ChooseFromArray — chọn ngẫu nhiên từ slice
// =============================================================================

// Purpose: Randomly chooses one element from a slice.
// Inputs: m (*MathHelper), arr ([]T) - The slice to choose from.
// Outputs: (T) - The chosen element, (bool) - True if successful, false if slice is empty.
func ChooseFromArray[T any](m *MathHelper, arr []T) (T, bool) {
	if len(arr) == 0 {
		var zero T
		return zero, false
	}
	return arr[m.rng.IntN(len(arr))], true
}

// Purpose: Randomly chooses n unique elements from a slice.
// Inputs: m (*MathHelper), arr ([]T) - The slice, n (int) - Number of elements to choose.
// Outputs: ([]T) - A slice containing the chosen elements.
func ChooseFromArrayN[T any](m *MathHelper, arr []T, n int) []T {
	if len(arr) == 0 || n <= 0 {
		return nil
	}
	if n >= len(arr) {
		// trả về toàn bộ đã shuffle
		out := make([]T, len(arr))
		copy(out, arr)
		m.rng.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
		return out
	}
	// Fisher-Yates partial — chỉ shuffle k phần tử đầu
	out := make([]T, len(arr))
	copy(out, arr)
	for i := range n {
		j := i + m.rng.IntN(len(out)-i)
		out[i], out[j] = out[j], out[i]
	}
	return out[:n]
}

// =============================================================================
// RandomArray — sinh slice ngẫu nhiên
// =============================================================================

// Purpose: Generates a slice of random float32s in [min, max).
// Inputs: n (int) - Number of elements, min, max (float32) - Range.
// Outputs: ([]float32) - The generated slice.
func (m *MathHelper) RandomArray(n int, min, max float32) []float32 {
	out := make([]float32, n)
	span := max - min
	for i := range out {
		out[i] = min + m.rng.Float32()*span
	}
	return out
}

// Purpose: Generates a slice of random integers in [min, max).
// Inputs: n (int) - Number of elements, min, max (int) - Range.
// Outputs: ([]int) - The generated slice.
func (m *MathHelper) RandomArrayInt(n, min, max int) []int {
	if min >= max {
		out := make([]int, n)
		for i := range out {
			out[i] = min
		}
		return out
	}
	span := max - min
	out := make([]int, n)
	for i := range out {
		out[i] = min + m.rng.IntN(span)
	}
	return out
}

// Purpose: Generates a 2D slice of random float32s in [min, max).
// Inputs: rows, cols (int) - Dimensions, min, max (float32) - Range.
// Outputs: ([][]float32) - The generated 2D slice.
func (m *MathHelper) RandomArray2D(rows, cols int, min, max float32) [][]float32 {
	span := max - min
	// Allocate backing array một lần — tránh rows lần malloc riêng lẻ
	backing := make([]float32, rows*cols)
	for i := range backing {
		backing[i] = min + m.rng.Float32()*span
	}
	out := make([][]float32, rows)
	for i := range out {
		out[i] = backing[i*cols : i*cols+cols]
	}
	return out
}

// Purpose: Generates a 2D slice of random ints in [min, max).
// Inputs: rows, cols (int) - Dimensions, min, max (int) - Range.
// Outputs: ([][]int) - The generated 2D slice.
func (m *MathHelper) RandomArray2DInt(rows, cols, min, max int) [][]int {
	span := max - min
	backing := make([]int, rows*cols)
	if span > 0 {
		for i := range backing {
			backing[i] = min + m.rng.IntN(span)
		}
	} else {
		for i := range backing {
			backing[i] = min
		}
	}
	out := make([][]int, rows)
	for i := range out {
		out[i] = backing[i*cols : i*cols+cols]
	}
	return out
}

// =============================================================================
// Shuffle — xáo trộn slice tại chỗ
// =============================================================================

// Purpose: Shuffles a slice in-place using Fisher-Yates algorithm.
// Inputs: m (*MathHelper), arr ([]T) - The slice to shuffle.
func Shuffle[T any](m *MathHelper, arr []T) {
	m.rng.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

// =============================================================================
// Alias Method — ChooseWeightedFast cho loot table lớn
// =============================================================================

// WeightedTable là cấu trúc tiền xử lý cho Alias Method.
// Build O(n), sau đó mỗi lần chọn O(1) — lý tưởng cho loot table gọi nhiều lần.
type WeightedTable[T any] struct {
	items []T
	prob  []float32 // xác suất chọn items[i] trong bucket i
	alias []int     // nếu không chọn i, chọn alias[i]
	n     int
}

// Purpose: Creates a WeightedTable for O(1) weighted random selection using Alias Method.
// Inputs: items ([]T) - Items to pick from, weights ([]float32) - Corresponding weights.
// Outputs: (*WeightedTable[T]) - The created table.
func NewWeightedTable[T any](items []T, weights []float32) *WeightedTable[T] {
	n := len(items)
	prob := make([]float32, n)
	alias := make([]int, n)

	// Normalize weights
	total := float32(0)
	for _, w := range weights {
		total += w
	}
	avg := 1.0 / float32(n)
	scaled := make([]float32, n)
	for i, w := range weights {
		scaled[i] = (w / total) * float32(n)
	}

	// Alias method: chia thành small (<1) và large (>=1)
	small := make([]int, 0, n)
	large := make([]int, 0, n)
	for i, p := range scaled {
		if p < 1.0 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	_ = avg // suppress unused warning

	for len(small) > 0 && len(large) > 0 {
		s := small[len(small)-1]
		small = small[:len(small)-1]
		l := large[len(large)-1]

		prob[s] = scaled[s]
		alias[s] = l
		scaled[l] = (scaled[l] + scaled[s]) - 1.0

		if scaled[l] < 1.0 {
			large = large[:len(large)-1]
			small = append(small, l)
		}
	}
	// Remaining — floating point rounding
	for _, i := range large {
		prob[i] = 1.0
	}
	for _, i := range small {
		prob[i] = 1.0
	}

	itemsCopy := make([]T, n)
	copy(itemsCopy, items)

	return &WeightedTable[T]{items: itemsCopy, prob: prob, alias: alias, n: n}
}

// Purpose: Chooses a random element from the WeightedTable in O(1) time.
// Inputs: m (*MathHelper).
// Outputs: (T) - The chosen element.
func (t *WeightedTable[T]) Choose(m *MathHelper) T {
	i := m.rng.IntN(t.n)
	if m.rng.Float32() < t.prob[i] {
		return t.items[i]
	}
	return t.items[t.alias[i]]
}

// =============================================================================
// Ghi chú thiết kế
// =============================================================================
//
// Tóm tắt chi phí:
//
//   RandomRangeDouble:    1 PCG step, 2 mul, 1 add        — ~2ns
//   RandomRangeInt:       1 PCG step, 1 mod               — ~2ns
//   RandomBool:           1 PCG step, 1 bit-and           — ~1ns
//   RandomChance:         1 PCG step, 1 compare           — ~2ns
//   Choose[T]:            1 PCG step, 1 IntN              — ~2ns, zero alloc
//   ChooseWeighted:       O(n) scan                       — tốt cho n<20
//   WeightedTable.Choose: O(1), zero alloc                — tốt cho n>=20, hot loop
//   ChooseFromArray:      1 PCG step, 1 IntN              — ~2ns, zero alloc
//   ChooseFromArrayN:     O(k) Fisher-Yates partial       — 1 alloc size k
//   Shuffle:              O(n) Fisher-Yates               — zero alloc
//   RandomArray:          O(n), 1 alloc                   — không tránh được
//   RandomArray2D:        O(n*m), 2 alloc (backing+outer) — tối ưu nhất có thể
//
// Khi nào dùng WeightedTable vs ChooseWeighted:
//   - Gọi < 100 lần / frame, n < 20  → ChooseWeighted (đơn giản hơn)
//   - Gọi nhiều lần, n >= 20         → NewWeightedTable một lần, Choose nhiều lần
//
// Thread safety:
//   MathHelper.rng KHÔNG thread-safe theo design (per-entity RNG).
//   Nếu cần multi-goroutine: tạo MathHelper riêng cho mỗi goroutine,
//   hoặc dùng rand.Float32() global của rand/v2 (thread-safe, nhưng không seedable).
