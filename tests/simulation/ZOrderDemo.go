//go:build ignore

// ZOrderDemo demonstrates Z-Order (Depth Sorting) for Custom Drawing and Text Align:
//   - 4 objects created in random order but drawn according to ZOrder.
//   - Uses t.Rect() to draw background.
//   - Each object demonstrates a different SetTextAlign mode.
//
// Run: go run .\tests\simulation\ZOrderDemo.go
package main

import (
	"fmt"
	"image/color"

	"github.com/Nguyen-Agn/N-Engine/modules/napi"
	"github.com/Nguyen-Agn/N-Engine/modules/napi/ncom"
)

// ─── TestObj ─────────────────────────────────────────────────────────────────

type TestObj struct {
	ncom.Object
	ncom.Pos
	ncom.Spr // Cần nhúng Spr để có thuộc tính ZOrder
	ncom.Drw
	colorName string
	boxColor  color.RGBA
	hAlign    string
	vAlign    string
}

func (t *TestObj) BeCreated() {
	// Không cần khởi tạo hình ảnh sprite vì chúng ta vẽ Rect thủ công trong Draw()
}

func (t *TestObj) Draw() {
	// Vẽ khung màu nền (Kích thước 100x100)
	t.Rect(t.X(), t.Y(), 100, 100, t.boxColor)

	// Test text align
	t.SetTextAlign(t.hAlign, t.vAlign)

	// Điểm gốc để test align là chính giữa khung Rect (X+50, Y+50)
	t.Text(fmt.Sprintf("%s\nZ: %d\n%s,%s", t.colorName, t.ZOrder(), t.hAlign, t.vAlign), t.X()+50, t.Y()+50, color.RGBA{255, 255, 255, 255})

	// Vẽ một điểm nhỏ tại (X+50, Y+50) để làm mốc tham chiếu trực quan cho việc căn lề
	t.Circle(t.X()+50, t.Y()+50, 3, color.RGBA{255, 0, 0, 255})
}

// Hàm hỗ trợ tạo Object
func NewTestObj(x, y float32, z int, col color.RGBA, name, hAlign, vAlign string) *TestObj {
	obj := &TestObj{
		colorName: name,
		boxColor:  col,
		hAlign:    hAlign,
		vAlign:    vAlign,
	}

	// Đăng ký entity với các component: Position, Sprite (để có ZOrder), Draw, Scene-Current
	napi.Obj.NewObject(obj, "test_"+name, "pos spr drw sce-cur")

	obj.SetX(x)
	obj.SetY(y)
	obj.SetZOrder(z) // Gắn thứ tự vẽ
	return obj
}

// ─── main ─────────────────────────────────────────────────────────────────────

func main() {
	// Khởi tạo Engine
	napi.Game.Init(napi.GameConfig{
		Title:  "ZOrder & Text Align Demo",
		Width:  400,
		Height: 400,
	})

	// Tạo Scene
	napi.Scene.NewSceneAndGo("main", "map-400-400")

	// Tạo 4 đối tượng với thứ tự khởi tạo khác với Z-Order

	// Đỏ khởi tạo đầu tiên, Z=4 (vẽ trên cùng), căn giữa ("center", "center")
	NewTestObj(60, 60, 4, color.RGBA{200, 50, 50, 255}, "Đỏ", "center", "center")

	// Xanh lá khởi tạo thứ 2, Z=1 (vẽ dưới cùng), căn trái-trên ("left", "top")
	NewTestObj(240, 240, 1, color.RGBA{50, 200, 50, 255}, "Xanh Lá", "left", "top")

	// Xanh dương khởi tạo thứ 3, Z=3 (nằm dưới Đỏ), căn phải-dưới ("right", "bottom")
	NewTestObj(180, 100, 3, color.RGBA{50, 50, 200, 255}, "Xanh Dương", "right", "bottom")

	// Vàng khởi tạo cuối cùng, Z=2 (nằm trên Xanh lá, dưới Xanh dương), căn trái-dưới ("left", "bottom")
	NewTestObj(100, 180, 2, color.RGBA{200, 200, 50, 255}, "Vàng", "left", "bottom")

	// Bắt đầu game loop
	napi.Game.GameStart()
}
