package domain

// ISpriteLoader định nghĩa giao ước tải hình ảnh từ file thành ISpriteLW.
type ISpriteLoader interface {
	// LoadSingle tải một ảnh đơn từ đường dẫn path, trả về ISpriteLW gồm 1 frame.
	LoadSingle(path string) (ISpriteLW, error)
	// LoadSheet tải một spritesheet và cắt thành nhiều frame theo lưới (cols x rows).
	// gapX, gapY là khoảng cách (pixel) giữa các frame trên sheet.
	LoadSheet(path string, frameW, frameH, cols, rows, gapX, gapY int) (ISpriteLW, error)
}

// IAudioLoader định nghĩa giao ước tải một file âm thanh từ đĩa thành IAudioLW.
type IAudioLoader interface {
	// Load tải file âm thanh tại đường dẫn path và trả về IAudioLW sẵn sàng sử dụng.
	Load(path string) (IAudioLW, error)
}

// IManifestLoader định nghĩa giao ước đọc file manifest (JSON) và
// nạp toàn bộ tài nguyên được khai báo trong đó vào một IGlobal store.
type IManifestLoader interface {
	// LoadFromFile đọc file manifest tại đường dẫn filePath,
	// parse nội dung, load từng tài nguyên và lưu vào store.
	// Trả về lỗi nếu file không đọc được hoặc có tài nguyên nào load thất bại.
	LoadFromFile(filePath string, store IGlobal) error
}
