# Hướng dẫn sử dụng các thư viện trong AutoWorld

Tài liệu này tóm tắt các Interface và cách sử dụng của hai thư viện chính: **Ebitengine** (Game Engine) và **Donburi** (ECS).

## 1. Ebitengine (v2)

Ebitengine hoạt động dựa trên một Interface trung tâm mà bạn phải triển khai để chạy game.

### Interface `ebiten.Game`

Mọi ứng dụng Ebitengine đều phải thực thi interface này:

```go
type Game interface {
    // Update được gọi mỗi khung hình (mặc định 60 lần/giây).
    // Đây là nơi xử lý logic: input, di chuyển, va chạm.
    Update() error

    // Draw được gọi mỗi khi màn hình cần vẽ lại.
    // Tham số screen là "canvas" chính để vẽ lên.
    Draw(screen *ebiten.Image)

    // Layout nhận vào kích thước cửa sổ và trả về kích thước logic của game.
    // Giúp game tự động co giãn (scaling).
    Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}
```

### Cách sử dụng cơ bản
```go
import "github.com/hajimehoshi/ebiten/v2"

type MyGame struct{}

func (g *MyGame) Update() error { return nil }
func (g *MyGame) Draw(screen *ebiten.Image) { /* Vẽ tại đây */ }
func (g *MyGame) Layout(w, h int) (int, int) { return 640, 480 }

func main() {
    ebiten.RunGame(&MyGame{})
}
```

---

## 2. Donburi (ECS)

Donburi không dựa nhiều vào Interface như Ebitengine mà tập trung vào cấu trúc dữ liệu hiệu năng cao. Tuy nhiên, các khái niệm cốt lõi vận hành như sau:

### Các thành phần chính

1.  **World**: Chứa toàn bộ thực thể và dữ liệu.
2.  **Component**: Các struct chứa dữ liệu thuần túy (không có logic).
3.  **Entry**: Một "tay cầm" (handle) để truy cập dữ liệu của một thực thể cụ thể.
4.  **Query**: Công cụ để lọc ra các thực thể có các component nhất định.

### Cách sử dụng cơ bản

#### Định nghĩa Component
```go
type Position struct { X, Y float64 }
var PositionTag = donburi.NewComponentType[Position]()
```

#### Tạo thực thể (Entity)
```go
world := donburi.NewWorld()
entity := world.Create(PositionTag)
entry := world.Entry(entity)

// Gán dữ liệu
donburi.SetValue(entry, PositionTag, Position{X: 10, Y: 20})
```

#### Truy vấn dữ liệu (Query)
```go
query := donburi.NewQuery(donburi.All(PositionTag))
query.Each(world, func(entry *donburi.Entry) {
    pos := donburi.Get[Position](entry, PositionTag)
    pos.X += 1 // Cập nhật vị trí
})
```

---

## 3. Kết hợp với Core Interface của AutoWorld

Chúng ta sẽ bọc (wrap) các thư viện này sau các interface của mình trong `core/interfaces.go` để đảm bảo tính linh hoạt.

- `IWorld` (AutoWorld) sẽ chứa một `donburi.World`.
- `Update()` của `IWorld` sẽ gọi các `ISystem.Update()`.
- `ISystem` sẽ sử dụng `donburi.Query` để xử lý logic.
- `Draw()` của `IWorld` sẽ được gọi từ `ebiten.Game.Draw()`.

---

> [!NOTE]
> Do chúng ta sử dụng thư mục `.libs`, bạn hãy import theo đường dẫn chuẩn (ví dụ: `github.com/hajimehoshi/ebiten/v2`). 
> Hệ thống Go sẽ tự động tìm thấy chúng nhờ cấu hình `replace` trong `go.mod`.
