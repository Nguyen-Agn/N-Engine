package domain

import (
	"image/color"

	"github.com/yohamta/donburi"
)

// ─── Core Object Interface ─────────────────────────────────────────────────────

// IObject là interface gốc của mọi thực thể (entity) trong Engine.
// Bất kỳ đối tượng nào tham gia vào vòng lặp game đều phải implement interface này.
type IObject interface {
	// OnCreate được gọi một lần duy nhất khi Object được khởi tạo vào Scene.
	// Dùng để thiết lập dữ liệu ban đầu, đăng ký event, v.v.
	OnCreate()

	// OnStep được gọi mỗi frame để cập nhật logic của Object.
	OnStep()

	// OnDestroy được gọi khi Object bị xóa khỏi Scene.
	// Dùng để giải phóng tài nguyên, hủy đăng ký event, v.v.
	OnDestroy()

	// OnSave được gọi khi game thực hiện lưu lại.
	// Nhận vào một map trống để developer gắn dữ liệu cần lưu.
	OnSave(data map[string]any)

	// OnLoad được gọi khi load game tương ứng cho Object.
	OnLoad(data map[string]any)

	// GetPool trả về Pool quản lý Object này (nếu có).
	GetPool() IPool
	// SetPool liên kết Object này với một Pool.
	SetPool(pool IPool)

	// Entry trả về ECS entry (donburi) của Object.
	// Dùng để đọc/ghi dữ liệu Component trong hệ thống ECS.
	Entry() *donburi.Entry

	// Using classes to set value of each component
	// Ex: pos-x-5 to set x = 5
	// Ex: spr-c-run to set current sprite to sprite name "run"
	// Format: <component's name>-<variable's name>-<value>
	// Note that: this way only suit for integer/string type value
	SetTokens(tokenClasses string)

	// Remove IObject from Game
	Remove()
}

// IPool định nghĩa giao diện chuẩn cho mọi hệ thống Object Pool của N-Engine.
type IPool interface {
	// Put cất IObject vào pool. Trả về true nếu thành công, false nếu pool đã đầy (để Engine xóa tận gốc).
	Put(obj IObject) bool
}

// ─── Component Interfaces ─────────────────────────────────────────────────────
// Mỗi interface dưới đây tương ứng với một Component Data trong ObjectData.go.
// Object implement interface nào nghĩa là nó sở hữu Component đó.

// IPosition cung cấp tọa độ 2D cho Object (ánh xạ tới PositionData).
type IPosition interface {
	// X trả về tọa độ ngang hiện tại (pixel).
	X() float32
	// Y trả về tọa độ dọc hiện tại (pixel).
	Y() float32
	// SetX thiết lập tọa độ ngang.
	SetX(x float32)
	// SetY thiết lập tọa độ dọc.
	SetY(y float32)
}

