//go:build ignore

// Demo tính năng Save/Load
// Run: go run .\tests\save_load_demo\SaveLoadDemo.go
package main

import (
	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// --- Khai báo Hero ---
type Hero struct {
	napi.IObject
	ncom.Pos
	ncom.Info
	ncom.Spr
	ncom.Deb
}

func NewHero() *Hero {
	h := &Hero{}
	// Tokens: pos (vị trí), drw (vẽ), inf (thông tin + savetag), sce-cur (scene hiện tại)
	napi.Obj.NewObject(h, "hero", "pos sce-cur spr deb")
	h.SetX(300)
	h.SetY(200)
	h.SetSaveTag("hero_1") // Bắt buộc phải set SaveTag để lưu
	return h
}

func (h *Hero) OnCreate() {
	h.AddSprite("normal", napi.Assert.GetSprite("character"))
	h.SetCurrentSprite("normal")

}

func (h *Hero) OnStep() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		h.SetY(h.Y() - 3)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		h.SetY(h.Y() + 3)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		h.SetX(h.X() - 3)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		h.SetX(h.X() + 3)
	}
	h.Log(fmt.Sprintf("%f-%f", h.X(), h.Y()))
}

// Hook lưu dữ liệu
func (h *Hero) OnSave(data map[string]any) {
	data["x"] = h.X()
	data["y"] = h.Y()
	log.Printf("Saving Hero... X: %v, Y: %v\n", h.X(), h.Y())
}

// Hook tải dữ liệu
func (h *Hero) OnLoad(data map[string]any) {
	if x, ok := data["x"].(float64); ok {
		h.SetX(float32(x))
	}
	if y, ok := data["y"].(float64); ok {
		h.SetY(float32(y))
	}
	log.Printf("Loaded Hero... X: %v, Y: %v\n", h.X(), h.Y())
}

// --- Khai báo Bouncing Ball ---
type BouncingBall struct {
	napi.IObject
	ncom.Pos
	ncom.Drw
	ncom.Info
	ncom.Deb
	ncom.Box
	VX, VY float32
	Color  color.RGBA
}

func NewBouncingBall(id string, x, y, vx, vy float32, clr color.RGBA) *BouncingBall {
	b := &BouncingBall{VX: vx, VY: vy, Color: clr}
	napi.Obj.NewObject(b, "ball_"+id, "pos drw inf sce-cur deb box")
	b.SetX(x)
	b.SetY(y)
	b.SetSaveTag("ball_" + id)
	b.Debug("box info")
	b.SetBoxH(15)
	b.SetBoxW(15)
	return b
}

func (b *BouncingBall) OnStep() {
	b.SetX(b.X() + b.VX)
	b.SetY(b.Y() + b.VY)

	// Bounce
	if b.X() <= b.BoxW() || b.X() >= 610 {
		b.VX = -b.VX
	}
	if b.Y() <= b.BoxH() || b.Y() >= 450 {
		b.VY = -b.VY
	}
}

func (b *BouncingBall) Draw() {
	b.Circle(b.X(), b.Y(), 15, b.Color)
}

func (b *BouncingBall) OnSave(data map[string]any) {
	data["x"] = b.X()
	data["y"] = b.Y()
	data["vx"] = b.VX
	data["vy"] = b.VY
}

func (b *BouncingBall) OnLoad(data map[string]any) {
	if x, ok := data["x"].(float64); ok {
		b.SetX(float32(x))
	}
	if y, ok := data["y"].(float64); ok {
		b.SetY(float32(y))
	}
	if vx, ok := data["vx"].(float64); ok {
		b.VX = float32(vx)
	}
	if vy, ok := data["vy"].(float64); ok {
		b.VY = float32(vy)
	}
}

// --- Khai báo Fading Star ---
type FadingStar struct {
	napi.IObject
	ncom.Pos
	ncom.Drw
	ncom.Info
	Alpha float32
	Dir   float32
}

func NewFadingStar(id string, x, y float32) *FadingStar {
	s := &FadingStar{Alpha: 255, Dir: -5}
	napi.Obj.NewObject(s, "star_"+id, "pos drw inf sce-cur")
	s.SetX(x)
	s.SetY(y)
	s.SetSaveTag("star_" + id)
	return s
}

