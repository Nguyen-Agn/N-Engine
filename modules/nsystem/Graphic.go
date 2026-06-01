package nsystem

import (
	"image"
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

// BuildDrawOptions chuyển đổi SpriteData và PositionData thành mảng DrawOptions sẵn sàng vẽ.
// camX, camY là toạ độ camera trong map space — được trừ khỏi vị trí entity để chuyển sang screen space.
// Trả về nil nếu không có sprite hợp lệ để vẽ.
func BuildDrawOptions(pos PositionData, spr SpriteData, camX, camY float32) []*DrawOptions {
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

	if spr.IsNineSlice {
		return buildNineSliceDrawOptions(img, pos, spr, camX, camY)
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

	return []*DrawOptions{
		{
			Image: img,
			Opts:  opts,
		},
	}
}

func buildNineSliceDrawOptions(img *ebiten.Image, pos PositionData, spr SpriteData, camX, camY float32) []*DrawOptions {
	var result []*DrawOptions

	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()

	top := spr.NineSlice.Top
	bottom := spr.NineSlice.Bottom
	left := spr.NineSlice.Left
	right := spr.NineSlice.Right

	// Kiểm tra an toàn cho các slice
	if left+right > srcW {
		left = srcW / 2
		right = srcW - left
	}
	if top+bottom > srcH {
		top = srcH / 2
		bottom = srcH - top
	}

	targetW := float64(srcW) * float64(spr.ScaleX)
	targetH := float64(srcH) * float64(spr.ScaleY)

	srcX := []int{0, left, srcW - right, srcW}
	srcY := []int{0, top, srcH - bottom, srcH}

	midW := targetW - float64(left+right)
	if midW < 0 {
		midW = 0
	}
	midH := targetH - float64(top+bottom)
	if midH < 0 {
		midH = 0
	}

	dstX := []float64{0, float64(left), float64(left) + midW, targetW}
	dstY := []float64{0, float64(top), float64(top) + midH, targetH}

	finalX := float64(pos.X) + float64(spr.OffsetX) - float64(camX)
	finalY := float64(pos.Y) + float64(spr.OffsetY) - float64(camY)

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			sw := srcX[col+1] - srcX[col]
			sh := srcY[row+1] - srcY[row]

			if sw <= 0 || sh <= 0 {
				continue
			}

			dw := dstX[col+1] - dstX[col]
			dh := dstY[row+1] - dstY[row]

			if dw <= 0 || dh <= 0 {
				continue
			}

			rect := image.Rect(srcX[col], srcY[row], srcX[col+1], srcY[row+1])
			subImg := img.SubImage(rect).(*ebiten.Image)

			opts := &ebiten.DrawImageOptions{}

			// Scale patch
			scaleX := dw / float64(sw)
			scaleY := dh / float64(sh)
			opts.GeoM.Scale(scaleX, scaleY)

			// Dịch patch về đúng tọa độ trong khung (0, 0)
			opts.GeoM.Translate(dstX[col], dstY[row])

			// Dịch tâm để có thể xoay quanh giữa khung
			opts.GeoM.Translate(-targetW/2, -targetH/2)

			// Xoay toàn khối
			opts.GeoM.Rotate(float64(spr.Rotation) * math.Pi / 180)

			// Dịch đến vị trí cuối
			opts.GeoM.Translate(finalX, finalY)

			// Áp dụng màu sắc
			opts.ColorScale.ScaleWithColor(toFloat32Color(spr.ImageColor))

			result = append(result, &DrawOptions{
				Image: subImg,
				Opts:  opts,
			})
		}
	}

	return result
}

// toFloat32Color chuyển color.RGBA (0-255) sang ebiten.ColorScale (0.0-1.0)
func toFloat32Color(c color.RGBA) color.RGBA {
	return c // ebiten.ColorScale.ScaleWithColor nhận color.RGBA trực tiếp
}