// ISprite cung cấp khả năng quản lý hình ảnh và animation cho Object (ánh xạ tới SpriteData).
type ISprite interface {
	// SpriteIdx trả về index frame ảnh hiện tại trong animation.
	SpriteIdx() int
	// SetSpriteIdx thiết lập index frame ảnh hiện tại.
	SetSpriteIdx(spriteIdx int)

	// ImageSpeed trả về tốc độ chạy animation (frame/tick).
	ImageSpeed() float32
	// SetImageSpeed thiết lập tốc độ chạy animation.
	SetImageSpeed(imageSpeed float32)

	// Rotation trả về góc xoay hiện tại của sprite (radian).
	Rotation() float32
	// SetRotation thiết lập góc xoay của sprite.
	SetRotation(rotation float32)

	// OffsetX trả về độ lệch ngang khi vẽ sprite so với tọa độ Object.
	OffsetX() float32
	// SetOffsetX thiết lập độ lệch ngang khi vẽ.
	SetOffsetX(offsetX float32)

	// OffsetY trả về độ lệch dọc khi vẽ sprite so với tọa độ Object.
	OffsetY() float32
	// SetOffsetY thiết lập độ lệch dọc khi vẽ.
	SetOffsetY(offsetY float32)

	// ImageColor trả về màu tô phủ lên sprite (RGBA).
	ImageColor() color.RGBA
	// SetImageColor thiết lập màu tô phủ lên sprite.
	SetImageColor(imageColor color.RGBA)

	// ScaleX trả về hệ số co giãn ngang của sprite. 1.0 = kích thước gốc.
	ScaleX() float32
	// SetScaleX thiết lập hệ số co giãn ngang.
	SetScaleX(scaleX float32)

	// ScaleY trả về hệ số co giãn dọc của sprite. 1.0 = kích thước gốc.
	ScaleY() float32
	// SetScaleY thiết lập hệ số co giãn dọc.
	SetScaleY(scaleY float32)

	// Sprite trả về ISpriteLW theo tên. Trả về nil nếu không tồn tại.
	Sprite(name string) ISpriteLW
	// SetSprite ghi đè sprite theo tên (tạo mới nếu chưa có).
	SetSprite(name string, sprite ISpriteLW)
	// AddSprite thêm sprite mới. Trả về false nếu tên đã tồn tại (không ghi đè).
	AddSprite(name string, sprite ISpriteLW) bool
	// RemoveSprite xóa sprite theo tên. Trả về false nếu không tìm thấy.
	RemoveSprite(name string) bool

	// SetCurrentSprite chọn sprite đang active theo tên để DrawSystem vẽ.
	SetCurrentSprite(name string)
	// GetCurrentSprite trả về ISpriteLW của sprite đang active. Trả về nil nếu chưa chọn.
	GetCurrentSprite() ISpriteLW

	// NextImage tăng ImageIndex lên 1, tự động wrap về 0 khi vượt quá Length.
	NextImage()
	// ImageIndex trả về index frame ảnh hiện tại trong sprite active.
	ImageIndex() int
	// SetImageIndex thiết lập index frame ảnh trong sprite active.
	SetImageIndex(imageIndex int)

	// Enable / Disable 9Slice Mode
	// String Ex: "5" ->5:all, "5 6" -> 5:top&bottom, 6:right&left, "1 2 3 4" -> each
	Set9Slice(turn bool, TopRightBottomLeft string)

	// ZOrder trả về giá trị thứ tự vẽ của sprite (số nhỏ vẽ trước, số lớn vẽ sau).
	ZOrder() int
	// SetZOrder thiết lập thứ tự vẽ của sprite và kích hoạt cờ cập nhật (dirty flag).
	SetZOrder(z int)
}

// IBox cung cấp hitbox hình học cho Object dùng trong va chạm (ánh xạ tới BoxData).
type IBox interface {
	// BoxW trả về chiều rộng của hitbox (pixel).
	BoxW() float32
	// SetBoxW thiết lập chiều rộng hitbox.
	SetBoxW(boxW float32)

	// BoxH trả về chiều cao của hitbox (pixel).
	BoxH() float32
	// SetBoxH thiết lập chiều cao hitbox.
	SetBoxH(boxH float32)

	// BoxX trả về offset ngang của hitbox so với tọa độ Object.
	BoxX() float32
	// SetBoxX thiết lập offset ngang của hitbox.
	SetBoxX(boxX float32)

	// BoxY trả về offset dọc của hitbox so với tọa độ Object.
	BoxY() float32
	// SetBoxY thiết lập offset dọc của hitbox.
	SetBoxY(boxY float32)

	// IsCollidable trả về trạng thái kích hoạt va chạm của Object.
	IsCollidable() bool
	// SetIsCollidable bật/tắt va chạm cho Object.
	SetIsCollidable(isCollidable bool)

	// Shape trả về hình dạng hitbox (rectangle hoặc circle).
	Shape() BoxShape
	// SetShape thiết lập hình dạng hitbox.
	SetShape(shape BoxShape)
}

