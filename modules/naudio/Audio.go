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
	// Được set khi IsLoopingData là true trong AudioData.
	isLooping bool

	// IsPauseRequested: true nếu âm thanh đang bị tạm dừng (hoặc yêu cầu tạm dừng).
	// Được set khi ShouldPauseData là true trong AudioData.
	isPauseRequested bool
}

// NewAudioLW tạo một thực thể âm thanh mới từ dữ liệu gốc
func NewAudioLW(ctx *audio.Context, data []byte) *AudioLW {
	return &AudioLW{
		context: ctx,
		buffer:  data,
	}
}

func (this *AudioLW) Play(name string, volume float32, pitch float32) {
	// Tạo player mới từ buffer dùng chung
	this.player = this.context.NewPlayerFromBytes(this.buffer)

	// Nếu đang phát thì không phát đè lên (hoặc bạn có thể Rewind() nếu muốn phát lại từ đầu)
	if !this.player.IsPlaying() {
		this.player.SetVolume(float64(volume))
		// Lưu ý: Ebitengine chuẩn không hỗ trợ Pitch trực tiếp trên Player,
		// cần Resampling nếu muốn đổi Pitch. Tạm thời ta tập trung vào Volume.
		this.player.Rewind()
		this.player.Play()
	}
}

func (this *AudioLW) Pause() {
	if this.player != nil {
		this.player.Pause()
	}
}

func (this *AudioLW) Stop() {
	if this.player != nil {
		this.player.Pause()
		this.player.Rewind()
	}
}

func (this *AudioLW) IsPlaying() bool {
	return this.player != nil && this.player.IsPlaying()
}

func (this *AudioLW) IsPaused() bool {
	return this.player != nil && !this.player.IsPlaying()
}
