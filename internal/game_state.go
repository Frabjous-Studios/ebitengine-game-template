package internal

import "github.com/DrJosh9000/yarn"

type GameState struct {
	// CurrentNode is the current storylet node which the player is working on.
	CurrentNode string

	// Vars is a list of all yarn variables.
	Vars yarn.MapVariableStorage
}