// IAudio cung cấp khả năng phát âm thanh cho Object (ánh xạ tới AudioData).
type IAudio interface {
	// Audio trả về IAudioLW của kênh âm thanh đang active. Trả về nil nếu chưa có.
	Audio() IAudioLW
	// SetAudio gắn một IAudioLW vào Object với tên kênh chỉ định.
	SetAudio(audioName string, audio IAudioLW)

	// AudioName trả về tên kênh âm thanh đang active.
	AudioName() string
	// SetAudioName chọn kênh âm thanh active theo tên.
	SetAudioName(audioName string)

	// AudioSpeed trả về tốc độ phát âm thanh (1.0 = bình thường).
	AudioSpeed() float32
	// SetAudioSpeed thiết lập tốc độ phát âm thanh.
	SetAudioSpeed(audioSpeed float32)

	// Volume trả về âm lượng hiện tại (0.0 - 1.0).
	Volume() float32
	// SetVolume thiết lập âm lượng.
	SetVolume(volume float32)

	// Pitch trả về cao độ hiện tại (1.0 = bình thường).
	Pitch() float32
	// SetPitch thiết lập cao độ âm thanh.
	SetPitch(pitch float32)

	// Play phát kênh âm thanh tên name với volume và pitch chỉ định.
	// AudioSystem sẽ xử lý lệnh này vào frame tiếp theo.
	Play(name string, volume float32, pitch float32)
	// PlayDefault phát kênh âm thanh tên name với volume và pitch mặc định.
	PlayDefault(name string)
	// StopAudio dừng âm thanh đang phát theo tên.
	StopAudio(name string)
	// PauseAudio tạm dừng âm thanh đang phát theo tên.
	PauseAudio(name string)
	// ResumeAudio tiếp tục phát âm thanh đang tạm dừng theo tên.
	ResumeAudio(name string)
	// SetLooping bật/tắt chế độ lặp lại tự động cho âm thanh theo tên.
	SetLooping(name string, loop bool)
	// IsLooping kiểm tra xem âm thanh chỉ định có đang ở chế độ lặp lại không.
	IsLooping(name string) bool
}

// IInfor cung cấp thông tin định danh cho Object (ánh xạ tới InforData).
type IInfor interface {
	// GetName trả về tên hiển thị hoặc nhãn đại diện của Object.
	GetName() string
	// GetId trả về mã định danh duy nhất (số nguyên tự tăng) của Object.
	GetId() int
	// AddTag thêm một tag định danh vào Object. Chuỗi sẽ được băm thành số nguyên uint64.
	AddTag(tag string)
	// HasTag kiểm tra xem Object có chứa tag này không (băm chuỗi đầu vào).
	HasTag(tag string) bool
	// HasTagHash kiểm tra xem Object có chứa mã băm tag này không.
	HasTagHash(hash uint64) bool
	// IsDead kiểm tra trạng thái của Object. True = không tham gia logic, chờ dọn dẹp.
	IsDead() bool
	// SetIsDead thiết lập trạng thái chết/sống của Object (dùng để hồi sinh từ Pool).
	SetIsDead(dead bool)
	// MarkDead đánh dấu Object là đã chết. Deferred destruction.
	MarkDead()
	// SaveTag trả về mã tag để phân biệt object khi lưu game.
	SaveTag() string
	// SetSaveTag thiết lập mã tag để lưu game.
	SetSaveTag(tag string)
}

// ICollision cung cấp khả năng xử lý va chạm.
type ICollision interface {
	// OnCollisionTag đăng ký callback, kích hoạt mỗi frame khi va chạm với Object có chứa tag chỉ định.
	OnCollisionTag(tag string, handler func(other IObject))
}

// IDirection cung cấp hướng di chuyển cho Object (ánh xạ tới DirectionData).
type IDirection interface {
	// Direction trả về góc hướng hiện tại của Object (đơn vị độ, 0–360).
	Direction() float32
	// SetDirection thiết lập góc hướng mới cho Object.
	SetDirection(dir float32)
	// Rotate xoay hướng Object thêm delta độ (dương = theo chiều kim đồng hồ).
	Rotate(delta float32)
}

