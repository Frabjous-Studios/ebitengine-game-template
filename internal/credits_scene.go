package internal

import (
	"github.com/Frabjous-Studios/ebitengine-game-template/internal/debug"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CreditsScene struct {
	Game *Game
	bg   *ebiten.Image
}

func NewCreditsScene(game *Game) (*CreditsScene, error) {
	return &CreditsScene{
		Game: game,
	}, nil
}

func (s *CreditsScene) Update() error {
	keys := inpututil.AppendPressedKeys(nil)

	if len(keys) > 0 || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		menu, err := NewMainMenuScene(s.Game)
		if err != nil {
			debug.Printf("error constructing main menu: %v", err)
			return err
		}
		s.Game.ChangeScene(menu)
	}
	return nil
}

func (s *CreditsScene) Draw(screen *ebiten.Image) {
	if s.bg != nil {
		screen.DrawImage(s.bg, nil)
	}
}
