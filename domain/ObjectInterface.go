package domain

import (
	"image/color"

	"github.com/yohamta/donburi"
)

// ─── Core Object Interface ─────────────────────────────────────────────────────

// IObject là interface gốc của mọi thực thể (entity) trong Engine.
// Bất kỳ đối tượng nào tham gia vào vòng lặp game đều phải implement interface này.
type IObject interface {
	// OnCreate được gọi một lần duy nhất khi Object được khởi tạo vào Scene.
	// Dùng để thiết lập dữ liệu ban đầu, đăng ký event, v.v.
	// Purpose: Called once when the object is initialized into the scene. Used for setup.
	// Inputs: None.
	// Outputs: None.
	OnCreate()

	// OnStep được gọi mỗi frame để cập nhật logic của Object.
	// Purpose: Called every frame to update the object's logic.
	// Inputs: None.
	// Outputs: None.
	OnStep()

	// OnDestroy được gọi khi Object bị xóa khỏi Scene.
	// Dùng để giải phóng tài nguyên, hủy đăng ký event, v.v.
	// Purpose: Called when the object is removed from the scene. Used for cleanup.
	// Inputs: None.
	// Outputs: None.
	OnDestroy()

	// OnSave được gọi khi game thực hiện lưu lại.
	// Nhận vào một map trống để developer gắn dữ liệu cần lưu.
	// Purpose: Called during a game save event to allow the object to persist its state.
	// Inputs: data map[string]any - A map to store key-value data to be saved.
	// Outputs: None.
	OnSave(data map[string]any)

	// OnLoad được gọi khi load game tương ứng cho Object.
	// Purpose: Called during a game load event to restore the object's state.
	// Inputs: data map[string]any - The restored map containing the saved key-value data.
	// Outputs: None.
	OnLoad(data map[string]any)

	// GetPool trả về Pool quản lý Object này (nếu có).
	// Purpose: Retrieves the object pool managing this instance, if any.
	// Inputs: None.
	// Outputs: IPool - The pool manager, or nil if not pooled.
	GetPool() IPool
	
	// SetPool liên kết Object này với một Pool.
	// Purpose: Links the object to a specific object pool manager.
	// Inputs: pool IPool - The pool manager.
	// Outputs: None.
	SetPool(pool IPool)

	// Entry trả về ECS entry (donburi) của Object.
	// Dùng để đọc/ghi dữ liệu Component trong hệ thống ECS.
	// Purpose: Retrieves the underlying ECS (donburi) entry for direct component access.
	// Inputs: None.
	// Outputs: *donburi.Entry - The ECS entity entry.
	Entry() *donburi.Entry

	// Using classes to set value of each component
	// Ex: pos-x-5 to set x = 5
	// Ex: spr-c-run to set current sprite to sprite name "run"
	// Format: <component's name>-<variable's name>-<value>
	// Note that: this way only suit for integer/string type value
	// Purpose: Sets multiple component values using a tokenized string format.
	// Inputs: tokenClasses string - A space-separated string of property tokens (e.g. "pos-x-5 spr-c-run").
	// Outputs: None.
	// Special requirements: Only supports primitive types like integer or string for values.
	SetTokens(tokenClasses string)

	// Remove IObject from Game
	// Purpose: Immediately removes the object from the game and schedules it for destruction.
	// Inputs: None.
	// Outputs: None.
	Remove()
}

// IPool định nghĩa giao diện chuẩn cho mọi hệ thống Object Pool của N-Engine.
type IPool interface {
	// Put cất IObject vào pool. Trả về true nếu thành công, false nếu pool đã đầy (để Engine xóa tận gốc).
	// Purpose: Returns the object to the pool for later reuse.
	// Inputs: obj IObject - The object to return to the pool.
	// Outputs: bool - True if successfully pooled, false if the pool is full (meaning the object should be destroyed).
	Put(obj IObject) bool
}

// ─── Component Interfaces ─────────────────────────────────────────────────────
// Mỗi interface dưới đây tương ứng với một Component Data trong ObjectData.go.
// Object implement interface nào nghĩa là nó sở hữu Component đó.

// IPosition cung cấp tọa độ 2D cho Object (ánh xạ tới PositionData).
type IPosition interface {
	// X trả về tọa độ ngang hiện tại (pixel).
	// Purpose: Retrieves the current X coordinate.
	// Inputs: None.
	// Outputs: float32 - The X position in pixels.
	X() float32
	
	// Y trả về tọa độ dọc hiện tại (pixel).
	// Purpose: Retrieves the current Y coordinate.
	// Inputs: None.
	// Outputs: float32 - The Y position in pixels.
	Y() float32
	
	// SetX thiết lập tọa độ ngang.
	// Purpose: Sets the X coordinate.
	// Inputs: x float32 - The new X position.
	// Outputs: None.
	SetX(x float32)
	
	// SetY thiết lập tọa độ dọc.
	// Purpose: Sets the Y coordinate.
	// Inputs: y float32 - The new Y position.
	// Outputs: None.
	SetY(y float32)
}

