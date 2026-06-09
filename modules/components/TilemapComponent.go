package components

import (
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

// Tilemap is the token for the Tilemap Component used in the ECS.
var Tilemap = enginetype.Tilemap

// init registers the Tilemap Component initializer with default values.
func init() {
	enginetype.RegisterComponentInitializer("til", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Tilemap, TilemapData{
			Grid:      make([]int, 0),
			IsVisible: true,
		})
	})
}

// TilemapComponent is a mixin to embed into Custom Objects.
// It requires the Object to have an IObject base to call Entry().
type TilemapComponent struct {
	IObject
	data *TilemapData
}

// BindComponent binds the base object and retrieves the TilemapData from the ECS.
// Inputs: base - the base IObject to bind.
func (t *TilemapComponent) BindComponent(base IObject) {
	t.IObject = base
	t.data = enginetype.GetComponent(base, Tilemap)
}

// Sprite returns the sprite associated with the tilemap.
// Outputs: ISpriteLW representing the tilemap's sprite.
func (t TilemapComponent) Sprite() ISpriteLW {
	if t.data == nil {
		return nil
	}
	return t.data.Sprite
}

// SetSprite sets the sprite for the tilemap.
// Inputs: s - the new sprite.
func (t *TilemapComponent) SetSprite(s ISpriteLW) {
	if t.data != nil {
		t.data.Sprite = s
	}
}

// TileWidth returns the width of a single tile.
// Outputs: int representing the tile width in pixels.
func (t TilemapComponent) TileWidth() int {
	if t.data == nil {
		return 0
	}
	return t.data.TileWidth
}

// SetTileWidth sets the width of a single tile.
// Inputs: w - the new tile width in pixels.
func (t *TilemapComponent) SetTileWidth(w int) {
	if t.data != nil {
		t.data.TileWidth = w
	}
}

// TileHeight returns the height of a single tile.
// Outputs: int representing the tile height in pixels.
func (t TilemapComponent) TileHeight() int {
	if t.data == nil {
		return 0
	}
	return t.data.TileHeight
}

// SetTileHeight sets the height of a single tile.
// Inputs: h - the new tile height in pixels.
func (t *TilemapComponent) SetTileHeight(h int) {
	if t.data != nil {
		t.data.TileHeight = h
	}
}

// Cols returns the number of columns in the tilemap grid.
// Outputs: int representing the number of columns.
func (t TilemapComponent) Cols() int {
	if t.data == nil {
		return 0
	}
	return t.data.Cols
}

// Rows returns the number of rows in the tilemap grid.
// Outputs: int representing the number of rows.
func (t TilemapComponent) Rows() int {
	if t.data == nil {
		return 0
	}
	return t.data.Rows
}

// Resize changes the dimensions of the tilemap grid and reinitializes it with empty tiles (-1).
// Inputs: cols - new number of columns, rows - new number of rows.
func (t *TilemapComponent) Resize(cols, rows int) {
	if t.data == nil {
		return
	}
	t.data.Cols = cols
	t.data.Rows = rows
	t.data.Grid = make([]int, cols*rows)
	// Initialize empty tiles with -1 by default
	for i := range t.data.Grid {
		t.data.Grid[i] = -1
	}
}

// GetTile retrieves the tile ID at the specified column and row.
// Inputs: col - the column index, row - the row index.
// Outputs: int representing the tile ID (-1 if empty or out of bounds).
func (t TilemapComponent) GetTile(col, row int) int {
	if t.data == nil {
		return -1
	}
	if col < 0 || col >= t.data.Cols || row < 0 || row >= t.data.Rows {
		return -1
	}
	return t.data.Grid[row*t.data.Cols+col]
}

// SetTile sets the tile ID at the specified column and row.
// Inputs: col - the column index, row - the row index, tileID - the new tile ID.
func (t *TilemapComponent) SetTile(col, row int, tileID int) {
	if t.data == nil {
		return
	}
	if col < 0 || col >= t.data.Cols || row < 0 || row >= t.data.Rows {
		return
	}
	t.data.Grid[row*t.data.Cols+col] = tileID
}

// SetGrid populates the entire tilemap grid from a 2D slice.
// Inputs: grid - 2D slice of tile IDs.
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

// IsVisible returns whether the tilemap is currently visible.
// Outputs: bool indicating visibility.
func (t TilemapComponent) IsVisible() bool {
	if t.data == nil {
		return false
	}
	return t.data.IsVisible
}

// SetIsVisible sets the visibility of the tilemap.
// Inputs: visible - true to show, false to hide.
func (t *TilemapComponent) SetIsVisible(visible bool) {
	if t.data != nil {
		t.data.IsVisible = visible
	}
}
