package naudio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// AudioLW thực thi IAudioLW.
// Nó là một "Flyweight" giữ player riêng nhưng dùng chung buffer dữ liệu.
type AudioLW struct {
	context *audio.Context
	buffer  []byte        // Dữ liệu âm thanh đã giải mã (Shared)
	player  *audio.Player // Trình phát riêng cho từng thực thể (Stateful)

	// --- Cờ Trạng Thái (State Flags) ---

	// IsLooping: true nếu âm thanh này cần lặp lại vô tận.
	// Được set khi IsLooping là true trong AudioData.
	isLooping bool

	// IsPauseRequested: true nếu âm thanh đang bị tạm dừng (hoặc yêu cầu tạm dừng).
	// Được set khi ShouldPause là true trong AudioData.
	isPauseRequested bool

	// isStopped: true nếu âm thanh đã được yêu cầu dừng hẳn.
	// Hữu ích để phân biệt với trường hợp trình phát kết thúc tự nhiên.
	isStopped bool
}

// Purpose: Creates a new lightweight audio instance from raw byte data.
// Inputs: ctx (*audio.Context) - The Ebitengine audio context, data ([]byte) - The decoded audio data.
// Outputs: (*AudioLW) - The newly created AudioLW object.
func NewAudioLW(ctx *audio.Context, data []byte) *AudioLW {
	return &AudioLW{
		context:   ctx,
		buffer:    data,
		isStopped: true, // Mặc định là đang dừng
	}
}

// Purpose: Starts playing the audio if it is not already playing. Creates the player lazily.
// Inputs: name (string) - Name of the audio (unused here), volume (float32) - Playback volume, pitch (float32) - Playback pitch (not implemented natively here).
func (this *AudioLW) Play(name string, volume float32, pitch float32) {
	// Chỉ tạo player một lần, tránh leak memory do tạo mới liên tục
	if this.player == nil {
		this.player = this.context.NewPlayerFromBytes(this.buffer)
	}

	this.isStopped = false
	this.isPauseRequested = false

	if !this.player.IsPlaying() {
		this.player.SetVolume(float64(volume))
		this.player.Rewind()
		this.player.Play()
	}
}

// Purpose: Pauses the currently playing audio.
func (this *AudioLW) Pause() {
	if this.player != nil && this.player.IsPlaying() {
		this.player.Pause()
		this.isPauseRequested = true
	}
}

// Purpose: Resumes playing the audio if it was paused.
func (this *AudioLW) Resume() {
	if this.player != nil && !this.player.IsPlaying() && this.isPauseRequested {
		this.player.Play()
		this.isPauseRequested = false
	}
}

// Purpose: Stops the playback completely and rewinds to the beginning.
func (this *AudioLW) Stop() {
	if this.player != nil {
		this.player.Pause()
		this.player.Rewind()
		this.isPauseRequested = false
		this.isStopped = true
	}
}

// Purpose: Checks if the audio is currently playing.
// Outputs: (bool) - True if playing, false otherwise.
func (this *AudioLW) IsPlaying() bool {
	return this.player != nil && this.player.IsPlaying()
}

// Purpose: Checks if the audio playback has been paused.
// Outputs: (bool) - True if a pause was requested, false otherwise.
func (this *AudioLW) IsPaused() bool {
	return this.isPauseRequested
}

// Purpose: Checks if the audio is fully stopped or hasn't started yet.
// Outputs: (bool) - True if stopped, false otherwise.
func (this *AudioLW) IsStopped() bool {
	return this.isStopped || this.player == nil
}

// Purpose: Sets whether the audio should loop continuously.
// Inputs: loop (bool) - True to enable looping, false to disable.
func (this *AudioLW) SetLooping(loop bool) {
	this.isLooping = loop
}

// Purpose: Checks if the audio is set to loop continuously.
// Outputs: (bool) - True if looping is enabled, false otherwise.
func (this *AudioLW) IsLooping() bool {
	return this.isLooping
}

// Purpose: Adjusts the playback volume of the audio.
// Inputs: volume (float32) - The new volume level.
func (this *AudioLW) SetVolume(volume float32) {
	if this.player != nil {
		this.player.SetVolume(float64(volume))
	}
}
