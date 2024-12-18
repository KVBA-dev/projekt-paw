package handlers

import (
	"github.com/labstack/echo/v4"
	// "html"
	"projekt-paw/data"
	// "projekt-paw/views"
)

func Login(ctx echo.Context, state *data.State) error {
	//userName := html.EscapeString(state.Ctx.Request().FormValue("username"))

	reply := struct {
		id int
	}{
		id: 69,
	}
	return ctx.JSON(200, reply)
}
