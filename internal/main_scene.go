package internal

import (
	"github.com/DrJosh9000/yarn"
	"github.com/Frabjous-Studios/ebitengine-game-template/internal/debug"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"log"
)

type BaseScene struct {
}

// Player is the player character in the game
type Player struct {
	Coords image.Point
}

// Move moves the player upwards
func (p *Player) Move() {
	p.Coords.Y--
}

type MainScene struct {
	Player *Player
	Game   *Game

	State  *GameState
	Runner *DialogueRunner
}

func NewMainScene(g *Game) *MainScene {
	runner, err := NewDialogueRunner()
	if err != nil {
		log.Fatal(err)
	}
	return &MainScene{
		Player: &Player{Coords: image.Pt(g.Width/2, g.Height/2)},
		Game:   g,
		Runner: runner,
		State: &GameState{
			CurrentNode: "Start",
			Vars:        make(yarn.MapVariableStorage),
		},
	}
}

func (m *MainScene) Update() error {
	if !m.Runner.running {
		go m.startRunner()
	}

	// Movement controls
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		m.Player.Move()
	}

	// XXX: Write game logic here

	return nil
}

func (m *MainScene) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(m.Player.Coords.X),
		float64(m.Player.Coords.Y),
		20,
		20,
		color.White,
	)
}

func (m *MainScene) startRunner() {
	if err := m.Runner.Start(m.State); err != nil {
		debug.Printf("error starting runner: %v", err)
		return
	}
}
