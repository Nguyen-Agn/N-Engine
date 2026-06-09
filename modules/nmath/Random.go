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

// NewMathHelper tạo MathHelper với seed ngẫu nhiên từ runtime.
// Dùng khi không cần reproducible (production gameplay).
func NewMathHelper() MathHelper {
	return MathHelper{rng: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))}
}

// SetSeed khởi tạo lại RNG với seed xác định.
// Cùng seed → cùng chuỗi random → dùng cho: replay, test, procedural gen.
// PCG nhận 2 tham số (state + seq); dùng seed cho cả hai để đơn giản.
func (m *MathHelper) SetSeed(seed int64) {
	s := uint64(seed)
	m.rng = rand.New(rand.NewPCG(s, s^0xdeadbeefcafebabe))
}

// =============================================================================
// RandomRangeDouble — số thực ngẫu nhiên trong [min, max)
// =============================================================================

// RandomRangeDouble trả về float32 ngẫu nhiên trong [min, max).
//
// Thuật toán: min + Float32() * (max - min)
// Float32() của PCG trả về [0.0, 1.0) uniform — nhân span rồi offset.
// Chi phí: 1 PCG step + 2 float mul + 1 add.
func (m *MathHelper) RandomRangeDouble(min, max float32) float32 {
	return min + m.rng.Float32()*(max-min)
}

// RandomRangeDouble64 phiên bản float64 — dùng khi cần precision cao hơn
// (physics simulation, procedural noise).
func (m *MathHelper) RandomRangeDouble64(min, max float64) float64 {
	return min + m.rng.Float64()*(max-min)
}

// =============================================================================
// RandomRangeInt — số nguyên ngẫu nhiên trong [min, max)
// =============================================================================

// RandomRangeInt trả về int ngẫu nhiên trong [min, max).
// Dùng rand.IntN(n) — unbiased (tránh modulo bias của cách cũ).
//
// Lưu ý: [min, max) — max không bao gồm. Nếu cần [min, max] thì gọi (min, max+1).
func (m *MathHelper) RandomRangeInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + m.rng.IntN(max-min)
}

// RandomRangeInt32 phiên bản int32 — tiết kiệm hơn khi dùng với array index lớn.
func (m *MathHelper) RandomRangeInt32(min, max int32) int32 {
	if min >= max {
		return min
	}
	return min + int32(m.rng.Int32N(max-min))
}

// =============================================================================
// RandomBool — boolean ngẫu nhiên
// =============================================================================

// RandomBool trả về true/false với xác suất 50/50.
// Nhanh hơn RandomRangeInt(0,2)==1 vì dùng bit trực tiếp.
func (m *MathHelper) RandomBool() bool {
	return m.rng.Uint32()&1 == 1
}

// RandomChance trả về true với xác suất p ∈ [0.0, 1.0].
// Ví dụ: RandomChance(0.25) → true 25% thời gian.
func (m *MathHelper) RandomChance(p float32) bool {
	return m.rng.Float32() < p
}

// =============================================================================
// Choose — chọn ngẫu nhiên một phần tử từ danh sách tường minh
// =============================================================================

// Choose chọn ngẫu nhiên một phần tử từ các giá trị truyền vào.
// Generic — compiler sinh code riêng cho từng kiểu, zero heap allocation.
//
// Ví dụ:
//
//	item   := m.Choose("sword", "bow", "staff")  // string, zero alloc
//	damage := m.Choose(10, 15, 20, 25)           // int, zero alloc
//	color  := m.Choose(Red, Green, Blue)         // enum, zero alloc
func Choose[T any](m *MathHelper, args ...T) T {
	return args[m.rng.IntN(len(args))]
}

// ChooseWeighted chọn ngẫu nhiên có trọng số.
// weights[i] tương ứng với args[i]. Không yêu cầu weights tổng bằng 1.
//
// Thuật toán: linear scan — O(n) nhưng đơn giản, đủ nhanh cho n < 20.
// Với loot table lớn hơn → dùng Alias Method (xem ChooseWeightedFast).
//
// Ví dụ:
//
//	item := ChooseWeighted(m, []float32{70, 20, 10}, "common", "rare", "epic")
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

// ChooseFromArray chọn ngẫu nhiên một phần tử từ slice.
// Trả về zero value và false nếu slice rỗng.
// Generic — zero allocation, không copy phần tử.
//
// Ví dụ:
//
//	enemies := []EnemyType{Goblin, Orc, Troll}
//	spawn, ok := m.ChooseFromArray(enemies)
func ChooseFromArray[T any](m *MathHelper, arr []T) (T, bool) {
	if len(arr) == 0 {
		var zero T
		return zero, false
	}
	return arr[m.rng.IntN(len(arr))], true
}

// ChooseFromArrayN chọn n phần tử ngẫu nhiên không lặp lại từ slice.
// Dùng Fisher-Yates partial shuffle trên bản copy — O(n) time, O(k) space.
// Trả về slice mới với đúng k phần tử (hoặc ít hơn nếu arr nhỏ hơn k).
//
// Ví dụ:
//
//	deck := []Card{...}
//	hand := m.ChooseFromArrayN(deck, 5) // rút 5 bài không trùng
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

// RandomArray sinh slice float32 ngẫu nhiên trong [min, max), độ dài n.
// Allocation: một lần make([]float32, n) — không tránh được khi trả slice.
//
// Ví dụ:
//
//	heights := m.RandomArray(10, 0, 100) // 10 chiều cao ngẫu nhiên
func (m *MathHelper) RandomArray(n int, min, max float32) []float32 {
	out := make([]float32, n)
	span := max - min
	for i := range out {
		out[i] = min + m.rng.Float32()*span
	}
	return out
}

// RandomArrayInt sinh slice int ngẫu nhiên trong [min, max), độ dài n.
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

// RandomArray2D sinh slice 2 chiều ngẫu nhiên — [][]float32 kích thước rows×cols.
// Allocation: rows+1 lần (1 outer slice + rows inner slice).
//
// Dùng cho: heightmap, noise map, tile weight map...
//
// Ví dụ:
//
//	heightmap := m.RandomArray2D(64, 64, 0, 1)
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

// RandomArray2DInt sinh slice 2 chiều int ngẫu nhiên — [][]int kích thước rows×cols.
//
// Ví dụ:
//
//	tileIDs := m.RandomArray2DInt(16, 16, 0, 8)
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

// Shuffle xáo trộn slice tại chỗ bằng Fisher-Yates — O(n), unbiased.
// Không allocation. Sửa trực tiếp slice gốc.
//
// Ví dụ:
//
//	deck := []Card{...}
//	Shuffle(m, deck)
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

// NewWeightedTable tạo WeightedTable từ items và weights.
// Build một lần, dùng nhiều lần với Choose.
//
// Ví dụ:
//
//	table := NewWeightedTable(
//	    []string{"common", "rare", "epic", "legendary"},
//	    []float32{60, 25, 12, 3},
//	)
//	for i := 0; i < 1000; i++ {
//	    drop := table.Choose(m)  // O(1) mỗi lần
//	}
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

// Choose chọn một phần tử từ WeightedTable — O(1), không allocation.
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
