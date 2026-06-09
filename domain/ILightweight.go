package domain

import "github.com/hajimehoshi/ebiten/v2"

// ─── Lightweight Wrapper Interfaces ───────────────────────────────────────────
// Các interface này bọc (wrap) tài nguyên nặng (ảnh, âm thanh) thành dạng
// nhẹ hơn để ECS component có thể tham chiếu mà không bị phụ thuộc thư viện.

// ISpriteLW là lightweight wrapper chứa tập hợp các frame ảnh (ebiten.Image).
// Được dùng trong SpriteData và TilemapData để lưu trữ tài nguyên đồ họa.
type ISpriteLW interface {
	// Image trả về frame ảnh tại vị trí index. Trả về nil nếu index ngoài biên.
	// Purpose: Retrieves the image frame at a specific index.
	// Inputs: index int - The zero-based frame index.
	// Outputs: *ebiten.Image - The corresponding image frame, or nil if the index is out of bounds.
	Image(index int) *ebiten.Image

	// Length trả về tổng số frame ảnh trong sprite.
	// Purpose: Retrieves the total number of frames in the sprite.
	// Inputs: None.
	// Outputs: int - Total number of frames.
	Length() int

	// Width trả về chiều rộng của mỗi frame (pixel).
	// Purpose: Retrieves the width of the sprite frames.
	// Inputs: None.
	// Outputs: int - The width in pixels.
	Width() int

	// Height trả về chiều cao của mỗi frame (pixel).
	// Purpose: Retrieves the height of the sprite frames.
	// Inputs: None.
	// Outputs: int - The height in pixels.
	Height() int

	// AddImage thêm một frame ảnh vào cuối danh sách.
	// Purpose: Appends a new image frame to the end of the sprite sequence.
	// Inputs: image *ebiten.Image - The image frame to add.
	// Outputs: None.
	AddImage(image *ebiten.Image)

	// RemoveImage xóa frame ảnh tại vị trí index. Không làm gì nếu index ngoài biên.
	// Purpose: Removes an image frame at the specified index.
	// Inputs: index int - The index to remove.
	// Outputs: None.
	// Special requirements: Does nothing if the index is out of bounds.
	RemoveImage(index int)
}

// IAudioLW là lightweight wrapper đại diện cho một đoạn âm thanh sẵn sàng phát.
// Được dùng trong AudioData để điều khiển phát/dừng âm thanh thông qua AudioSystem.
type IAudioLW interface {
	// Play phát âm thanh với tên kênh, âm lượng (volume) và cao độ (pitch) chỉ định.
	// Purpose: Starts playing the audio with specific playback parameters.
	// Inputs:
	//   - name string: An optional channel name to categorize the audio.
	//   - volume float32: Playback volume (typically 0.0 to 1.0).
	//   - pitch float32: Playback speed/pitch multiplier (1.0 is normal).
	// Outputs: None.
	Play(name string, volume float32, pitch float32)

	// Pause tạm dừng âm thanh đang phát.
	// Purpose: Pauses the currently playing audio.
	// Inputs: None.
	// Outputs: None.
	Pause()

	// Stop dừng hoàn toàn và reset về đầu âm thanh.
	// Purpose: Stops playback completely and resets the audio position to the beginning.
	// Inputs: None.
	// Outputs: None.
	Stop()

	// IsPlaying trả về true nếu âm thanh đang được phát.
	// Purpose: Checks if the audio is currently playing.
	// Inputs: None.
	// Outputs: bool - True if playing.
	IsPlaying() bool

	// IsPaused trả về true nếu âm thanh đang tạm dừng.
	// Purpose: Checks if the audio is currently paused.
	// Inputs: None.
	// Outputs: bool - True if paused.
	IsPaused() bool

	// Resume tiếp tục phát âm thanh đang bị tạm dừng.
	// Purpose: Resumes playback of a paused audio track.
	// Inputs: None.
	// Outputs: None.
	Resume()

	// IsStopped trả về true nếu âm thanh đã dừng hoàn toàn (chưa phát hoặc đã gọi Stop).
	// Purpose: Checks if the audio is completely stopped.
	// Inputs: None.
	// Outputs: bool - True if stopped or never played.
	IsStopped() bool

	// SetLooping bật/tắt chế độ lặp lại vô tận cho âm thanh này.
	// Purpose: Toggles infinite looping behavior for the audio track.
	// Inputs: loop bool - True to loop indefinitely, false to play once.
	// Outputs: None.
	SetLooping(loop bool)

	// IsLooping kiểm tra xem âm thanh này có đang ở chế độ lặp lại không.
	// Purpose: Checks if the audio is set to loop indefinitely.
	// Inputs: None.
	// Outputs: bool - True if looping is enabled.
	IsLooping() bool

	// SetVolume thiết lập âm lượng cho trình phát âm thanh (0.0 đến 1.0).
	// Purpose: Sets the playback volume for the audio track.
	// Inputs: volume float32 - The volume level (typically between 0.0 and 1.0).
	// Outputs: None.
	SetVolume(volume float32)
}
