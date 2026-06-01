package layout

import (
	"testing"
)

// ─── Helpers ───────────────────────────────────────────────────────────────────

// newDiv tạo Div đơn giản tại (0,0) với kích thước cho trước.
func newDiv(name string, x, y, w, h int) *Div {
	return NewDivWithNameAndWidthAndHeight(name, x, y, w, h, "px")
}

// newA tạo thẻ A không có Game Object (object = nil).
// Dùng để test layout thuần mà không cần ECS.
func newA(name string, w, h int) *A {
	return NewA(name, nil, w, h)
}

// ─── Div: Constructor ──────────────────────────────────────────────────────────

// TestDivConstructor kiểm tra Div được khởi tạo đúng giá trị.
func TestDivConstructor(t *testing.T) {
	d := newDiv("root", 10, 20, 400, 300)

	if d.Name() != "root" {
		t.Errorf("Name: want %q, got %q", "root", d.Name())
	}
	if d.BoxX() != 10 || d.BoxY() != 20 {
		t.Errorf("Position: want (10,20), got (%d,%d)", d.BoxX(), d.BoxY())
	}
	if d.BoxW() != 400 || d.BoxH() != 300 {
		t.Errorf("Size: want (400,300), got (%d,%d)", d.BoxW(), d.BoxH())
	}
	// Config mặc định: DirRow | AlignStart | JustifyStart
	wantConfig := DirRow | AlignStart | JustifyStart
	if d.LayoutConfig() != wantConfig {
		t.Errorf("LayoutConfig: want %d, got %d", wantConfig, d.LayoutConfig())
	}
	if !d.IsVisible() {
		t.Error("IsVisible: want true by default")
	}
}

// ─── Div: Children management ─────────────────────────────────────────────────

// TestDivAddRemoveChildren kiểm tra AddChildren và RemoveChildren hoạt động đúng.
func TestDivAddRemoveChildren(t *testing.T) {
	parent := newDiv("parent", 0, 0, 800, 600)
	c1 := newDiv("c1", 0, 0, 100, 100)
	c2 := newDiv("c2", 0, 0, 100, 100)

	parent.AddChildren(c1)
	parent.AddChildren(c2)

	if len(parent.Childrens()) != 2 {
		t.Fatalf("AddChildren: want 2 children, got %d", len(parent.Childrens()))
	}

	parent.RemoveChildren(c1)
	if len(parent.Childrens()) != 1 {
		t.Fatalf("RemoveChildren: want 1 child, got %d", len(parent.Childrens()))
	}
	if parent.Childrens()[0] != c2 {
		t.Error("RemoveChildren: wrong child was removed")
	}
}

// TestDivGetChildrenByName kiểm tra tìm kiếm đệ quy theo tên.
func TestDivGetChildrenByName(t *testing.T) {
	root := newDiv("root", 0, 0, 800, 600)
	child := newDiv("child", 0, 0, 200, 200)
	grandchild := newDiv("grandchild", 0, 0, 100, 100)

	child.AddChildren(grandchild)
	root.AddChildren(child)

	// Tìm trực tiếp
	if root.GetChildrenByName("child") != child {
		t.Error("GetChildrenByName: failed to find direct child")
	}
	// Tìm đệ quy (grandchild nằm trong child, không trực tiếp trong root)
	if root.GetChildrenByName("grandchild") != grandchild {
		t.Error("GetChildrenByName: failed to find grandchild recursively")
	}
	// Không tồn tại
	if root.GetChildrenByName("ghost") != nil {
		t.Error("GetChildrenByName: should return nil for missing name")
	}
}

// ─── Div: DirRow layouts ───────────────────────────────────────────────────────

