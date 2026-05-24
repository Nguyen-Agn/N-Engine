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

		// Xử lý lệnh DỪNG trước (ưu tiên cao hơn)
		if data.ShouldStop {
			audioLW.Stop()
			data.ShouldStop = false // Reset cờ sau khi xử lý
			return
		}

		// Xử lý lệnh PHÁT
		if data.ShouldPlay {
			// Chỉ phát nếu chưa đang phát (tránh phát đè)
			if !audioLW.IsPlaying() {
				audioLW.Play(data.AudioName, data.Volume, data.Pitch)
			}
			data.ShouldPlay = false // Reset cờ sau khi xử lý
		}
	})
}
