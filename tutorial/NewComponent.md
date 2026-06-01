# Tạo Component Tùy biến / Custom Component

> **Vision**: 5 bước đơn giản để thêm component tùy biến vào bất kỳ object nào.
> **Vision**: 5 simple steps to add a custom component to any object.

---

## 1. Giải thích / Explanation

Engine cung cấp sẵn các component cốt lõi (Position, Sprite, Box...). Tuy nhiên, bạn thường xuyên cần lưu trữ dữ liệu riêng của mình trong ECS (Ví dụ: `HealthData`, `InventoryData`).
Custom component cho phép bạn làm điều đó một cách an toàn kiểu (type-safe) và hiệu năng cao thông qua `napi.GenericComponent`.

**5 bước / 5 steps:**
1. Định nghĩa struct dữ liệu / Define data struct
2. Khai báo Component Type bằng token / Create ComponentType token
3. Tạo wrapper struct kế thừa `napi.GenericComponent[T]` / Create wrapper
4. Nhúng (embed) vào Object / Embed into Object
5. Khởi tạo bằng `napi.NewGenericComponent` trước khi gọi `NewObject`.

---

## 2. Ví dụ Đầy đủ / Full Example

### Bước 1: Định nghĩa dữ liệu / Step 1: Define data struct
Đây là struct lưu trữ dữ liệu thực sự trong hệ thống ECS của Donburi.
```go
package game

// HealthData - dữ liệu máu của entity
type HealthData struct {
	Current int
	Max     int
	Regen   float32
}
```

### Bước 2: Tạo Component Type / Step 2: Create Component Type
Sử dụng hàm từ module `enginetype` hoặc `napi.NewComponentType` (đại diện cho API nội bộ). Bằng `autoworld/modules/enginetype`:
```go
import "autoworld/modules/napi"

// Token (chuỗi 3 ký tự) dùng cho hàm napi.Obj.NewObject() sau này
var StatsComp = napi.NewComponentType[HealthData]("sta")
```
> *Lưu ý: API nội bộ cho NewComponentType có thể nằm ở `napi.NewComponentType` hoặc `enginetype.NewComponentType`. Dưới góc độ user, ta tham chiếu qua module tương ứng.*

### Bước 3: Tạo Wrapper Struct / Step 3: Create Wrapper Struct
Tạo một Mixin bao bọc lấy generic component để viết thêm các methods tùy biến tiện ích.
```go
// StatsComponent - wrapper component với methods tiện ích
type StatsComponent struct {
	napi.GenericComponent[HealthData]
}

// Hàm nhận sát thương (sẽ tự động hiển thị cho Object nào nhúng StatsComponent)
func (s *StatsComponent) TakeDamage(dmg int) {
	// s.Get() trả về pointer đến HealthData trong ECS
	data := s.Get()
	data.Current -= dmg
	if data.Current < 0 {
		data.Current = 0
	}
}

func (s *StatsComponent) IsAlive() bool {
	return s.Get().Current > 0
}
```

### Bước 4 & 5: Nhúng và Khởi tạo / Step 4 & 5: Embed & Init

```go
type Hero struct {
	napi.IObject
	napi.Pos
	napi.Spr
	StatsComponent // Bước 4: Nhúng mixin
}

func NewHero() *Hero {
	h := &Hero{
		// Bước 5: Phải gán GenericComponent từ StatsComp trước khi gọi NewObject
		StatsComponent: StatsComponent{
			GenericComponent: napi.NewGenericComponent(StatsComp),
		},
	}

	// Đưa chuỗi "sta" vào để Engine sinh Component cho object này
	napi.Obj.NewObject(h, "MainHero", "pos spr sta sce-main")

	// Set data ban đầu
	h.Get().Max = 100
	h.Get().Current = 100
	h.Get().Regen = 1.5

	return h
}

func (h *Hero) StepUpdate() {
	h.TakeDamage(1) // Gọi method custom đã tạo ở Bước 3
	if !h.IsAlive() {
		h.Destroy() // Biến mất khỏi game
	}
}
```

Với mô hình này, dữ liệu nằm sâu dưới Donburi để đảm bảo hiệu suất, nhưng bạn vẫn code game bằng phong cách OOP (Gọi `h.TakeDamage()`) vô cùng thân thiện.