// ISprite cung cấp khả năng quản lý hình ảnh và animation cho Object (ánh xạ tới SpriteData).
type ISprite interface {
	// SpriteIdx trả về index frame ảnh hiện tại trong animation.
	// Purpose: Retrieves the current animation frame index.
	// Inputs: None.
	// Outputs: int - The zero-based frame index.
	SpriteIdx() int
	
	// SetSpriteIdx thiết lập index frame ảnh hiện tại.
	// Purpose: Sets the animation frame index manually.
	// Inputs: spriteIdx int - The frame index to jump to.
	// Outputs: None.
	SetSpriteIdx(spriteIdx int)

	// ImageSpeed trả về tốc độ chạy animation (frame/tick).
	// Purpose: Retrieves the animation playback speed.
	// Inputs: None.
	// Outputs: float32 - Speed multiplier (frames per tick).
	ImageSpeed() float32
	
	// SetImageSpeed thiết lập tốc độ chạy animation.
	// Purpose: Sets the animation playback speed.
	// Inputs: imageSpeed float32 - Speed multiplier.
	// Outputs: None.
	SetImageSpeed(imageSpeed float32)

	// Rotation trả về góc xoay hiện tại của sprite (radian).
	// Purpose: Retrieves the sprite's rotation angle.
	// Inputs: None.
	// Outputs: float32 - Rotation angle in radians.
	Rotation() float32
	
	// SetRotation thiết lập góc xoay của sprite.
	// Purpose: Sets the sprite's rotation angle.
	// Inputs: rotation float32 - Angle in radians.
	// Outputs: None.
	SetRotation(rotation float32)

	// OffsetX trả về độ lệch ngang khi vẽ sprite so với tọa độ Object.
	// Purpose: Retrieves the horizontal draw offset.
	// Inputs: None.
	// Outputs: float32 - X offset in pixels.
	OffsetX() float32
	
	// SetOffsetX thiết lập độ lệch ngang khi vẽ.
	// Purpose: Sets the horizontal draw offset.
	// Inputs: offsetX float32 - X offset in pixels.
	// Outputs: None.
	SetOffsetX(offsetX float32)

	// OffsetY trả về độ lệch dọc khi vẽ sprite so với tọa độ Object.
	// Purpose: Retrieves the vertical draw offset.
	// Inputs: None.
	// Outputs: float32 - Y offset in pixels.
	OffsetY() float32
	
	// SetOffsetY thiết lập độ lệch dọc khi vẽ.
	// Purpose: Sets the vertical draw offset.
	// Inputs: offsetY float32 - Y offset in pixels.
	// Outputs: None.
	SetOffsetY(offsetY float32)

	// ImageColor trả về màu tô phủ lên sprite (RGBA).
	// Purpose: Retrieves the color tint applied to the sprite.
	// Inputs: None.
	// Outputs: color.RGBA - The tint color.
	ImageColor() color.RGBA
	
	// SetImageColor thiết lập màu tô phủ lên sprite.
	// Purpose: Sets the color tint for the sprite.
	// Inputs: imageColor color.RGBA - The tint color to apply.
	// Outputs: None.
	SetImageColor(imageColor color.RGBA)

	// ScaleX trả về hệ số co giãn ngang của sprite. 1.0 = kích thước gốc.
	// Purpose: Retrieves the horizontal scale factor.
	// Inputs: None.
	// Outputs: float32 - The X scale factor (1.0 is original size).
	ScaleX() float32
	
	// SetScaleX thiết lập hệ số co giãn ngang.
	// Purpose: Sets the horizontal scale factor.
	// Inputs: scaleX float32 - The X scale factor.
	// Outputs: None.
	SetScaleX(scaleX float32)

	// ScaleY trả về hệ số co giãn dọc của sprite. 1.0 = kích thước gốc.
	// Purpose: Retrieves the vertical scale factor.
	// Inputs: None.
	// Outputs: float32 - The Y scale factor (1.0 is original size).
	ScaleY() float32
	
	// SetScaleY thiết lập hệ số co giãn dọc.
	// Purpose: Sets the vertical scale factor.
	// Inputs: scaleY float32 - The Y scale factor.
	// Outputs: None.
	SetScaleY(scaleY float32)

	// Sprite trả về ISpriteLW theo tên. Trả về nil nếu không tồn tại.
	// Purpose: Retrieves a specific sprite by its name from the object's local registry.
	// Inputs: name string - The name of the sprite.
	// Outputs: ISpriteLW - The loaded sprite, or nil if not found.
	Sprite(name string) ISpriteLW
	
	// SetSprite ghi đè sprite theo tên (tạo mới nếu chưa có).
	// Purpose: Assigns or overwrites a sprite in the local registry under a specific name.
	// Inputs:
	//   - name string: The identifier for the sprite.
	//   - sprite ISpriteLW: The sprite to assign.
	// Outputs: None.
	SetSprite(name string, sprite ISpriteLW)
	
	// AddSprite thêm sprite mới. Trả về false nếu tên đã tồn tại (không ghi đè).
	// Purpose: Adds a new sprite to the local registry without overwriting existing ones.
	// Inputs:
	//   - name string: The identifier for the sprite.
	//   - sprite ISpriteLW: The sprite to add.
	// Outputs: bool - True if successfully added, false if a sprite with that name already exists.
	AddSprite(name string, sprite ISpriteLW) bool
	
	// RemoveSprite xóa sprite theo tên. Trả về false nếu không tìm thấy.
	// Purpose: Removes a sprite from the local registry by name.
	// Inputs: name string - The identifier for the sprite to remove.
	// Outputs: bool - True if successfully removed, false if not found.
	RemoveSprite(name string) bool

	// SetCurrentSprite chọn sprite đang active theo tên để DrawSystem vẽ.
	// Purpose: Sets the currently active sprite for rendering by name.
	// Inputs: name string - The name of the sprite to switch to.
	// Outputs: None.
	SetCurrentSprite(name string)
	
	// GetCurrentSprite trả về ISpriteLW của sprite đang active. Trả về nil nếu chưa chọn.
	// Purpose: Retrieves the currently active sprite being rendered.
	// Inputs: None.
	// Outputs: ISpriteLW - The active sprite, or nil if none is set.
	GetCurrentSprite() ISpriteLW

	// NextImage tăng ImageIndex lên 1, tự động wrap về 0 khi vượt quá Length.
	// Purpose: Advances the animation to the next frame, wrapping to 0 at the end.
	// Inputs: None.
	// Outputs: None.
	NextImage()
	
	// ImageIndex trả về index frame ảnh hiện tại trong sprite active.
	// Purpose: Retrieves the current frame index (same as SpriteIdx).
	// Inputs: None.
	// Outputs: int - The frame index.
	ImageIndex() int
	
	// SetImageIndex thiết lập index frame ảnh trong sprite active.
	// Purpose: Sets the current frame index manually (same as SetSpriteIdx).
	// Inputs: imageIndex int - The frame index.
	// Outputs: None.
	SetImageIndex(imageIndex int)

	// Enable / Disable 9Slice Mode
	// String Ex: "5" ->5:all, "5 6" -> 5:top&bottom, 6:right&left, "1 2 3 4" -> each
	// Purpose: Toggles 9-slice rendering mode and sets margin values via a string configuration.
	// Inputs:
	//   - turn bool: True to enable 9-slice rendering, false to disable.
	//   - TopRightBottomLeft string: Space-separated margin values (e.g., "5 5 5 5").
	// Outputs: None.
	Set9Slice(turn bool, TopRightBottomLeft string)

	// ZOrder trả về giá trị thứ tự vẽ của sprite (số nhỏ vẽ trước, số lớn vẽ sau).
	// Purpose: Retrieves the rendering Z-order.
	// Inputs: None.
	// Outputs: int - The Z-order index.
	ZOrder() int
	
	// SetZOrder thiết lập thứ tự vẽ của sprite và kích hoạt cờ cập nhật (dirty flag).
	// Purpose: Sets the rendering Z-order. Lower numbers draw first.
	// Inputs: z int - The new Z-order index.
	// Outputs: None.
	SetZOrder(z int)
}

