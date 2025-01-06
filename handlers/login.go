package handlers

import (
	"fmt"
	"projekt-paw/data"

	"github.com/labstack/echo/v4"
	// "projekt-paw/views"
)

type LoginReply struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func LoginAnonymous(ctx echo.Context, state *data.State) error {
	userName := ctx.QueryParam("username")
	ctx.Logger().Print(fmt.Sprintf("user %s logged in anonymously", userName))
	reply := &LoginReply{
		Id:   69,
		Name: userName,
	}
	return ctx.JSON(200, reply)
}