// TestDivRowJustifyStart kiểm tra layout Row + JustifyStart (mặc định).
// Các con xếp từ trái qua phải bắt đầu từ x của Div.
func TestDivRowJustifyStart(t *testing.T) {
	root := newDiv("root", 50, 50, 400, 200)
	root.SetLayoutConfig(DirRow | AlignStart | JustifyStart)

	c1 := newA("c1", 80, 40)
	c2 := newA("c2", 80, 40)
	c3 := newA("c3", 80, 40)
	root.AddChildren(c1)
	root.AddChildren(c2)
	root.AddChildren(c3)

	root.ComputeLayout(50, 50, 400, 200)

	// c1 bắt đầu ở x=50 (x của root)
	if c1.BoxX() != 50 {
		t.Errorf("c1.BoxX: want 50, got %d", c1.BoxX())
	}
	// c2 bắt đầu ở x = 50 + 80 = 130
	if c2.BoxX() != 130 {
		t.Errorf("c2.BoxX: want 130, got %d", c2.BoxX())
	}
	// c3 bắt đầu ở x = 130 + 80 = 210
	if c3.BoxX() != 210 {
		t.Errorf("c3.BoxX: want 210, got %d", c3.BoxX())
	}
	// Tất cả AlignStart → y = y của root = 50
	for _, c := range []ILayout{c1, c2, c3} {
		if c.BoxY() != 50 {
			t.Errorf("%s.BoxY: want 50, got %d", c.Name(), c.BoxY())
		}
	}
}

// TestDivRowJustifyCenter kiểm tra layout Row + JustifyCenter.
// Các con được căn giữa theo trục ngang.
func TestDivRowJustifyCenter(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 200)
	root.SetLayoutConfig(DirRow | AlignStart | JustifyCenter)

	// 2 item, mỗi item 100px wide → tổng 200px → startX = (400-200)/2 = 100
	c1 := newA("c1", 100, 50)
	c2 := newA("c2", 100, 50)
	root.AddChildren(c1)
	root.AddChildren(c2)

	root.ComputeLayout(0, 0, 400, 200)

	if c1.BoxX() != 100 {
		t.Errorf("c1.BoxX: want 100, got %d", c1.BoxX())
	}
	if c2.BoxX() != 200 {
		t.Errorf("c2.BoxX: want 200, got %d", c2.BoxX())
	}
}

// TestDivRowJustifyEnd kiểm tra layout Row + JustifyEnd.
// Các con được đẩy về bên phải.
func TestDivRowJustifyEnd(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 200)
	root.SetLayoutConfig(DirRow | AlignStart | JustifyEnd)

	// 2 item × 100px = 200px tổng → startX = 400 - 200 = 200
	c1 := newA("c1", 100, 50)
	c2 := newA("c2", 100, 50)
	root.AddChildren(c1)
	root.AddChildren(c2)

	root.ComputeLayout(0, 0, 400, 200)

	if c1.BoxX() != 200 {
		t.Errorf("c1.BoxX: want 200, got %d", c1.BoxX())
	}
	if c2.BoxX() != 300 {
		t.Errorf("c2.BoxX: want 300, got %d", c2.BoxX())
	}
}

// TestDivRowAlignCenter kiểm tra layout Row + AlignCenter.
// Các con được căn giữa theo trục dọc.
func TestDivRowAlignCenter(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 200)
	root.SetLayoutConfig(DirRow | AlignCenter | JustifyStart)

	// Item cao 40px, root cao 200px → y = (200-40)/2 = 80
	c1 := newA("c1", 100, 40)
	root.AddChildren(c1)

	root.ComputeLayout(0, 0, 400, 200)

	if c1.BoxY() != 80 {
		t.Errorf("c1.BoxY: want 80, got %d", c1.BoxY())
	}
}

// TestDivRowAlignEnd kiểm tra layout Row + AlignEnd.
// Các con được đẩy xuống dưới theo trục dọc.
func TestDivRowAlignEnd(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 200)
	root.SetLayoutConfig(DirRow | AlignEnd | JustifyStart)

	// Item cao 40px, root cao 200px → y = 200 - 40 = 160
	c1 := newA("c1", 100, 40)
	root.AddChildren(c1)

	root.ComputeLayout(0, 0, 400, 200)

	if c1.BoxY() != 160 {
		t.Errorf("c1.BoxY: want 160, got %d", c1.BoxY())
	}
}

