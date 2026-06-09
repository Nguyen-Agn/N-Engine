package domain

import "image/color"

// ─── Feature Interfaces ────────────────────────────────────────────────────────
// Các interface này định nghĩa hành vi của các tính năng render đặc biệt:
// hình nền (Background) và bản đồ gạch (Tilemap).
// Chúng không phải Component thông thường mà là object chuyên biệt gắn với Scene.

// IBackground định nghĩa giao diện tương tác với hình nền hoặc màu nền của Scene.
type IBackground interface {
	// Color trả về màu nền hiện tại của Scene (dùng khi không có ảnh nền).
	// Purpose: Retrieves the current background color of the Scene.
	// Inputs: None.
	// Outputs: color.RGBA - The current background color.
	Color() color.RGBA
	
	// SetColor thiết lập màu nền mới cho Scene.
	// Purpose: Sets a new background color for the Scene.
	// Inputs: c color.RGBA - The new background color to apply.
	// Outputs: None.
	SetColor(c color.RGBA)

	// Sprite trả về ảnh nền hiện tại. Trả về nil nếu chỉ dùng màu nền.
	// Purpose: Retrieves the current background sprite image.
	// Inputs: None.
	// Outputs: ISpriteLW - The current background sprite, or nil if only color is used.
	Sprite() ISpriteLW
	
	// SetSprite thiết lập ảnh nền mới.
	// Purpose: Sets a new background sprite image.
	// Inputs: s ISpriteLW - The new sprite to use as the background.
	// Outputs: None.
	SetSprite(s ISpriteLW)

	// RepeatX trả về trạng thái lặp ảnh nền theo chiều ngang.
	// Purpose: Checks if the background image repeats horizontally.
	// Inputs: None.
	// Outputs: bool - True if horizontal repeat is enabled, false otherwise.
	RepeatX() bool
	
	// SetRepeatX bật/tắt lặp ảnh nền theo chiều ngang.
	// Purpose: Toggles horizontal repeating of the background image.
	// Inputs: repeatX bool - True to enable horizontal repeating, false to disable.
	// Outputs: None.
	SetRepeatX(repeatX bool)

	// RepeatY trả về trạng thái lặp ảnh nền theo chiều dọc.
	// Purpose: Checks if the background image repeats vertically.
	// Inputs: None.
	// Outputs: bool - True if vertical repeat is enabled, false otherwise.
	RepeatY() bool
	
	// SetRepeatY bật/tắt lặp ảnh nền theo chiều dọc.
	// Purpose: Toggles vertical repeating of the background image.
	// Inputs: repeatY bool - True to enable vertical repeating, false to disable.
	// Outputs: None.
	SetRepeatY(repeatY bool)

	// Stretch trả về trạng thái co giãn ảnh nền.
	// Purpose: Checks if the background image is stretched to fit the screen.
	// Inputs: None.
	// Outputs: bool - True if stretching is enabled, false otherwise.
	Stretch() bool
	
	// SetStretch bật/tắt co giãn ảnh nền theo kích thước màn hình.
	// Purpose: Toggles background image stretching to fit the screen dimensions.
	// Inputs: stretch bool - True to stretch the image, false otherwise.
	// Outputs: None.
	SetStretch(stretch bool)

	// ScrollSpeedX trả về tốc độ tự cuộn ngang (pixel/frame). 0 = không cuộn.
	// Purpose: Retrieves the horizontal auto-scroll speed.
	// Inputs: None.
	// Outputs: float32 - The scrolling speed in pixels per frame (0 means no scrolling).
	ScrollSpeedX() float32
	
	// SetScrollSpeedX thiết lập tốc độ tự cuộn ngang.
	// Purpose: Sets the horizontal auto-scroll speed.
	// Inputs: speed float32 - The scrolling speed in pixels per frame.
	// Outputs: None.
	SetScrollSpeedX(speed float32)

	// ScrollSpeedY trả về tốc độ tự cuộn dọc (pixel/frame). 0 = không cuộn.
	// Purpose: Retrieves the vertical auto-scroll speed.
	// Inputs: None.
	// Outputs: float32 - The scrolling speed in pixels per frame (0 means no scrolling).
	ScrollSpeedY() float32
	
	// SetScrollSpeedY thiết lập tốc độ tự cuộn dọc.
	// Purpose: Sets the vertical auto-scroll speed.
	// Inputs: speed float32 - The scrolling speed in pixels per frame.
	// Outputs: None.
	SetScrollSpeedY(speed float32)

	// OffsetX trả về offset cuộn ngang hiện tại (pixel).
	// Purpose: Retrieves the current horizontal scrolling offset.
	// Inputs: None.
	// Outputs: float32 - The horizontal offset in pixels.
	OffsetX() float32
	
	// SetOffsetX thiết lập offset cuộn ngang.
	// Purpose: Sets a specific horizontal scrolling offset.
	// Inputs: offset float32 - The horizontal offset to apply in pixels.
	// Outputs: None.
	SetOffsetX(offset float32)

	// OffsetY trả về offset cuộn dọc hiện tại (pixel).
	// Purpose: Retrieves the current vertical scrolling offset.
	// Inputs: None.
	// Outputs: float32 - The vertical offset in pixels.
	OffsetY() float32
	
	// SetOffsetY thiết lập offset cuộn dọc.
	// Purpose: Sets a specific vertical scrolling offset.
	// Inputs: offset float32 - The vertical offset to apply in pixels.
	// Outputs: None.
	SetOffsetY(offset float32)

	// IsVisible trả về trạng thái hiển thị của Background.
	// Purpose: Checks if the background is currently visible.
	// Inputs: None.
	// Outputs: bool - True if visible, false otherwise.
	IsVisible() bool
	
	// SetIsVisible bật/tắt hiển thị Background.
	// Purpose: Toggles the visibility of the background.
	// Inputs: visible bool - True to show the background, false to hide it.
	// Outputs: None.
	SetIsVisible(visible bool)
}

