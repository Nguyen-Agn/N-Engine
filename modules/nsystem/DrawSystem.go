package nsystem

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// DrawSystem chịu trách nhiệm vẽ tất cả các thực thể: Background, Tilemap, và Sprite.
// Nó chỉ đọc dữ liệu từ Component, không thay đổi trạng thái game.
//
// Camera offset (camX, camY) được trừ khỏi tọa độ entity trước khi vẽ,
// chuyển từ map space sang screen space. Truyền 0, 0 cho GUI layer.
type DrawSystem struct {
	// query lọc thực thể có đủ Sprite + Position
	query *donburi.Query
	// bgQuery lọc thực thể Background
	bgQuery *donburi.Query
	// tilemapQuery lọc thực thể Tilemap
	tilemapQuery *donburi.Query
	// screen là canvas vẽ, được set mỗi frame
	screen *ebiten.Image
}

// NewDrawSystem khởi tạo DrawSystem với các bộ lọc Component.
func NewDrawSystem() *DrawSystem {
	return &DrawSystem{
		query:        donburi.NewQuery(filter.Contains(Sprite, Position)),
		bgQuery:      donburi.NewQuery(filter.Contains(Background)),
		tilemapQuery: donburi.NewQuery(filter.Contains(Tilemap)),
	}
}

// SetScreen thiết lập canvas đích cho frame hiện tại.
// Phải gọi hàm này trước Draw() mỗi frame từ ebiten.Game.Draw().
func (ds *DrawSystem) SetScreen(screen *ebiten.Image) {
	ds.screen = screen
}

// Draw duyệt qua các thực thể và vẽ theo thứ tự: Background → Tilemap → Sprite.
// camX, camY là tọa độ camera trong map space — được trừ khỏi tọa độ entity.
// Truyền 0, 0 cho GUI layer để vẽ ở screen space không bị offset.
func (ds *DrawSystem) Draw(w donburi.World, camX, camY float32) {
	if ds.screen == nil {
		return
	}

	screenW := float32(ds.screen.Bounds().Dx())
	screenH := float32(ds.screen.Bounds().Dy())

	// ─── 1. Background ────────────────────────────────────────────────────────
	// Background luôn fill screen — không áp dụng camera offset
	ds.bgQuery.Each(w, func(entry *donburi.Entry) {
		bgData := donburi.Get[BackgroundData](entry, Background)
		if !bgData.IsVisible {
			return
		}

		if bgData.Color.A > 0 {
			ds.screen.Fill(bgData.Color)
		}

		if bgData.Sprite == nil || bgData.Sprite.Length() == 0 {
			return
		}

		img := bgData.Sprite.Image(0)
		if img == nil {
			return
		}
		imgW := float64(img.Bounds().Dx())
		imgH := float64(img.Bounds().Dy())
		sw := float64(ds.screen.Bounds().Dx())
		sh := float64(ds.screen.Bounds().Dy())

		if bgData.RepeatX || bgData.RepeatY {
			bgData.OffsetX += bgData.ScrollSpeedX
			bgData.OffsetY += bgData.ScrollSpeedY

			startX, endX, stepX := 0.0, sw, imgW
			if bgData.RepeatX {
				ox := float64(int(bgData.OffsetX) % int(imgW))
				if ox > 0 {
					ox -= imgW
				}
				startX = ox
			} else {
				startX = float64(bgData.OffsetX)
				endX = startX + imgW
				stepX = sw + 1
			}

			startY, endY, stepY := 0.0, sh, imgH
			if bgData.RepeatY {
				oy := float64(int(bgData.OffsetY) % int(imgH))
				if oy > 0 {
					oy -= imgH
				}
				startY = oy
			} else {
				startY = float64(bgData.OffsetY)
				endY = startY + imgH
				stepY = sh + 1
			}

			for x := startX; x < endX; x += stepX {
				for y := startY; y < endY; y += stepY {
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(x, y)
					ds.screen.DrawImage(img, opts)
				}
			}
		} else if bgData.Stretch {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(sw/imgW, sh/imgH)
			ds.screen.DrawImage(img, opts)
		} else {
			bgData.OffsetX += bgData.ScrollSpeedX
			bgData.OffsetY += bgData.ScrollSpeedY
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(bgData.OffsetX), float64(bgData.OffsetY))
			ds.screen.DrawImage(img, opts)
		}
	})

	// ─── 2. Tilemap ───────────────────────────────────────────────────────────
	// Tile-level culling: chỉ vẽ tile nằm trong viewport
	ds.tilemapQuery.Each(w, func(entry *donburi.Entry) {
		tilemapData := donburi.Get[TilemapData](entry, Tilemap)
		if !tilemapData.IsVisible || tilemapData.Sprite == nil || len(tilemapData.Grid) == 0 {
			return
		}

		var originX, originY float32
		if entry.HasComponent(Position) {
			pos := donburi.Get[PositionData](entry, Position)
			originX = pos.X
			originY = pos.Y
		}

		tw := tilemapData.TileWidth
		th := tilemapData.TileHeight

		// Tính range tile cần vẽ dựa trên viewport (culling)
		startCol := int((camX-originX)/float32(tw)) - 1
		startRow := int((camY-originY)/float32(th)) - 1
		endCol := int((camX-originX+screenW)/float32(tw)) + 1
		endRow := int((camY-originY+screenH)/float32(th)) + 1

		if startCol < 0 {
			startCol = 0
		}
		if startRow < 0 {
			startRow = 0
		}
		if endCol > tilemapData.Cols {
			endCol = tilemapData.Cols
		}
		if endRow > tilemapData.Rows {
			endRow = tilemapData.Rows
		}

		for r := startRow; r < endRow; r++ {
			for c := startCol; c < endCol; c++ {
				tileID := tilemapData.Grid[r*tilemapData.Cols+c]
				if tileID < 0 {
					continue
				}
				tileImg := tilemapData.Sprite.Image(tileID)
				if tileImg == nil {
					continue
				}
				screenX := originX + float32(c*tw) - camX
				screenY := originY + float32(r*th) - camY
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(screenX), float64(screenY))
				ds.screen.DrawImage(tileImg, opts)
			}
		}
	})

	// ─── 3. Sprite ────────────────────────────────────────────────────────────
	// AABB culling: bỏ qua entity nằm hoàn toàn ngoài viewport
	ds.query.Each(w, func(entry *donburi.Entry) {
		sprData := donburi.Get[SpriteData](entry, Sprite)
		if !sprData.IsVisible || sprData.CurrentSprite == "" {
			return
		}

		posData := donburi.Get[PositionData](entry, Position)

		// Chuyển tọa độ sang screen space
		screenX := posData.X - camX
		screenY := posData.Y - camY

		// Lấy kích thước sprite để culling chính xác
		currentSpr := sprData.Sprite[sprData.CurrentSprite]
		if currentSpr == nil {
			return
		}
		sprW := float32(currentSpr.Width()) * sprData.ScaleX
		sprH := float32(currentSpr.Height()) * sprData.ScaleY

		// AABB culling: bỏ qua nếu nằm hoàn toàn ngoài viewport
		if screenX+sprW < 0 || screenY+sprH < 0 || screenX > screenW || screenY > screenH {
			return
		}

		drawOpts := BuildDrawOptions(*posData, *sprData, camX, camY)
		if drawOpts == nil {
			return
		}
		ds.screen.DrawImage(drawOpts.Image, drawOpts.Opts)

		// Chuyển frame animation
		if currentSpr.Length() > 0 {
			sprData.SpriteIdx += int(sprData.ImageSpeed)
			sprData.SpriteIdx %= currentSpr.Length()
		}
	})
}
