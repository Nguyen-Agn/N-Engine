package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// Khai báo token cho Tilemap Component
var Tilemap = enginetype.Tilemap

func init() {
	enginetype.RegisterComponentInitializer("til", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Tilemap, TilemapData{
			Grid:      make([]int, 0),
			IsVisible: true,
		})
	})
}

// TilemapComponent là Mixin để nhúng vào Custom Object.
// Yêu cầu Object phải có base IObject để gọi Entry().
type TilemapComponent struct {
	IObject
	data *TilemapData
}

func (t *TilemapComponent) BindComponent(base IObject) {
	t.IObject = base
	t.data = enginetype.GetComponent(base, Tilemap)
}

func (t TilemapComponent) Sprite() ISpriteLW {
	if t.data == nil {
		return nil
	}
	return t.data.Sprite
}

func (t *TilemapComponent) SetSprite(s ISpriteLW) {
	if t.data != nil {
		t.data.Sprite = s
	}
}

func (t TilemapComponent) TileWidth() int {
	if t.data == nil {
		return 0
	}
	return t.data.TileWidth
}

func (t *TilemapComponent) SetTileWidth(w int) {
	if t.data != nil {
		t.data.TileWidth = w
	}
}

func (t TilemapComponent) TileHeight() int {
	if t.data == nil {
		return 0
	}
	return t.data.TileHeight
}

func (t *TilemapComponent) SetTileHeight(h int) {
	if t.data != nil {
		t.data.TileHeight = h
	}
}

func (t TilemapComponent) Cols() int {
	if t.data == nil {
		return 0
	}
	return t.data.Cols
}

func (t TilemapComponent) Rows() int {
	if t.data == nil {
		return 0
	}
	return t.data.Rows
}

func (t *TilemapComponent) Resize(cols, rows int) {
	if t.data == nil {
		return
	}
	t.data.Cols = cols
	t.data.Rows = rows
	t.data.Grid = make([]int, cols*rows)
	// Khởi tạo các ô rỗng mặc định bằng -1
	for i := range t.data.Grid {
		t.data.Grid[i] = -1
	}
}

func (t TilemapComponent) GetTile(col, row int) int {
	if t.data == nil {
		return -1
	}
	if col < 0 || col >= t.data.Cols || row < 0 || row >= t.data.Rows {
		return -1
	}
	return t.data.Grid[row*t.data.Cols+col]
}

func (t *TilemapComponent) SetTile(col, row int, tileID int) {
	if t.data == nil {
		return
	}
	if col < 0 || col >= t.data.Cols || row < 0 || row >= t.data.Rows {
		return
	}
	t.data.Grid[row*t.data.Cols+col] = tileID
}

func (t *TilemapComponent) SetGrid(grid [][]int) {
	if t.data == nil {
		return
	}
	rows := len(grid)
	if rows == 0 {
		t.Resize(0, 0)
		return
	}
	cols := len(grid[0])
	t.Resize(cols, rows)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			t.SetTile(c, r, grid[r][c])
		}
	}
}

func (t TilemapComponent) IsVisible() bool {
	if t.data == nil {
		return false
	}
	return t.data.IsVisible
}

func (t *TilemapComponent) SetIsVisible(visible bool) {
	if t.data != nil {
		t.data.IsVisible = visible
	}
}
