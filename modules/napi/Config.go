package napi

// ─── Global Config Accessors ──────────────────────────────────────────────────
// Đọc giá trị từ IGlobalConfig của Engine (key-value store có kiểu).
// Các key thường được gán lúc khởi động trong NewGame (game-title, game-width, v.v.)
// hoặc bởi game code qua engine().Config.SetValue(key, value).

// VarInt trả về giá trị kiểu int từ global config theo key. Trả về 0 nếu không có.
func VarInt(key string) int {
	return engine().Config.GetInt(key)
}

// VarInt64 trả về giá trị kiểu int64 từ global config theo key. Trả về 0 nếu không có.
func VarInt64(key string) int64 {
	return engine().Config.GetInt64(key)
}

// VarFloat32 trả về giá trị kiểu float32 từ global config theo key. Trả về 0.0 nếu không có.
func VarFloat32(key string) float32 {
	return engine().Config.GetFloat32(key)
}

// VarFloat64 trả về giá trị kiểu float64 từ global config theo key. Trả về 0.0 nếu không có.
func VarFloat64(key string) float64 {
	return engine().Config.GetFloat64(key)
}

// VarString trả về giá trị kiểu string từ global config theo key. Trả về "" nếu không có.
func VarString(key string) string {
	return engine().Config.GetString(key)
}

// VarBool trả về giá trị kiểu bool từ global config theo key. Trả về false nếu không có.
func VarBool(key string) bool {
	return engine().Config.GetBool(key)
}
