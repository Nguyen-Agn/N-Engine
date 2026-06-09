package domain

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// base
type PositionData struct {
	X, Y float32
}

// NineSlice holds 9-slice configuration for a sprite.
type NineSlice struct {
	Top, Right, Bottom, Left int
}

// graphic
type SpriteData struct {
	Sprite           map[string]ISpriteLW
	CurrentSprite    string
	SpriteIdx        int
	IsVisible        bool // false = skip rendering this entity
	ImageSpeed       float32
	Rotation         float32
	OffsetX, OffsetY float32
	ImageColor       color.RGBA
	ScaleX, ScaleY   float32
	IsNineSlice      bool
	NineSlice        NineSlice
	ZOrder           int
	IsZOrderDirty    bool
}

// BackgroundData holds background image or color data for a Scene.
type BackgroundData struct {
	Sprite       ISpriteLW  // background image (nil = color only)
	Color        color.RGBA // background fill color
	RepeatX      bool       // tile image horizontally
	RepeatY      bool       // tile image vertically
	Stretch      bool       // stretch image to screen size
	ScrollSpeedX float32    // auto-scroll speed X (pixels/frame)
	ScrollSpeedY float32    // auto-scroll speed Y (pixels/frame)
	OffsetX      float32    // current scroll offset X
	OffsetY      float32    // current scroll offset Y
	IsVisible    bool       // visibility flag
}

// TilemapData holds tilemap grid data.
type TilemapData struct {
	Sprite     ISpriteLW // tileset sprite sheet
	TileWidth  int       // tile width in pixels
	TileHeight int       // tile height in pixels
	Cols       int       // number of columns
	Rows       int       // number of rows
	Grid       []int     // flat array of tile IDs (Cols*Rows); -1 = empty
	IsVisible  bool      // visibility flag
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
	Id      int
	Name    string
	Tags    []uint64
	IsDead  bool
	SaveTag string
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

// KeyBinding liên kết một tập hợp phím (OR logic) với một handler function.
// EventType xác định loại sự kiện cần bắt (giữ, vừa nhấn, vừa thả).
// Handler nhận tên phím đã kích hoạt dưới dạng chuỗi (e.g. "w", "space").
type KeyBinding struct {
	Keys      []Key
	EventType EventType
	Handler   func(key string) // nhận tên phím đã trigger
}

// MouseBinding liên kết một nút chuột với một handler function.
// EventType xác định loại sự kiện (giữ, vừa nhấn, vừa thả).
type MouseBinding struct {
	Button    MouseButton
	EventType EventType
	Handler   func(button string) // nhận tên nút chuột đã trigger
}

// InputData là ECS data cho Input Component (keyboard).
// Chứa toàn bộ KeyBindings mà object muốn lắng nghe.
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
	// ShouldPlay signals AudioSystem to play audio. Reset to false after processing.
	ShouldPlay bool
	// ShouldStop signals AudioSystem to stop current audio.
	ShouldStop bool
	// ShouldPause signals AudioSystem to pause current audio.
	ShouldPause bool
	// ShouldResume signals AudioSystem to resume paused audio.
	ShouldResume bool
	// IsLooping signals that the current audio should automatically loop.
	IsLooping bool
}

// draw

// DrawData is the ECS data for the Draw Component (token "drw").
// Screen, CamX, CamY are transient — injected by DrawSystem before each Draw() call.
// Dev does not interact with this struct directly.
type DrawData struct {
	Screen *ebiten.Image // current render canvas — set by DrawSystem each frame
	CamX   float32       // camera offset X — set by DrawSystem each frame
	CamY   float32       // camera offset Y — set by DrawSystem each frame
}

type BoxShape string

const (
	BSRectangle BoxShape = "rectangle"
	BSCircle    BoxShape = "circle"
)

// DebugData lưu trữ cấu hình hiển thị overlay debug và log tạm thời của Object.
type DebugData struct {
	ShowBox      bool       // Có vẽ hitbox không
	ShowPos      bool       // Có vẽ origin crosshair không
	ShowInfo     bool       // Có hiển thị info (ID, Name) không
	Color        color.RGBA // Màu vẽ debug overlay
	CustomLog    string     // Log text từ lệnh Log()
}