// IBox cung cấp hitbox hình học cho Object dùng trong va chạm (ánh xạ tới BoxData).
type IBox interface {
	// BoxW trả về chiều rộng của hitbox (pixel).
	// Purpose: Retrieves the collision box width.
	// Inputs: None.
	// Outputs: float32 - Width in pixels.
	BoxW() float32
	
	// SetBoxW thiết lập chiều rộng hitbox.
	// Purpose: Sets the collision box width.
	// Inputs: boxW float32 - New width in pixels.
	// Outputs: None.
	SetBoxW(boxW float32)

	// BoxH trả về chiều cao của hitbox (pixel).
	// Purpose: Retrieves the collision box height.
	// Inputs: None.
	// Outputs: float32 - Height in pixels.
	BoxH() float32
	
	// SetBoxH thiết lập chiều cao hitbox.
	// Purpose: Sets the collision box height.
	// Inputs: boxH float32 - New height in pixels.
	// Outputs: None.
	SetBoxH(boxH float32)

	// BoxX trả về offset ngang của hitbox so với tọa độ Object.
	// Purpose: Retrieves the horizontal offset of the hitbox relative to the object's position.
	// Inputs: None.
	// Outputs: float32 - X offset in pixels.
	BoxX() float32
	
	// SetBoxX thiết lập offset ngang của hitbox.
	// Purpose: Sets the horizontal offset of the hitbox.
	// Inputs: boxX float32 - New X offset in pixels.
	// Outputs: None.
	SetBoxX(boxX float32)

	// BoxY trả về offset dọc của hitbox so với tọa độ Object.
	// Purpose: Retrieves the vertical offset of the hitbox relative to the object's position.
	// Inputs: None.
	// Outputs: float32 - Y offset in pixels.
	BoxY() float32
	
	// SetBoxY thiết lập offset dọc của hitbox.
	// Purpose: Sets the vertical offset of the hitbox.
	// Inputs: boxY float32 - New Y offset in pixels.
	// Outputs: None.
	SetBoxY(boxY float32)

	// IsCollidable trả về trạng thái kích hoạt va chạm của Object.
	// Purpose: Checks if collision detection is active for this object.
	// Inputs: None.
	// Outputs: bool - True if collidable.
	IsCollidable() bool
	
	// SetIsCollidable bật/tắt va chạm cho Object.
	// Purpose: Toggles collision detection for this object.
	// Inputs: isCollidable bool - True to enable, false to disable.
	// Outputs: None.
	SetIsCollidable(isCollidable bool)

	// Shape trả về hình dạng hitbox (rectangle hoặc circle).
	// Purpose: Retrieves the geometric shape of the collision hitbox.
	// Inputs: None.
	// Outputs: BoxShape - The hitbox shape (e.g., rectangle, circle).
	Shape() BoxShape
	
	// SetShape thiết lập hình dạng hitbox.
	// Purpose: Sets the geometric shape of the collision hitbox.
	// Inputs: shape BoxShape - The new shape to use.
	// Outputs: None.
	SetShape(shape BoxShape)
}

