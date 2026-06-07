package core

import (
	"autoworld/domain"
	"autoworld/modules/nsystem"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Camera Ä‘á»‹nh nghÄ©a má»™t viewport nhÃ¬n vÃ o Physical Map.
// Camera chá»‹u trÃ¡ch nhiá»‡m: xÃ¡c Ä‘á»‹nh vÃ¹ng nhÃ¬n, follow target, vÃ  render qua DrawSystem.
//
// Tá»a Ä‘á»™ (x, y) lÃ  gÃ³c trÃªn-trÃ¡i cá»§a viewport trong **map space**.
// DrawSystem sáº½ trá»« (x, y) khá»i tá»a Ä‘á»™ entity Ä‘á»ƒ chuyá»ƒn sang screen space.
type Camera struct {
	x, y          float32 // vá»‹ trÃ­ camera trong map space
	width, height int     // kÃ­ch thÆ°á»›c viewport (pixel)
	target        IObject // optional: follow target (nil = khÃ´ng follow)
	drawSystem    IDrawSystem
	screen        *ebiten.Image
}

// NewCamera khá»Ÿi táº¡o Camera má»›i vá»›i kÃ­ch thÆ°á»›c viewport cho trÆ°á»›c.
// viewW, viewH thÆ°á»ng báº±ng kÃ­ch thÆ°á»›c cá»­a sá»• game.
func NewCamera(viewW, viewH int) *Camera {
	return &Camera{
		x:          0,
		y:          0,
		width:      viewW,
		height:     viewH,
		drawSystem: nsystem.NewDrawSystem(),
	}
}

// â”€â”€â”€ Getters / Setters â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// X tráº£ vá» tá»a Ä‘á»™ ngang cá»§a camera trong map space.
func (c *Camera) X() float32 { return c.x }

// Y tráº£ vá» tá»a Ä‘á»™ dá»c cá»§a camera trong map space.
func (c *Camera) Y() float32 { return c.y }

// SetX dá»‹ch chuyá»ƒn camera Ä‘áº¿n tá»a Ä‘á»™ ngang má»›i trong map space.
func (c *Camera) SetX(x float32) { c.x = x }

// SetY dá»‹ch chuyá»ƒn camera Ä‘áº¿n tá»a Ä‘á»™ dá»c má»›i trong map space.
func (c *Camera) SetY(y float32) { c.y = y }

// Width tráº£ vá» chiá»u rá»™ng viewport (pixel).
func (c *Camera) Width() int { return c.width }

// Height tráº£ vá» chiá»u cao viewport (pixel).
func (c *Camera) Height() int { return c.height }

// SetTarget Ä‘áº·t IObject lÃ m má»¥c tiÃªu Ä‘á»ƒ camera tá»± Ä‘á»™ng theo dÃµi má»—i frame.
// Truyá»n nil Ä‘á»ƒ táº¯t follow.
func (c *Camera) SetTarget(obj IObject) { c.target = obj }

// â”€â”€â”€ Update â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// Update cáº­p nháº­t vá»‹ trÃ­ camera má»—i frame.
// Náº¿u cÃ³ target, camera Ä‘Æ°á»£c Ä‘áº·t sao cho target náº±m á»Ÿ giá»¯a viewport.
// mapW, mapH dÃ¹ng Ä‘á»ƒ clamp camera trong biÃªn báº£n Ä‘á»“ (0 = khÃ´ng giá»›i háº¡n).
func (c *Camera) Update(mapW, mapH int) {
	if c.target == nil {
		return
	}

	// Láº¥y vá»‹ trÃ­ target tá»« ECS (cáº§n PositionData)
	entry := c.target.Entry()
	if entry == nil || !entry.HasComponent(Position) {
		return
	}
	posData := donburi.Get[PositionData](entry, Position)

	// Äáº·t camera sao cho target á»Ÿ giá»¯a viewport
	c.x = posData.X - float32(c.width)/2
	c.y = posData.Y - float32(c.height)/2

	// Clamp trong biÃªn báº£n Ä‘á»“
	if mapW > 0 {
		if c.x < 0 {
			c.x = 0
		}
		if c.x+float32(c.width) > float32(mapW) {
			c.x = float32(mapW - c.width)
		}
	}
	if mapH > 0 {
		if c.y < 0 {
			c.y = 0
		}
		if c.y+float32(c.height) > float32(mapH) {
			c.y = float32(mapH - c.height)
		}
	}
}

// â”€â”€â”€ Render â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

// SetScreen thiáº¿t láº­p canvas Ä‘Ã­ch (ebiten.Image) Ä‘á»ƒ render lÃªn.
// Pháº£i gá»i má»—i frame trÆ°á»›c Draw() â€” thá»±c hiá»‡n bá»Ÿi EbitenGame.
func (c *Camera) SetScreen(screen *ebiten.Image) {
	c.screen = screen
	c.drawSystem.SetScreen(screen)
}

// Draw render scene lÃªn mÃ n hÃ¬nh theo thá»© tá»±:
//  1. Physical Map (physWorld) vá»›i camera offset (camX, camY) â€” entity ngoÃ i viewport bá»‹ skip.
//  2. GUI Map (guiWorld) vá»›i offset 0,0 â€” luÃ´n Ä‘Ã¨ lÃªn trÃªn, khÃ´ng bá»‹ culling camera.
//
// guiWorld cÃ³ thá»ƒ nil náº¿u Scene khÃ´ng cÃ³ GUI Map.
func (c *Camera) Draw(physWorld donburi.World, guiWorld donburi.World) {
	if c.screen == nil {
		return
	}
	// Váº½ Physical Map vá»›i camera offset
	c.drawSystem.Draw(physWorld, c.x, c.y)

	// Váº½ GUI Map khÃ´ng cÃ³ camera offset (screen space)
	if guiWorld != nil {
		c.drawSystem.Draw(guiWorld, 0, 0)
	}
}



// GetDrawSystem returns the DrawSystem as IDrawObjectRegistry so Scene can wire it to Map.
// This enables Map.AddObject to auto-register IDraw objects without importing nsystem directly.
func (c *Camera) GetDrawSystem() domain.IDrawObjectRegistry {
	if reg, ok := c.drawSystem.(domain.IDrawObjectRegistry); ok {
		return reg
	}
	return nil
}
