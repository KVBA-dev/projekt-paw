package handlers

import (
	"projekt-paw/data"
	"projekt-paw/views"
)

func ShowList(state *data.State, page uint) error {
	offset := (page - 1) * 10
	shownGames := state.Games[offset:]
	if len(shownGames) >= 10 {
		shownGames = shownGames[:10]
	} else {
		shownGames = shownGames[:len(shownGames)]
	}
	return RenderView(200, state.Ctx, views.List("", shownGames))
}

func CreateGame(state *data.State) error {
	var game data.Game
	game.Started = false
	state.Games = append(state.Games, game)
	return state.Ctx.HTML(200, "dupa")
}
