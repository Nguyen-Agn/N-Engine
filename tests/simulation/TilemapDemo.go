package main

import (
	"image/color"
	"log"

	"autoworld/modules/napi"
)

// BackgroundObject là Custom Object đại diện cho nền game
type BackgroundObject struct {
	napi.IObject
	napi.IBackground
}

func NewBackgroundObject() *BackgroundObject {
	bg := &BackgroundObject{}
	napi.NewObject(bg, "bg_obj", "bg")

	// Đặt màu nền sẫm (Xanh lục sẫm)
	bg.SetColor(color.RGBA{10, 30, 20, 255})

	// Gán ảnh nền, cấu hình lặp và tự động cuộn
	bg.SetSprite(napi.GetSprite("bg-image"))
	bg.SetRepeatX(true)
	bg.SetRepeatY(true)
	bg.SetScrollSpeedX(0.3)
	bg.SetScrollSpeedY(0.1)

	napi.Register(bg, "main")
	return bg
}

// TilemapObject là Custom Object quản lý địa hình bản đồ
type TilemapObject struct {
	napi.IObject
	napi.ITilemap
	napi.IPosition
}

func NewTilemapObject() *TilemapObject {
	tm := &TilemapObject{}
	napi.NewObject(tm, "tilemap_obj", "pos til")

	tm.SetSprite(napi.GetSprite("tileset"))
	tm.SetTileWidth(32)
	tm.SetTileHeight(32)

	// Tạo ma trận gạch 2D
	grid := [][]int{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, 0, 1, 2, 3, -1, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
	}
	tm.SetGrid(grid)

	// Đặt vị trí bắt đầu vẽ bản đồ
	tm.SetX(0)
	tm.SetY(100)

	napi.Register(tm, "main")
	return tm
}

// BouncingPlayer là đối tượng di chuyển qua lại để kiểm tra lớp phủ
type BouncingPlayer struct {
	napi.IObject
	napi.IPosition
	napi.ISprite
	napi.IInput
	napi.IVelocity
}

func NewBouncingPlayer(x, y float32) *BouncingPlayer {
	p := &BouncingPlayer{}
	napi.NewObject(p, "player_obj", "pos spr inp vel")
	p.SetX(x)
	p.SetY(y)

	p.ListenOn("space", p.OnSpace)

	napi.Register(p, "main")
	return p
}

func (p *BouncingPlayer) Create() {
	p.SetSprite("normal", napi.GetSprite("player-sprite"))
	p.SetCurrentSprite("normal")

	p.SetVelocityX(1)
	p.SetVelocityY(0)

}

func (this *BouncingPlayer) OnSpace() {
	this.SetVelocityX(-this.VelocityX())
}

func (p *BouncingPlayer) StepUpdate() {
	x := p.X()
	w := float32(napi.VarInt("game-width"))
	if x < 0 || x > w {
		p.SetVelocityX(-p.VelocityX())
	}

}

func main() {
	// 1. Khởi tạo cấu hình engine
	cfg := napi.GameConfig{
		Title:  "AutoWorld - Tilemap & Background Demo",
		Width:  640,
		Height: 480,
	}
	napi.Init(cfg)

	// 2. Load manifest tài nguyên
	napi.LoadFromFile("./tilemap_manifest.toml")

	// 3. Khởi tạo Scene
	_, err := napi.NewSceneAndGo("main", "map-640-480")
	if err != nil {
		log.Fatalf("Không thể khởi tạo Scene: %v", err)
	}

	// 4. Tạo các thực thể game
	NewBackgroundObject()
	NewTilemapObject()
	NewBouncingPlayer(50, 150)

	// 5. Chạy game loop
	napi.GameStart()
}
