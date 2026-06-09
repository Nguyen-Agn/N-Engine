package napi

// ─── Global Config Accessors ──────────────────────────────────────────────────
// Đọc giá trị từ IGlobalConfig của Engine (key-value store có kiểu).
// Các key thường được gán lúc khởi động trong NewGame (game-title, game-width, v.v.)
// hoặc bởi game code qua engine().Config.SetValue(key, value).

type storeGroup struct{}

// Store là nhóm hàm truy xuất tài nguyên, âm thanh và biến toàn cục.
// Store is the function group for accessing assets, audio, and global variables.
var Store = &storeGroup{}

// =============================================================================
// Generic Wrapper cho các biến toàn cục (Thay thế toàn bộ Int, String, Float...)
// =============================================================================

// StoreGet retrieves a global value by key and safely casts it to the requested generic type.
//
// Purpose: Provides type-safe access to the engine's global configuration map.
//
// Inputs:
// - key (string): The identifier for the global variable.
//
// Outputs:
// - T: The value cast to type T, or the zero value of T if the key does not exist or casting fails.
func StoreGet[T any](key string) T {
	val, _ := engine().Config.GetValue(key).(T)
	return val
}

// Value assigns or updates a global configuration value.
//
// Purpose: Stores arbitrary data in the global configuration map using a string key.
//
// Inputs:
// - key (string): The identifier for the stored value.
// - value (any): The data to store.
func (c *storeGroup) Value(key string, value any) {
	engine().Config.SetValue(key, value)
}

// NewConst attempts to set a global configuration value permanently.
//
// Purpose: Registers a constant value in the global configuration map. If the key is already set as a constant, the assignment will fail.
//
// Inputs:
// - key (string): The identifier for the constant.
// - value (any): The immutable data to store.
//
// Outputs:
// - bool: True if the constant was successfully created, false if it already exists.
func (c *storeGroup) NewConst(key string, value any) bool {
	return engine().Config.NewConst(key, value)
}
