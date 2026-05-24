# Hướng dẫn sử dụng
## Tạo đối tượng
```go
// định nghĩa 
type Player struct {
    napi.IObject // triển khai object
    napi.ISprite // triển khai hình ảnh (cần vị trí)
    napi.IPosition // triển khai vị trí
    napi.IInput // triển khai nhận sự kiện  điều khiển
}
```

```go
// tạo Contrustor 
func NewPlayer()  *PLayer {
    p := &PLayer{} // tạo 1 con trỏ để lưu đối tượng

    // hiện thực hóa đối tượng PLayer thông qua thư viện
    // NewObject nhận vào 3 biến
    // 1: là con trỏ lưu trữ đối tượng
    // 2: tên đối tượng (có thể là "")
    // 3: danh sách component cần được triển khai
    // Ví dụ: spr: Sprite, pos: Postion, inp: Input
    napi.NewObject(p, "PLayer","spr pos inp");

    // đăng ký Object vào scene
    // đầu thứ 2 là tên của scene mà bạn muốn
    napi.Register(p, "scene")
}
```
```go
// lập trình hành viện
// Đây là hàm thức duy nhất 1 lần khi Object được tạo 
// Thực thi sau Contrustor và trong cùng nhịp sự kiện của
func (this *Player) Create() {
    this.SetX(0);
    this.SetY(0);

    // Lắng nghe phím "space" -> nếu có -> chạy hàm OnSpace (tự định nghĩa)
    this.ListenOn("space", p.OnSpace)
}

// Đây là hàm sẽ lặp lại liên tục nhiều lần mỗi giây
func (this *PLayer) StepUpdate() {

}

// Đây là hàm sẽ chạy trước khi Object bị hủy 
func (this *PLayer) Destroy() {

}

// Đây là hàm tự định nghĩa, ví dụ khi phím Space được nhấn
func (this *PLayer) OnSpace() {

}

```
## Game begin 
```go
    func main() {
	// 1. Khởi tạo cấu hình engine
	cfg := napi.GameConfig{
		Title:  "Example",
		Width:  640,
		Height: 480,
	}
    // Nạp cấu hình game
    // có thể gôm chung bước 1
	napi.Init(cfg)

	// 2. Load tài nguyên từ từ tệp danh  sách .toml
	napi.LoadFromFile("./tilemap_manifest.toml")

	// 3. Khởi tạo Scene
    // đầu vào 1: là tên scene 
    // Đầu vào 2: là các lớp thành phần (map,gui,map-w-h,gui-w-h)
    // map-w-h: là tạo bàn đồ có kích thước w*h
	_, err := napi.NewSceneAndGo("main", "map-640-480")
	if err != nil {
		log.Fatalf("Không thể khởi tạo Scene: %v", err)
	}

	// 4. Tạo các thực thể game
	_ := NewPlayer() // ví dụ PLayer

	// 5. Chạy game loop
	napi.GameStart()
}
```
## Background
Bản chất là 1 Object với component napi.IBackground
```go
type BackgroundObject struct {
	napi.IObject
	napi.IBackground
}

func NewBackgroundObject() *BackgroundObject {
	bg := &BackgroundObject{}
	napi.NewObject(bg, "bg_obj", "bg")

	// Đặt màu nền sẫm (Xanh lục sẫm)
	bg.SetColor(color.RGBA{10, 30, 20, 255})

	
    // napi.Gétprite(tên ảnh được dặt trước đó) từ store sau khi LoadFromFile
    img := napi.GetSprite("bg-image")
    // Gán ảnh nền
	bg.SetSprite(img)

    //  Đặt cấu hình cho ảnh nền
	bg.SetRepeatX(true)
	bg.SetRepeatY(true)

	bg.SetScrollSpeedX(0.3)
	bg.SetScrollSpeedY(0.1)

    // đăng ký như object bình thường
	napi.Register(bg, "main")
	return bg
}
```

## Tile Map
Bản chất là 1 object với napi.ITilemap
```go
type TilemapObject struct {
	napi.IObject
	napi.ITilemap
	napi.IPosition // cần cho TileMap để xác định vị trí bắt đầu
}

func NewTilemapObject() *TilemapObject {
	tm := &TilemapObject{}
	napi.NewObject(tm, "tilemap_obj", "pos tile")

    // tương tự gán ảnh của background
	tm.SetSprite(napi.GetSprite("tileset"))

    // kích thước
	tm.SetTileWidth(32)
	tm.SetTileHeight(32)

	// Tạo ma trận gạch 2D
    // giá trị n chính tile thứ n tính theo hàng trước cột sau.
    // Giá trị bắt đầu từ 0 --> n-1; -1 là rỗng
	grid := [][]int{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1,  0,  1,  2,  3, -1, -1, -1, -1, -1, -1, -1, -1,  0,  1,  2,  3, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
	}
    // gán ma trận cho tm
	tm.SetGrid(grid)

	// Đặt vị trí bắt đầu vẽ bản đồ
	tm.SetX(0)
	tm.SetY(100)


    // Đăng ký object vào scene cần dùng
	napi.Register(tm, "main")
	return tm
}
```

```go
func main() {
    ... // scene created
    NewBackgroundObject()
    NewTilemapObject()
    ... 
}
```

## Tự định nghĩa Compoenet
```go

```