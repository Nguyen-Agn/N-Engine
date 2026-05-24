# Module enginetype — Engine Component Type Registry

## Mục tiêu

Module `enginetype` là **trung tâm đăng ký Component Type** của toàn bộ Engine.
Nó giải quyết vấn đề thứ tự khởi tạo (init order) của Go: thay vì để mỗi module tự tạo `donburi.ComponentType`, tất cả Component Type được khai báo tại đây để đảm bảo không có race condition hay nil pointer trong giai đoạn startup.

---

## Chức năng chính

### 1. `Register.go` — Khai báo Component Types tập trung

Nơi duy nhất định nghĩa các biến Component Type toàn cục cho Engine:

| Biến | Token | Mô tả |
|------|-------|-------|
| `Position` | `"pos"` | Tọa độ X, Y của đối tượng |
| `Sprite` | `"spr"` | Hình ảnh, animation |
| `Box` | `"box"` | Hitbox, va chạm |
| `Audio` | `"aud"` | Âm thanh |
| `Infor` | `"inf"` | ID và tên đối tượng (bắt buộc) |
| `Direction` | `"dir"` | Góc hướng di chuyển |
| `Input` | `"inp"` | Input callback (lắng nghe phím) |

### 2. `ComponentsType.go` — Factory và Registry

- **`NewComponentType[T](token)`**: Tạo một `donburi.ComponentType[T]` mới và tự động đăng ký vào registry nếu token khác `""`. Cho phép Dev tạo custom component riêng và sử dụng trong chuỗi `componentCode`.
- **`GetComponentType(token)`**: Tra cứu Component Type theo token string.
- **`RegisterComponentInitializer(token, fn)`**: Đăng ký hàm khởi tạo mặc định cho một component.
- **`InitializeComponent(token, entry)`**: Gọi hàm khởi tạo tương ứng khi Object được tạo.

### 3. `ComponentHelper.go` — Helpers Type-Safe

- **`SetComponent[T](obj, comp, value)`**: Gán giá trị cho custom component trên object.
- **`GetComponent[T](obj, comp)`**: Lấy con trỏ đến data của custom component. Trả về `nil` nếu object không có component này.
- **`AddComponentType[T](obj, comp)`**: Thêm component type vào ECS entry của object sau khi object đã được tạo.

---

## Quy tắc quan trọng cho Agent

1. **Mọi Component Type mới** phải được khai báo tại `Register.go` với `NewComponentType[T](token)`.
2. **Mọi giá trị khởi tạo mặc định** phải được đăng ký trong `init()` của module `components` tương ứng, không phải tại đây.
3. Không import `core`, `napi` hay bất kỳ module cấp cao hơn — enginetype chỉ phụ thuộc vào `domain` và `donburi`.
