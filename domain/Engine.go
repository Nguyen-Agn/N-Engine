package domain



type Engine struct {
	Scene    ISceneManager
	Config   IGlobalConfig
	Input    IInputManager
	Store    IGlobal        // Global resource store (Singleton từ nglobal)
	// Chiều rộng cửa sổ (pixel)
	Width int
	// Chiều cao cửa sổ (pixel)
	Height int
	// Sample rate cho âm thanh (thường là 44100)
	SampleRate int
	// SaveDir là thư mục lưu trữ file save (mặc định "./saves")
	SaveDir string
	// AutoSaveVars cho biết có tự động lưu toàn bộ biến từ nsys/Global không
	AutoSaveVars bool
}
