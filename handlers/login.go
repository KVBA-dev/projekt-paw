package handlers

import (
	"github.com/labstack/echo/v4"
	"html"
	"projekt-paw/data"
	// "projekt-paw/views"
)

type LoginReply struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func Login(ctx echo.Context, state *data.State) error {
	userName := html.EscapeString(ctx.Request().FormValue("username"))
	ctx.Logger().Print("user logged in anonymously")
	reply := &LoginReply{
		Id:   69,
		Name: userName,
	}
	return ctx.JSON(200, reply)
}
