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

// NewAudioLW tạo một thực thể âm thanh mới từ dữ liệu gốc
func NewAudioLW(ctx *audio.Context, data []byte) *AudioLW {
	return &AudioLW{
		context:   ctx,
		buffer:    data,
		isStopped: true, // Mặc định là đang dừng
	}
}

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

func (this *AudioLW) Pause() {
	if this.player != nil && this.player.IsPlaying() {
		this.player.Pause()
		this.isPauseRequested = true
	}
}

func (this *AudioLW) Resume() {
	if this.player != nil && !this.player.IsPlaying() && this.isPauseRequested {
		this.player.Play()
		this.isPauseRequested = false
	}
}

func (this *AudioLW) Stop() {
	if this.player != nil {
		this.player.Pause()
		this.player.Rewind()
		this.isPauseRequested = false
		this.isStopped = true
	}
}

func (this *AudioLW) IsPlaying() bool {
	return this.player != nil && this.player.IsPlaying()
}

func (this *AudioLW) IsPaused() bool {
	return this.isPauseRequested
}

func (this *AudioLW) IsStopped() bool {
	return this.isStopped || this.player == nil
}

func (this *AudioLW) SetLooping(loop bool) {
	this.isLooping = loop
}

func (this *AudioLW) IsLooping() bool {
	return this.isLooping
}

func (this *AudioLW) SetVolume(volume float32) {
	if this.player != nil {
		this.player.SetVolume(float64(volume))
	}
}
