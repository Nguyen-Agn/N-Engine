package ncom

import (
	"autoworld/domain"
	"autoworld/modules/components"

	"github.com/yohamta/donburi"
)

// ─── Draw ───────────────────────────────────────────────────────────────────
// Token: No
// Provide: (you can override them)
// - OnCreate(): once after object created
// - OnDestroy(): once after object destroyed
// - OnStep():every (mid) frame
// - OnSave(map[string]any): When napi.Game.SaveGame called
// - OnLoad(map[string]any): When napi.Game.LoadGame called
// - SetTokens(string): Convient function, set value(s) by 1 string of classes
type Object = domain.IObject

// ─── Position ───────────────────────────────────────────────────────────────
// Token: "pos"
// Provide: X(), Y(), SetX(), SetY()
type Pos = components.PositionComponent

// ─── Sprite ─────────────────────────────────────────────────────────────────
// Token: "spr"
// Provide: AddSprite(), SetCurrentSprite(), SetScaleX/Y(), SetImageSpeed()...
type Spr = components.SpriteComponent

// ─── Box (Hitbox) ────────────────────────────────────────────────────────────
// Token: "box" (tự động thêm khi dùng "col")
// Provide: BoxW(), BoxH(), SetBoxW/H(), IsCollidable()...
type Box = components.BoxComponent

// ─── Audio ──────────────────────────────────────────────────────────────────
// Token: "aud"
// Provide: Play(), Stop(), SetVolume(), SetPitch()...
type Aud = components.AudioComponent

// ─── Infor (Identity) ───────────────────────────────────────────────────────
// Token: "inf" (tự động thêm vào mọi object)
// Provide: GetName(), GetId(), SaveTag(), SetSaveTag(), IsDead(), MarkDead()...
type Info = components.InforComponent

// ─── Direction ──────────────────────────────────────────────────────────────
// Token: "dir"
// Provide: Direction(), SetDirection(), Rotate()
type Dir = components.DirectionComponent

// ─── Input (Keyboard) ───────────────────────────────────────────────────────
// Token: "inp"
// Provide: ListenOn(key, eventType, handler)
//   - handler nhận tên phím đã trigger (e.g. "w", "space", "a")
//   - eventType: domain.EventPressed / EventJustPressed / EventJustReleased
type Inp = components.InputComponent

// ─── Mouse ──────────────────────────────────────────────────────────────────
// Token: không cần (không lưu ECS data)
// Provide:
//   - MouseX() int, MouseY() int         — tọa độ con trỏ chuột (pixel)
//   - WheelX() float64, WheelY() float64 — tốc độ cuộn trục X/Y
//   - ListenMouseOn(button, eventType, handler)
//     button: "left", "right", "middle"
//     eventType: domain.EventPressed / EventJustPressed / EventJustReleased
//     handler nhận tên nút chuột đã trigger (e.g. "left")
type Mouse = components.MouseComponent

// ─── Background ─────────────────────────────────────────────────────────────
// Token: "back"
// Provide: SetBackground(), SetRepeat(), SetScrollSpeed()...
type Back = components.BackgroundComponent

// ─── Tilemap ────────────────────────────────────────────────────────────────
// Token: "til"
// Provide: SetTilemap(), SetGrid()...
type Tile = components.TilemapComponent

// ─── Alarm ──────────────────────────────────────────────────────────────────
// Token: "alr"
// Provide: SetAlarm(id, frames, callback), GetAlarm(id), StopAlarm(id)
type Alrm = components.AlarmComponent

// ─── Velocity ───────────────────────────────────────────────────────────────
// Token: "vel"
// Provide: VelocityX/Y(), SetVelocity(), Friction(), MaxSpeed()...
type Velo = components.VelocityComponent

// ─── Tween ──────────────────────────────────────────────────────────────────
// Token: "twn" (tự động thêm "spr")
// Provide: TweenMove(), TweenScale(), TweenAlpha()
type Twn = components.TweenComponent

// ─── Collision ──────────────────────────────────────────────────────────────
// Token: "col" (tự động thêm "box")
// Provide: AddTag(), OnCollisionTag(tag, handler)
type Col = components.CollisionComponent

// ─── Draw ───────────────────────────────────────────────────────────────────
// Token: "drw" (tự động thêm "pos")
// Provide: Rect(), Circle(), Line(), Text(), Image()...
// Object cũng cần implement IDraw interface (hàm Draw()).
type Drw = components.DrawComponent

// ─── Debug ──────────────────────────────────────────────────────────────────
// Token: "deb"
// Provide: Log(), DrawDebug()...
type Deb = components.DebugComponent

// ─── Generic Component ──────────────────────────────────────────────────────
// Dùng để tạo Custom Component với kiểu dữ liệu riêng.
// Kết hợp với napi.NewComponentType[T]() để khai báo ComponentType.
//
// Ví dụ:
//
//	var StatsComp = napi.NewComponentType[StatsData]("sta")
//
//	type StatsComponent struct {
//	    ncom.Generic[StatsData]
//	}
//
//	func (s *StatsComponent) HP() int { return s.Get().Health }
type Generic[T any] = components.GenericComponent[T]

// NewGeneric tạo Generic component đã gắn sẵn ComponentType.
// Phải gọi khi khởi tạo struct, trước khi gọi napi.Obj.NewObject().
//
// Ví dụ:
//
//	h := &Hero{
//	    StatsComponent: StatsComponent{
//	        Generic: ncom.NewGeneric(StatsComp),
//	    },
//	}
func NewGeneric[T any](comp *donburi.ComponentType[T]) components.GenericComponent[T] {
	return components.NewGenericComponent(comp)
}