// ─── Div: DirColumn layouts ────────────────────────────────────────────────────

// TestDivColumnJustifyStart kiểm tra layout Column + JustifyStart.
// Các con xếp từ trên xuống dưới bắt đầu từ y của Div.
func TestDivColumnJustifyStart(t *testing.T) {
	root := newDiv("root", 0, 10, 400, 300)
	root.SetLayoutConfig(DirColumn | AlignStart | JustifyStart)

	c1 := newA("c1", 100, 50)
	c2 := newA("c2", 100, 60)
	root.AddChildren(c1)
	root.AddChildren(c2)

	root.ComputeLayout(0, 10, 400, 300)

	if c1.BoxY() != 10 {
		t.Errorf("c1.BoxY: want 10, got %d", c1.BoxY())
	}
	if c2.BoxY() != 60 { // 10 + 50
		t.Errorf("c2.BoxY: want 60, got %d", c2.BoxY())
	}
}

// TestDivColumnJustifyCenter kiểm tra layout Column + JustifyCenter.
func TestDivColumnJustifyCenter(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 300)
	root.SetLayoutConfig(DirColumn | AlignStart | JustifyCenter)

	// 2 item × 50px = 100px tổng → startY = (300-100)/2 = 100
	c1 := newA("c1", 100, 50)
	c2 := newA("c2", 100, 50)
	root.AddChildren(c1)
	root.AddChildren(c2)

	root.ComputeLayout(0, 0, 400, 300)

	if c1.BoxY() != 100 {
		t.Errorf("c1.BoxY: want 100, got %d", c1.BoxY())
	}
	if c2.BoxY() != 150 {
		t.Errorf("c2.BoxY: want 150, got %d", c2.BoxY())
	}
}

// TestDivColumnJustifyEnd kiểm tra layout Column + JustifyEnd.
func TestDivColumnJustifyEnd(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 300)
	root.SetLayoutConfig(DirColumn | AlignStart | JustifyEnd)

	// 2 item × 50px = 100px → startY = 300 - 100 = 200
	c1 := newA("c1", 100, 50)
	c2 := newA("c2", 100, 50)
	root.AddChildren(c1)
	root.AddChildren(c2)

	root.ComputeLayout(0, 0, 400, 300)

	if c1.BoxY() != 200 {
		t.Errorf("c1.BoxY: want 200, got %d", c1.BoxY())
	}
	if c2.BoxY() != 250 {
		t.Errorf("c2.BoxY: want 250, got %d", c2.BoxY())
	}
}

// TestDivColumnAlignCenter kiểm tra Column + AlignCenter (căn giữa trục ngang).
func TestDivColumnAlignCenter(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 300)
	root.SetLayoutConfig(DirColumn | AlignCenter | JustifyStart)

	// Item rộng 100px, root rộng 400px → x = (400-100)/2 = 150
	c1 := newA("c1", 100, 50)
	root.AddChildren(c1)

	root.ComputeLayout(0, 0, 400, 300)

	if c1.BoxX() != 150 {
		t.Errorf("c1.BoxX: want 150, got %d", c1.BoxX())
	}
}

// TestDivColumnAlignEnd kiểm tra Column + AlignEnd (căn phải trục ngang).
func TestDivColumnAlignEnd(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 300)
	root.SetLayoutConfig(DirColumn | AlignEnd | JustifyStart)

	// Item rộng 100px, root rộng 400px → x = 400 - 100 = 300
	c1 := newA("c1", 100, 50)
	root.AddChildren(c1)

	root.ComputeLayout(0, 0, 400, 300)

	if c1.BoxX() != 300 {
		t.Errorf("c1.BoxX: want 300, got %d", c1.BoxX())
	}
}

