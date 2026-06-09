package nsprite

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite LightWeight
type SpriteLW struct {
	images        []*ebiten.Image
	width, height int
}

// Image returns the ebiten.Image at the given index.
// Inputs: index - position of the image in the sprite sequence.
// Outputs: pointer to the ebiten.Image, or nil if the index is out of bounds.
func (this *SpriteLW) Image(index int) *ebiten.Image {
	if index < 0 || index >= len(this.images) {
		return nil
	}
	return this.images[index]
}

// Length returns the total number of images in this sprite.
// Outputs: integer count of stored images.
func (this *SpriteLW) Length() int {
	return len(this.images)
}

// Width returns the base width of the sprite.
// Outputs: width in pixels.
func (this *SpriteLW) Width() int {
	return this.width
}

// Height returns the base height of the sprite.
// Outputs: height in pixels.
func (this *SpriteLW) Height() int {
	return this.height
}

// AddImage appends a new image to the end of the sprite's sequence.
// Inputs: image - pointer to the ebiten.Image to add.
func (this *SpriteLW) AddImage(image *ebiten.Image) {
	this.images = append(this.images, image)
}

// RemoveImage removes an image from the sprite at the specified index.
// Inputs: index - position of the image to remove. Does nothing if out of bounds.
func (this *SpriteLW) RemoveImage(index int) {
	if index < 0 || index >= len(this.images) {
		return
	}
	this.images = append(this.images[:index], this.images[index+1:]...)
}
