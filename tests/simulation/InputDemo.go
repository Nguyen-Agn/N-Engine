//go:build ignore

// InputDemo demonstrates the full input system:
//   - Keyboard: Pressed (hold), JustPressed (edge ↓), JustReleased (edge ↑)
//   - Key groups: "alpha" (text-like input), "number", "all"
//   - Mouse: button events (left/right/middle), cursor position, scroll wheel
//
// Events are printed to the console AND displayed on screen via DrawComponent.
//
// Run: go run .\tests\simulation\InputDemo.go
package main

import (
	"fmt"
	"image/color"
	"strings"

	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
)

// ─── InputLogger — ghi lại các sự kiện gần nhất ──────────────────────────────

// maxLog là số dòng log tối đa hiển thị trên màn hình.
const maxLog = 16

// logLines lưu các sự kiện để render lên màn hình.
var logLines []string

// pushLog thêm một dòng log, giới hạn maxLog dòng.
func pushLog(line string) {
	fmt.Println(line) // in ra terminal
	logLines = append(logLines, line)
	if len(logLines) > maxLog {
		logLines = logLines[len(logLines)-maxLog:]
	}
}

// ─── InputTester — object chính lắng nghe tất cả sự kiện ────────────────────

type InputTester struct {
	napi.IObject
	ncom.Pos   // SetX/Y, X/Y — kéo theo "pos" token
	ncom.Drw   // token "drw"
	ncom.Inp   // token "inp" — keyboard bindings
	ncom.Mouse // không cần ECS token — mouse bindings
	ncom.Deb   // token "deb" — hiển thị debug
	ncom.Box   // token "box" — hiển thị debug
}

func (t *InputTester) OnCreate() {
	t.SetX(100)
	t.SetY(100)

	// Bật debug
	t.Debug("pos box info")
	t.Log("Demo Object")

	// ── Keyboard: phím đơn & phím gộp bằng khoảng trắng ─────────────────

	// Nhiều phím cách nhau bằng dấu cách
	t.ListenOn("w a s d up down left right", "", func(key string) {
		pushLog(fmt.Sprintf("[HOLD  ] (multi) → \"%s\"", key))
	})

	// Nhấn 1 lần (JustPressed)
	t.ListenOn("space enter", "pressed", func(key string) {
		pushLog(fmt.Sprintf("[DOWN  ] (multi) → \"%s\"", key))
	})

	// Escape — JustReleased
	t.ListenOn("escape", "released", func(key string) {
		pushLog(fmt.Sprintf("[UP    ] escape → \"%s\"", key))
	})

	// ── Keyboard: nhóm phím ───────────────────────────────────────────────

	// "alpha" — nhập văn bản: bắt JustPressed để tránh lặp, giống nhập chữ
	t.ListenOn("alpha", "pressed", func(key string) {
		pushLog(fmt.Sprintf("[ALPHA ] JustPressed → \"%s\"  (typed letter)", key))
	})

	// "number" — nhập số: JustPressed
	t.ListenOn("number", "pressed", func(key string) {
		pushLog(fmt.Sprintf("[NUMBER] JustPressed → \"%s\"  (typed digit)", key))
	})

	// "all" — JustReleased: bắt tất cả phím khi thả, demo tổng quát
	t.ListenOn("all", "released", func(key string) {
		pushLog(fmt.Sprintf("[ALL   ] JustReleased → \"%s\"", key))
	})

	// ── Mouse: nút chuột ─────────────────────────────────────────────────

	// Trái — JustPressed
	t.ListenMouseOn("left", "pressed", func(btn string) {
		mx, my := t.MouseX(), t.MouseY()
		pushLog(fmt.Sprintf("[MOUSE ] left  DOWN  @ (%d, %d)", mx, my))
	})

	// Trái — JustReleased
	t.ListenMouseOn("left", "released", func(btn string) {
		mx, my := t.MouseX(), t.MouseY()
		pushLog(fmt.Sprintf("[MOUSE ] left  UP    @ (%d, %d)", mx, my))
	})

	// Nhiều nút chuột cách nhau bằng dấu cách
	t.ListenMouseOn("right middle", "pressed", func(btn string) {
		pushLog(fmt.Sprintf("[MOUSE ] %s DOWN", btn))
	})

	// Giữ chuột trái — Pressed (mỗi frame)
	t.ListenMouseOn("left", "", func(btn string) {
		mx, my := t.MouseX(), t.MouseY()
		pushLog(fmt.Sprintf("[MOUSE ] left  HOLD  @ (%d, %d)", mx, my))
	})
}