// IAudio cung cấp khả năng phát âm thanh cho Object (ánh xạ tới AudioData).
type IAudio interface {
	// Audio trả về IAudioLW của kênh âm thanh đang active. Trả về nil nếu chưa có.
	// Purpose: Retrieves the currently active audio object.
	// Inputs: None.
	// Outputs: IAudioLW - The active audio, or nil if none is selected.
	Audio() IAudioLW
	
	// SetAudio gắn một IAudioLW vào Object với tên kênh chỉ định.
	// Purpose: Registers a new audio track to the object under a specific name.
	// Inputs:
	//   - audioName string: The identifier for the audio.
	//   - audio IAudioLW: The audio object.
	// Outputs: None.
	SetAudio(audioName string, audio IAudioLW)

	// AudioName trả về tên kênh âm thanh đang active.
	// Purpose: Retrieves the name of the currently active audio track.
	// Inputs: None.
	// Outputs: string - The active audio name.
	AudioName() string
	
	// SetAudioName chọn kênh âm thanh active theo tên.
	// Purpose: Sets the active audio track by name.
	// Inputs: audioName string - The name of the track to activate.
	// Outputs: None.
	SetAudioName(audioName string)

	// AudioSpeed trả về tốc độ phát âm thanh (1.0 = bình thường).
	// Purpose: Retrieves the playback speed multiplier.
	// Inputs: None.
	// Outputs: float32 - The speed multiplier (1.0 is normal).
	AudioSpeed() float32
	
	// SetAudioSpeed thiết lập tốc độ phát âm thanh.
	// Purpose: Sets the playback speed multiplier.
	// Inputs: audioSpeed float32 - The new speed multiplier.
	// Outputs: None.
	SetAudioSpeed(audioSpeed float32)

	// Volume trả về âm lượng hiện tại (0.0 - 1.0).
	// Purpose: Retrieves the playback volume.
	// Inputs: None.
	// Outputs: float32 - The volume level (0.0 to 1.0).
	Volume() float32
	
	// SetVolume thiết lập âm lượng.
	// Purpose: Sets the playback volume.
	// Inputs: volume float32 - The new volume level (0.0 to 1.0).
	// Outputs: None.
	SetVolume(volume float32)

	// Pitch trả về cao độ hiện tại (1.0 = bình thường).
	// Purpose: Retrieves the playback pitch multiplier.
	// Inputs: None.
	// Outputs: float32 - The pitch multiplier (1.0 is normal).
	Pitch() float32
	
	// SetPitch thiết lập cao độ âm thanh.
	// Purpose: Sets the playback pitch multiplier.
	// Inputs: pitch float32 - The new pitch multiplier.
	// Outputs: None.
	SetPitch(pitch float32)

	// Play phát kênh âm thanh tên name với volume và pitch chỉ định.
	// AudioSystem sẽ xử lý lệnh này vào frame tiếp theo.
	// Purpose: Queues an audio track to play with specific volume and pitch.
	// Inputs:
	//   - name string: The track name.
	//   - volume float32: Volume level (0.0 to 1.0).
	//   - pitch float32: Pitch multiplier.
	// Outputs: None.
	Play(name string, volume float32, pitch float32)
	
	// PlayDefault phát kênh âm thanh tên name với volume và pitch mặc định.
	// Purpose: Queues an audio track to play with default volume and pitch.
	// Inputs: name string - The track name.
	// Outputs: None.
	PlayDefault(name string)
	
	// StopAudio dừng âm thanh đang phát theo tên.
	// Purpose: Queues a stop command for a specific audio track.
	// Inputs: name string - The track name.
	// Outputs: None.
	StopAudio(name string)
	
	// PauseAudio tạm dừng âm thanh đang phát theo tên.
	// Purpose: Queues a pause command for a specific audio track.
	// Inputs: name string - The track name.
	// Outputs: None.
	PauseAudio(name string)
	
	// ResumeAudio tiếp tục phát âm thanh đang tạm dừng theo tên.
	// Purpose: Queues a resume command for a paused audio track.
	// Inputs: name string - The track name.
	// Outputs: None.
	ResumeAudio(name string)
	
	// SetLooping bật/tắt chế độ lặp lại tự động cho âm thanh theo tên.
	// Purpose: Toggles automatic looping for a specific audio track.
	// Inputs:
	//   - name string: The track name.
	//   - loop bool: True to enable looping.
	// Outputs: None.
	SetLooping(name string, loop bool)
	
	// IsLooping kiểm tra xem âm thanh chỉ định có đang ở chế độ lặp lại không.
	// Purpose: Checks if a specific audio track is set to loop.
	// Inputs: name string - The track name.
	// Outputs: bool - True if looping is enabled.
	IsLooping(name string) bool
}

