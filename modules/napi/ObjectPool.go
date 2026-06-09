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

// NewObjectPool creates a new object storage pool.
//
// Purpose: Initializes an object pool based on the provided configuration, allowing objects to be reused to minimize garbage collection.
//
// Inputs:
// - config (PoolConfig[T]): The configuration object detailing how to instantiate and reset objects.
//
// Outputs:
// - *ObjectPool[T]: The newly created object pool.
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

// Get retrieves an object from the pool, creates a new one if empty, resets its state, and registers it to the scene.
//
// Purpose: Provides an active object instance ready for use in the game.
//
// Inputs:
// - sceneName (string): The scene to register the retrieved object to.
//
// Outputs:
// - T: The active game object.
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

// Put stores an object back into the pool for future reuse.
//
// Purpose: Recycles an object. Returns true if stored successfully, or false if the pool has reached MaxSize (indicating the engine should permanently delete it). Note: Game developers usually do not need to call this manually; Engine handles it on object removal.
//
// Inputs:
// - obj (domain.IObject): The object to store.
//
// Outputs:
// - bool: True if pooled, false if the pool was full.
func (p *ObjectPool[T]) Put(obj domain.IObject) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.items) >= p.config.MaxSize {
		return false // Pool đầy -> True Deletion
	}

	p.items = append(p.items, obj.(T))
	return true
}

// Release safely removes the object from the active scene and schedules it for return to the pool.
//
// Purpose: A convenience method functionally identical to napi.Obj.Remove(obj). Marks the object for deferred routing back to this pool at the end of the frame.
//
// Inputs:
// - obj (T): The object to release.
func (p *ObjectPool[T]) Release(obj T) {
	Obj.Remove(obj)
}
