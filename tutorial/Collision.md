# Va chạm / Collision & Hitbox (ncom.Col & ncom.Box)

> **Vision**: Định nghĩa vùng va chạm và bắt sự kiện qua tag nhanh chóng.
> **Vision**: Define collision boxes and catch events via tags quickly.

---

## 1. Giải thích / Explanation

Engine sử dụng 2 component đi kèm với nhau để xử lý va chạm:
The Engine uses 2 interconnected components for collision:
- `ncom.Box`: Xác định hình khối, kích thước (Hình chữ nhật, Hình tròn) và các thuộc tính vật lý (Solid, Collidable).
- `ncom.Col`: Định nghĩa các "Tag" cho Object và lắng nghe sự kiện va chạm với các Tag khác.

Lưu ý: Khi nhúng `ncom.Col` (token `"col"`), engine sẽ tự động thêm token `"box"`.
Note: Embedding `ncom.Col` (token `"col"`) automatically adds the `"box"` token.

---

## 2. Ví dụ / Code Example

```go
package objects

import (
	"fmt"

	"autoworld/modules/napi"
	"autoworld/modules/napi/ncom"
)

type Bullet struct {
	ncom.Object
	ncom.Pos
	ncom.Spr
	ncom.Col // Kéo theo ncom.Box
}

func NewBullet(x, y float32) *Bullet {
	b := &Bullet{}
	napi.Obj.NewObject(b, "Bullet", "pos spr col sce-main")
	
	b.SetX(x)
	b.SetY(y)
	return b
}

func (b *Bullet) Create() {
	// --- Cấu hình Hitbox (ncom.Box) ---
	b.SetBoxW(10) // Rộng 10px / Width 10px
	b.SetBoxH(10) // Cao 10px / Height 10px
	
	// Dịch chuyển tâm hitbox (Tùy chọn) / Offset hitbox (Optional)
	// b.SetBoxOffsetX(-5)
	// b.SetBoxOffsetY(-5)

	// Đặt hình dáng (Mặc định là napi.BSRectangle) / Set shape (Default is Rectangle)
	// Có thể đổi sang: b.SetShape(napi.BSCircle)

	// --- Cấu hình Va chạm (ncom.Col) ---
	// Gắn tag cho viên đạn này là "projectile"
	// Tag this bullet as "projectile"
	b.AddTag("projectile")

	// Lắng nghe sự kiện: Nếu chạm vào Object có tag "enemy" -> Gọi hàm onHitEnemy
	// Listen event: If hitting an Object with "enemy" tag -> Call onHitEnemy
	b.OnCollisionTag("enemy", b.onHitEnemy)
}

// Hàm callback nhận vào Object mà nó vừa va chạm
// Callback receives the Object it collided with
func (b *Bullet) onHitEnemy(other ncom.Object) {
	fmt.Println("Viên đạn đã trúng kẻ thù! / Bullet hit enemy!")
	
	// Tiêu diệt viên đạn
	b.Destroy()

	// Bạn có thể ép kiểu 'other' để thao tác (VD: enemy.TakeDamage())
}
```

## 3. Các thuộc tính hữu ích / Useful Properties
- `SetSolid(true)`: Ngăn không cho 2 Object có hitbox đè lên nhau (Dùng cho Tường, Đất). / Prevent hitboxes from overlapping (Used for Walls, Ground).
- `SetCollidable(false)`: Vô hiệu hóa va chạm tạm thời mà không cần hủy Object. / Temporarily disable collision without destroying the Object.
