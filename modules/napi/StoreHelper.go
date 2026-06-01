package napi

import (
	"log"
)

type AssertGroup struct{}

// Store là nhóm hàm truy xuất tài nguyên, âm thanh và biến toàn cục.
// Store is the function group for accessing assets, audio, and global variables.
var Assert = &AssertGroup{}

// LoadManifest đọc file TOML manifest và load toàn bộ tài nguyên vào global store.
// Lỗi audio (ví dụ WAV bit depth không phù hợp) chỉ log cảnh báo, không dừng chương trình.
// func LoadManifest(path string) {
// 	e := engine()
// 	spriteLoader := nasset.NewSpriteLoader()
// 	audioLoader := nasset.NewAudioLoader(e.AudioCtx, 44100)
// 	manifestLoader := nasset.NewManifestLoader(spriteLoader, audioLoader)

// 	if err := manifestLoader.LoadFromFile(path, e.Store); err != nil {
// 		log.Printf("[napi] LoadManifest '%s': %v", path, err)
// 	}
// }

// GetSprite lấy sprite đã load từ global store bằng key. Trả về nil nếu không tìm thấy.
func (s *AssertGroup) GetSprite(key string) ISpriteLW {
	return engine().Store.GetSprite(key)
}

// GetAudio lấy audio đã load từ global store bằng key. Trả về nil nếu không tìm thấy.
func (s *AssertGroup) GetAudio(key string) IAudioLW {
	return engine().Store.GetAudio(key)
}

// GetConst lấy giá trị hằng số từ global store bằng key. Trả về nil nếu không có.
func (s *AssertGroup) GetConst(key string) any {
	return engine().Store.GetConst(key)
}

// GetObject lấy IObject đã đặt tên từ global store.
// Chỉ có object được tạo bằng NewObject(name, ...) mới được lưu tự động.
// Trả về nil nếu không tìm thấy.
func (s *AssertGroup) GetObject(key string) IObject {
	return engine().Store.GetObject(key)
}

// Play phát audio từ global store với volume và pitch mặc định (1.0).
func (s *AssertGroup) Play(key string) {
	s.PlayAt(key, 1.0, 1.0)
}

// PlayAt phát audio với volume và pitch tùy chỉnh.
func (s *AssertGroup) PlayAt(key string, volume, pitch float32) {
	audio := engine().Store.GetAudio(key)
	if audio == nil {
		log.Printf("[napi] Play: không tìm thấy audio '%s'", key)
		return
	}
	audio.Play(key, volume, pitch)
}

// Stop dừng audio từ global store.
func (s *AssertGroup) Stop(key string) {
	audio := engine().Store.GetAudio(key)
	if audio == nil {
		return
	}
	audio.Stop()
}
