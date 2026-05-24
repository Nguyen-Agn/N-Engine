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

func VarInt(key string) int {
	return store.GetInt(key)
}

func VarInt64(key string) int64 {
	return store.GetInt64(key)
}

func VarFloat32(key string) float32 {
	return store.GetFloat32(key)
}

func VarFloat64(key string) float64 {
	return store.GetFloat64(key)
}

func VarString(key string) string {
	return store.GetString(key)
}

func VarBool(key string) bool {
	return store.GetBool(key)
}
