# AutoWorld — nLayout Module Reference

**Ý Nghĩa**:
    Module `nlayout` cung cấp một hệ thống Auto-Layout tương tự như Flexbox đơn giản hóa, dùng để tự động thiết lập và căn chỉnh vị trí của các UI hoặc Game Object trên màn hình một cách linh hoạt mà không cần đặt tọa độ cứng (hard-code).

**Chức năng chính**:
    1. **Register.go**:
        File duy nhất kết nối với `domain`, làm nhiệm vụ Type Alias (ví dụ: `type ILayout = domain.ILayout`) nhằm nội bộ hóa (localize) các tham chiếu và tuân thủ chặt chẽ kiến trúc của Engine.
    2. **Div (layoutDiv.go)**:
        Đóng vai trò là vùng chứa (Container). Một `Div` có thể chứa các `Div` hoặc `A` khác. Nó sở hữu thuật toán tính toán và phân bổ vị trí tự động cho các thẻ con dựa trên các cờ cấu hình (`layoutConfig`).
    3. **A (layoutA.go)**:
        Đóng vai trò là thẻ Adapter (cầu nối). Nó tuân thủ `ILayout` nhưng mục đích chính là bọc (wrap) một `domain.IObject` của game. Bất cứ khi nào thẻ `A` được layout đặt vào vị trí mới, nó sẽ tự động đồng bộ tọa độ đó xuống cho Game Object.
    4. **Bitwise Configuration**:
        Sử dụng kỹ thuật Bitwise flag (thông qua `layoutConfig`) để thiết lập layout rất gọn nhẹ. Hỗ trợ:
        - *Hướng (Direction)*: `DirRow`, `DirColumn`.
        - *Căn chỉnh trục phụ (AlignItems)*: `AlignStart`, `AlignCenter`, `AlignEnd`.
        - *Dàn trải trục chính (JustifyContent)*: `JustifyStart`, `JustifyCenter`, `JustifyEnd`.

**Quy trình hoạt động**:
    1. **Khởi tạo và Xây dựng Cây Layout**:
        Tạo đối tượng `Div` gốc. Cấu hình layout thông qua `SetLayoutConfig` (ví dụ: `div.SetLayoutConfig(layout.DirColumn | layout.AlignCenter | layout.JustifyCenter)`).
    2. **Thêm Đối Tượng Game**:
        Gói các `domain.IObject` vào thẻ `A` tương ứng, và thêm (`AddChildren`) các thẻ `A` (hoặc các `Div` lồng nhau) vào `Div` gốc.
    3. **Tính Toán Vị Trí (Compute)**:
        Gọi hàm `ComputeLayout(x, y, w, h)` từ `Div` gốc. Hệ thống sẽ đệ quy quét qua toàn bộ cấu trúc cây.
    4. **Áp dụng Tọa Độ**:
        Trong bước đệ quy, hàm `ComputeLayout` sẽ tính toán và truyền tọa độ chính xác xuống thẻ `A`. Thẻ `A` sau đó sẽ gọi `SetX`, `SetY`, `SetBoxX`, `SetBoxY` để cập nhật trạng thái vật lý/hình ảnh cho đối tượng Game.
