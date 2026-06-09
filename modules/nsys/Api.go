package nsys

// ConstInt retrieves a constant value as an integer using its key.
// Inputs: key - string identifier for the constant.
// Outputs: integer value if found and strictly typed, otherwise 0.
func ConstInt(key string) int {
	if val, ok := store.GetConst(key).(int); ok {
		return val
	}
	return 0
}

// ConstInt64 retrieves a constant value as an int64 using its key.
// Inputs: key - string identifier for the constant.
// Outputs: int64 value if found and strictly typed, otherwise 0.
func ConstInt64(key string) int64 {
	if val, ok := store.GetConst(key).(int64); ok {
		return val
	}
	return 0
}

// ConstFloat32 retrieves a constant value as a float32 using its key.
// Inputs: key - string identifier for the constant.
// Outputs: float32 value if found and strictly typed, otherwise 0.
func ConstFloat32(key string) float32 {
	if val, ok := store.GetConst(key).(float32); ok {
		return val
	}
	return 0
}

// ConstFloat64 retrieves a constant value as a float64 using its key.
// Inputs: key - string identifier for the constant.
// Outputs: float64 value if found and strictly typed, otherwise 0.
func ConstFloat64(key string) float64 {
	if val, ok := store.GetConst(key).(float64); ok {
		return val
	}
	return 0
}

// ConstString retrieves a constant value as a string using its key.
// Inputs: key - string identifier for the constant.
// Outputs: string value if found and strictly typed, otherwise an empty string.
func ConstString(key string) string {
	if val, ok := store.GetConst(key).(string); ok {
		return val
	}
	return ""
}

// ConstBool retrieves a constant value as a boolean using its key.
// Inputs: key - string identifier for the constant.
// Outputs: boolean value if found and strictly typed, otherwise false.
func ConstBool(key string) bool {
	if val, ok := store.GetConst(key).(bool); ok {
		return val
	}
	return false
}

// Get safely retrieves a mutable variable from the old system and casts it to type T.
// Inputs: key - string identifier.
// Outputs: the value typed as T, or zero-value of T if not found or cast fails.
func Get[T any](key string) T {
	val, _ := store.GetValue(key).(T)
	return val
}

// ConstGet safely retrieves an immutable constant from the old system and casts it to type T.
// Inputs: key - string identifier.
// Outputs: the value typed as T, or zero-value of T if not found or cast fails.
func ConstGet[T any](key string) T {
	val, _ := store.GetConst(key).(T)
	return val
}
