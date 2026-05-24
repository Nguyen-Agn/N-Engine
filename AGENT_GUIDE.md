# AutoWorld Engine — Agent Guide

Tài liệu này được thiết kế đặc biệt dành cho các AI Agent (hoặc developer mới) để nắm bắt nhanh toàn bộ kiến trúc, triết lý thiết kế và các quy tắc bắt buộc của engine AutoWorld.

**Vui lòng đọc kỹ tài liệu này trước khi thực hiện bất kỳ thay đổi nào trong source code.**

---

## 1. Tổng quan Kiến trúc

AutoWorld là một 2D Game Engine xây dựng trên nền tảng **Ebitengine** (render/audio/input) và **Donburi** (ECS - Entity Component System). Kiến trúc được chia thành các layer kiểm soát chặt chẽ dependency:

### 1.1. Các Layer chính (Từ thấp đến cao)

1. **`domain`**: Chứa toàn bộ **Interfaces** và **Data Structs** cơ bản (PositionData, SpriteData...). Tuyệt đối không chứa logic. Các module khác giao tiếp với nhau thông qua interface ở đây.
2. **`modules/enginetype`**: Định nghĩa và đăng ký các Component Type toàn cục cho hệ thống ECS (Donburi).
3. **`modules/nobject`, `modules/nsprite`, `modules/naudio`**: Các module thực thi hệ thống con (LogicSystem, DrawSystem, AudioSystem).
4. **`modules/core`**: "Trái tim" của engine. Chứa `Engine`, `SceneManager`, `Map`, `Camera` và loop chính kết nối với Ebitengine. Core sử dụng các system từ nsprite/naudio/nobject để vận hành game.
5. **`modules/napi`**: **API Layer cao nhất**. Che giấu hoàn toàn sự phức tạp của ECS, Ebiten, và Core. Game developer **chỉ import `napi`** để viết game logic. `napi` không bao giờ được import ngược lại bởi các module core.

---

## 2. Các Khái niệm Cốt lõi (Core Concepts)

### 2.1. Kiến trúc Scene - Map - Camera
- **Scene**: Đại diện cho 1 màn chơi. Chứa Physical Map, GUI Map, và Camera.
- **Physical Map**: Chứa ECS World (`donburi.World`), Object list, và các System logic/audio/input. Tọa độ tính theo **map space** (tuyệt đối).
- **GUI Map**: Map riêng dành cho HUD/UI. Render ở **screen space** (không bị ảnh hưởng bởi camera).
- **Camera**: Quản lý viewport, tự động bám theo (follow) mục tiêu (target) có giới hạn theo biên bản đồ (bounds). DrawSystem sử dụng tọa độ Camera để trừ offset khi vẽ Physical Map.

### 2.2. Hệ thống ECS Hybrid (ECS + OOP Mix)
Để thân thiện với developer quen OOP, AutoWorld kết hợp ECS bên dưới với OOP bề mặt qua `napi`:
- Dữ liệu thực sự sống trong **donburi ECS** (Entity-Component).
- Code game định nghĩa struct nhúng các **Component Mixin** (`napi.IPosition`, `napi.ISprite`).
- Hàm `napi.NewObject` tự động tạo ECS Entity và sử dụng reflection (`napi.bind`) để nối component data từ ECS vào struct của user.

Ví dụ tạo Custom Object:
```go
type Player struct {
	napi.IObject
	napi.IPosition
	napi.ISprite
}

func NewPlayer(x, y float32) *Player {
	p := &Player{}
	// "pos spr" sẽ yêu cầu engine tạo PositionData và SpriteData trong ECS
	napi.NewObject(p, "player", "pos spr") 
	p.SetX(x)
	p.SetY(y)
	napi.Register(p, false) // Đăng ký vào Physical Map hiện tại
	return p
}
```

---

## 3. Các Quy tắc Bắt buộc (CRITICAL RULES)

Khi Agent thực hiện task, **phải tuân thủ tuyệt đối** các quy tắc sau (do user định nghĩa):

1. **Interface First**: Các class/module không nói chuyện trực tiếp với nhau mà phải thông qua Interface (trong `domain`). Luôn giả định interface cung cấp dữ liệu chính xác như comment mô tả.
2. **Comment rõ ràng**: Mỗi phương thức trong Interface hoặc struct công khai phải có comment giải thích rõ ràng.
3. **Open/Close Principle**: Hạn chế chỉnh sửa các module không liên quan trực tiếp đến task. Thiết kế sao cho dễ mở rộng thay vì sửa đổi code cũ.
4. **Đọc `explanation.md`**: Khi vào một module/folder mới, **luôn tìm và đọc file `explanation.md` đầu tiên** thay vì đọc source code, để hiểu triết lý của riêng module đó.
5. **Cập nhật `explanation.md`**: Nếu cấu trúc của một module thay đổi, bạn **bắt buộc** phải cập nhật lại `explanation.md` của module đó cho đồng bộ.
6. **KHÔNG dùng lệnh terminal (cmd/powershell)**: Tuyệt đối không dùng bash/cmd để thay đổi file (không dùng `rm`, `mv`, `cat`...). Hãy sử dụng các tool API (như `write_to_file`, `replace_file_content`, `multi_replace_file_content`) để quản lý mã nguồn.
7. **Bảo toàn comment hiện có**: Không xóa các comment hoặc docstring không liên quan đến phần thay đổi của bạn trừ khi được user chỉ định.

---

## 4. Troubleshooting / Cần lưu ý

- **Type-assert bị hạn chế**: Tránh dùng type-assert để bypass interface. Hãy thêm method cần thiết vào interface (nếu hợp lý) và tuân thủ thiết kế.
- **DrawSystem Culling**: Entity ngoài camera viewport sẽ bị culling (không vẽ), nhưng LogicSystem vẫn update bình thường.
- **Component Token**: Trong `napi.NewObject`, chuỗi component thường dùng là:
  - `pos`: PositionData (X,Y)
  - `spr`: SpriteData (Image, Scale, Rotate)
  - `box`: BoxData (Hitbox, Collision)
  - `aud`: AudioData
  - `inf`: InforData (luôn tự động được add)
  - `dir`: DirectionData
  - `bg`: BackgroundData
  - `tile`: TilemapData
