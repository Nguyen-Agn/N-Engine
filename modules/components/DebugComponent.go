package components

import (
	"image/color"
	"strings"

	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Debug = enginetype.Debug

func init() {
	enginetype.RegisterComponentInitializer("deb", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Debug, domain.DebugData{
			ShowBox:  true,
			ShowPos:  true,
			ShowInfo: true,
			Color:    color.RGBA{255, 255, 255, 255}, // Mặc định màu trắng
		})
	})
}

// DebugComponent là Mixin để nhúng vào Custom Object (token "deb").
// Cung cấp các hàm vẽ overlay debug cho object (hitbox, tâm tọa độ, thông tin ID/Tên)
// cũng như hàm Log(msg) in chữ tự do hiển thị trên Object.
type DebugComponent struct {
	IObject
	data *DebugData
}

func (c *DebugComponent) BindComponent(base IObject) {
	c.IObject = base
	c.data = enginetype.GetComponent(base, Debug)
}

// SetDebugColor thay đổi màu sắc của khung debug (mặc định là trắng).
func (c *DebugComponent) SetDebugColor(col color.RGBA) {
	if c.data != nil {
		c.data.Color = col
	}
}

// Debug cấu hình các thành phần hiển thị trên debug overlay thông qua chuỗi flags.
// Tham số là một chuỗi chứa các tuỳ chọn cách nhau bằng dấu cách:
// "box": vẽ khung va chạm
// "pos": vẽ tâm gốc tọa độ
// "info": in ra [ID] Name của object
// Gọi Debug("pos box") sẽ tắt info và chỉ vẽ pos và box.
func (c *DebugComponent) Debug(flags string) {
	if c.data == nil {
		return
	}
	c.data.ShowBox = false
	c.data.ShowPos = false
	c.data.ShowInfo = false

	tokens := strings.Fields(strings.ToLower(flags))
	for _, t := range tokens {
		switch t {
		case "box":
			c.data.ShowBox = true
		case "pos":
			c.data.ShowPos = true
		case "info":
			c.data.ShowInfo = true
		}
	}
}

// Log ghi lại một chuỗi văn bản để hiển thị ngay cạnh object trên màn hình (dùng cho debug).
// Nếu object không có vị trí (không có Pos component), text sẽ được in tĩnh
// ở góc trên trái màn hình, và tự động xếp chồng tránh đè chữ.
func (c *DebugComponent) Log(msg string) {
	if c.data != nil {
		c.data.CustomLog = msg
	}
}
