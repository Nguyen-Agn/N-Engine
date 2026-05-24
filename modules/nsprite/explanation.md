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
        Hệ thống render chính. Nó sử dụng Query của donburi để lọc ra các thực thể cần vẽ và thực hiện lệnh `DrawImage` lên canvas của Ebitengine.

**Quy trình hoạt động**:
    1. **Render Flow**:
       `Scene` gọi `DrawSystem.SetScreen(screen)` sau đó gọi `DrawSystem.Draw(world)`.
    2. **Query & Filter**:
       `DrawSystem` lọc các thực thể có đủ cả `Sprite` và `Position` component.
    3. **Build Options**:
       Với mỗi thực thể, `Graphic.BuildDrawOptions` tính toán tọa độ, độ xoay, tỉ lệ và màu sắc.
    4. **Rendering**:
       Hình ảnh được vẽ lên màn hình thông qua thư viện Ebitengine.