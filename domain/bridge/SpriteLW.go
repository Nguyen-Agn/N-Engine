package bridge

import "github.com/hajimehoshi/ebiten/v2"

// SpriteLW is a lightweight wrapper around ebiten.Image frames.
type SpriteLW struct {
	images        []*ebiten.Image
	width, height int
}

func NewSpriteLW(width, height int) *SpriteLW {
	return &SpriteLW{
		images: make([]*ebiten.Image, 0),
		width:  width,
		height: height,
	}
}

// return value of image (return nil if not exist)
func (this *SpriteLW) Image(index int) *ebiten.Image {
	if index < 0 || index >= len(this.images) {
		return nil
	}
	return this.images[index]
}

// return number of images (return 0 if not exist)
func (this *SpriteLW) Length() int {
	return len(this.images)
}

// return value of width (return 0 if not exist)
func (this *SpriteLW) Width() int {
	return this.width
}

// return value of height (return 0 if not exist)
func (this *SpriteLW) Height() int {
	return this.height
}

// add image (if exist, do nothing)
func (this *SpriteLW) AddImage(image *ebiten.Image) {
	this.images = append(this.images, image)
}

// remove image (if not exist, do nothing)
func (this *SpriteLW) RemoveImage(index int) {
	if index < 0 || index >= len(this.images) {
		return
	}
	this.images = append(this.images[:index], this.images[index+1:]...)
}
