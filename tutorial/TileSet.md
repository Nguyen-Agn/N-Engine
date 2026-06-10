# Bản đồ Ô gạch / Tilemap

> **Vision**: Khai báo grid đơn giản, engine vẽ tilemap tự động.
> **Vision**: Declare a simple grid, the engine draws the tilemap automatically.

---

## 1. Giải thích / Explanation

Tilemap object nhúng mixin `ncom.Tile` (`TilemapComponent`) để hiển thị bản đồ dạng lưới ô. Việc vẽ hàng nghìn ô gạch được Engine culling (lọc) và tối ưu hóa tự động bằng `DrawSystem`.
Tilemap object embeds `ncom.Tile` to display a grid-based map.

**Cách hoạt động / How it works:**
1. Khai báo sprite sheet chứa các tile / Declare a sprite sheet with tiles. (Cắt lưới tự động bằng Manifest `cols`, `rows`)
2. Đặt kích thước từng ô tile / Set tile width & height.
3. Khai báo grid (mảng 2D số nguyên) / Declare grid (2D integer array).
4. Khai báo Vị trí bắt đầu của Tilemap.
5. Engine tự vẽ / Engine renders automatically.

Giá trị trong mảng grid:
- `-1` hoặc `<0`: Ô trống (Không vẽ gì cả).
- `0, 1, 2...`: Index tương ứng của frame cắt ra từ Sprite Sheet.

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

// TilemapObject - bản đồ ô gạch / tile map object
type TilemapObject struct {
	ncom.Object // Lifecycle
	ncom.Pos     // Vị trí bản đồ / Map position
	ncom.Tile    // Tilemap component
}

// NewTilemapObject - tạo tilemap / create tilemap
func NewTilemapObject() *TilemapObject {
	t := &TilemapObject{}

	// Dùng "pos" cho vị trí và "til" cho Tilemap component
	napi.Obj.NewObject(t, "Level1Tilemap", "pos til sce-main")
	return t
}

// Create - thiết lập tilemap / setup tilemap
func (t *TilemapObject) Create() {
	// Đặt sprite sheet chứa các tile / Set sprite sheet with tiles
	t.SetSprite(napi.Assert.GetSprite("tileset_grass"))

	// Kích thước mỗi ô tile / Size of each tile cell (Ví dụ: 32x32 pixel)
	t.SetTileWidth(32)
	t.SetTileHeight(32)

	// Vị trí gốc vẽ Tilemap (top-left) trên bản đồ vật lý
	t.SetX(0)
	t.SetY(0)

	// Khai báo grid (Mảng 2D int)
	t.SetGrid([][]int{
		{ 0,  1,  1,  1,  1,  1,  1,  2},
		{ 8, -1, -1, -1, -1, -1, -1, 10},
		{ 8, -1,  4, -1, -1,  5, -1, 10},
		{16, 17, 17, 17, 17, 17, 17, 18},
	})
}
```

> **Lưu ý:** Việc va chạm (Collision) không gắn liền với Tilemap mặc định. Nếu bạn muốn nhân vật va chạm với tường, bạn nên tạo riêng các tảng "WallObject" nhúng `ncom.Box` đặt chồng lên tại các tọa độ tương ứng của Tilemap.