// IInfor cung cấp thông tin định danh cho Object (ánh xạ tới InforData).
type IInfor interface {
	// GetName trả về tên hiển thị hoặc nhãn đại diện của Object.
	// Purpose: Retrieves the object's human-readable name or label.
	// Inputs: None.
	// Outputs: string - The object name.
	GetName() string
	
	// GetId trả về mã định danh duy nhất (số nguyên tự tăng) của Object.
	// Purpose: Retrieves the object's unique auto-incrementing ID.
	// Inputs: None.
	// Outputs: int - The unique integer ID.
	GetId() int
	
	// AddTag thêm một tag định danh vào Object. Chuỗi sẽ được băm thành số nguyên uint64.
	// Purpose: Adds a categorical tag to the object (hashed internally).
	// Inputs: tag string - The tag string.
	// Outputs: None.
	AddTag(tag string)
	
	// HasTag kiểm tra xem Object có chứa tag này không (băm chuỗi đầu vào).
	// Purpose: Checks if the object has a specific tag.
	// Inputs: tag string - The tag string to check.
	// Outputs: bool - True if the object has the tag.
	HasTag(tag string) bool
	
	// HasTagHash kiểm tra xem Object có chứa mã băm tag này không.
	// Purpose: Checks if the object has a specific pre-hashed tag.
	// Inputs: hash uint64 - The hashed tag to check.
	// Outputs: bool - True if the object has the hashed tag.
	HasTagHash(hash uint64) bool
	
	// IsDead kiểm tra trạng thái của Object. True = không tham gia logic, chờ dọn dẹp.
	// Purpose: Checks if the object is marked as dead and awaiting removal.
	// Inputs: None.
	// Outputs: bool - True if dead.
	IsDead() bool
	
	// SetIsDead thiết lập trạng thái chết/sống của Object (dùng để hồi sinh từ Pool).
	// Purpose: Explicitly sets the dead/alive state, useful when reviving pooled objects.
	// Inputs: dead bool - The new dead state.
	// Outputs: None.
	SetIsDead(dead bool)
	
	// MarkDead đánh dấu Object là đã chết. Deferred destruction.
	// Purpose: Flags the object as dead, scheduling it for deferred removal.
	// Inputs: None.
	// Outputs: None.
	MarkDead()
	
	// SaveTag trả về mã tag để phân biệt object khi lưu game.
	// Purpose: Retrieves the object's save tag for serialization.
	// Inputs: None.
	// Outputs: string - The save tag.
	SaveTag() string
	
	// SetSaveTag thiết lập mã tag để lưu game.
	// Purpose: Sets the object's save tag for serialization.
	// Inputs: tag string - The new save tag.
	// Outputs: None.
	SetSaveTag(tag string)
}

// ICollision cung cấp khả năng xử lý va chạm.
type ICollision interface {
	// OnCollisionTag đăng ký callback, kích hoạt mỗi frame khi va chạm với Object có chứa tag chỉ định.
	// Purpose: Registers a callback triggered when colliding with another object holding a specific tag.
	// Inputs:
	//   - tag string: The target tag to check collisions against.
	//   - handler func(other IObject): The callback executed upon collision.
	// Outputs: None.
	OnCollisionTag(tag string, handler func(other IObject))
}

