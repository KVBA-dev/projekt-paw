package handlers

import (
	"github.com/labstack/echo/v4"
	"projekt-paw/data"
	"projekt-paw/views"
)

func ShowList(ctx echo.Context, state *data.State, page uint) error {
	offset := (page - 1) * 10
	shownGames := state.Games[offset:]
	if len(shownGames) >= 10 {
		shownGames = shownGames[:10]
	}
	return RenderView(200, ctx, views.List("", shownGames))
}

func CreateGame(ctx echo.Context, state *data.State) error {
	var game data.Game
	game.Started = false
	state.Games = append(state.Games, game)
	ctx.Logger().Print("game created")
	return ctx.HTML(200, "test")
}
