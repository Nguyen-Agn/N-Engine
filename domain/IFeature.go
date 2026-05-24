package domain

import "image/color"

// ─── Feature Interfaces ────────────────────────────────────────────────────────
// Các interface này định nghĩa hành vi của các tính năng render đặc biệt:
// hình nền (Background) và bản đồ gạch (Tilemap).
// Chúng không phải Component thông thường mà là object chuyên biệt gắn với Scene.

// IBackground định nghĩa giao diện tương tác với hình nền hoặc màu nền của Scene.
type IBackground interface {
	// Color trả về màu nền hiện tại của Scene (dùng khi không có ảnh nền).
	Color() color.RGBA
	// SetColor thiết lập màu nền mới cho Scene.
	SetColor(c color.RGBA)

	// Sprite trả về ảnh nền hiện tại. Trả về nil nếu chỉ dùng màu nền.
	Sprite() ISpriteLW
	// SetSprite thiết lập ảnh nền mới.
	SetSprite(s ISpriteLW)

	// RepeatX trả về trạng thái lặp ảnh nền theo chiều ngang.
	RepeatX() bool
	// SetRepeatX bật/tắt lặp ảnh nền theo chiều ngang.
	SetRepeatX(repeatX bool)

	// RepeatY trả về trạng thái lặp ảnh nền theo chiều dọc.
	RepeatY() bool
	// SetRepeatY bật/tắt lặp ảnh nền theo chiều dọc.
	SetRepeatY(repeatY bool)

	// Stretch trả về trạng thái co giãn ảnh nền.
	Stretch() bool
	// SetStretch bật/tắt co giãn ảnh nền theo kích thước màn hình.
	SetStretch(stretch bool)

	// ScrollSpeedX trả về tốc độ tự cuộn ngang (pixel/frame). 0 = không cuộn.
	ScrollSpeedX() float32
	// SetScrollSpeedX thiết lập tốc độ tự cuộn ngang.
	SetScrollSpeedX(speed float32)

	// ScrollSpeedY trả về tốc độ tự cuộn dọc (pixel/frame). 0 = không cuộn.
	ScrollSpeedY() float32
	// SetScrollSpeedY thiết lập tốc độ tự cuộn dọc.
	SetScrollSpeedY(speed float32)

	// OffsetX trả về offset cuộn ngang hiện tại (pixel).
	OffsetX() float32
	// SetOffsetX thiết lập offset cuộn ngang.
	SetOffsetX(offset float32)

	// OffsetY trả về offset cuộn dọc hiện tại (pixel).
	OffsetY() float32
	// SetOffsetY thiết lập offset cuộn dọc.
	SetOffsetY(offset float32)

	// IsVisible trả về trạng thái hiển thị của Background.
	IsVisible() bool
	// SetIsVisible bật/tắt hiển thị Background.
	SetIsVisible(visible bool)
}

// ITilemap định nghĩa giao diện tương tác với lưới bản đồ Tilemap.
// Tilemap là một lưới 2 chiều (Cols x Rows), mỗi ô lưu ID của tile trong tileset.
type ITilemap interface {
	// Sprite trả về Tileset sprite sheet đang dùng (chứa tất cả frame tile).
	Sprite() ISpriteLW
	// SetSprite thiết lập Tileset sprite sheet mới.
	SetSprite(s ISpriteLW)

	// TileWidth trả về chiều rộng của một ô tile (pixel).
	TileWidth() int
	// SetTileWidth thiết lập chiều rộng của một ô tile.
	SetTileWidth(w int)

	// TileHeight trả về chiều cao của một ô tile (pixel).
	TileHeight() int
	// SetTileHeight thiết lập chiều cao của một ô tile.
	SetTileHeight(h int)

	// Cols trả về số cột của lưới Tilemap.
	Cols() int
	// Rows trả về số hàng của lưới Tilemap.
	Rows() int

	// Resize thay đổi kích thước lưới bản đồ và reset toàn bộ Grid về giá trị -1 (trống).
	Resize(cols, rows int)

	// GetTile trả về ID tile tại tọa độ (col, row). Trả về -1 nếu tọa độ ngoài biên.
	GetTile(col, row int) int
	// SetTile thiết lập ID tile tại tọa độ (col, row). Không làm gì nếu tọa độ ngoài biên.
	SetTile(col, row int, tileID int)

	// SetGrid thiết lập toàn bộ lưới bản đồ từ mảng 2 chiều và tự động Resize tương ứng.
	// grid[row][col] = tileID. -1 = ô trống (không vẽ).
	SetGrid(grid [][]int)

	// IsVisible trả về trạng thái hiển thị của Tilemap.
	IsVisible() bool
	// SetIsVisible bật/tắt hiển thị Tilemap.
	SetIsVisible(visible bool)
}
