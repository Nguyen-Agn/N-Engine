package nsystem

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// AudioSystem chịu trách nhiệm duy nhất là điều phối việc phát âm thanh.
// Nó đọc cờ (flag) từ AudioData và thực hiện lệnh phát/dừng âm thanh.
// Điều này giúp tách biệt hoàn toàn logic game khỏi việc gọi thư viện audio.
type AudioSystem struct {
	// query là bộ truy vấn được tạo một lần và tái sử dụng, rất hiệu năng.
	query *donburi.Query
}

// NewAudioSystem khởi tạo AudioSystem với bộ lọc chỉ lấy thực thể có AudioData.
func NewAudioSystem() *AudioSystem {
	return &AudioSystem{
		query: donburi.NewQuery(filter.Contains(Audio)),
	}
}

// Update duyệt qua tất cả thực thể có AudioData và xử lý các cờ phát/dừng.
// Chỉ những thực thể có Component Audio mới được xử lý (nhờ Query).
func (this *AudioSystem) Update(w donburi.World) {
	this.query.Each(w, func(entry *donburi.Entry) {
		data := donburi.Get[AudioData](entry, Audio)

		// Bỏ qua nếu chưa có âm thanh nào được đặt tên
		if data.AudioName == "" {
			return
		}

		// Lấy IAudioLW từ map theo tên hiện tại
		audioLW, ok := data.Audio[data.AudioName]
		if !ok || audioLW == nil {
			return
		}

		// Cập nhật trạng thái lặp lại cho IAudioLW
		audioLW.SetLooping(data.IsLooping)

		// Cập nhật volume liên tục nếu đang phát (để hỗ trợ fade/tăng giảm volume real-time)
		if audioLW.IsPlaying() {
			audioLW.SetVolume(data.Volume)
		}

		// Xử lý lệnh DỪNG trước (ưu tiên cao nhất)
		if data.ShouldStop {
			audioLW.Stop()
			data.ShouldStop = false // Reset cờ sau khi xử lý
			return
		}

		// Xử lý lệnh TẠM DỪNG
		if data.ShouldPause {
			audioLW.Pause()
			data.ShouldPause = false
			return
		}

		// Xử lý lệnh TIẾP TỤC
		if data.ShouldResume {
			audioLW.Resume()
			data.ShouldResume = false
			return
		}

		// Xử lý lệnh PHÁT MỚI
		if data.ShouldPlay {
			// Chỉ phát nếu chưa đang phát (tránh phát đè)
			if !audioLW.IsPlaying() {
				audioLW.Play(data.AudioName, data.Volume, data.Pitch)
			}
			data.ShouldPlay = false // Reset cờ sau khi xử lý
		} else {
			// Xử lý tự động lặp lại (Loop)
			if data.IsLooping {
				// Nếu đang không phát, không bị pause, và KHÔNG bị dừng cưỡng bức
				if !audioLW.IsPlaying() && !audioLW.IsPaused() && !audioLW.IsStopped() {
					audioLW.Play(data.AudioName, data.Volume, data.Pitch)
				}
			}
		}
	})
}
