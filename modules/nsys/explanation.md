# Module nGlobal

## Mục tiêu
Quản lý các tài nguyên toàn cục (Global resources) của game bao gồm:
- **Sprites**: Hình ảnh, hoạt ảnh.
- **Audios**: Âm thanh, nhạc nền.
- **Global Objects**: Các đối tượng tồn tại xuyên suốt trò chơi.
- **Variables**: Các biến số cấu hình/trạng thái thay đổi được.
- **Constants**: Các hằng số cài đặt ít bị thay đổi.

## Kiến trúc
- Tuân thủ nguyên tắc **Interface First**. Giao tiếp thông qua `domain.IGlobal`.
- Sử dụng **Singleton Pattern** để đảm bảo chỉ có duy nhất một kho lưu trữ tài nguyên toàn cục trong vòng đời của Game. Lấy instance thông qua `nglobal.GetInstance()`.
- Sử dụng `sync.RWMutex` để đảm bảo luồng an toàn (Thread-safe) trong các thao tác thêm/sửa/đọc tài nguyên từ nhiều goroutine.

## Lưu ý
- Các hàm `Add` và `Update` hiện tại đều thực hiện thao tác ghi đè (overwrite) nếu key đã tồn tại.
- Để giữ cho mọi thứ đơn giản và an toàn, module này không hỗ trợ chức năng `Remove` để tránh trường hợp tài nguyên đang được sử dụng ở nơi khác bị xóa gây lỗi.
