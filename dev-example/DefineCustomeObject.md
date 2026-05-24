# How to Define a Custom Object

Hướng dẫn dành cho game developer khi cần tạo Object với type riêng (Player, Enemy, Boss...) và có logic tùy biến.

---

## Pattern: Embed `napi.IObject` & Component Mixins

Kể từ phiên bản mới, `domain.IObject` chỉ còn các phương thức quản lý logic vòng đời cốt lõi (`StepUpdate`, `Destroy`). 
Để sử dụng các tính năng như thay đổi tọa độ (`SetX`, `SetY`), vẽ hình (`SetSprite`), phát nhạc (`Play`), bạn cần nhúng các **Component Mixins** tương ứng từ package `components`.

### Cách triển khai một Custom Object:

```go
package your_game

import (
    "github.com/user/autoworld/modules/napi"
    "github.com/user/autoworld/modules/components"
)

// 1. Định nghĩa struct — embed napi.IObject và các Component Mixins cần dùng
type CustomObject struct {
    napi.IObject                  // Vòng đời cốt lõi (StepUpdate, Destroy)
    components.PositionComponent  // Khả năng thay đổi vị trí (X, Y, SetX, SetY)
    components.SpriteComponent    // Khả năng hiển thị hình ảnh (SetSprite, Scale...)
    components.AudioComponent     // Khả năng phát âm thanh (Play, StopAudio...)
    
    health     int
    customVar2 string
}

// 2. (Tùy chọn) Định nghĩa interface riêng cho CustomObject của bạn
type ICustomObject interface {
    napi.IObject
    components.PositionComponent
    Health() int
    TakeDamage(amount int)
}

// 3. Constructor — Tạo ECS entity, liên kết các component mixin, rồi đăng ký
func NewCustomObject(health int, customVar2 string) *CustomObject {
    // NewBaseObject tạo ECS entity với các component được chỉ định qua chuỗi (chưa đăng ký scene)
    base := napi.NewBaseObject("player", "pos spr aud")
    
    obj := &CustomObject{
        IObject:            base,
        PositionComponent:  components.PositionComponent{IObject: base},
        SpriteComponent:    components.SpriteComponent{IObject: base},
        AudioComponent:     components.AudioComponent{IObject: base},
        health:             health,
        customVar2:         customVar2,
    }

    // Thiết lập giá trị ban đầu qua các component mixin đã nhúng
    obj.SetX(100)
    obj.SetY(200)
    obj.SetSprite("idle", napi.GetSprite("hero_idle"))

    // Register đưa object vào scene để bắt đầu nhận sự kiện StepUpdate
    // global=false → current scene | global=true → Global Hidden Scene (persistent)
    napi.Register(obj, false)

    return obj
}

// 4. Override StepUpdate — Chạy mỗi frame
func (c *CustomObject) StepUpdate() {
    // Thực hiện game logic
    c.SetX(c.X() + 2) // Di chuyển sang phải
}

// 5. Override Destroy — Dọn dẹp khi object bị xóa
func (c *CustomObject) Destroy() {
    // Cleanup logic nếu cần
}

// 6. Các phương thức nghiệp vụ của riêng object
func (c *CustomObject) Health() int { return c.health }

func (c *CustomObject) TakeDamage(amount int) {
    c.health -= amount
    if c.health <= 0 {
        c.Destroy()
    }
}
```

---

## Custom Component (Định nghĩa dữ liệu riêng)

Nếu game của bạn cần lưu thêm các thuộc tính tùy biến vào ECS để hệ thống (system) khác đọc được, hãy sử dụng `napi.NewComponentType`:

```go
package your_game

import (
    "github.com/user/autoworld/modules/napi"
)

// 1. Định nghĩa data struct
type StatsData struct {
    MaxHealth int
    Mana      int
    Strength  int
}

// 2. Khai báo Component Type ở package level kèm theo tên token là "sta"
var StatsComp = napi.NewComponentType[StatsData]("sta")

// 3. Trong Constructor của Hero:
func NewHero() *Hero {
    obj := &Hero{
        // Chỉ rõ token "sta" trong chuỗi componentCode để ECS tự động đính kèm component này
        IObject: napi.NewBaseObject("hero", "pos spr sta"),
    }
    
    // Gán dữ liệu mặc định cho component
    napi.SetComponent(obj, StatsComp, StatsData{
        MaxHealth: 200,
        Mana:      100,
        Strength:  15,
    })
    
    napi.Register(obj, false)
    return obj
}

// 4. Đọc/Ghi component ở bất kỳ đâu
func (h *Hero) StepUpdate() {
    stats := napi.GetComponent(h, StatsComp) // Trả về *StatsData
    if stats != nil {
        stats.Mana += 1 // Regenerate mana mỗi frame
    }
}
```

---

## Bảng chọn lựa tính năng

| Nhu cầu | Giải pháp |
|---------|-----------|
| Object đơn giản, chỉ chứa ảnh/vị trí tĩnh, không có logic chạy mỗi frame | `napi.NewObject("bullet", "pos spr")` |
| Object có logic, có phương thức `StepUpdate()` riêng | Embed `napi.IObject` + embed các mixin (`PositionComponent`, `SpriteComponent`,...) + `napi.Register` |
| Muốn lưu thêm dữ liệu ECS tùy biến để truy xuất động | `napi.NewComponentType[T]("token")` + `SetComponent/GetComponent` |
| Object tồn tại xuyên suốt cảnh chơi (Persistent) | Gọi `napi.Register(obj, true)` hoặc có token `glo` |