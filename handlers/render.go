package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderView(status int, ctx echo.Context, comp templ.Component) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx.Response().Status = status
	return comp.Render(ctx.Request().Context(), ctx.Response().Writer)
}
