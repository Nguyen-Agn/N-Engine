# Âm thanh / Audio (ncom.Aud)

> **Vision**: Quản lý âm thanh đơn giản, phát và dừng hiệu ứng dễ dàng.
> **Vision**: Simple audio management, easy playback of effects.

---

## 1. Giải thích / Explanation

Audio Component (`ncom.Aud`) cung cấp khả năng phát âm thanh (nhạc nền hoặc hiệu ứng) trực tiếp từ một Object. Hệ thống sẽ tự động xử lý vòng đời phát âm thanh thông qua `AudioSystem`.
The Audio Component (`ncom.Aud`) provides the ability to play audio (BGM or SFX) directly from an Object. The system automatically handles the audio lifecycle via the `AudioSystem`.

**Tính năng / Features:**
- Phát, Dừng, Tạm dừng, Tiếp tục / Play, Stop, Pause, Resume
- Cấu hình âm lượng (Volume) và cao độ (Pitch) / Configure Volume and Pitch
- Chế độ lặp lại (Looping) cho nhạc nền / Looping mode for BGM

---

## 2. Cấu hình / Configuration

Trong `manifest.toml`, khai báo các file âm thanh:
In `manifest.toml`, declare your audio files:

```toml
[[audios]]
name = "bgm_main"
path = "assets/audio/theme.ogg"
type = "ogg"

[[audios]]
name = "sfx_jump"
path = "assets/audio/jump.wav"
type = "wav"
```

## 3. Ví dụ / Code Example

```go
package objects

import (
	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
)

type MusicPlayer struct {
	ncom.Object
	ncom.Aud // Thêm Audio component / Add Audio component
}

func NewMusicPlayer() *MusicPlayer {
	m := &MusicPlayer{}
	// "aud" là token để Engine cấp phát ECS Audio Data
	napi.Obj.NewObject(m, "GlobalMusicPlayer", "aud sce-main")
	return m
}

func (m *MusicPlayer) Create() {
	// Lấy audio từ hệ thống (đã load từ manifest)
	// Get audio from system (loaded from manifest)
	bgm := napi.Assert.GetAudio("bgm_main")
	jumpSfx := napi.Assert.GetAudio("sfx_jump")

	// Đăng ký âm thanh vào Object với tên tự chọn
	// Register audio tracks to the Object with a custom name
	m.SetAudio("bgm", bgm)
	m.SetAudio("jump", jumpSfx)

	// Bật chế độ lặp cho nhạc nền
	// Enable looping for background music
	m.SetLooping("bgm", true)

	// Cấu hình âm lượng mặc định (0.0 đến 1.0)
	// Set default volume (0.0 to 1.0)
	m.SetVolume(0.8)

	// Bắt đầu phát nhạc nền
	// Start playing background music
	m.PlayDefault("bgm")
}

func (m *MusicPlayer) Jump() {
	// Phát hiệu ứng nhảy với volume 1.0 và pitch 1.0
	// Play jump effect with 1.0 volume and 1.0 pitch
	m.Play("jump", 1.0, 1.0)
}
```

## 4. API Bổ sung / Additional APIs
- `m.StopAudio("bgm")`: Dừng phát. / Stop playing.
- `m.PauseAudio("bgm")`: Tạm dừng. / Pause playback.
- `m.ResumeAudio("bgm")`: Tiếp tục phát. / Resume playback.
- `m.SetPitch(1.5)`: Chỉnh cao độ để thay đổi độ thanh/trầm (1.0 là mặc định). / Set pitch multiplier (1.0 is default).
