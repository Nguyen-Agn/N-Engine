package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// EbitenGame is adapter between loop of Ebitengine -> N-engine
type EbitenGame struct {
	engine *Engine
}

// NewEbitenGame creates a new EbitenGame adapter instance.
//
// Purpose: Initializes the adapter that connects the underlying Ebitengine game loop with the N-engine.
//
// Inputs:
// - engine (*Engine): A pointer to the main N-engine instance.
//
// Outputs:
// - *EbitenGame: The constructed adapter ready to be run by Ebiten.
func NewEbitenGame(engine *Engine) *EbitenGame {
	return &EbitenGame{engine: engine}
}

// Update is called every frame to process game logic.
//
// Purpose: Drives the engine's logic loop by updating the InputManager first, followed by the SceneManager.
//
// Outputs:
// - error: Returns an error if the SceneManager encounters one, which may halt the game loop; otherwise, returns nil.
func (g *EbitenGame) Update() error {
	g.engine.Input.Update()
	return g.engine.Scene.Update()
}

// Draw is called every frame to render graphics.
//
// Purpose: Passes the rendering surface to the current scene's camera and instructs the scene to draw itself.
//
// Inputs:
// - screen (*ebiten.Image): The Ebitengine screen buffer to draw onto.
//
// Special Requirements:
// - Does nothing if there is no active scene.
func (g *EbitenGame) Draw(screen *ebiten.Image) {
	currentScene := g.engine.Scene.GetCurrentScene()
	if currentScene == nil {
		return
	}

	camera := currentScene.GetCamera()
	if camera != nil {
		camera.SetScreen(screen)
	}

	currentScene.Draw()
}

// Layout determines the logical screen size.
//
// Purpose: Called by Ebitengine when the window resizes to determine the logical dimensions of the game screen, delegating the decision to the SceneManager.
//
// Inputs:
// - outsideWidth (int): The current window width.
// - outsideHeight (int): The current window height.
//
// Outputs:
// - (int, int): The logical width and height of the game screen.
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Scene.Layout(outsideWidth, outsideHeight)
}
