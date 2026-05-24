package domain

type Key int

const (
	KeyUnknown Key = iota
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyApostrophe
	KeyBackslash
	KeyBackspace
	KeyCapsLock
	KeyComma
	KeyDelete
	KeyDown
	KeyEnd
	KeyEnter
	KeyEqual
	KeyEscape
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyGraveAccent
	KeyHome
	KeyInsert
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPAdd
	KeyKPDecimal
	KeyKPDivide
	KeyKPEnter
	KeyKPEqual
	KeyKPMultiply
	KeyKPSubtract
	KeyLeft
	KeyLeftAlt
	KeyLeftBracket
	KeyLeftControl
	KeyLeftShift
	KeyMenu
	KeyMinus
	KeyNumLock
	KeyPageDown
	KeyPageUp
	KeyPause
	KeyPeriod
	KeyPrintScreen
	KeyRight
	KeyRightAlt
	KeyRightBracket
	KeyRightControl
	KeyRightShift
	KeyScrollLock
	KeySemicolon
	KeySlash
	KeySpace
	KeyTab
	KeyUp
)

type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

// IInputManager quản lý input từ keyboard và mouse, không phụ thuộc thư viện ngoài.
type IInputManager interface {
	Update()

	// Keyboard
	IsKeyPressed(key Key) bool
	IsKeyJustPressed(key Key) bool
	IsKeyJustReleased(key Key) bool

	// Mouse
	CursorPosition() (int, int)
	IsMouseButtonPressed(button MouseButton) bool
	IsMouseButtonJustPressed(button MouseButton) bool
	IsMouseButtonJustReleased(button MouseButton) bool
	WheelOffset() (float64, float64)

	// Virtual Actions (Action Mapping)
	BindAction(action string, keys ...Key)
	IsActionPressed(action string) bool
	IsActionJustPressed(action string) bool
	IsActionJustReleased(action string) bool
}

// KeyNameMap cho phép tra cứu Key theo tên chuỗi dễ nhớ.
// Dùng trong ListenOn để chuyển "w", "space", "enter" sang hằng số Key.
var KeyNameMap = map[string]Key{
	"a": KeyA, "b": KeyB, "c": KeyC, "d": KeyD, "e": KeyE,
	"f": KeyF, "g": KeyG, "h": KeyH, "i": KeyI, "j": KeyJ,
	"k": KeyK, "l": KeyL, "m": KeyM, "n": KeyN, "o": KeyO,
	"p": KeyP, "q": KeyQ, "r": KeyR, "s": KeyS, "t": KeyT,
	"u": KeyU, "v": KeyV, "w": KeyW, "x": KeyX, "y": KeyY,
	"z": KeyZ,
	"0": Key0, "1": Key1, "2": Key2, "3": Key3, "4": Key4,
	"5": Key5, "6": Key6, "7": Key7, "8": Key8, "9": Key9,
	"space":     KeySpace,
	"enter":     KeyEnter,
	"escape":    KeyEscape,
	"backspace": KeyBackspace,
	"delete":    KeyDelete,
	"tab":       KeyTab,
	"up":        KeyUp,
	"down":      KeyDown,
	"left":      KeyLeft,
	"right":     KeyRight,
	"shift":     KeyLeftShift,
	"ctrl":      KeyLeftControl,
	"alt":       KeyLeftAlt,
	"f1":        KeyF1, "f2": KeyF2, "f3": KeyF3, "f4": KeyF4,
	"f5":        KeyF5, "f6": KeyF6, "f7": KeyF7, "f8": KeyF8,
	"f9":        KeyF9, "f10": KeyF10, "f11": KeyF11, "f12": KeyF12,
}

// KeyGroupMap ánh xạ tên nhóm đặc biệt sang danh sách các Key.
// Dùng trong ListenOn để hỗ trợ các trường hợp nhập liệu, v.d: "number", "alpha".
var KeyGroupMap = map[string][]Key{
	// Nhóm chữ cái a-z: dùng để nhập văn bản
	"alpha": {
		KeyA, KeyB, KeyC, KeyD, KeyE, KeyF, KeyG, KeyH, KeyI, KeyJ,
		KeyK, KeyL, KeyM, KeyN, KeyO, KeyP, KeyQ, KeyR, KeyS, KeyT,
		KeyU, KeyV, KeyW, KeyX, KeyY, KeyZ,
	},
	// Nhóm số 0-9: dùng để nhập số
	"number": {Key0, Key1, Key2, Key3, Key4, Key5, Key6, Key7, Key8, Key9},
	// Nhóm phím mũi tên
	"arrows": {KeyUp, KeyDown, KeyLeft, KeyRight},
	// Nhóm WASD: di chuyển kiểu game PC
	"wasd": {KeyW, KeyA, KeyS, KeyD},
	// Nhóm tất cả phím: dùng khi muốn bắt bất kỳ phím nào
	"all": {
		KeyA, KeyB, KeyC, KeyD, KeyE, KeyF, KeyG, KeyH, KeyI, KeyJ,
		KeyK, KeyL, KeyM, KeyN, KeyO, KeyP, KeyQ, KeyR, KeyS, KeyT,
		KeyU, KeyV, KeyW, KeyX, KeyY, KeyZ,
		Key0, Key1, Key2, Key3, Key4, Key5, Key6, Key7, Key8, Key9,
		KeySpace, KeyEnter, KeyEscape, KeyBackspace, KeyDelete, KeyTab,
		KeyUp, KeyDown, KeyLeft, KeyRight,
		KeyLeftShift, KeyLeftControl, KeyLeftAlt,
		KeyF1, KeyF2, KeyF3, KeyF4, KeyF5, KeyF6,
		KeyF7, KeyF8, KeyF9, KeyF10, KeyF11, KeyF12,
	},
}
