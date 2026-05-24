package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// EbitenGame là adapter giữa vòng lặp Ebitengine và Engine của AutoWorld.
// Implement interface ebiten.Game — truyền vào ebiten.RunGame() để bắt đầu game.
//
// Không còn cần type-assert sang *Scene:
// IScene expose GetCamera() → EbitenGame gọi camera.SetScreen(screen) trực tiếp.
type EbitenGame struct {
	engine *Engine
}

// NewEbitenGame tạo EbitenGame bọc quanh Engine đã khởi tạo.
func NewEbitenGame(engine *Engine) *EbitenGame {
	return &EbitenGame{engine: engine}
}

// Update được Ebitengine gọi mỗi frame để xử lý logic.
// Cập nhật InputManager trước, sau đó ủy quyền cho SceneManager.
func (g *EbitenGame) Update() error {
	g.engine.Input.Update()
	return g.engine.Scene.Update()
}

// Draw được Ebitengine gọi mỗi frame sau Update để render.
// Set screen vào Camera của Scene hiện tại, sau đó gọi Draw().
func (g *EbitenGame) Draw(screen *ebiten.Image) {
	currentScene := g.engine.Scene.GetCurrentScene()
	if currentScene == nil {
		return
	}
	// Truyền screen vào Camera — không cần type-assert sang *Scene
	camera := currentScene.GetCamera()
	if camera != nil {
		camera.SetScreen(screen)
	}
	currentScene.Draw()
}

// Layout được Ebitengine gọi để lấy kích thước màn hình logic.
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Scene.Layout(outsideWidth, outsideHeight)
}
