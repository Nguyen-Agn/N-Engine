package nsys

func ConstInt(key string) int {
	if val, ok := store.GetConst(key).(int); ok {
		return val
	}
	return 0
}

func ConstInt64(key string) int64 {
	if val, ok := store.GetConst(key).(int64); ok {
		return val
	}
	return 0
}

func ConstFloat32(key string) float32 {
	if val, ok := store.GetConst(key).(float32); ok {
		return val
	}
	return 0
}

func ConstFloat64(key string) float64 {
	if val, ok := store.GetConst(key).(float64); ok {
		return val
	}
	return 0
}

func ConstString(key string) string {
	if val, ok := store.GetConst(key).(string); ok {
		return val
	}
	return ""
}

func ConstBool(key string) bool {
	if val, ok := store.GetConst(key).(bool); ok {
		return val
	}
	return false
}

// Generic Wrapper cho biến toàn cục của hệ thống nsys cũ
func Get[T any](key string) T {
	val, _ := store.GetValue(key).(T)
	return val
}

// Generic Wrapper cho Const của hệ thống nsys cũ
func ConstGet[T any](key string) T {
	val, _ := store.GetConst(key).(T)
	return val
}