func (s *FadingStar) OnStep() {
	s.Alpha += s.Dir
	if s.Alpha <= 0 {
		s.Alpha = 0
		s.Dir = 5
	} else if s.Alpha >= 255 {
		s.Alpha = 255
		s.Dir = -5
	}
}

func (s *FadingStar) Draw() {
	c := color.RGBA{255, 255, 0, uint8(s.Alpha)}
	// Vẽ hình thoi đơn giản để giả làm sao
	s.Rect(s.X(), s.Y(), 20, 20, c)
	s.Text("Star", s.X()-5, s.Y()-10, c)
}

func (s *FadingStar) OnSave(data map[string]any) {
	data["alpha"] = s.Alpha
	data["dir"] = s.Dir
}

func (s *FadingStar) OnLoad(data map[string]any) {
	if a, ok := data["alpha"].(float64); ok {
		s.Alpha = float32(a)
	}
	if d, ok := data["dir"].(float64); ok {
		s.Dir = float32(d)
	}
}

// --- Khai báo Dynamic Box & Spawner ---
type DynamicBox struct {
	napi.IObject
	ncom.Pos
	ncom.Drw
	ncom.Info
}

func NewDynamicBox(id string, x, y float32) *DynamicBox {
	d := &DynamicBox{}
	napi.Obj.NewObject(d, "dyn_"+id, "pos drw inf sce-cur")
	d.SetX(x)
	d.SetY(y)
	d.SetSaveTag("dyn_" + id)
	return d
}

func (d *DynamicBox) OnStep() {}

func (d *DynamicBox) Draw() {
	d.Rect(d.X(), d.Y(), 20, 20, color.RGBA{200, 100, 50, 255})
	d.Text("Dyn", d.X(), d.Y()-5, color.RGBA{255, 255, 255, 255})
}

func (d *DynamicBox) OnSave(data map[string]any) {
	data["x"] = d.X()
	data["y"] = d.Y()
}

func (d *DynamicBox) OnLoad(data map[string]any) {
	if x, ok := data["x"].(float64); ok {
		d.SetX(float32(x))
	}
	if y, ok := data["y"].(float64); ok {
		d.SetY(float32(y))
	}
}

type Spawner struct {
	napi.IObject
	ncom.Info
	Boxes   []*DynamicBox
	Counter int
}

func NewSpawner() *Spawner {
	s := &Spawner{}
	napi.Obj.NewObject(s, "spawner", "inf sce-cur")
	s.SetSaveTag("spawner_sys")
	return s
}

func (s *Spawner) OnStep() {
	// Nhấn Space để tạo thêm 1 đối tượng
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.Counter++
		bx := NewDynamicBox(fmt.Sprintf("%d", s.Counter), 50+float32(len(s.Boxes)*30), 350)
		s.Boxes = append(s.Boxes, bx)
		log.Printf("Spawned Box #%d. Total: %d\n", s.Counter, len(s.Boxes))
	}
	// Nhấn Delete để xóa đối tượng cuối cùng
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		if len(s.Boxes) > 0 {
			last := s.Boxes[len(s.Boxes)-1]
			napi.Obj.Remove(last) // Remove thật khỏi Engine (deferred, cuối frame)
			s.Boxes = s.Boxes[:len(s.Boxes)-1]
			log.Printf("Deleted a Box. Total: %d\n", len(s.Boxes))
		}
	}
}

// boxEntry lưu toạ độ của mỗi box để Spawner tự restore khi Load.
type boxEntry struct {
	X, Y float32
}

func (s *Spawner) OnSave(data map[string]any) {
	positions := make([]map[string]any, len(s.Boxes))
	for i, b := range s.Boxes {
		positions[i] = map[string]any{"x": b.X(), "y": b.Y()}
	}
	data["count"] = len(s.Boxes)
	data["positions"] = positions
}

