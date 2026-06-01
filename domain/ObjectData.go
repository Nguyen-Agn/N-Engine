package domain

import (
	"image/color"
)

// base
type PositionData struct {
	X, Y float32
}

// NineSlice chứa cấu hình 9-slice cho sprite.
type NineSlice struct {
	Top, Right, Bottom, Left int
}

// graphic
type SpriteData struct {
	Sprite           map[string]ISpriteLW
	CurrentSprite    string
	SpriteIdx        int
	IsVisible        bool // false = không vẽ thực thể này
	ImageSpeed       float32
	Rotation         float32
	OffsetX, OffsetY float32
	ImageColor       color.RGBA
	ScaleX, ScaleY   float32
	IsNineSlice      bool
	NineSlice        NineSlice
}

// BackgroundData lưu thông tin hình nền hoặc màu nền của Scene.
type BackgroundData struct {
	Sprite       ISpriteLW  // Ảnh nền (nil nếu chỉ dùng màu nền)
	Color        color.RGBA // Màu nền để fill màn hình (thường dùng làm màu gốc)
	RepeatX      bool       // Lặp lại ảnh nền theo chiều ngang
	RepeatY      bool       // Lặp lại ảnh nền theo chiều dọc
	Stretch      bool       // Co giãn ảnh nền vừa kích thước màn hình
	ScrollSpeedX float32    // Tốc độ tự động cuộn ngang (auto scroll)
	ScrollSpeedY float32    // Tốc độ tự động cuộn dọc (auto scroll)
	OffsetX      float32    // Offset cuộn X hiện tại
	OffsetY      float32    // Offset cuộn Y hiện tại
	IsVisible    bool       // Trạng thái hiển thị
}

// TilemapData lưu trữ thông tin bản đồ gạch.
type TilemapData struct {
	Sprite     ISpriteLW // Tileset sprite sheet (chứa các frame ảnh tile)
	TileWidth  int       // Chiều rộng mỗi tile (pixel)
	TileHeight int       // Chiều cao mỗi tile (pixel)
	Cols       int       // Số cột của bản đồ
	Rows       int       // Số hàng của bản đồ
	Grid       []int     // Mảng 1 chiều lưu ID của Tile kích thước Cols * Rows (-1 nghĩa là trống)
	IsVisible  bool      // Trạng thái hiển thị
}

// physics
type BoxData struct {
	BoxW, BoxH   float32
	BoxX, BoxY   float32
	Shape        BoxShape
	IsCollidable bool
}

// information

type InforData struct {
	Id     int
	Name   string
	Tags   []uint64
	IsDead bool
}

// collision
type CollisionData struct {
	Handlers     map[uint64]func(other IObject)
	IsCollidable bool
}

// alarm
type Alarm struct {
	Id         string
	FramesLeft int
	IsActive   bool
	Callback   func()
}

type AlarmData struct {
	Alarms []Alarm
}

// velocity
type VelocityData struct {
	Vx, Vy   float32
	Friction float32
	MaxSpeed float32
}

// tween
type Tween struct {
	Id         string
	IsActive   bool
	TargetType string // "move", "scale", "alpha"
	Duration   int    // frames
	Elapsed    int    // frames
	StartX     float32
	StartY     float32
	EndX       float32
	EndY       float32
	StartColor color.RGBA
	EndColor   color.RGBA
}

type TweenData struct {
	Tweens []Tween
}

// direction
type DirectionData struct {
	Direction float32
}

// input

// KeyBinding liên kết một nhóm phím (Keys) với một hàm xử lý (Handler).
// OR logic: bất kỳ phím nào trong Keys được nhấn thì Handler được kích hoạt.
type KeyBinding struct {
	Keys    []Key
	Handler func()
}

// InputData là dữ liệu ECS của Input Component.
// Lưu danh sách các KeyBinding mà Object muốn lắng nghe.
type InputData struct {
	Bindings []KeyBinding
}

// audio
type AudioData struct {
	Audio      map[string]IAudioLW
	AudioName  string
	AudioSpeed float32
	Volume     float32
	Pitch      float32
	// ShouldPlay là cờ ra lệnh cho AudioSystem phát âm thanh.
	// Object đặt thành true, AudioSystem sẽ đọc và reset về false.
	ShouldPlay bool
	// ShouldStop là cờ ra lệnh cho AudioSystem dừng âm thanh hiện tại.
	ShouldStop bool
}

type BoxShape string

const (
	BSRectangle BoxShape = "rectangle"
	BSCircle    BoxShape = "circle"
)

// var (
// 	Position = donburi.NewComponentType[PositionData]()
// 	Sprite   = donburi.NewComponentType[SpriteData]()
// 	Box      = donburi.NewComponentType[BoxData]()
// 	Audio    = donburi.NewComponentType[AudioData]()
// )
