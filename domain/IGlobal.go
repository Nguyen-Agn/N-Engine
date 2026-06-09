package domain

// IGlobal định nghĩa giao diện cho Global Resource Store — nơi lưu trữ tập trung
// các tài nguyên (Sprite, Audio, Object) và biến chia sẻ (Variable, Constant)
// có thể truy cập từ bất kỳ đâu trong toàn bộ hệ thống Game.
type IGlobal interface {
	// GetSprite trả về ISpriteLW theo key. Trả về nil nếu key không tồn tại.
	// Purpose: Retrieves a lightweight sprite resource by its unique key.
	// Inputs: key string - The identifier for the sprite.
	// Outputs: ISpriteLW - The loaded sprite, or nil if the key does not exist.
	GetSprite(key string) ISpriteLW
	
	// AddSprite thêm mới hoặc ghi đè một Sprite theo key.
	// Purpose: Registers a new lightweight sprite or overwrites an existing one under the specified key.
	// Inputs:
	//   - key string: The unique identifier.
	//   - sprite ISpriteLW: The sprite object to store.
	// Outputs: None.
	AddSprite(key string, sprite ISpriteLW)
	
	// UpdateSprite cập nhật Sprite đã có theo key.
	// Hoạt động tương tự Add nếu key chưa tồn tại.
	// Purpose: Updates an existing lightweight sprite under the specified key, or adds it if it doesn't exist.
	// Inputs:
	//   - key string: The unique identifier.
	//   - sprite ISpriteLW: The sprite object to store.
	// Outputs: None.
	UpdateSprite(key string, sprite ISpriteLW)

	// GetAudio trả về IAudioLW theo key. Trả về nil nếu key không tồn tại.
	// Purpose: Retrieves an audio resource by its unique key.
	// Inputs: key string - The identifier for the audio.
	// Outputs: IAudioLW - The loaded audio, or nil if the key does not exist.
	GetAudio(key string) IAudioLW
	
	// AddAudio thêm mới hoặc ghi đè một Audio theo key.
	// Purpose: Registers a new audio resource or overwrites an existing one under the specified key.
	// Inputs:
	//   - key string: The unique identifier.
	//   - audio IAudioLW: The audio object to store.
	// Outputs: None.
	AddAudio(key string, audio IAudioLW)
	
	// UpdateAudio cập nhật Audio đã có theo key.
	// Purpose: Updates an existing audio resource under the specified key.
	// Inputs:
	//   - key string: The unique identifier.
	//   - audio IAudioLW: The audio object to store.
	// Outputs: None.
	UpdateAudio(key string, audio IAudioLW)

	// GetObject trả về IObject theo key. Trả về nil nếu key không tồn tại.
	// Purpose: Retrieves a registered global game object by its unique key.
	// Inputs: key string - The identifier for the object.
	// Outputs: IObject - The game object, or nil if the key does not exist.
	GetObject(key string) IObject
	
	// AddObject thêm mới hoặc ghi đè một Object theo key.
	// Purpose: Registers a new global game object or overwrites an existing one under the specified key.
	// Inputs:
	//   - key string: The unique identifier.
	//   - object IObject: The game object to store.
	// Outputs: None.
	AddObject(key string, object IObject)
	
	// UpdateObject cập nhật Object đã có theo key.
	// Purpose: Updates an existing global game object under the specified key.
	// Inputs:
	//   - key string: The unique identifier.
	//   - object IObject: The game object to store.
	// Outputs: None.
	UpdateObject(key string, object IObject)

	// ShareGlobal nhúng giao diện chia sẻ biến key-value có kiểu.
	ShareGlobal

	// UpdateConst cập nhật giá trị hằng số đã có theo key.
	// Purpose: Updates the value of an existing constant.
	// Inputs:
	//   - key string: The identifier of the constant.
	//   - value any: The new value to set.
	// Outputs: None.
	UpdateConst(key string, value any)
}

// ShareGlobal định nghĩa giao diện đọc/ghi biến key-value với kiểu dữ liệu cụ thể.
// Được nhúng vào IGlobal và IGlobalConfig để chia sẻ khả năng lưu trữ biến động.
type ShareGlobal interface {
	// SetValue lưu một giá trị bất kỳ theo key. Ghi đè nếu key đã tồn tại.
	// Purpose: Stores an arbitrary value under a specific key, overwriting any existing value.
	// Inputs:
	//   - key string: The identifier for the variable.
	//   - value any: The data to store.
	// Outputs: None.
	SetValue(key string, value any)

	// GetValue trả về giá trị dưới dạng any theo key.
	// Sử dụng Generic Wrapper ở API tầng trên để ép kiểu.
	// Purpose: Retrieves an arbitrary value by its key.
	// Inputs: key string - The identifier for the variable.
	// Outputs: any - The stored value, un-typed. Expects the caller to cast it.
	GetValue(key string) any


	// GetConst trả về giá trị hằng số theo key. Trả về nil nếu key không tồn tại.
	// Purpose: Retrieves a constant value by its key.
	// Inputs: key string - The identifier for the constant.
	// Outputs: any - The stored constant value, or nil if it doesn't exist.
	GetConst(key string) any
	
	// NewConst khai báo hoặc ghi đè một hằng số theo key.
	// Purpose: Declares a new constant or overwrites an existing one.
	// Inputs:
	//   - key string: The identifier for the constant.
	//   - value any: The value to store.
	// Outputs: bool - True if successfully created, false if creation failed.
	NewConst(key string, value any) bool

	// DumpVariables trả về bản sao của toàn bộ biến để lưu vào file save.
	// Purpose: Returns a copy of all current variables for saving state.
	// Inputs: None.
	// Outputs: map[string]any - A map containing all stored variables.
	DumpVariables() map[string]any
	
	// RestoreVariables khôi phục toàn bộ biến từ bản sao (file save).
	// Purpose: Restores variables from a previously dumped state map.
	// Inputs: data map[string]any - The saved map data to restore from.
	// Outputs: None.
	RestoreVariables(data map[string]any)
}
