# Bố cục Giao diện / UI Layout (nLayout)

> **Vision**: Layout flexbox-like — khai báo UI bằng code, không cần editor.
> **Vision**: Flexbox-like layout — declare UI in code, no editor needed.

---

## 1. Giải thích / Explanation

N-Engine dùng hệ thống **nLayout** để bố trí UI theo dạng flexbox, sau đó bạn có thể gắn Game Objects trực tiếp vào các tọa độ đã được tính toán.
N-Engine uses the **nLayout** system for flexbox-style UI layout, and you can attach Game Objects to the calculated coordinates.

**Hai loại node / Two node types:**
- `layout.Div` — container chứa nhiều children / container holding multiple children. Nó tính toán vị trí theo thuật toán Flexbox.
- `layout.A` — leaf node (nút lá), đại diện cho kích thước thực của một object / leaf node representing real object size.

**Hướng layout / Layout direction:**
- `layout.DirRow` — ngang (mặc định) / horizontal
- `layout.DirColumn` — dọc / vertical

**Cấu hình Layout:** 
Kết hợp hướng (Dir), Canh lề (Align), và Dàn đều (Justify) bằng toán tử Bitwise OR `|`.

---

## 2. Quy trình sử dụng (Workflow)

1. Tạo một container (`Div`) làm gốc.
2. Cấu hình khoảng cách (`Gap`) và kiểu dàn trang (`LayoutConfig`).
3. Tạo các phần tử con (`A`), chỉ định chiều rộng/cao của chúng, thêm vào `Div`.
4. Kích hoạt tính toán: `div.ComputeLayout(x, y, w, h)`.
5. Tạo các Game Object thực sự của bạn và đặt vị trí của chúng bằng tọa độ đã lấy từ `A.BoxX()` và `A.BoxY()`.

---

## 3. Ví dụ / Code Example

Đoạn code sau chia một vùng 400x300 pixel, bố trí theo chiều dọc, căn giữa tất cả các thành phần:

```go
package main

import (
	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
	layout "autoworld/modules/nlayout"
)

func createUIMenu() {
	// 1. Tạo node gốc (Root Div)
	// Vị trí: (100, 100), Kích thước: 400x300
	root := layout.NewDivWithNameAndWidthAndHeight("MenuRoot", 100, 100, 400, 300, "px")
	
	// 2. Bố trí Dọc, Căn Giữa ngang (AlignCenter), Dàn Giữa dọc (JustifyCenter)
	root.SetLayoutConfig(layout.DirColumn | layout.AlignCenter | layout.JustifyCenter)
	root.SetGap(20) // Khoảng cách giữa các nút là 20px

	// 3. Khai báo các Node lá (đại diện cho vị trí của các Nút bấm)
	btnStart := layout.NewA("btn-start", nil, 200, 50)
	btnOption := layout.NewA("btn-option", nil, 200, 50)
	btnExit := layout.NewA("btn-exit", nil, 200, 50)

	// Thêm vào root
	root.AddChildren(btnStart, btnOption, btnExit)

	// 4. Bắt đầu tính toán toàn bộ cành cây layout
	root.ComputeLayout(100, 100, 400, 300)

	// 5. Spawn các Object thực sự của Game Engine tại tọa độ đã tính
	NewMyCustomButton("Start", btnStart.BoxX(), btnStart.BoxY(), 200, 50)
	NewMyCustomButton("Option", btnOption.BoxX(), btnOption.BoxY(), 200, 50)
	NewMyCustomButton("Exit", btnExit.BoxX(), btnExit.BoxY(), 200, 50)
}

// Khai báo một Custom Object Button
type MyCustomButton struct {
	ncom.Object
	ncom.Pos
	ncom.Spr
}

func NewMyCustomButton(name string, x, y, w, h int) *MyCustomButton {
	btn := &MyCustomButton{}
	// Sử dụng sce-cur để tự động add vào Current Scene
	napi.Obj.NewObject(btn, name, "pos spr sce-cur")
	
	btn.SetX(float32(x))
	btn.SetY(float32(y))
	// Set Sprite v.v...
	
	return btn
}
```

> **Tham khảo chi tiết:** Bạn có thể tham khảo file `tests/simulation/LayoutDemo.go` trong source code để xem trực quan các bố trí nâng cao (Row/Column + Start/Center/End).
