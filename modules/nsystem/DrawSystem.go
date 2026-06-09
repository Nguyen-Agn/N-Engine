package nsystem

import (
	"fmt"
	"sort"

	"autoworld/domain"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/font/basicfont"
)

// DrawSystem is responsible for rendering all entities: Background, Tilemap, Sprite,
// and custom IDraw objects. It only reads Component data — never mutates game state.
//
// Camera offset (camX, camY) is subtracted from entity positions to convert
// from map space to screen space. Pass 0, 0 for GUI layers (screen space).
type DrawSystem struct {
	// query filters entities with Sprite + Position components
	query *donburi.Query
	// bgQuery filters Background entities
	bgQuery *donburi.Query
	// tilemapQuery filters Tilemap entities
	tilemapQuery *donburi.Query
	// drawQuery filters entities with the Draw component (to inject screen/cam before Draw())
	drawQuery *donburi.Query
	// debugQuery filters entities with the Debug component
	debugQuery *donburi.Query
	// drawObjects holds objects that implement IDraw; populated by AddDrawObject
	drawObjects []domain.IObject
	// sortedSprites contains valid Sprite entries sorted by ZOrder
	sortedSprites []*donburi.Entry
	// screen is the render target; set every frame via SetScreen
	screen *ebiten.Image
}

// NewDrawSystem initialises DrawSystem with all required component filters.
// Outputs: Returns a pointer to a newly initialized DrawSystem.
func NewDrawSystem() *DrawSystem {
	return &DrawSystem{
		query:         donburi.NewQuery(filter.Contains(Sprite, Position)),
		bgQuery:       donburi.NewQuery(filter.Contains(Background)),
		tilemapQuery:  donburi.NewQuery(filter.Contains(Tilemap)),
		drawQuery:     donburi.NewQuery(filter.Contains(Draw)),
		debugQuery:    donburi.NewQuery(filter.Contains(Debug)),
		sortedSprites: make([]*donburi.Entry, 0),
	}
}

// SetScreen sets the render target for the current frame.
// Inputs: screen (*ebiten.Image) - The Ebitengine image to draw to.
// Purpose: Must be called before Draw() every frame (usually done by the Camera) to supply the target rendering surface.
func (ds *DrawSystem) SetScreen(screen *ebiten.Image) {
	ds.screen = screen
}

// AddDrawObject registers an object whose Draw() method will be called each frame.
// Inputs: obj (domain.IObject) - The object implementing IDraw to be registered.
// Purpose: Called automatically by Map.AddObject when it detects an IDraw implementation to ensure it gets rendered.
func (ds *DrawSystem) AddDrawObject(obj domain.IObject) {
	ds.drawObjects = append(ds.drawObjects, obj)
}

// RemoveDrawObject unregisters an object from the draw loop.
// Inputs: obj (domain.IObject) - The object to be removed from the draw list.
// Purpose: Safely removes the object even if it was never registered, preventing it from being drawn anymore.
func (ds *DrawSystem) RemoveDrawObject(obj domain.IObject) {
	list := ds.drawObjects[:0]
	for _, o := range ds.drawObjects {
		if o != obj {
			list = append(list, o)
		}
	}
	ds.drawObjects = list
}