// ITilemap định nghĩa giao diện tương tác với lưới bản đồ Tilemap.
// Tilemap là một lưới 2 chiều (Cols x Rows), mỗi ô lưu ID của tile trong tileset.
type ITilemap interface {
	// Sprite trả về Tileset sprite sheet đang dùng (chứa tất cả frame tile).
	// Purpose: Retrieves the current tileset sprite sheet used for the tilemap.
	// Inputs: None.
	// Outputs: ISpriteLW - The tileset sprite sheet containing all tile frames.
	Sprite() ISpriteLW
	
	// SetSprite thiết lập Tileset sprite sheet mới.
	// Purpose: Sets a new tileset sprite sheet for the tilemap.
	// Inputs: s ISpriteLW - The new sprite sheet to use.
	// Outputs: None.
	SetSprite(s ISpriteLW)

	// TileWidth trả về chiều rộng của một ô tile (pixel).
	// Purpose: Retrieves the width of a single tile in the tilemap.
	// Inputs: None.
	// Outputs: int - The width of a tile in pixels.
	TileWidth() int
	
	// SetTileWidth thiết lập chiều rộng của một ô tile.
	// Purpose: Sets the width of a single tile.
	// Inputs: w int - The new tile width in pixels.
	// Outputs: None.
	SetTileWidth(w int)

	// TileHeight trả về chiều cao của một ô tile (pixel).
	// Purpose: Retrieves the height of a single tile in the tilemap.
	// Inputs: None.
	// Outputs: int - The height of a tile in pixels.
	TileHeight() int
	
	// SetTileHeight thiết lập chiều cao của một ô tile.
	// Purpose: Sets the height of a single tile.
	// Inputs: h int - The new tile height in pixels.
	// Outputs: None.
	SetTileHeight(h int)

	// Cols trả về số cột của lưới Tilemap.
	// Purpose: Retrieves the total number of columns in the tilemap grid.
	// Inputs: None.
	// Outputs: int - Number of columns.
	Cols() int
	
	// Rows trả về số hàng của lưới Tilemap.
	// Purpose: Retrieves the total number of rows in the tilemap grid.
	// Inputs: None.
	// Outputs: int - Number of rows.
	Rows() int

	// Resize thay đổi kích thước lưới bản đồ và reset toàn bộ Grid về giá trị -1 (trống).
	// Purpose: Resizes the tilemap grid dimensions and resets all cells to -1 (empty).
	// Inputs:
	//   - cols int: The new number of columns.
	//   - rows int: The new number of rows.
	// Outputs: None.
	Resize(cols, rows int)

	// GetTile trả về ID tile tại tọa độ (col, row). Trả về -1 nếu tọa độ ngoài biên.
	// Purpose: Gets the tile ID at the specified grid coordinates.
	// Inputs:
	//   - col int: The column index.
	//   - row int: The row index.
	// Outputs: int - The tile ID, or -1 if the coordinates are out of bounds.
	GetTile(col, row int) int
	
	// SetTile thiết lập ID tile tại tọa độ (col, row). Không làm gì nếu tọa độ ngoài biên.
	// Purpose: Sets the tile ID at the specified grid coordinates.
	// Inputs:
	//   - col int: The column index.
	//   - row int: The row index.
	//   - tileID int: The ID of the tile to place.
	// Outputs: None.
	// Special requirements: Does nothing if the coordinates are out of bounds.
	SetTile(col, row int, tileID int)

	// SetGrid thiết lập toàn bộ lưới bản đồ từ mảng 2 chiều và tự động Resize tương ứng.
	// grid[row][col] = tileID. -1 = ô trống (không vẽ).
	// Purpose: Applies a 2D array to set the entire tilemap grid, automatically resizing to match.
	// Inputs: grid [][]int - A 2D array where grid[row][col] is the tile ID (-1 for empty).
	// Outputs: None.
	SetGrid(grid [][]int)

	// IsVisible trả về trạng thái hiển thị của Tilemap.
	// Purpose: Checks if the tilemap is currently visible.
	// Inputs: None.
	// Outputs: bool - True if visible, false otherwise.
	IsVisible() bool
	
	// SetIsVisible bật/tắt hiển thị Tilemap.
	// Purpose: Toggles the visibility of the tilemap.
	// Inputs: visible bool - True to show the tilemap, false to hide it.
	// Outputs: None.
	SetIsVisible(visible bool)
}
