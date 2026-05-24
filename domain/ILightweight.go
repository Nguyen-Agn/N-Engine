package domain

import "github.com/hajimehoshi/ebiten/v2"

// ─── Lightweight Wrapper Interfaces ───────────────────────────────────────────
// Các interface này bọc (wrap) tài nguyên nặng (ảnh, âm thanh) thành dạng
// nhẹ hơn để ECS component có thể tham chiếu mà không bị phụ thuộc thư viện.

// ISpriteLW là lightweight wrapper chứa tập hợp các frame ảnh (ebiten.Image).
// Được dùng trong SpriteData và TilemapData để lưu trữ tài nguyên đồ họa.
type ISpriteLW interface {
	// Image trả về frame ảnh tại vị trí index. Trả về nil nếu index ngoài biên.
	Image(index int) *ebiten.Image

	// Length trả về tổng số frame ảnh trong sprite.
	Length() int

	// Width trả về chiều rộng của mỗi frame (pixel).
	Width() int

	// Height trả về chiều cao của mỗi frame (pixel).
	Height() int

	// AddImage thêm một frame ảnh vào cuối danh sách.
	AddImage(image *ebiten.Image)

	// RemoveImage xóa frame ảnh tại vị trí index. Không làm gì nếu index ngoài biên.
	RemoveImage(index int)
}

// IAudioLW là lightweight wrapper đại diện cho một đoạn âm thanh sẵn sàng phát.
// Được dùng trong AudioData để điều khiển phát/dừng âm thanh thông qua AudioSystem.
type IAudioLW interface {
	// Play phát âm thanh với tên kênh, âm lượng (volume) và cao độ (pitch) chỉ định.
	Play(name string, volume float32, pitch float32)

	// Pause tạm dừng âm thanh đang phát.
	Pause()

	// Stop dừng hoàn toàn và reset về đầu âm thanh.
	Stop()

	// IsPlaying trả về true nếu âm thanh đang được phát.
	IsPlaying() bool

	// IsPaused trả về true nếu âm thanh đang tạm dừng.
	IsPaused() bool
}