// IDirection cung cấp hướng di chuyển cho Object (ánh xạ tới DirectionData).
type IDirection interface {
	// Direction trả về góc hướng hiện tại của Object (đơn vị độ, 0–360).
	// Purpose: Retrieves the object's current directional angle.
	// Inputs: None.
	// Outputs: float32 - Angle in degrees (0-360).
	Direction() float32
	
	// SetDirection thiết lập góc hướng mới cho Object.
	// Purpose: Sets the object's directional angle.
	// Inputs: dir float32 - New angle in degrees.
	// Outputs: None.
	SetDirection(dir float32)
	
	// Rotate xoay hướng Object thêm delta độ (dương = theo chiều kim đồng hồ).
	// Purpose: Rotates the object's direction by a given delta.
	// Inputs: delta float32 - Degrees to rotate (positive is clockwise).
	// Outputs: None.
	Rotate(delta float32)
}

// IInput cung cấp khả năng đăng ký lắng nghe phím bấm cho Object (ánh xạ tới InputData).
type IInput interface {
	// ListenOn đăng ký một handler được gọi khi phím (hoặc nhóm phím) kích hoạt.
	// Purpose: Registers an input callback for specific keys or key groups.
	// Inputs:
	//   - key string: Key name or group (e.g., "w a s d", "alpha", "all").
	//   - eventType string: Trigger condition ("" for hold, "pressed"/"p", "released"/"r").
	//   - handler func(key string): Callback executed with the triggered key name.
	// Outputs: None.
	ListenOn(key string, eventType string, handler func(key string))
}

// IMouse cung cấp khả năng đọc trạng thái chuột và đăng ký lắng nghe nút chuột cho Object.
type IMouse interface {
	// MouseX trả về tọa độ ngang của con trỏ chuột (pixel, tính từ góc trên-trái màn hình).
	// Purpose: Retrieves the mouse cursor's X coordinate in screen space.
	// Inputs: None.
	// Outputs: int - X coordinate.
	MouseX() int
	
	// MouseY trả về tọa độ dọc của con trỏ chuột (pixel, tính từ góc trên-trái màn hình).
	// Purpose: Retrieves the mouse cursor's Y coordinate in screen space.
	// Inputs: None.
	// Outputs: int - Y coordinate.
	MouseY() int

	// WheelX trả về độ cuộn bánh xe chuột theo trục X trong frame hiện tại.
	// Dương = cuộn phải, âm = cuộn trái.
	// Purpose: Retrieves horizontal scroll wheel offset.
	// Inputs: None.
	// Outputs: float64 - X offset.
	WheelX() float64
	
	// WheelY trả về độ cuộn bánh xe chuột theo trục Y trong frame hiện tại.
	// Dương = cuộn xuống, âm = cuộn lên.
	// Purpose: Retrieves vertical scroll wheel offset.
	// Inputs: None.
	// Outputs: float64 - Y offset.
	WheelY() float64

	// ListenMouseOn đăng ký một handler được gọi theo loại sự kiện khi nút chuột kích hoạt.
	// Purpose: Registers an input callback for specific mouse buttons.
	// Inputs:
	//   - button string: Button name (e.g., "left", "right").
	//   - eventType string: Trigger condition ("" for hold, "pressed"/"p", "released"/"r").
	//   - handler func(button string): Callback executed with the triggered button name.
	// Outputs: None.
	ListenMouseOn(button string, eventType string, handler func(button string))
}

// IAlarm cung cấp bộ đếm thời gian để thực thi callback sau N frames (ánh xạ tới AlarmData).
type IAlarm interface {
	// SetAlarm thiết lập một timer với ID cụ thể. Sau khi đếm ngược hết số frames, callback sẽ được gọi.
	// Nếu gọi lại SetAlarm với cùng ID, timer cũ sẽ bị ghi đè.
	// Purpose: Schedules a callback to execute after a specified number of frames.
	// Inputs:
	//   - id string: Unique identifier for the alarm.
	//   - frames int: Number of frames to wait.
	//   - callback func(): The function to execute.
	// Outputs: None.
	// Special requirements: Overwrites any existing alarm with the same ID.
	SetAlarm(id string, frames int, callback func())
	
	// GetAlarm trả về số frames còn lại của alarm chỉ định. Trả về 0 nếu không tìm thấy hoặc đã hoàn thành.
	// Purpose: Retrieves the remaining frames for a specific alarm.
	// Inputs: id string - The alarm identifier.
	// Outputs: int - Remaining frames (0 if done or not found).
	GetAlarm(id string) int
	
	// StopAlarm hủy bỏ alarm chỉ định (callback sẽ không được gọi).
	// Purpose: Cancels a pending alarm before it triggers.
	// Inputs: id string - The alarm identifier.
	// Outputs: None.
	StopAlarm(id string)
}