func (t *InputTester) OnStep() {
	// Cuộn bánh xe — kiểm tra mỗi frame
	wx, wy := t.WheelX(), t.WheelY()
	if wx != 0 || wy != 0 {
		pushLog(fmt.Sprintf("[WHEEL ] X=%.2f  Y=%.2f", wx, wy))
	}
}

func (t *InputTester) Draw() {
	mx, my := t.MouseX(), t.MouseY()

	// ── Tiêu đề ──────────────────────────────────────────────────────────
	t.Rect(0, 0, 640, 22, color.RGBA{20, 20, 40, 230})
	t.Text("INPUT DEMO  —  alpha / number / all / mouse", 8, 16, color.RGBA{120, 200, 255, 255})

	// ── Hướng dẫn ────────────────────────────────────────────────────────
	t.Rect(0, 22, 640, 38, color.RGBA{10, 10, 25, 200})
	hint := "WASD=hold  SPACE=↓  ENTER=↑  A-Z=alpha↓  0-9=num↓  ANY=↑  LMB/RMB/MMB  Scroll"
	t.Text(hint, 6, 36, color.RGBA{180, 180, 180, 200})
	t.Text(fmt.Sprintf("Mouse: (%d, %d)", mx, my), 490, 36, color.RGBA{255, 220, 80, 255})

	// ── Vùng log ─────────────────────────────────────────────────────────
	t.Rect(0, 60, 640, float32(maxLog*18+4), color.RGBA{8, 12, 22, 220})

	for i, line := range logLines {
		y := float32(60 + 4 + i*18)
		// Màu theo loại sự kiện
		c := lineColor(line)
		t.Text(line, 8, y+12, c)
	}

	// ── Con trỏ chuột (dấu thập) ─────────────────────────────────────────
	cx, cy := float32(mx), float32(my)
	t.Line(cx-10, cy, cx+10, cy, color.RGBA{255, 80, 80, 200}, 1.5)
	t.Line(cx, cy-10, cx, cy+10, color.RGBA{255, 80, 80, 200}, 1.5)
}

// lineColor trả về màu dựa trên prefix của dòng log.
func lineColor(line string) color.RGBA {
	switch {
	case strings.HasPrefix(line, "[HOLD"):
		return color.RGBA{100, 200, 100, 255} // xanh lá — giữ phím
	case strings.HasPrefix(line, "[DOWN"):
		return color.RGBA{100, 180, 255, 255} // xanh dương — nhấn
	case strings.HasPrefix(line, "[UP"):
		return color.RGBA{255, 200, 80, 255} // vàng — thả
	case strings.HasPrefix(line, "[ALPHA"):
		return color.RGBA{200, 120, 255, 255} // tím — chữ cái
	case strings.HasPrefix(line, "[NUMBER"):
		return color.RGBA{255, 160, 60, 255} // cam — chữ số
	case strings.HasPrefix(line, "[ALL"):
		return color.RGBA{160, 160, 160, 255} // xám — all
	case strings.HasPrefix(line, "[MOUSE"):
		return color.RGBA{255, 100, 140, 255} // hồng — chuột
	case strings.HasPrefix(line, "[WHEEL"):
		return color.RGBA{80, 230, 200, 255} // cyan — cuộn
	default:
		return color.RGBA{200, 200, 200, 255}
	}
}

func NewInputTester() *InputTester {
	t := &InputTester{}
	napi.Obj.NewObject(t, "input_tester", "pos drw inp deb box sce-cur")
	return t
}

// ─── main ─────────────────────────────────────────────────────────────────────

func main() {
	napi.Game.Init(napi.GameConfig{
		Title:  "Input Demo — alpha / number / all / mouse",
		Width:  640,
		Height: maxLog*18 + 70,
	})

	napi.Scene.NewSceneAndGo("main", "map-640-400")

	NewInputTester()

	fmt.Println("=== Input Demo ===")
	fmt.Println("WASD      → [HOLD  ] mỗi frame")
	fmt.Println("Space     → [DOWN  ] JustPressed")
	fmt.Println("Enter     → [UP    ] JustReleased")
	fmt.Println("A-Z       → [ALPHA ] JustPressed (giống nhập văn bản)")
	fmt.Println("0-9       → [NUMBER] JustPressed")
	fmt.Println("Any key   → [ALL   ] JustReleased")
	fmt.Println("LMB/RMB   → [MOUSE ] các sự kiện chuột")
	fmt.Println("Scroll    → [WHEEL ] tốc độ cuộn")
	fmt.Println("==================")

	napi.Game.GameStart()
}
