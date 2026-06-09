package bridge

import "github.com/hajimehoshi/ebiten/v2"

// SpriteLW is a lightweight wrapper around ebiten.Image frames.
type SpriteLW struct {
	images        []*ebiten.Image
	width, height int
}

// Purpose: Constructs and initializes a new empty SpriteLW.
// Inputs:
//   - width int: The width of the frames.
//   - height int: The height of the frames.
// Outputs: *SpriteLW - The newly created sprite wrapper.
func NewSpriteLW(width, height int) *SpriteLW {
	return &SpriteLW{
		images: make([]*ebiten.Image, 0),
		width:  width,
		height: height,
	}
}

// return value of image (return nil if not exist)
// Purpose: Retrieves the underlying ebiten.Image frame at the specified index.
// Inputs: index int - The zero-based frame index.
// Outputs: *ebiten.Image - The frame image, or nil if the index is invalid.
func (this *SpriteLW) Image(index int) *ebiten.Image {
	if index < 0 || index >= len(this.images) {
		return nil
	}
	return this.images[index]
}

// return number of images (return 0 if not exist)
// Purpose: Retrieves the total number of frames in this sprite sequence.
// Inputs: None.
// Outputs: int - The frame count.
func (this *SpriteLW) Length() int {
	return len(this.images)
}

// return value of width (return 0 if not exist)
// Purpose: Retrieves the width of the sprite frames.
// Inputs: None.
// Outputs: int - Width in pixels.
func (this *SpriteLW) Width() int {
	return this.width
}

// return value of height (return 0 if not exist)
// Purpose: Retrieves the height of the sprite frames.
// Inputs: None.
// Outputs: int - Height in pixels.
func (this *SpriteLW) Height() int {
	return this.height
}

// add image (if exist, do nothing)
// Purpose: Appends a new frame image to the sprite sequence.
// Inputs: image *ebiten.Image - The image to append.
// Outputs: None.
func (this *SpriteLW) AddImage(image *ebiten.Image) {
	this.images = append(this.images, image)
}

// remove image (if not exist, do nothing)
// Purpose: Removes a frame image from the sequence by index.
// Inputs: index int - The index to remove.
// Outputs: None.
func (this *SpriteLW) RemoveImage(index int) {
	if index < 0 || index >= len(this.images) {
		return
	}
	this.images = append(this.images[:index], this.images[index+1:]...)
}