// IVelocity cung cấp vật lý cơ bản cho Object (vận tốc, ma sát) (ánh xạ tới VelocityData).
type IVelocity interface {
	// Purpose: Retrieves the current X velocity.
	// Inputs: None.
	// Outputs: float32 - Velocity in X direction.
	VelocityX() float32
	// Purpose: Retrieves the current Y velocity.
	// Inputs: None.
	// Outputs: float32 - Velocity in Y direction.
	VelocityY() float32
	// Purpose: Sets the X velocity.
	// Inputs: vx float32 - New X velocity.
	// Outputs: None.
	SetVelocityX(vx float32)
	// Purpose: Sets the Y velocity.
	// Inputs: vy float32 - New Y velocity.
	// Outputs: None.
	SetVelocityY(vy float32)
	// Purpose: Sets both X and Y velocities simultaneously.
	// Inputs:
	//   - vx float32: New X velocity.
	//   - vy float32: New Y velocity.
	// Outputs: None.
	SetVelocity(vx, vy float32)
	// Purpose: Adds to the current X and Y velocities.
	// Inputs:
	//   - vx float32: Value to add to X velocity.
	//   - vy float32: Value to add to Y velocity.
	// Outputs: None.
	AddVelocity(vx, vy float32)

	// Purpose: Retrieves the friction applied per frame.
	// Inputs: None.
	// Outputs: float32 - Friction coefficient.
	Friction() float32
	// Purpose: Sets the friction applied per frame.
	// Inputs: f float32 - New friction coefficient.
	// Outputs: None.
	SetFriction(f float32)

	// Purpose: Retrieves the maximum allowed movement speed.
	// Inputs: None.
	// Outputs: float32 - Max speed limit.
	MaxSpeed() float32
	// Purpose: Sets the maximum allowed movement speed.
	// Inputs: speed float32 - New max speed limit.
	// Outputs: None.
	SetMaxSpeed(speed float32)
}

// ITween cung cấp khả năng nội suy giá trị (Lerp) mượt mà theo thời gian (ánh xạ tới TweenData).
type ITween interface {
	// TweenMove di chuyển mượt mà Object từ vị trí hiện tại đến (targetX, targetY) trong duration (frames).
	// Purpose: Interpolates the object's position smoothly over a set duration.
	// Inputs:
	//   - targetX float32: Target X coordinate.
	//   - targetY float32: Target Y coordinate.
	//   - duration int: Time to reach target in frames.
	// Outputs: None.
	TweenMove(targetX, targetY float32, duration int)
	
	// TweenScale co giãn mượt mà Object từ Scale hiện tại đến (targetScaleX, targetScaleY) trong duration.
	// Purpose: Interpolates the object's scale smoothly over a set duration.
	// Inputs:
	//   - targetScaleX float32: Target X scale.
	//   - targetScaleY float32: Target Y scale.
	//   - duration int: Time to reach target in frames.
	// Outputs: None.
	TweenScale(targetScaleX, targetScaleY float32, duration int)
	
	// TweenAlpha làm mờ/rõ mượt mà Object từ Alpha hiện tại đến targetAlpha (0-255) trong duration.
	// Purpose: Interpolates the object's transparency smoothly over a set duration.
	// Inputs:
	//   - targetAlpha uint8: Target alpha value (0-255).
	//   - duration int: Time to reach target in frames.
	// Outputs: None.
	TweenAlpha(targetAlpha uint8, duration int)
}

// IDraw is implemented by Objects that want to perform custom drawing each frame.
// The Object must also have the DrawComponent (token "drw") embedded.
// DrawSystem automatically detects IDraw and calls Draw() after all Sprite entities are rendered.
type IDraw interface {
	// Draw is called every frame by DrawSystem after Sprite rendering.
	// Call drawing methods (Rect, Circle, Text...) from the embedded DrawComponent inside here.
	// Purpose: Custom draw callback executed each frame.
	// Inputs: None.
	// Outputs: None.
	Draw()
}

