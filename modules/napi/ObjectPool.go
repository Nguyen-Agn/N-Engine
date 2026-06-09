package napi

import (
	"autoworld/domain"
	"sync"
)

// PoolConfig chứa các hàm khởi tạo và reset cho ObjectPool.
type PoolConfig[T domain.IObject] struct {
	// New khởi tạo Object mới hoàn toàn (thường chứa napi.Obj.NewObject).
	New func() T
	// Reset thiết lập lại trạng thái của Object (Máu, Tọa độ, IsVisible, v.v.).
	Reset func(T)
	// MaxSize giới hạn số lượng tối đa Object lưu trong kho. Vượt quá sẽ bị Engine xóa sạch.
	// Nếu <= 0, mặc định là 100.
	MaxSize int
}

// ObjectPool là hệ thống kho lưu trữ Object ECS (Zero-Allocation).
type ObjectPool[T domain.IObject] struct {
	items  []T
	config PoolConfig[T]
	mu     sync.Mutex
}

// NewObjectPool tạo một kho chứa Object mới.
func NewObjectPool[T domain.IObject](config PoolConfig[T]) *ObjectPool[T] {
	if config.New == nil {
		panic("[ObjectPool] PoolConfig.New cannot be nil")
	}
	if config.MaxSize <= 0 {
		config.MaxSize = 100
	}
	return &ObjectPool[T]{
		items:  make([]T, 0, config.MaxSize),
		config: config,
	}
}

// Get rút một Object từ kho (hoặc tạo mới nếu kho rỗng), reset trạng thái,
// và đăng ký lại vào Vòng lặp cập nhật của Scene.
func (p *ObjectPool[T]) Get(sceneName string) T {
	p.mu.Lock()
	var obj T
	n := len(p.items)
	found := false
	if n > 0 {
		obj = p.items[n-1]
		p.items = p.items[:n-1]
		found = true
	}
	p.mu.Unlock()

	// Nếu kho rỗng, tạo mới
	if !found {
		obj = p.config.New()
	}

	// Đảm bảo Object được gán cờ nhận diện từ Pool này
	if obj.GetPool() == nil {
		obj.SetPool(p)
	}

	// Dev tự reset logic game
	if p.config.Reset != nil {
		p.config.Reset(obj)
	}

	// Hồi sinh ECS Entity (bật lại cờ sống)
	if deadObj, ok := any(obj).(interface{ SetIsDead(bool) }); ok {
		deadObj.SetIsDead(false)
	}

	// Nạp lại vào LogicSystem của Scene
	Obj.Register(obj, sceneName)

	return obj
}

// Put cất Object vào kho. Hàm này trả về true nếu cất thành công.
// Trả về false nếu kho đã đầy, ra hiệu cho Engine xóa tận gốc (Tránh Memory Leak).
// Dev KHÔNG CẦN GỌI HÀM NÀY. Engine tự động gọi ngầm trong Map.go khi Dev gọi napi.Obj.Remove.
func (p *ObjectPool[T]) Put(obj domain.IObject) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.items) >= p.config.MaxSize {
		return false // Pool đầy -> True Deletion
	}

	p.items = append(p.items, obj.(T))
	return true
}

// Release là hàm thay thế cho napi.Obj.Remove(obj) (hoạt động giống y hệt).
// Dev có thể gọi pool.Release(obj) hoặc napi.Obj.Remove(obj) đều có kết quả Auto-Routing như nhau.
func (p *ObjectPool[T]) Release(obj T) {
	Obj.Remove(obj)
}