// IInput cung cấp khả năng đăng ký lắng nghe phím bấm cho Object (ánh xạ tới InputData).
type IInput interface {
	// ListenOn đăng ký một handler được gọi khi phím (hoặc nhóm phím) kích hoạt.
	// key: tên phím đơn ("w", "space", "enter"...) hoặc tên nhóm đặc biệt.
	// Hỗ trợ truyền nhiều phím/nhóm phím cách nhau bằng dấu cách (VD: "w a s d alpha").
	//   "alpha"  — bất kỳ phím chữ nào (a-z)
	//   "number" — bất kỳ phím số nào (0-9)
	//   "arrows" — bất kỳ phím mũi tên nào (↑↓←→)
	//   "wasd"   — bất kỳ phím W/A/S/D nào
	//   "all"    — bất kỳ phím nào (dùng cho nhập liệu tổng quát)
	// eventType xác định khi nào handler được gọi, là một chuỗi:
	//   ""       — mỗi frame khi phím đang được GIỮ
	//   "pressed"/"p"  — duy nhất 1 lần khi phím vừa được NHẤN XUỐNG
	//   "released"/"r" — duy nhất 1 lần khi phím vừa được THẢ RA
	// handler nhận tên phím đã trigger (ví dụ: "w", "space", "a"...).
	ListenOn(key string, eventType string, handler func(key string))
}

// IMouse cung cấp khả năng đọc trạng thái chuột và đăng ký lắng nghe nút chuột cho Object.
type IMouse interface {
	// MouseX trả về tọa độ ngang của con trỏ chuột (pixel, tính từ góc trên-trái màn hình).
	MouseX() int
	// MouseY trả về tọa độ dọc của con trỏ chuột (pixel, tính từ góc trên-trái màn hình).
	MouseY() int

	// WheelX trả về độ cuộn bánh xe chuột theo trục X trong frame hiện tại.
	// Dương = cuộn phải, âm = cuộn trái.
	WheelX() float64
	// WheelY trả về độ cuộn bánh xe chuột theo trục Y trong frame hiện tại.
	// Dương = cuộn xuống, âm = cuộn lên.
	WheelY() float64

	// ListenMouseOn đăng ký một handler được gọi theo loại sự kiện khi nút chuột kích hoạt.
	// button: "left"/"l", "right"/"r", "middle"/"m" (hoặc nhiều nút cách nhau bằng dấu cách, VD: "left right")
	// eventType: "" (giữ), "pressed"/"p" (vừa nhấn), "released"/"r" (vừa thả)
	// handler nhận tên nút chuột đã trigger (ví dụ: "left").
	ListenMouseOn(button string, eventType string, handler func(button string))
}

// IAlarm cung cấp bộ đếm thời gian để thực thi callback sau N frames (ánh xạ tới AlarmData).
type IAlarm interface {
	// SetAlarm thiết lập một timer với ID cụ thể. Sau khi đếm ngược hết số frames, callback sẽ được gọi.
	// Nếu gọi lại SetAlarm với cùng ID, timer cũ sẽ bị ghi đè.
	SetAlarm(id string, frames int, callback func())
	// GetAlarm trả về số frames còn lại của alarm chỉ định. Trả về 0 nếu không tìm thấy hoặc đã hoàn thành.
	GetAlarm(id string) int
	// StopAlarm hủy bỏ alarm chỉ định (callback sẽ không được gọi).
	StopAlarm(id string)
}

// IVelocity cung cấp vật lý cơ bản cho Object (vận tốc, ma sát) (ánh xạ tới VelocityData).
type IVelocity interface {
	VelocityX() float32
	VelocityY() float32
	SetVelocityX(vx float32)
	SetVelocityY(vy float32)
	SetVelocity(vx, vy float32)
	AddVelocity(vx, vy float32)

	Friction() float32
	SetFriction(f float32)

	MaxSpeed() float32
	SetMaxSpeed(speed float32)
}

