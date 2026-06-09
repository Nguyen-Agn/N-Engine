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
	// Purpose: Sets the layout node's height.
	// Inputs: h int - The height in pixels.
	// Outputs: None.
	SetBoxH(h int)
	
	// SetBoxW thiết lập chiều rộng của Layout node (pixel).
	// Purpose: Sets the layout node's width.
	// Inputs: w int - The width in pixels.
	// Outputs: None.
	SetBoxW(w int)
	
	// SetBoxX thiết lập tọa độ ngang của Layout node (thường do cha tính toán).
	// Purpose: Sets the local X coordinate of the layout node.
	// Inputs: x int - The X position in pixels.
	// Outputs: None.
	SetBoxX(x int)
	
	// SetBoxY thiết lập tọa độ dọc của Layout node (thường do cha tính toán).
	// Purpose: Sets the local Y coordinate of the layout node.
	// Inputs: y int - The Y position in pixels.
	// Outputs: None.
	SetBoxY(y int)

	// BoxH trả về chiều cao hiện tại của Layout node (pixel).
	// Purpose: Retrieves the current height of the layout node.
	// Inputs: None.
	// Outputs: int - The height in pixels.
	BoxH() int
	
	// BoxW trả về chiều rộng hiện tại của Layout node (pixel).
	// Purpose: Retrieves the current width of the layout node.
	// Inputs: None.
	// Outputs: int - The width in pixels.
	BoxW() int
	
	// BoxX trả về tọa độ ngang hiện tại của Layout node (pixel).
	// Purpose: Retrieves the current local X coordinate.
	// Inputs: None.
	// Outputs: int - The X position in pixels.
	BoxX() int
	
	// BoxY trả về tọa độ dọc hiện tại của Layout node (pixel).
	// Purpose: Retrieves the current local Y coordinate.
	// Inputs: None.
	// Outputs: int - The Y position in pixels.
	BoxY() int

	// AddChildren thêm một ILayout con vào node này.
	// Purpose: Adds a child layout node to this parent node.
	// Inputs: children ILayout - The child node to add.
	// Outputs: None.
	AddChildren(children ILayout)
	
	// RemoveChildren xóa một ILayout con khỏi node này.
	// Purpose: Removes a specific child layout node from this parent.
	// Inputs: children ILayout - The child node to remove.
	// Outputs: None.
	RemoveChildren(children ILayout)
	
	// Childrens trả về danh sách tất cả các Layout con.
	// Purpose: Retrieves a list of all child layout nodes.
	// Inputs: None.
	// Outputs: []ILayout - Slice containing all children.
	Childrens() []ILayout
	
	// GetChildrenByName tìm Layout con theo tên. Trả về nil nếu không tìm thấy.
	// Purpose: Finds and returns a child layout node by its assigned name.
	// Inputs: name string - The name of the child node.
	// Outputs: ILayout - The child node, or nil if not found.
	GetChildrenByName(name string) ILayout

	// Name trả về tên định danh của Layout node.
	// Purpose: Retrieves the assigned name of the layout node.
	// Inputs: None.
	// Outputs: string - The node's name.
	Name() string
	
	// SetName thiết lập tên định danh của Layout node.
	// Purpose: Assigns a name to the layout node for identification.
	// Inputs: name string - The new name.
	// Outputs: None.
	SetName(name string)

	// IsVisible trả về trạng thái hiển thị của Layout node.
	// Purpose: Checks if the layout node and its children are visible.
	// Inputs: None.
	// Outputs: bool - True if visible, false otherwise.
	IsVisible() bool
	
	// SetIsVisible bật/tắt hiển thị Layout node (ẩn node cũng ẩn toàn bộ cây con).
	// Purpose: Toggles visibility of this node and its entire subtree.
	// Inputs: isVisible bool - True to show, false to hide.
	// Outputs: None.
	SetIsVisible(isVisible bool)

	// SetLayoutConfig thiết lập config bitfield kết hợp Direction | Align | Justify.
	// Ví dụ: DirColumn | AlignCenter | JustifyStart
	// Purpose: Sets the layout configuration using a bitfield of Dir, Align, and Justify masks.
	// Inputs: config int - The bitfield configuration value.
	// Outputs: None.
	SetLayoutConfig(config int)
	
	// LayoutConfig trả về config bitfield hiện tại.
	// Purpose: Retrieves the current layout configuration bitfield.
	// Inputs: None.
	// Outputs: int - The bitfield configuration value.
	LayoutConfig() int

	// SetGap thiết lập khoảng cách giữa các Layout con (pixel).
	// Purpose: Sets the gap spacing between child nodes in a flex-like container.
	// Inputs: gap int - The gap size in pixels.
	// Outputs: None.
	SetGap(gap int)
	
	// Gap trả về khoảng cách giữa các Layout con (pixel).
	// Purpose: Retrieves the gap spacing between child nodes.
	// Inputs: None.
	// Outputs: int - The gap size in pixels.
	Gap() int

	// ComputeLayout tính toán lại vị trí và kích thước của node này và toàn bộ cây con.
	// parentX, parentY là tọa độ góc trên-trái của node cha.
	// parentW, parentH là kích thước khả dụng từ node cha.
	// Purpose: Recursively calculates positions and dimensions for this node and its children based on flex-like rules.
	// Inputs:
	//   - parentX int: Parent's absolute X coordinate.
	//   - parentY int: Parent's absolute Y coordinate.
	//   - parentW int: Parent's available width.
	//   - parentH int: Parent's available height.
	// Outputs: None.
	ComputeLayout(parentX, parentY int, parentW, parentH int)
}
