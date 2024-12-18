package data

import "github.com/labstack/echo/v4"

type User struct {
	Name string
	Id   int64
}

type Game struct {
	Player1          int64
	Player2          int64
	Player1Name      string
	Player2Name      string
	Player1Character int8
	Player2Character int8
	Turn             bool
	Started          bool
}

type State struct {
	Games []Game
	Ctx   echo.Context
}