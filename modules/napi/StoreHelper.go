package napi

import (
	"log"
)

type assertGroup struct{}

// Store là nhóm hàm truy xuất tài nguyên, âm thanh và biến toàn cục.
// Store is the function group for accessing assets, audio, and global variables.
var Assert = &assertGroup{}

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

// GetSprite retrieves a loaded sprite from the global store by its key.
//
// Outputs:
// - ISpriteLW: The sprite interface, or nil if not found.
func (s *assertGroup) GetSprite(key string) ISpriteLW {
	return engine().Store.GetSprite(key)
}

// GetAudio retrieves loaded audio from the global store by its key.
//
// Outputs:
// - IAudioLW: The audio interface, or nil if not found.
func (s *assertGroup) GetAudio(key string) IAudioLW {
	return engine().Store.GetAudio(key)
}

// GetConst retrieves a constant value from the global store by its key.
//
// Outputs:
// - any: The constant value, or nil if it doesn't exist.
func (s *assertGroup) GetConst(key string) any {
	return engine().Store.GetConst(key)
}

// GetObject retrieves a named IObject from the global store.
//
// Purpose: Allows lookup of objects created via NewObject with a specific name.
//
// Outputs:
// - IObject: The game object, or nil if not found.
func (s *assertGroup) GetObject(key string) IObject {
	return engine().Store.GetObject(key)
}

// Play starts playback of an audio asset from the global store at default volume and pitch (1.0).
//
// Inputs:
// - key (string): The identifier of the audio asset.
func (s *assertGroup) Play(key string) {
	s.PlayAt(key, 1.0, 1.0)
}

// PlayAt starts playback of an audio asset with custom volume and pitch.
//
// Inputs:
// - key (string): The identifier of the audio asset.
// - volume (float32): The playback volume multiplier.
// - pitch (float32): The playback pitch multiplier.
func (s *assertGroup) PlayAt(key string, volume, pitch float32) {
	audio := engine().Store.GetAudio(key)
	if audio == nil {
		log.Printf("[napi] Play: không tìm thấy audio '%s'", key)
		return
	}
	audio.Play(key, volume, pitch)
}

// Stop halts playback of an audio asset.
//
// Inputs:
// - key (string): The identifier of the audio asset to stop.
func (s *assertGroup) Stop(key string) {
	audio := engine().Store.GetAudio(key)
	if audio == nil {
		return
	}
	audio.Stop()
}
