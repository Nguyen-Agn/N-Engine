package components

import (
	"image/color"
	"strings"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Debug = enginetype.Debug

// init initializes the default data for the debug component.
// It registers the "deb" component token with default display flags (box, pos, info) and white color.
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

// DebugComponent is a Mixin to embed into Custom Objects (using the "deb" token).
// It provides functionality to draw debug overlays for the object (hitbox, origin point, ID/Name info)
// as well as a Log(msg) function to display arbitrary text next to the object on screen.
type DebugComponent struct {
	IObject
	data *DebugData
}

// BindComponent binds the base object and its ECS data to this component.
// Inputs:
//   - base: The base IObject to bind to.
func (c *DebugComponent) BindComponent(base IObject) {
	c.IObject = base
	c.data = enginetype.GetComponent(base, Debug)
}

// SetDebugColor changes the color of the debug overlay graphics.
// Inputs:
//   - col: The new color.RGBA to use (default is white).
func (c *DebugComponent) SetDebugColor(col color.RGBA) {
	if c.data != nil {
		c.data.Color = col
	}
}

// Debug configures the display elements on the debug overlay using a space-separated string of flags.
// Inputs:
//   - flags: A string containing options:
//     "box": draws the collision bounding box.
//     "pos": draws the origin position.
//     "info": prints the [ID] and Name of the object.
//     For example, calling Debug("pos box") will disable info and only draw pos and box.
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

// Log records a text string to display near the object on screen (used for debugging).
// Purpose: Provides on-screen logging. If the object lacks a position component, the text
// will be drawn statically at the top-left of the screen and automatically stacked to avoid overlapping.
// Inputs:
//   - msg: The string message to display.
func (c *DebugComponent) Log(msg string) {
	if c.data != nil {
		c.data.CustomLog = msg
	}
}
