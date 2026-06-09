# AutoWorld — nSprite Module Reference

**Ý Nghĩa**:
    Module `nsprite` chịu trách nhiệm quản lý tài nguyên hình ảnh (Sprite) và thực hiện việc render (vẽ) các thực thể lên màn hình dựa trên dữ liệu từ ECS (donburi).

**Chức năng chính**:
    1. **Register.go**:
        Cửa ngõ duy nhất import từ `domain`. Định nghĩa các Type Alias cho dữ liệu Component (`PositionData`, `SpriteData`, `BoxData`) và liên kết trực tiếp các `ComponentType` dùng chung từ `enginetype` (`Position`, `Sprite`, `Box`) để đảm bảo tính nhất quán trên toàn engine.
    2. **Sprite (Sprite.go)**:
        Lớp lưu trữ tập hợp các khung hình (ebiten.Image) và thông số kỹ thuật của một tập tin Sprite. Được thiết kế theo mô hình Lightweight để tiết kiệm bộ nhớ.
    3. **Graphic (Graphic.go)**:
        Module nội bộ chịu trách nhiệm tính toán ma trận biến đổi (DrawOptions). Nó kết hợp dữ liệu vị trí và sprite để tạo ra các tham số sẵn sàng cho việc vẽ.
    4. **DrawSystem (DrawSystem.go)**:
        Hệ thống render chính. Nó sử dụng Query của donburi để lọc ra các thực thể cần vẽ và thực hiện lệnh `DrawImage` lên canvas của Ebitengine. Hỗ trợ cả sprite thông thường và sprite 9-slice.

**Tính năng 9-slice**:
    Khi thiết lập `IsNineSlice = true` và truyền thông số `NineSlice {Top, Right, Bottom, Left}` vào `SpriteData`, `Graphic.go` sẽ chia nhỏ hình ảnh thành 9 mảnh (patch) và điều chỉnh Scale theo `ScaleX` và `ScaleY` để tạo thành một đối tượng không bị biến dạng ở 4 góc. Kích thước đích được tính bằng `Kích_thước_gốc * Scale`. Điều này cực kỳ hữu ích để vẽ UI, hộp thoại (Box, Dialog).

**Quy trình hoạt động**:
    1. **Render Flow**:
       `Scene` gọi `DrawSystem.SetScreen(screen)` sau đó gọi `DrawSystem.Draw(world)`.
    2. **Query & Filter**:
       `DrawSystem` lọc các thực thể có đủ cả `Sprite` và `Position` component.
    3. **Build Options**:
       Với mỗi thực thể, `Graphic.BuildDrawOptions` tính toán tọa độ, độ xoay, tỉ lệ và màu sắc.
    4. **Rendering**:
       Hình ảnh được vẽ lên màn hình thông qua thư viện Ebitengine.

**Tính năng Z-Order (Depth Sorting)**:
    - Để giải quyết vấn đề thứ tự vẽ chồng lấn giữa các Sprite, N-Engine sử dụng thuộc tính `ZOrder` (kiểu `int`) lưu trong `SpriteData`.
    - `DrawSystem` duy trì một mảng nội bộ các thực thể hợp lệ (`sortedSprites`). Mảng này được tự động phân giải (cập nhật nếu có thêm bớt thực thể).
    - Để tối ưu hiệu năng theo phương châm "chỉ sort khi có thay đổi", `DrawSystem` lắng nghe cờ `IsZOrderDirty` (được bật tự động khi gọi `SetZOrder()`). Nếu cờ này bật hoặc có thực thể mới thêm vào, hệ thống mới tiến hành chạy thuật toán `sort.SliceStable` trên mảng. Do mảng đã gần như được sắp xếp sẵn, `sort.SliceStable` thực thi cực kỳ nhanh (thời gian chạy tiệm cận `O(N)`).
    - Quá trình vẽ cuối cùng diễn ra bằng cách duyệt mảng `sortedSprites` sau khi đã được sắp xếp và loại trừ AABB Culling.