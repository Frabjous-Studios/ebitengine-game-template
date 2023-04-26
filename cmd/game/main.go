// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"github.com/Frabjous-Studios/ebitengine-game-template/internal"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var err error
	gameWidth, gameHeight := 640, 480

	ebiten.SetWindowSize(gameWidth, gameHeight)
	ebiten.SetWindowTitle("ebitengine-game-template")

	game := &internal.Game{
		Width:  gameWidth,
		Height: gameHeight,
	}
	game.CurrScene, err = internal.NewMainMenuScene(game)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