// ─── Div: Gap ─────────────────────────────────────────────────────────────────

// TestDivRowGap kiểm tra gap giữa các item trong Row.
func TestDivRowGap(t *testing.T) {
	root := newDiv("root", 0, 0, 500, 200)
	root.SetLayoutConfig(DirRow | AlignStart | JustifyStart)
	root.SetGap(20)

	// c1 tại x=0, c2 tại x=0+100+20=120, c3 tại x=120+100+20=240
	c1 := newA("c1", 100, 50)
	c2 := newA("c2", 100, 50)
	c3 := newA("c3", 100, 50)
	root.AddChildren(c1)
	root.AddChildren(c2)
	root.AddChildren(c3)

	root.ComputeLayout(0, 0, 500, 200)

	if c1.BoxX() != 0 {
		t.Errorf("c1.BoxX: want 0, got %d", c1.BoxX())
	}
	if c2.BoxX() != 120 {
		t.Errorf("c2.BoxX: want 120, got %d", c2.BoxX())
	}
	if c3.BoxX() != 240 {
		t.Errorf("c3.BoxX: want 240, got %d", c3.BoxX())
	}
}

// TestDivColumnGap kiểm tra gap giữa các item trong Column.
func TestDivColumnGap(t *testing.T) {
	root := newDiv("root", 0, 0, 200, 500)
	root.SetLayoutConfig(DirColumn | AlignStart | JustifyStart)
	root.SetGap(15)

	c1 := newA("c1", 100, 60)
	c2 := newA("c2", 100, 60)
	root.AddChildren(c1)
	root.AddChildren(c2)

	root.ComputeLayout(0, 0, 200, 500)

	if c1.BoxY() != 0 {
		t.Errorf("c1.BoxY: want 0, got %d", c1.BoxY())
	}
	if c2.BoxY() != 75 { // 60 + 15
		t.Errorf("c2.BoxY: want 75, got %d", c2.BoxY())
	}
}

// ─── Div: Edge cases ──────────────────────────────────────────────────────────

// TestDivComputeNoChildren kiểm tra không có gì xảy ra khi Div trống.
func TestDivComputeNoChildren(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 300)
	// Phải không panic khi không có children
	root.ComputeLayout(0, 0, 400, 300)
}

// TestDivSingleChild kiểm tra JustifyCenter với 1 child (không có gap được cộng thêm).
func TestDivSingleChild(t *testing.T) {
	root := newDiv("root", 0, 0, 400, 200)
	root.SetLayoutConfig(DirRow | AlignCenter | JustifyCenter)

	c1 := newA("c1", 100, 100)
	root.AddChildren(c1)

	root.ComputeLayout(0, 0, 400, 200)

	// x = (400-100)/2 = 150, y = (200-100)/2 = 50
	if c1.BoxX() != 150 {
		t.Errorf("c1.BoxX: want 150, got %d", c1.BoxX())
	}
	if c1.BoxY() != 50 {
		t.Errorf("c1.BoxY: want 50, got %d", c1.BoxY())
	}
}

// ─── A: Basic behaviour ───────────────────────────────────────────────────────

// TestAConstructor kiểm tra A được khởi tạo đúng.
func TestAConstructor(t *testing.T) {
	a := newA("btn", 120, 40)

	if a.Name() != "btn" {
		t.Errorf("Name: want %q, got %q", "btn", a.Name())
	}
	if a.BoxW() != 120 || a.BoxH() != 40 {
		t.Errorf("Size: want (120,40), got (%d,%d)", a.BoxW(), a.BoxH())
	}
	if a.Object() != nil {
		t.Error("Object: want nil for NewA(nil, ...)")
	}
}

