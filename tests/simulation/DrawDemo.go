//go:build ignore

// DrawDemo demonstrates DrawComponent (token "drw"):
//   - Text rendering with default font
//   - Filled rectangle
//   - Rectangle outline (stroke)
//   - Filled circle
//   - Circle outline
//
// Run: go run .\tests\simulation\DrawDemo.go
package main

import (
	"fmt"
	"image/color"

	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
)

// ─── HUD: text + rect background ─────────────────────────────────────────────

type HUD struct {
	ncom.Object
	ncom.Pos
	ncom.Drw
	frame int
}

func (h *HUD) BeCreated() {
	h.SetX(10)
	h.SetY(10)
}

func (h *HUD) StepUpdate() {
	h.frame++
}

func (h *HUD) Draw() {
	// Semi-transparent black background panel
	h.Rect(8, 8, 220, 45, color.RGBA{0, 0, 0, 160})
	//	h.RectStroke(8, 8, 220, 45, color.RGBA{100, 200, 255, 255}, 2)
	// Frame counter text
	h.Text(fmt.Sprintf("Frame: %d", h.frame), 16, 30, color.RGBA{255, 230, 80, 255})
}

func NewHUD() *HUD {
	h := &HUD{}
	napi.Obj.NewObject(h, "hud", "pos drw sce-cur")
	return h
}

// ─── ShapeDemo: circles and rects ────────────────────────────────────────────

type ShapeDemo struct {
	ncom.Object
	ncom.Pos
	ncom.Drw
}

func (s *ShapeDemo) BeCreated() {
	s.SetX(160)
	s.SetY(160)
}

func (s *ShapeDemo) Draw() {
	cx, cy := s.X(), s.Y()

	// Filled rect
	s.Rect(cx-80, cy-40, 160, 80, color.RGBA{30, 100, 200, 180})
	// Rect stroke
	s.RectStroke(cx-80, cy-40, 160, 80, color.RGBA{255, 255, 255, 200}, 2)

	// Filled circle
	s.Circle(cx, cy+110, 45, color.RGBA{200, 60, 60, 200})
	// Circle stroke
	s.CircleStroke(cx, cy+110, 45, color.RGBA{255, 200, 200, 255}, 3)

	// Label
	s.Text("DrawComponent Demo", cx-75, cy-55, color.RGBA{200, 255, 200, 255})
}

func NewShapeDemo() *ShapeDemo {
	sd := &ShapeDemo{}
	napi.Obj.NewObject(sd, "shape_demo", "pos drw sce-cur")
	return sd
}

// ─── main ─────────────────────────────────────────────────────────────────────

func main() {
	napi.Game.Init(napi.GameConfig{
		Title:  "DrawComponent Demo",
		Width:  320,
		Height: 320,
	})

	napi.Scene.NewSceneAndGo("main", "map-320-320")
	NewHUD()
	NewShapeDemo()

	napi.Game.GameStart()
}
