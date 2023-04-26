package internal

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"sync"
)

// TPS is the number of ticks per second, read once when the game starts.
var TPS float64
var TPSOnce sync.Once

// Game represents the main game state
type Game struct {
	Width     int
	Height    int
	CurrScene Scene
}

type Scene interface {
	Update() error
	Draw(*ebiten.Image)
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Width, g.Height
}

// Update calculates game logic
func (g *Game) Update() error {
	TPSOnce.Do(func() {
		ebiten.SetTPS(60)
		TPS = float64(ebiten.TPS())
	})

	// Pressing Q any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("game quit by player")
	}

	// Pressing F toggles full-screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	return g.CurrScene.Update()

}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrScene.Draw(screen)
}

// ChangeScene sets the current scene to the provided Scene.
func (g *Game) ChangeScene(s Scene) {
	g.CurrScene = s
}
