# AutoWorld — nMath Module Reference

**Ý Nghĩa**:
    Module `nmath` đóng vai trò là thư viện tiện ích (Helper) cung cấp các hàm toán học đặc thù cho việc lập trình game 2D. 
    Các hàm trong đây chủ yếu sử dụng `float32` để tương thích hoàn toàn với hệ thống tọa độ và các thông số của N-Engine, giúp lập trình viên không cần mất công ép kiểu liên tục từ `float64` (mặc định của thư viện `math` trong Go).

**Chức năng chính**:
    1. Cung cấp các hàm kiểm soát giá trị: `Clamp`, `Lerp`, `Approach`, `Wrap`, `Sign`.
    2. Các hàm tính toán tọa độ vector và góc (Degree): `LengthDirX`, `LengthDirY`, `Distance`, `Angle`, `AngleDif`, `AngleAdd`.
    3. Trình tạo ngẫu nhiên (Random) chuyên dụng cho cơ chế game: `RandomRangeDouble`, `RandomRangeInt`, `Choose`, `SetSeed`.
    4. Căn chỉnh hệ trục tọa độ lưới: `Snap`.

**Quy trình hoạt động**:
    Module này không chứa trạng thái nội tại (stateless) ngoại trừ hàm Random (được quản lý seed). Nó cung cấp bộ khung hàm rỗng dạng `MathHelper` và được khởi tạo thành một biến toàn cục `napi.Math`.
    
    Từ phía người dùng (Dev), chỉ cần gọi thông qua `napi.Math.[TênHàm]` ở bất kỳ đâu trong game logic.

*(Lưu ý: Các thuật toán lõi bên trong đang được thảo luận và tối ưu, hiện tại đây là cấu trúc rỗng.)*
