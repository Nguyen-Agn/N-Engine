package domain

// ISpriteLoader định nghĩa giao ước tải hình ảnh từ file thành ISpriteLW.
type ISpriteLoader interface {
	// Purpose: Loads a single image file from disk and returns a lightweight sprite containing 1 frame.
	// Inputs: path string - The file path to the image.
	// Outputs:
	//   1. ISpriteLW - The loaded lightweight sprite object.
	//   2. error - Returns an error if the file cannot be found or decoded.
	LoadSingle(path string) (ISpriteLW, error)

	// Purpose: Loads a spritesheet image and slices it into multiple frames based on grid dimensions.
	// Inputs:
	//   - path string: The file path to the spritesheet.
	//   - frameW int: Width of a single frame in pixels.
	//   - frameH int: Height of a single frame in pixels.
	//   - cols int: Number of columns in the spritesheet.
	//   - rows int: Number of rows in the spritesheet.
	//   - gapX int: Horizontal gap (in pixels) between frames.
	//   - gapY int: Vertical gap (in pixels) between frames.
	// Outputs:
	//   1. ISpriteLW - The loaded lightweight sprite containing multiple frames.
	//   2. error - Returns an error if loading or slicing fails.
	LoadSheet(path string, frameW, frameH, cols, rows, gapX, gapY int) (ISpriteLW, error)
}

// IAudioLoader định nghĩa giao ước tải một file âm thanh từ đĩa thành IAudioLW.
type IAudioLoader interface {
	// Purpose: Loads an audio file from disk into memory.
	// Inputs: path string - The file path to the audio file.
	// Outputs:
	//   1. IAudioLW - The loaded lightweight audio object.
	//   2. error - Returns an error if the file cannot be read or decoded.
	Load(path string) (IAudioLW, error)
}

// IManifestLoader định nghĩa giao ước đọc file manifest (JSON) và
// nạp toàn bộ tài nguyên được khai báo trong đó vào một IGlobal store.
type IManifestLoader interface {
	// Purpose: Parses a JSON manifest file and loads all declared resources into the provided global store.
	// Inputs:
	//   - filePath string: The file path to the manifest JSON.
	//   - store IGlobal: The global store where loaded resources will be registered.
	// Outputs:
	//   error - Returns an error if the manifest cannot be read, parsed, or if any resource fails to load.
	LoadFromFile(filePath string, store IGlobal) error
}