// IDrawComponent defines the primitive drawing methods provided by DrawComponent.
// Dev receives these methods via embedding napi.Drw in a Custom Object.
// All coordinates are in map space — camera offset is applied automatically.
//
// Advanced path-based methods (PathFill, PathStroke) and coordinate helpers
// (ScreenX, ScreenY) are available directly on napi.Drw but are not part of
// this interface to keep domain free of ebiten/vector dependencies.
type IDrawComponent interface {
	// --- Filled shapes ---

	// Rect draws a filled rectangle at (x, y) with size (w, h).
	// Purpose: Draws a solid colored rectangle.
	// Inputs:
	//   - x float32: Top-left X coordinate.
	//   - y float32: Top-left Y coordinate.
	//   - w float32: Width.
	//   - h float32: Height.
	//   - c color.RGBA: Fill color.
	// Outputs: None.
	Rect(x, y, w, h float32, c color.RGBA)

	// Circle draws a filled circle centered at (x, y) with radius r.
	// Purpose: Draws a solid colored circle.
	// Inputs:
	//   - x float32: Center X coordinate.
	//   - y float32: Center Y coordinate.
	//   - r float32: Radius.
	//   - c color.RGBA: Fill color.
	// Outputs: None.
	Circle(x, y, r float32, c color.RGBA)

	// --- Stroke (outline) shapes ---

	// RectStroke draws a rectangle outline at (x, y) with size (w, h).
	// strokeWidth controls the border thickness in pixels.
	// Purpose: Draws a rectangular outline.
	// Inputs:
	//   - x float32: Top-left X coordinate.
	//   - y float32: Top-left Y coordinate.
	//   - w float32: Width.
	//   - h float32: Height.
	//   - c color.RGBA: Stroke color.
	//   - strokeWidth float32: Thickness of the border.
	// Outputs: None.
	RectStroke(x, y, w, h float32, c color.RGBA, strokeWidth float32)

	// CircleStroke draws a circle outline centered at (x, y) with radius r.
	// strokeWidth controls the border thickness in pixels.
	// Purpose: Draws a circular outline.
	// Inputs:
	//   - x float32: Center X coordinate.
	//   - y float32: Center Y coordinate.
	//   - r float32: Radius.
	//   - c color.RGBA: Stroke color.
	//   - strokeWidth float32: Thickness of the border.
	// Outputs: None.
	CircleStroke(x, y, r float32, c color.RGBA, strokeWidth float32)

	// Line draws a straight line from (x0, y0) to (x1, y1).
	// strokeWidth controls the line thickness in pixels.
	// Purpose: Draws a straight line segment.
	// Inputs:
	//   - x0 float32: Start X.
	//   - y0 float32: Start Y.
	//   - x1 float32: End X.
	//   - y1 float32: End Y.
	//   - c color.RGBA: Line color.
	//   - strokeWidth float32: Line thickness.
	// Outputs: None.
	Line(x0, y0, x1, y1 float32, c color.RGBA, strokeWidth float32)

	// --- Text ---

	// SetTextAlign cấu hình căn lề cho các lệnh vẽ chữ tiếp theo.
	// Hỗ trợ truyền "left"/"l", "center"/"c", "right"/"r", "justify"/"j".
	// Purpose: Configures text alignment for subsequent text drawing calls.
	// Inputs:
	//   - hAlign string: Horizontal alignment (e.g. "center").
	//   - vAlign string: Vertical alignment.
	// Outputs: None.
	SetTextAlign(hAlign, vAlign string)

	// SetTextOverflow cấu hình giới hạn kích thước khung chữ và cách xử lý tràn.
	// Truyền 0 cho maxWidth/maxHeight nếu không muốn giới hạn.
	// mode hỗ trợ "visible"/"v", "hidden"/"h", "scale"/"s".
	// Purpose: Configures a bounding box and overflow behavior for text rendering.
	// Inputs:
	//   - maxWidth float32: Maximum allowed width (0 for unlimited).
	//   - maxHeight float32: Maximum allowed height (0 for unlimited).
	//   - mode string: Overflow mode ("visible", "hidden", "scale").
	// Outputs: None.
	SetTextOverflow(maxWidth, maxHeight float32, mode string)

	// Text draws a string at (x, y) with the default font and color c.
	// Use napi.SetDefaultFont() to change the engine-wide font.
	// Purpose: Draws a text string at a specified location.
	// Inputs:
	//   - text string: The string to draw.
	//   - x float32: X coordinate.
	//   - y float32: Y coordinate.
	//   - c color.RGBA: Text color.
	// Outputs: None.
	Text(text string, x, y float32, c color.RGBA)

	// TextEx draws a string at (x, y) with a uniform scale applied.
	// scale 1.0 = default size, 2.0 = double size.
	// Purpose: Draws a text string with an additional scaling factor.
	// Inputs:
	//   - text string: The string to draw.
	//   - x float32: X coordinate.
	//   - y float32: Y coordinate.
	//   - c color.RGBA: Text color.
	//   - scale float64: Uniform scale multiplier.
	// Outputs: None.
	TextEx(text string, x, y float32, c color.RGBA, scale float64)

	// --- Image ---

	// Image draws frame idx of the given ISpriteLW at (x, y).
	// Allows manual sprite rendering without a SpriteComponent.
	// Purpose: Manually draws a specific frame of a sprite.
	// Inputs:
	//   - sprite ISpriteLW: The sprite to render.
	//   - idx int: The frame index.
	//   - x float32: X coordinate.
	//   - y float32: Y coordinate.
	// Outputs: None.
	Image(sprite ISpriteLW, idx int, x, y float32)
}