// TestASetBoxXY kiểm tra SetBoxX/SetBoxY cập nhật tọa độ nội bộ khi object == nil.
func TestASetBoxXY(t *testing.T) {
	a := newA("a", 50, 50)
	a.SetBoxX(77)
	a.SetBoxY(88)

	if a.BoxX() != 77 {
		t.Errorf("BoxX: want 77, got %d", a.BoxX())
	}
	if a.BoxY() != 88 {
		t.Errorf("BoxY: want 88, got %d", a.BoxY())
	}
}

// TestAComputeLayoutNoObject kiểm tra ComputeLayout không panic khi object == nil.
func TestAComputeLayoutNoObject(t *testing.T) {
	a := newA("a", 50, 50)
	a.SetBoxX(10)
	a.SetBoxY(20)
	// Phải không panic
	a.ComputeLayout(10, 20, 50, 50)
}

// TestAGetChildrenByName kiểm tra A trả về nil cho mọi tên (A không có children đệ quy).
func TestAGetChildrenByName(t *testing.T) {
	a := newA("a", 50, 50)
	result := a.GetChildrenByName("anything")
	if result != nil {
		t.Error("GetChildrenByName: want nil on A with no children")
	}
}

// ─── Integration: Nested Div ──────────────────────────────────────────────────

// TestNestedDivLayout kiểm tra layout lồng nhau: Div chứa Div khác chứa A.
// Mô phỏng cấu trúc: root (Column,Center,Center) → row (Row,Center,Center) → các nút.
func TestNestedDivLayout(t *testing.T) {
	// root 800×600 tại (0,0), Column, Center, Center
	root := newDiv("root", 0, 0, 800, 600)
	root.SetLayoutConfig(DirColumn | AlignCenter | JustifyCenter)

	// row 400×100, Row, Center, Center
	row := newDiv("row", 0, 0, 400, 100)
	row.SetLayoutConfig(DirRow | AlignCenter | JustifyCenter)

	btn1 := newA("btn1", 100, 60)
	btn2 := newA("btn2", 100, 60)
	row.AddChildren(btn1)
	row.AddChildren(btn2)
	root.AddChildren(row)

	// root ComputeLayout → đặt row vào giữa root
	root.ComputeLayout(0, 0, 800, 600)

	// row phải được căn giữa: x = (800-400)/2 = 200, y = (600-100)/2 = 250
	if row.BoxX() != 200 {
		t.Errorf("row.BoxX: want 200, got %d", row.BoxX())
	}
	if row.BoxY() != 250 {
		t.Errorf("row.BoxY: want 250, got %d", row.BoxY())
	}

	// row tính layout con của mình: 2 btn × 100px = 200px → startX = row.x + (400-200)/2 = 200+100 = 300
	// btn1 tại x=300, btn2 tại x=400
	// y con = row.y + (100-60)/2 = 250 + 20 = 270
	if btn1.BoxX() != 300 {
		t.Errorf("btn1.BoxX: want 300, got %d", btn1.BoxX())
	}
	if btn2.BoxX() != 400 {
		t.Errorf("btn2.BoxX: want 400, got %d", btn2.BoxX())
	}
	if btn1.BoxY() != 270 {
		t.Errorf("btn1.BoxY: want 270, got %d", btn1.BoxY())
	}
}

// TestOffsetRootPosition kiểm tra layout với root không ở (0,0).
// Mọi tọa độ con phải được tính tương đối so với vị trí root.
func TestOffsetRootPosition(t *testing.T) {
	root := newDiv("root", 100, 200, 400, 200)
	root.SetLayoutConfig(DirRow | AlignStart | JustifyStart)

	c1 := newA("c1", 80, 50)
	root.AddChildren(c1)

	root.ComputeLayout(100, 200, 400, 200)

	// c1 phải bắt đầu tại x=100 (x của root), y=200 (y của root)
	if c1.BoxX() != 100 {
		t.Errorf("c1.BoxX: want 100, got %d", c1.BoxX())
	}
	if c1.BoxY() != 200 {
		t.Errorf("c1.BoxY: want 200, got %d", c1.BoxY())
	}
}