// ITween cung cấp khả năng nội suy giá trị (Lerp) mượt mà theo thời gian (ánh xạ tới TweenData).
type ITween interface {
	// TweenMove di chuyển mượt mà Object từ vị trí hiện tại đến (targetX, targetY) trong duration (frames).
	TweenMove(targetX, targetY float32, duration int)
	// TweenScale co giãn mượt mà Object từ Scale hiện tại đến (targetScaleX, targetScaleY) trong duration.
	TweenScale(targetScaleX, targetScaleY float32, duration int)
	// TweenAlpha làm mờ/rõ mượt mà Object từ Alpha hiện tại đến targetAlpha (0-255) trong duration.
	TweenAlpha(targetAlpha uint8, duration int)
}

// IDraw is implemented by Objects that want to perform custom drawing each frame.
// The Object must also have the DrawComponent (token "drw") embedded.
// DrawSystem automatically detects IDraw and calls Draw() after all Sprite entities are rendered.
type IDraw interface {
	// Draw is called every frame by DrawSystem after Sprite rendering.
	// Call drawing methods (Rect, Circle, Text...) from the embedded DrawComponent inside here.
	Draw()
}

// IDrawComponent defines the primitive drawing methods provided by DrawComponent.
// Dev receives these methods via embedding napi.Drw in a Custom Object.
// All coordinates are in map space — camera offset is applied automatically.
//
// Advanced path-based methods (PathFill, PathStroke) and coordinate helpers
// (ScreenX, ScreenY) are available directly on napi.Drw but are not part of
// this interface to keep domain free of ebiten/vector dependencies.
type IDrawComponent interface {
	// --- Filled shapes ---

	// Rect draws a filled rectangle at (x, y) with size (w, h).
	Rect(x, y, w, h float32, c color.RGBA)

	// Circle draws a filled circle centered at (x, y) with radius r.
	Circle(x, y, r float32, c color.RGBA)

	// --- Stroke (outline) shapes ---

	// RectStroke draws a rectangle outline at (x, y) with size (w, h).
	// strokeWidth controls the border thickness in pixels.
	RectStroke(x, y, w, h float32, c color.RGBA, strokeWidth float32)

	// CircleStroke draws a circle outline centered at (x, y) with radius r.
	// strokeWidth controls the border thickness in pixels.
	CircleStroke(x, y, r float32, c color.RGBA, strokeWidth float32)

	// Line draws a straight line from (x0, y0) to (x1, y1).
	// strokeWidth controls the line thickness in pixels.
	Line(x0, y0, x1, y1 float32, c color.RGBA, strokeWidth float32)

	// --- Text ---

	// SetTextAlign cấu hình căn lề cho các lệnh vẽ chữ tiếp theo.
	// Hỗ trợ truyền "left"/"l", "center"/"c", "right"/"r", "justify"/"j".
	SetTextAlign(hAlign, vAlign string)

	// SetTextOverflow cấu hình giới hạn kích thước khung chữ và cách xử lý tràn.
	// Truyền 0 cho maxWidth/maxHeight nếu không muốn giới hạn.
	// mode hỗ trợ "visible"/"v", "hidden"/"h", "scale"/"s".
	SetTextOverflow(maxWidth, maxHeight float32, mode string)

	// Text draws a string at (x, y) with the default font and color c.
	// Use napi.SetDefaultFont() to change the engine-wide font.
	Text(text string, x, y float32, c color.RGBA)

	// TextEx draws a string at (x, y) with a uniform scale applied.
	// scale 1.0 = default size, 2.0 = double size.
	TextEx(text string, x, y float32, c color.RGBA, scale float64)

	// --- Image ---

	// Image draws frame idx of the given ISpriteLW at (x, y).
	// Allows manual sprite rendering without a SpriteComponent.
	Image(sprite ISpriteLW, idx int, x, y float32)
}
