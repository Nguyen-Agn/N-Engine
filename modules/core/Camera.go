package core

import (
	"autoworld/domain"
	"autoworld/modules/nmath"
	"autoworld/modules/nsystem"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Camera
type Camera struct {
	x, y          float32
	width, height int
	target        IObject
	drawSystem    IDrawSystem
	screen        *ebiten.Image
}

// NewCamera creates a new Camera instance with the specified view dimensions.
//
// Purpose: Initializes a camera that determines what part of the game world is visible.
//
// Inputs:
// - viewW (int): The width of the camera's viewport in pixels.
// - viewH (int): The height of the camera's viewport in pixels.
//
// Outputs:
// - *Camera: A pointer to the newly created Camera instance.
func NewCamera(viewW, viewH int) *Camera {
	return &Camera{
		x:          0,
		y:          0,
		width:      viewW,
		height:     viewH,
		drawSystem: nsystem.NewDrawSystem(),
	}
}

// X returns the current X coordinate of the camera.
//
// Purpose: Retrieves the horizontal position of the top-left corner of the camera in the game world.
//
// Outputs:
// - float32: The X coordinate.
func (c *Camera) X() float32 { return c.x }

// Y returns the current Y coordinate of the camera.
//
// Purpose: Retrieves the vertical position of the top-left corner of the camera in the game world.
//
// Outputs:
// - float32: The Y coordinate.
func (c *Camera) Y() float32 { return c.y }

// SetX updates the X coordinate of the camera.
//
// Purpose: Manually sets the horizontal position of the camera.
//
// Inputs:
// - x (float32): The new X coordinate.
func (c *Camera) SetX(x float32) { c.x = x }

// SetY updates the Y coordinate of the camera.
//
// Purpose: Manually sets the vertical position of the camera.
//
// Inputs:
// - y (float32): The new Y coordinate.
func (c *Camera) SetY(y float32) { c.y = y }

// Width returns the width of the camera's viewport.
//
// Purpose: Retrieves the viewport width.
//
// Outputs:
// - int: The width in pixels.
func (c *Camera) Width() int { return c.width }

// Height returns the height of the camera's viewport.
//
// Purpose: Retrieves the viewport height.
//
// Outputs:
// - int: The height in pixels.
func (c *Camera) Height() int { return c.height }

// SetTarget assigns an object for the camera to follow.
//
// Purpose: Binds the camera to a specific game object so it centers on the object during updates.
//
// Inputs:
// - obj (IObject): The target object to follow. Must have a Position component to be tracked effectively.
func (c *Camera) SetTarget(obj IObject) { c.target = obj }

// Update adjusts the camera's position to follow its target, keeping it within map bounds.
//
// Purpose: Centers the camera on the target object every frame, ensuring the camera does not pan outside the map boundaries.
//
// Inputs:
// - mapW (int): The total width of the map in pixels.
// - mapH (int): The total height of the map in pixels.
//
// Special Requirements:
// - If no target is set, or the target lacks a Position component, the camera does nothing.
func (c *Camera) Update(mapW, mapH int) {
	if c.target == nil {
		return
	}

	entry := c.target.Entry()
	if entry == nil || !entry.HasComponent(Position) {
		return
	}
	targetData := donburi.Get[PositionData](entry, Position)

	// New raw position
	c.x = targetData.X - float32(c.width)/2
	c.y = targetData.Y - float32(c.height)/2

	// Final position
	c.x = nmath.GetInstance().Clamp(c.x, 0, float32(mapW-c.width))
	c.y = nmath.GetInstance().Clamp(c.y, 0, float32(mapH-c.height))
}

// Draw renders both the physical world and the GUI world onto the screen.
//
// Purpose: Renders entities using the DrawSystem. It applies a viewport offset for the physical world and renders the GUI world statically.
//
// Inputs:
// - physWorld (donburi.World): The ECS world containing game map entities.
// - guiWorld (donburi.World): The ECS world containing GUI entities.
//
// Special Requirements:
// - If the screen has not been set via SetScreen, the method does nothing.
func (c *Camera) Draw(physWorld donburi.World, guiWorld donburi.World) {
	if c.screen == nil {
		return
	}

	c.drawSystem.Draw(physWorld, c.x, c.y)

	if guiWorld != nil {
		c.drawSystem.Draw(guiWorld, 0, 0)
	}
}

// GetDrawSystem retrieves the registry interface for drawing objects.
//
// Purpose: Enables adding drawable objects to the DrawSystem without requiring external modules to import the nsystem directly.
//
// Outputs:
// - domain.IDrawObjectRegistry: The drawing registry interface, or nil if not implemented.
func (c *Camera) GetDrawSystem() domain.IDrawObjectRegistry {
	if reg, ok := c.drawSystem.(domain.IDrawObjectRegistry); ok {
		return reg
	}
	return nil
}

// SetScreen assigns the Ebiten screen image for the camera's draw system.
//
// Purpose: Updates the render target for the camera and its associated draw system.
//
// Inputs:
// - screen (*ebiten.Image): The new screen buffer to render graphics onto.
func (c *Camera) SetScreen(screen *ebiten.Image) {
	c.screen = screen
	c.drawSystem.SetScreen(screen)
}
