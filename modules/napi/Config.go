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

// StoreGet lấy giá trị toàn cục theo key và tự ép kiểu. Trả về zero value nếu không tìm thấy.
// Ví dụ: napi.StoreGet[int]("score")
func StoreGet[T any](key string) T {
	val, _ := engine().Config.GetValue(key).(T)
	return val
}

// Setvalue cho global value
func (c *storeGroup) Value(key string, value any) {
	engine().Config.SetValue(key, value)
}

// new const
func (c *storeGroup) NewConst(key string, value any) bool {
	return engine().Config.NewConst(key, value)
}
