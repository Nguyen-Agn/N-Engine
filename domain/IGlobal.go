package domain

// IGlobal định nghĩa giao diện cho Global Resource Store — nơi lưu trữ tập trung
// các tài nguyên (Sprite, Audio, Object) và biến chia sẻ (Variable, Constant)
// có thể truy cập từ bất kỳ đâu trong toàn bộ hệ thống Game.
type IGlobal interface {
	// GetSprite trả về ISpriteLW theo key. Trả về nil nếu key không tồn tại.
	GetSprite(key string) ISpriteLW
	// AddSprite thêm mới hoặc ghi đè một Sprite theo key.
	AddSprite(key string, sprite ISpriteLW)
	// UpdateSprite cập nhật Sprite đã có theo key.
	// Hoạt động tương tự Add nếu key chưa tồn tại.
	UpdateSprite(key string, sprite ISpriteLW)

	// GetAudio trả về IAudioLW theo key. Trả về nil nếu key không tồn tại.
	GetAudio(key string) IAudioLW
	// AddAudio thêm mới hoặc ghi đè một Audio theo key.
	AddAudio(key string, audio IAudioLW)
	// UpdateAudio cập nhật Audio đã có theo key.
	UpdateAudio(key string, audio IAudioLW)

	// GetObject trả về IObject theo key. Trả về nil nếu key không tồn tại.
	GetObject(key string) IObject
	// AddObject thêm mới hoặc ghi đè một Object theo key.
	AddObject(key string, object IObject)
	// UpdateObject cập nhật Object đã có theo key.
	UpdateObject(key string, object IObject)

	// ShareGlobal nhúng giao diện chia sẻ biến key-value có kiểu.
	ShareGlobal

	// UpdateConst cập nhật giá trị hằng số đã có theo key.
	UpdateConst(key string, value any)
}

// ShareGlobal định nghĩa giao diện đọc/ghi biến key-value với kiểu dữ liệu cụ thể.
// Được nhúng vào IGlobal và IGlobalConfig để chia sẻ khả năng lưu trữ biến động.
type ShareGlobal interface {
	// SetValue lưu một giá trị bất kỳ theo key. Ghi đè nếu key đã tồn tại.
	SetValue(key string, value any)

	// GetInt trả về giá trị kiểu int theo key. Trả về 0 nếu key không tồn tại.
	GetInt(key string) int

	// GetInt64 trả về giá trị kiểu int64 theo key. Trả về 0 nếu key không tồn tại.
	GetInt64(key string) int64

	// GetString trả về giá trị kiểu string theo key. Trả về "" nếu key không tồn tại.
	GetString(key string) string

	// GetFloat32 trả về giá trị kiểu float32 theo key. Trả về 0.0 nếu key không tồn tại.
	GetFloat32(key string) float32

	// GetFloat64 trả về giá trị kiểu float64 theo key. Trả về 0.0 nếu key không tồn tại.
	GetFloat64(key string) float64

	// GetBool trả về giá trị kiểu bool theo key. Trả về false nếu key không tồn tại.
	GetBool(key string) bool

	// GetConst trả về giá trị hằng số theo key. Trả về nil nếu key không tồn tại.
	GetConst(key string) any
	// NewConst khai báo hoặc ghi đè một hằng số theo key.
	NewConst(key string, value any) bool
}
