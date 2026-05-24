package nsystem

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawOptions là kết quả đã được tính toán sẵn, sẵn sàng để vẽ.
// Graphic.go chịu trách nhiệm tạo ra struct này từ SpriteData + PositionData.
type DrawOptions struct {
	// Image là frame ảnh cụ thể sẽ được vẽ
	Image *ebiten.Image
	// Opts là tất cả các thông số biến đổi hình học và màu sắc
	Opts *ebiten.DrawImageOptions
}

// BuildDrawOptions chuyển đổi SpriteData và PositionData thành DrawOptions sẵn sàng vẽ.
// camX, camY là toạ độ camera trong map space — được trừ khỏi vị trí entity để chuyển sang screen space.
// Trả về nil nếu không có sprite hợp lệ để vẽ.
func BuildDrawOptions(pos PositionData, spr SpriteData, camX, camY float32) *DrawOptions {
	// Lấy sprite đang được chọn
	spriteLW, ok := spr.Sprite[spr.CurrentSprite]
	if !ok || spriteLW == nil {
		return nil
	}

	// Lấy frame ảnh hiện tại theo chỉ số
	img := spriteLW.Image(spr.SpriteIdx)
	if img == nil {
		return nil
	}

	opts := &ebiten.DrawImageOptions{}

	// 1. Dịch về gốc tọa độ (0,0) để phép xoay và scale hoạt động đúng
	w := float64(spriteLW.Width())
	h := float64(spriteLW.Height())
	opts.GeoM.Translate(-w/2, -h/2)

	// 2. Áp dụng Scale (ScaleX âm = lật ngang, ScaleY âm = lật dọc)
	opts.GeoM.Scale(float64(spr.ScaleX), float64(spr.ScaleY))

	// 3. Áp dụng Rotation (chuyển từ độ sang radian)
	opts.GeoM.Rotate(float64(spr.Rotation) * math.Pi / 180)

	// 4. Dịch về vị trí cuối cùng = (map pos + offset) - camera offset = screen pos
	finalX := float64(pos.X) + float64(spr.OffsetX) - float64(camX)
	finalY := float64(pos.Y) + float64(spr.OffsetY) - float64(camY)
	opts.GeoM.Translate(finalX, finalY)

	// 5. Áp dụng màu sắc (Color Scale)
	opts.ColorScale.ScaleWithColor(toFloat32Color(spr.ImageColor))

	return &DrawOptions{
		Image: img,
		Opts:  opts,
	}
}

// toFloat32Color chuyển color.RGBA (0-255) sang ebiten.ColorScale (0.0-1.0)
func toFloat32Color(c color.RGBA) color.RGBA {
	return c // ebiten.ColorScale.ScaleWithColor nhận color.RGBA trực tiếp
}