func (s *Spawner) OnLoad(data map[string]any) {
	targetCount := 0
	if c, ok := data["count"].(float64); ok {
		targetCount = int(c)
	}
	log.Printf("Spawner: Restoring %d boxes\n", targetCount)

	// Xóa bớt box hiện tại nếu thừa (bằng Remove thật)
	for len(s.Boxes) > targetCount {
		last := s.Boxes[len(s.Boxes)-1]
		napi.Obj.Remove(last)
		s.Boxes = s.Boxes[:len(s.Boxes)-1]
	}

	// Đọc vị trí đã lưu
	var savedPositions []map[string]any
	if raw, ok := data["positions"].([]any); ok {
		for _, item := range raw {
			if m, ok := item.(map[string]any); ok {
				savedPositions = append(savedPositions, m)
			}
		}
	}

	// Spawn thêm nếu thiếu, đặt đúng tọa độ đã lưu
	for len(s.Boxes) < targetCount {
		idx := len(s.Boxes)
		s.Counter++
		var x, y float32
		if idx < len(savedPositions) {
			if xv, ok := savedPositions[idx]["x"].(float64); ok {
				x = float32(xv)
			}
			if yv, ok := savedPositions[idx]["y"].(float64); ok {
				y = float32(yv)
			}
		} else {
			x = 50 + float32(idx*30)
			y = 350
		}
		bx := NewDynamicBox(fmt.Sprintf("%d", s.Counter), x, y)
		s.Boxes = append(s.Boxes, bx)
	}
}

// --- Khai báo Nút bấm (Button) ---
type Button struct {
	napi.IObject
	ncom.Pos
	ncom.Drw

	Label   string
	OnClick func()
}

func NewButton(label string, x, y float32, onClick func()) *Button {
	b := &Button{Label: label, OnClick: onClick}
	napi.Obj.NewObject(b, "button_"+label, "pos drw sce-cur")
	b.SetX(x)
	b.SetY(y)
	return b
}

func (b *Button) OnStep() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		// Box check đơn giản cho nút
		if float32(mx) >= b.X() && float32(mx) <= b.X()+120 &&
			float32(my) >= b.Y() && float32(my) <= b.Y()+40 {
			if b.OnClick != nil {
				b.OnClick()
			}
		}
	}
}

func (b *Button) Draw() {
	b.Rect(b.X(), b.Y(), 120, 40, color.RGBA{80, 100, 220, 255})
	b.Text(b.Label, b.X()+15, b.Y()+25, color.RGBA{255, 255, 255, 255})
}

// --- Hàm Main ---
func main() {
	napi.Game.Init(napi.GameConfig{
		Title:   "Save/Load Demo",
		Width:   640,
		Height:  480,
		SaveDir: "saves", // Nơi lưu trữ file
	})

	napi.Scene.NewSceneAndGo("main", "map-640-480")

	napi.Game.LoadFromFile("./SharedObject/Config.toml")
	// Tạo đối tượng Hero di chuyển
	NewHero()

	// Tạo nhiều đối tượng BouncingBall
	NewBouncingBall("1", 100, 100, 2, 3, color.RGBA{255, 0, 0, 255})
	NewBouncingBall("2", 300, 150, -3, 2, color.RGBA{0, 0, 255, 255})
	NewBouncingBall("3", 500, 300, 4, -2, color.RGBA{255, 0, 255, 255})

	// Tạo vài FadingStar
	NewFadingStar("1", 150, 400)
	NewFadingStar("2", 450, 80)

	// Tạo hệ thống Spawner (quản lý thêm/xóa)
	NewSpawner()

	// Tạo nút Save
	NewButton("SAVE (Click)", 20, 20, func() {
		path := "./save_load_demo/demo_slot.json"
		err := napi.Game.SaveGame(path)
		if err != nil {
			log.Println("Save Failed:", err)
		} else {
			log.Println("Save Successful to", path)
		}
	})

	// Tạo nút Load
	NewButton("LOAD (Click)", 160, 20, func() {
		path := "./save_load_demo/demo_slot.json"
		err := napi.Game.LoadGame(path)
		if err != nil {
			log.Println("Load Failed:", err)
		} else {
			log.Println("Load Successful from", path)
		}
	})

	napi.Game.GameStart()
}
