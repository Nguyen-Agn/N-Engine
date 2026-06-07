# AutoWorld — nSave Module Reference

**Ý Nghĩa**:
    Module `nsave` chịu trách nhiệm quản lý hệ thống lưu và tải trò chơi (Save/Load System). Module này hiện thực hóa interface `domain.ISaveManager`, thu thập dữ liệu từ các đối tượng và biến toàn cục, tuần tự hóa thành định dạng JSON và lưu xuống ổ đĩa, cũng như quy trình ngược lại khi tải game.

**Chức năng chính**:
    1. **FileAdapter (FileAdapter.go)**:
        Xử lý I/O tập tin. Đọc, ghi, kiểm tra sự tồn tại, xóa và liệt kê các slot save (các file JSON trong thư mục lưu trữ). Tách biệt logic đọc/ghi đĩa với logic thu thập dữ liệu.
    2. **Collector (Collector.go)**:
        Thu thập dữ liệu từ Scene hiện tại (thông qua `GetObjects`) và các biến Global (thông qua `DumpVariables`). Chỉ thu thập các Object có `SaveTag` hợp lệ và trả về khác `nil` từ hàm `OnSave()`. **Bỏ qua Object có `IsDead() == true`** (đã bị queue xóa, không được lưu vào file).
    3. **Applier (Applier.go)**:
        Phân phối dữ liệu từ `SaveSnapshot` trở lại cho các Object hiện có (qua `OnLoad()`) và biến Global (qua `RestoreVariables()`). **Bỏ qua Object có `IsDead() == true`** — ngăn các object đã bị Remove "sống lại" khi load.
    4. **SaveManager (SaveManager.go)**:
        Trung tâm điều phối. Lấy dữ liệu từ Collector, chuyển cho FileAdapter ghi xuống, và lấy dữ liệu từ FileAdapter, chuyển cho Applier để nạp lại trạng thái.
        Implement `domain.ISaveManager`.

**Quy tắc**:
    - **Không phụ thuộc ngược**: `nsave` chỉ import `domain`, không import `core`, `nsys` hay `napi`.
    - **Stable ID**: Phân biệt các object khi save/load dựa trên `SaveTag()`. Nếu object không có SaveTag, nó sẽ bị bỏ qua.
    - **Dead Guard**: Object có `IsDead() == true` bị skip ở cả Collector lẫn Applier — đây là cơ chế phối hợp với `Map.RemoveObject()` để tránh ghost resurrection.