// Draw renders all entities in order: Background → Tilemap → Sprite → IDraw.
// Inputs: 
//   w (donburi.World) - The ECS world containing entities to draw.
//   camX, camY (float32) - Camera coordinates in map space used as screen offset. Pass 0, 0 for GUI layers.
// Purpose: Manages the complete rendering pipeline. Performs viewport culling and properly handles z-ordering. Also handles debug rendering if the debug component is present.
func (ds *DrawSystem) Draw(w donburi.World, camX, camY float32) {
	if ds.screen == nil {
		return
	}

	screenW := float32(ds.screen.Bounds().Dx())
	screenH := float32(ds.screen.Bounds().Dy())

	// ─── 1. Background ────────────────────────────────────────────────────────
	// Background always fills the screen — no camera offset applied.
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
	// Tile-level culling: only render tiles inside the viewport.
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
	// Sync and Sort Sprites based on ZOrder
	currentMap := make(map[donburi.Entity]bool)
	hasDirty := false

	ds.query.Each(w, func(entry *donburi.Entry) {
		currentMap[entry.Entity()] = true
		sprData := donburi.Get[domain.SpriteData](entry, Sprite)
		if sprData.IsZOrderDirty {
			hasDirty = true
			sprData.IsZOrderDirty = false
		}
	})

	// Remove dead entries
	alive := ds.sortedSprites[:0]
	for _, entry := range ds.sortedSprites {
		if entry.Valid() && currentMap[entry.Entity()] {
			alive = append(alive, entry)
		} else {
			hasDirty = true
		}
	}
	ds.sortedSprites = alive

	// Add new entries
	if len(ds.sortedSprites) != len(currentMap) {
		existingMap := make(map[donburi.Entity]bool)
		for _, entry := range ds.sortedSprites {
			existingMap[entry.Entity()] = true
		}
		ds.query.Each(w, func(entry *donburi.Entry) {
			if !existingMap[entry.Entity()] {
				ds.sortedSprites = append(ds.sortedSprites, entry)
				hasDirty = true
			}
		})
	}

	// Sort if needed
	if hasDirty {
		sort.SliceStable(ds.sortedSprites, func(i, j int) bool {
			zI := donburi.Get[domain.SpriteData](ds.sortedSprites[i], Sprite).ZOrder
			zJ := donburi.Get[domain.SpriteData](ds.sortedSprites[j], Sprite).ZOrder
			return zI < zJ
		})
	}

	// Render sorted sprites with AABB culling
	for _, entry := range ds.sortedSprites {
		sprData := donburi.Get[domain.SpriteData](entry, Sprite)
		if !sprData.IsVisible || sprData.CurrentSprite == "" {
			continue
		}

		posData := donburi.Get[domain.PositionData](entry, Position)

		screenX := posData.X - camX
		screenY := posData.Y - camY

		currentSpr := sprData.Sprite[sprData.CurrentSprite]
		if currentSpr == nil {
			continue
		}

		sprW := float32(currentSpr.Width()) * sprData.ScaleX
		sprH := float32(currentSpr.Height()) * sprData.ScaleY

		// AABB cull
		if screenX+sprW < 0 || screenY+sprH < 0 || screenX > screenW || screenY > screenH {
			continue
		}

		drawOptsList := BuildDrawOptions(*posData, *sprData, camX, camY)
		for _, drawOpts := range drawOptsList {
			if drawOpts == nil {
				continue
			}
			ds.screen.DrawImage(drawOpts.Image, drawOpts.Opts)
		}

		// Advance animation frame
		if currentSpr.Length() > 0 {
			sprData.SpriteIdx += int(sprData.ImageSpeed)
			sprData.SpriteIdx %= currentSpr.Length()
		}
	}

	// ─── 4. IDraw objects ─────────────────────────────────────────────────────
	// Inject current screen and camera offset into DrawData for each Draw entity,
	// then call Draw() on objects that implement IDraw.
	ds.drawQuery.Each(w, func(entry *donburi.Entry) {
		drawData := donburi.Get[DrawData](entry, Draw)
		drawData.Screen = ds.screen
		drawData.CamX = camX
		drawData.CamY = camY
	})
	
	// Sort drawObjects by ZOrder if they also implement ISprite
	sort.SliceStable(ds.drawObjects, func(i, j int) bool {
		zI, zJ := 0, 0
		if sprI, ok := ds.drawObjects[i].(domain.ISprite); ok {
			zI = sprI.ZOrder()
		}
		if sprJ, ok := ds.drawObjects[j].(domain.ISprite); ok {
			zJ = sprJ.ZOrder()
		}
		return zI < zJ
	})

	for _, obj := range ds.drawObjects {
		if drawer, ok := obj.(domain.IDraw); ok {
			drawer.Draw()
		}
	}

	// ─── 5. Debug ─────────────────────────────────────────────────────────────
	var logYOffset float32 = 0
	basicFace := text.NewGoXFace(basicfont.Face7x13)

	ds.debugQuery.Each(w, func(entry *donburi.Entry) {
		debugData := donburi.Get[domain.DebugData](entry, Debug)
		if !debugData.ShowBox && !debugData.ShowPos && !debugData.ShowInfo && debugData.CustomLog == "" {
			return
		}

		hasPos := entry.HasComponent(Position)
		var screenX, screenY float32
		if hasPos {
			pos := donburi.Get[domain.PositionData](entry, Position)
			screenX = pos.X - camX
			screenY = pos.Y - camY
		} else {
			screenX = 5
			screenY = logYOffset
			logYOffset += 16
		}

		// Vẽ chữ Info / CustomLog
		var textToDraw string
		if debugData.ShowInfo && entry.HasComponent(Infor) {
			info := donburi.Get[domain.InforData](entry, Infor)
			textToDraw = fmt.Sprintf("[%d] %s", info.Id, info.Name)
		}
		if debugData.CustomLog != "" {
			if textToDraw != "" {
				textToDraw += " | " + debugData.CustomLog
			} else {
				textToDraw = debugData.CustomLog
			}
		}

		if textToDraw != "" {
			op := &text.DrawOptions{}
			if hasPos {
				op.PrimaryAlign = text.AlignCenter
				
				if entry.HasComponent(Box) {
					box := donburi.Get[domain.BoxData](entry, Box)
					op.GeoM.Translate(float64(screenX+box.BoxX+box.BoxW/2), float64(screenY)-16)
				} else {
					op.GeoM.Translate(float64(screenX), float64(screenY)-16)
				}
			} else {
				op.PrimaryAlign = text.AlignStart
				op.GeoM.Translate(float64(screenX), float64(screenY))
			}
			op.ColorScale.ScaleWithColor(debugData.Color)
			text.Draw(ds.screen, textToDraw, basicFace, op)
		}

		// Chỉ vẽ pos/box khi entity CÓ Position thực sự trên map
		if !hasPos {
			return
		}

		if debugData.ShowPos {
			vector.StrokeLine(ds.screen, screenX-4, screenY, screenX+4, screenY, 1, debugData.Color, false)
			vector.StrokeLine(ds.screen, screenX, screenY-4, screenX, screenY+4, 1, debugData.Color, false)
		}

		if debugData.ShowBox && entry.HasComponent(Box) {
			box := donburi.Get[domain.BoxData](entry, Box)
			boxScreenX := screenX + box.BoxX
			boxScreenY := screenY + box.BoxY
			if box.Shape == domain.BSCircle {
				radius := box.BoxW / 2
				vector.StrokeCircle(ds.screen, boxScreenX+radius, boxScreenY+radius, radius, 1, debugData.Color, false)
			} else {
				vector.StrokeRect(ds.screen, boxScreenX, boxScreenY, box.BoxW, box.BoxH, 1, debugData.Color, false)
			}
		}
	})
}
