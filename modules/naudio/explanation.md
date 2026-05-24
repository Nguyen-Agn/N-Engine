# AutoWorld — nAudio Module Reference

**Ý Nghĩa**:
    Module `naudio` chịu trách nhiệm quản lý và phát âm thanh trong game, tách biệt hoàn toàn logic kích hoạt âm thanh khỏi hạ tầng phát âm thanh của Ebitengine.

**Chức năng chính**:
    1. **Register.go**:
        Thực hiện kết nối duy nhất với `domain`. Định nghĩa Type Alias cho `AudioData`, `IAudioLW` và liên kết trực tiếp `ComponentType` cho Audio từ `enginetype` (`Audio`) để đảm bảo tính nhất quán trên toàn engine.
    2. **Audio (Audio.go)**:
        Lớp bọc (wrapper) cho một tệp âm thanh cụ thể, triển khai các lệnh phát, dừng, tạm dừng cơ bản.
    3. **AudioSystem (AudioSystem.go)**:
        Hệ thống điều phối âm thanh trong ECS. Nó lắng nghe các thay đổi trạng thái (cờ hiệu) trong `AudioData` để thực hiện các hành động tương ứng.

**Quy trình hoạt động**:
    1. **Update Flow**:
       `Scene` gọi `AudioSystem.Update(world)` trong mỗi khung hình.
    2. **Check Flags**:
       Hệ thống kiểm tra các cờ `ShouldPlay` hoặc `ShouldStop` trong `AudioData` của từng thực thể.
    3. **Execute Command**:
       Nếu cờ được bật, `AudioSystem` gọi phương thức tương ứng trên đối tượng `Audio` và sau đó reset cờ về `false`.
    4. **Output**:
       Âm thanh được phát ra loa thông qua thư viện audio của Ebitengine.