package domain

// ─── Layout Constants ──────────────────────────────────────────────────────────
// Các hằng số dưới đây được encode vào một int duy nhất (layoutConfig) dùng bitfield.
// Sử dụng các hàm bitwise để kết hợp: vd. DirColumn | AlignCenter | JustifyEnd

// Direction — chiều sắp xếp các con trong Layout.
const (
	DirRow    = 0 << 0 // Sắp xếp con theo hàng ngang (mặc định)
	DirColumn = 1 << 0 // Sắp xếp con theo cột dọc
	DirMask   = 1 << 0 // Bitmask để đọc trường Direction
)

// Align — căn chỉnh các con theo trục phụ (cross-axis).
// Khi DirRow: căn theo chiều dọc. Khi DirColumn: căn theo chiều ngang.
const (
	AlignStart  = 0 << 1 // Căn đầu (trái hoặc trên)
	AlignCenter = 1 << 1 // Căn giữa
	AlignEnd    = 2 << 1 // Căn cuối (phải hoặc dưới)
	AlignMask   = 7 << 1 // Bitmask để đọc trường Align
)

// Justify — phân bổ không gian các con theo trục chính (main-axis).
// Khi DirRow: theo chiều ngang. Khi DirColumn: theo chiều dọc.
const (
	JustifyStart  = 0 << 4 // Tất cả con dồn về đầu
	JustifyCenter = 1 << 4 // Tất cả con tập trung giữa
	JustifyEnd    = 2 << 4 // Tất cả con dồn về cuối
	JustifyMask   = 7 << 4 // Bitmask để đọc trường Justify
)

// ILayout định nghĩa giao diện cho một node trong cây Layout UI.
// Layout tính toán vị trí và kích thước của các phần tử UI theo kiểu flexbox đơn giản.
type ILayout interface {
	// SetBoxH thiết lập chiều cao của Layout node (pixel).
	SetBoxH(h int)
	// SetBoxW thiết lập chiều rộng của Layout node (pixel).
	SetBoxW(w int)
	// SetBoxX thiết lập tọa độ ngang của Layout node (thường do cha tính toán).
	SetBoxX(x int)
	// SetBoxY thiết lập tọa độ dọc của Layout node (thường do cha tính toán).
	SetBoxY(y int)

	// BoxH trả về chiều cao hiện tại của Layout node (pixel).
	BoxH() int
	// BoxW trả về chiều rộng hiện tại của Layout node (pixel).
	BoxW() int
	// BoxX trả về tọa độ ngang hiện tại của Layout node (pixel).
	BoxX() int
	// BoxY trả về tọa độ dọc hiện tại của Layout node (pixel).
	BoxY() int

	// AddChildren thêm một ILayout con vào node này.
	AddChildren(children ILayout)
	// RemoveChildren xóa một ILayout con khỏi node này.
	RemoveChildren(children ILayout)
	// Childrens trả về danh sách tất cả các Layout con.
	Childrens() []ILayout
	// GetChildrenByName tìm Layout con theo tên. Trả về nil nếu không tìm thấy.
	GetChildrenByName(name string) ILayout

	// Name trả về tên định danh của Layout node.
	Name() string
	// SetName thiết lập tên định danh của Layout node.
	SetName(name string)

	// IsVisible trả về trạng thái hiển thị của Layout node.
	IsVisible() bool
	// SetIsVisible bật/tắt hiển thị Layout node (ẩn node cũng ẩn toàn bộ cây con).
	SetIsVisible(isVisible bool)

	// SetLayoutConfig thiết lập config bitfield kết hợp Direction | Align | Justify.
	// Ví dụ: DirColumn | AlignCenter | JustifyStart
	SetLayoutConfig(config int)
	// LayoutConfig trả về config bitfield hiện tại.
	LayoutConfig() int

	// SetGap thiết lập khoảng cách giữa các Layout con (pixel).
	SetGap(gap int)
	// Gap trả về khoảng cách giữa các Layout con (pixel).
	Gap() int

	// ComputeLayout tính toán lại vị trí và kích thước của node này và toàn bộ cây con.
	// parentX, parentY là tọa độ góc trên-trái của node cha.
	// parentW, parentH là kích thước khả dụng từ node cha.
	ComputeLayout(parentX, parentY int, parentW, parentH int)
}
