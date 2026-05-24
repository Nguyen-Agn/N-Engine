package core

import (
	"autoworld/domain"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ─── Key Mapping Table ────────────────────────────────────────────────────────
// Ánh xạ domain.Key (độc lập thư viện) sang ebiten.Key (ebiten-specific).
// Bảng này là cầu nối duy nhất giữa định nghĩa phím trừu tượng và thư viện ebiten.
var ebitenKeyMap = map[domain.Key]ebiten.Key{
	domain.KeyA:            ebiten.KeyA,
	domain.KeyB:            ebiten.KeyB,
	domain.KeyC:            ebiten.KeyC,
	domain.KeyD:            ebiten.KeyD,
	domain.KeyE:            ebiten.KeyE,
	domain.KeyF:            ebiten.KeyF,
	domain.KeyG:            ebiten.KeyG,
	domain.KeyH:            ebiten.KeyH,
	domain.KeyI:            ebiten.KeyI,
	domain.KeyJ:            ebiten.KeyJ,
	domain.KeyK:            ebiten.KeyK,
	domain.KeyL:            ebiten.KeyL,
	domain.KeyM:            ebiten.KeyM,
	domain.KeyN:            ebiten.KeyN,
	domain.KeyO:            ebiten.KeyO,
	domain.KeyP:            ebiten.KeyP,
	domain.KeyQ:            ebiten.KeyQ,
	domain.KeyR:            ebiten.KeyR,
	domain.KeyS:            ebiten.KeyS,
	domain.KeyT:            ebiten.KeyT,
	domain.KeyU:            ebiten.KeyU,
	domain.KeyV:            ebiten.KeyV,
	domain.KeyW:            ebiten.KeyW,
	domain.KeyX:            ebiten.KeyX,
	domain.KeyY:            ebiten.KeyY,
	domain.KeyZ:            ebiten.KeyZ,
	domain.Key0:            ebiten.Key0,
	domain.Key1:            ebiten.Key1,
	domain.Key2:            ebiten.Key2,
	domain.Key3:            ebiten.Key3,
	domain.Key4:            ebiten.Key4,
	domain.Key5:            ebiten.Key5,
	domain.Key6:            ebiten.Key6,
	domain.Key7:            ebiten.Key7,
	domain.Key8:            ebiten.Key8,
	domain.Key9:            ebiten.Key9,
	domain.KeyApostrophe:   ebiten.KeyApostrophe,
	domain.KeyBackslash:    ebiten.KeyBackslash,
	domain.KeyBackspace:    ebiten.KeyBackspace,
	domain.KeyCapsLock:     ebiten.KeyCapsLock,
	domain.KeyComma:        ebiten.KeyComma,
	domain.KeyDelete:       ebiten.KeyDelete,
	domain.KeyDown:         ebiten.KeyDown,
	domain.KeyEnd:          ebiten.KeyEnd,
	domain.KeyEnter:        ebiten.KeyEnter,
	domain.KeyEqual:        ebiten.KeyEqual,
	domain.KeyEscape:       ebiten.KeyEscape,
	domain.KeyF1:           ebiten.KeyF1,
	domain.KeyF2:           ebiten.KeyF2,
	domain.KeyF3:           ebiten.KeyF3,
	domain.KeyF4:           ebiten.KeyF4,
	domain.KeyF5:           ebiten.KeyF5,
	domain.KeyF6:           ebiten.KeyF6,
	domain.KeyF7:           ebiten.KeyF7,
	domain.KeyF8:           ebiten.KeyF8,
	domain.KeyF9:           ebiten.KeyF9,
	domain.KeyF10:          ebiten.KeyF10,
	domain.KeyF11:          ebiten.KeyF11,
	domain.KeyF12:          ebiten.KeyF12,
	domain.KeyGraveAccent:  ebiten.KeyGraveAccent,
	domain.KeyHome:         ebiten.KeyHome,
	domain.KeyInsert:       ebiten.KeyInsert,
	domain.KeyKP0:          ebiten.KeyKP0,
	domain.KeyKP1:          ebiten.KeyKP1,
	domain.KeyKP2:          ebiten.KeyKP2,
	domain.KeyKP3:          ebiten.KeyKP3,
	domain.KeyKP4:          ebiten.KeyKP4,
	domain.KeyKP5:          ebiten.KeyKP5,
	domain.KeyKP6:          ebiten.KeyKP6,
	domain.KeyKP7:          ebiten.KeyKP7,
	domain.KeyKP8:          ebiten.KeyKP8,
	domain.KeyKP9:          ebiten.KeyKP9,
	domain.KeyKPAdd:        ebiten.KeyKPAdd,
	domain.KeyKPDecimal:    ebiten.KeyKPDecimal,
	domain.KeyKPDivide:     ebiten.KeyKPDivide,
	domain.KeyKPEnter:      ebiten.KeyKPEnter,
	domain.KeyKPEqual:      ebiten.KeyKPEqual,
	domain.KeyKPMultiply:   ebiten.KeyKPMultiply,
	domain.KeyKPSubtract:   ebiten.KeyKPSubtract,
	domain.KeyLeft:         ebiten.KeyLeft,
	domain.KeyLeftAlt:      ebiten.KeyAlt,
	domain.KeyLeftBracket:  ebiten.KeyLeftBracket,
	domain.KeyLeftControl:  ebiten.KeyControl,
	domain.KeyLeftShift:    ebiten.KeyShift,
	domain.KeyMenu:         ebiten.KeyMenu,
	domain.KeyMinus:        ebiten.KeyMinus,
	domain.KeyNumLock:      ebiten.KeyNumLock,
	domain.KeyPageDown:     ebiten.KeyPageDown,
	domain.KeyPageUp:       ebiten.KeyPageUp,
	domain.KeyPause:        ebiten.KeyPause,
	domain.KeyPeriod:       ebiten.KeyPeriod,
	domain.KeyPrintScreen:  ebiten.KeyPrintScreen,
	domain.KeyRight:        ebiten.KeyRight,
	domain.KeyRightAlt:     ebiten.KeyAlt,
	domain.KeyRightBracket: ebiten.KeyRightBracket,
	domain.KeyRightControl: ebiten.KeyControl,
	domain.KeyRightShift:   ebiten.KeyShift,
	domain.KeyScrollLock:   ebiten.KeyScrollLock,
	domain.KeySemicolon:    ebiten.KeySemicolon,
	domain.KeySlash:        ebiten.KeySlash,
	domain.KeySpace:        ebiten.KeySpace,
	domain.KeyTab:          ebiten.KeyTab,
	domain.KeyUp:           ebiten.KeyUp,
}

// ebitenMouseMap ánh xạ domain.MouseButton sang ebiten.MouseButton.
var ebitenMouseMap = map[domain.MouseButton]ebiten.MouseButton{
	domain.MouseButtonLeft:   ebiten.MouseButtonLeft,
	domain.MouseButtonRight:  ebiten.MouseButtonRight,
	domain.MouseButtonMiddle: ebiten.MouseButtonMiddle,
}

// ─── InputManager ─────────────────────────────────────────────────────────────

// InputManager triển khai domain.IInputManager bằng thư viện Ebitengine.
// Là bridge duy nhất giữa hệ thống input trừu tượng (domain) và ebiten input API.
type InputManager struct {
	// actions lưu ánh xạ tên action → danh sách Key (Virtual Action Mapping).
	actions map[string][]domain.Key
}

// NewInputManager khởi tạo InputManager mới với bảng action rỗng.
func NewInputManager() *InputManager {
	return &InputManager{
		actions: make(map[string][]domain.Key),
	}
}

// ─── Keyboard ─────────────────────────────────────────────────────────────────

// Update được gọi mỗi frame. Ebiten tự cập nhật trạng thái input nội bộ,
// nên method này hiện tại không cần làm gì thêm.
func (im *InputManager) Update() {}

// IsKeyPressed trả về true nếu phím đang được giữ trong frame hiện tại.
func (im *InputManager) IsKeyPressed(key domain.Key) bool {
	if eKey, ok := ebitenKeyMap[key]; ok {
		return ebiten.IsKeyPressed(eKey)
	}
	return false
}

// IsKeyJustPressed trả về true nếu phím vừa được nhấn xuống tại frame này (edge trigger).
func (im *InputManager) IsKeyJustPressed(key domain.Key) bool {
	if eKey, ok := ebitenKeyMap[key]; ok {
		return inpututil.IsKeyJustPressed(eKey)
	}
	return false
}

// IsKeyJustReleased trả về true nếu phím vừa được thả ra tại frame này (edge trigger).
func (im *InputManager) IsKeyJustReleased(key domain.Key) bool {
	if eKey, ok := ebitenKeyMap[key]; ok {
		return inpututil.IsKeyJustReleased(eKey)
	}
	return false
}

// ─── Mouse ────────────────────────────────────────────────────────────────────

// CursorPosition trả về tọa độ con trỏ chuột hiện tại (pixel, tính từ góc trên-trái).
func (im *InputManager) CursorPosition() (int, int) {
	return ebiten.CursorPosition()
}

// IsMouseButtonPressed trả về true nếu nút chuột đang được giữ.
func (im *InputManager) IsMouseButtonPressed(button domain.MouseButton) bool {
	if eBtn, ok := ebitenMouseMap[button]; ok {
		return ebiten.IsMouseButtonPressed(eBtn)
	}
	return false
}

// IsMouseButtonJustPressed trả về true nếu nút chuột vừa được nhấn tại frame này.
func (im *InputManager) IsMouseButtonJustPressed(button domain.MouseButton) bool {
	if eBtn, ok := ebitenMouseMap[button]; ok {
		return inpututil.IsMouseButtonJustPressed(eBtn)
	}
	return false
}

// IsMouseButtonJustReleased trả về true nếu nút chuột vừa được thả tại frame này.
func (im *InputManager) IsMouseButtonJustReleased(button domain.MouseButton) bool {
	if eBtn, ok := ebitenMouseMap[button]; ok {
		return inpututil.IsMouseButtonJustReleased(eBtn)
	}
	return false
}

// WheelOffset trả về độ cuộn bánh xe chuột theo trục X và Y trong frame hiện tại.
func (im *InputManager) WheelOffset() (float64, float64) {
	return ebiten.Wheel()
}

// ─── Virtual Action Mapping ───────────────────────────────────────────────────

// BindAction gán một tên action logic với một hoặc nhiều Key vật lý.
// Sau khi bind, dùng IsActionPressed/JustPressed/JustReleased thay vì kiểm tra từng phím.
// Ví dụ: BindAction("jump", domain.KeySpace, domain.KeyW)
func (im *InputManager) BindAction(action string, keys ...domain.Key) {
	im.actions[action] = keys
}

// IsActionPressed trả về true nếu bất kỳ phím nào trong action đang được giữ (OR logic).
func (im *InputManager) IsActionPressed(action string) bool {
	keys, ok := im.actions[action]
	if !ok {
		return false
	}
	return slices.ContainsFunc(keys, im.IsKeyPressed)
}

// IsActionJustPressed trả về true nếu bất kỳ phím nào trong action vừa được nhấn (OR logic).
func (im *InputManager) IsActionJustPressed(action string) bool {
	keys, ok := im.actions[action]
	if !ok {
		return false
	}
	return slices.ContainsFunc(keys, im.IsKeyJustPressed)
}

// IsActionJustReleased trả về true nếu bất kỳ phím nào trong action vừa được thả (OR logic).
func (im *InputManager) IsActionJustReleased(action string) bool {
	keys, ok := im.actions[action]
	if !ok {
		return false
	}
	return slices.ContainsFunc(keys, im.IsKeyJustReleased)
}
