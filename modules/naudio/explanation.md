# AutoWorld — nAudio Module Reference

**Ý Nghĩa**:
    Module `naudio` chịu trách nhiệm quản lý và phát âm thanh trong game, tách biệt hoàn toàn logic kích hoạt âm thanh khỏi hạ tầng phát âm thanh của Ebitengine.

**Chức năng chính**:
    1. **Register.go**:
        Thực hiện kết nối duy nhất với `domain`. Định nghĩa Type Alias cho `AudioData`, `IAudioLW` và liên kết trực tiếp `ComponentType` cho Audio từ `enginetype` (`Audio`) để đảm bảo tính nhất quán trên toàn engine.
    2. **Audio (Audio.go)**:
        Lớp bọc (wrapper) cho một tệp âm thanh cụ thể, triển khai các lệnh phát, dừng, tạm dừng, tiếp tục và hỗ trợ tự động lặp (looping). Nó cũng quản lý bộ nhớ để tránh rò rỉ (chỉ tạo 1 instance Player).
    3. **AudioSystem (AudioSystem.go)**:
        Hệ thống điều phối âm thanh trong ECS. Nó lắng nghe các thay đổi trạng thái (cờ hiệu) trong `AudioData` để thực hiện các hành động tương ứng (Play, Pause, Resume, Stop) và đồng bộ Volume real-time.

**Quy trình hoạt động**:
    1. **Update Flow**:
       `Scene` gọi `AudioSystem.Update(world)` trong mỗi khung hình.
    2. **Check Flags**:
       Hệ thống kiểm tra các cờ điều khiển (`ShouldPlay`, `ShouldStop`, `ShouldPause`, `ShouldResume`) và trạng thái lặp (`IsLooping`) trong `AudioData` của từng thực thể.
    3. **Execute Command**:
       Nếu cờ điều khiển được bật, `AudioSystem` gọi phương thức tương ứng trên đối tượng `AudioLW` và sau đó reset cờ về `false`.
       Hệ thống cũng tự động thiết lập lại Volume theo real-time để cho phép hiệu ứng mờ dần (fade out/fade in).
       Nếu cờ `IsLooping` được bật và âm thanh đã phát xong một cách tự nhiên (không bị Stop cưỡng bức), nó sẽ tự động phát lại.
    4. **Output**:
       Âm thanh được phát ra loa thông qua thư viện audio của Ebitengine.