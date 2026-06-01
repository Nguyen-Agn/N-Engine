package napi

// ─── Global Config Accessors ──────────────────────────────────────────────────
// Đọc giá trị từ IGlobalConfig của Engine (key-value store có kiểu).
// Các key thường được gán lúc khởi động trong NewGame (game-title, game-width, v.v.)
// hoặc bởi game code qua engine().Config.SetValue(key, value).

type storeGroup struct{}

// Store là nhóm hàm truy xuất tài nguyên, âm thanh và biến toàn cục.
// Store is the function group for accessing assets, audio, and global variables.
var Store = &storeGroup{}

// VarInt trả về giá trị kiểu int từ global config theo key. Trả về 0 nếu không có.
func (c *storeGroup) Int(key string) int {
	return engine().Config.GetInt(key)
}

// VarInt64 trả về giá trị kiểu int64 từ global config theo key. Trả về 0 nếu không có.
func (c *storeGroup) Int64(key string) int64 {
	return engine().Config.GetInt64(key)
}

// VarFloat32 trả về giá trị kiểu float32 từ global config theo key. Trả về 0.0 nếu không có.
func (c *storeGroup) Float32(key string) float32 {
	return engine().Config.GetFloat32(key)
}

// VarFloat64 trả về giá trị kiểu float64 từ global config theo key. Trả về 0.0 nếu không có.
func (c *storeGroup) Float64(key string) float64 {
	return engine().Config.GetFloat64(key)
}

// VarString trả về giá trị kiểu string từ global config theo key. Trả về "" nếu không có.
func (c *storeGroup) String(key string) string {
	return engine().Config.GetString(key)
}

// VarBool trả về giá trị kiểu bool từ global config theo key. Trả về false nếu không có.
func (c *storeGroup) VarBool(key string) bool {
	return engine().Config.GetBool(key)
}

// Setvalue cho global value
func (c *storeGroup) Value(key string, value any) {
	engine().Config.SetValue(key, value)
}

// new const
func (c *storeGroup) NewConst(key string, value any) bool {
	return engine().Config.NewConst(key, value)
}
