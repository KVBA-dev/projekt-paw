package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"projekt-paw/data"
	"projekt-paw/handlers"
	"projekt-paw/views"
)

func main() {
	fmt.Println("Opening database...")
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	fmt.Println("success!")
	defer db.Close()

	state := &data.State{Games: make([]*data.Game, 0), Db: db}

	e := echo.New()
	e.Static("/css", "css")
	e.Static("/img", "img")
	e.Static("/js", "js")
	//define your routes and endpoints here
	e.GET("/", func(c echo.Context) error {
		return handlers.RenderView(200, c, views.Index())
	})
	e.GET("/list", func(c echo.Context) error {
		return handlers.ShowList(c, state, 1)
	})
	e.GET("/game", func(c echo.Context) error {
		return handlers.GetGamePage(c, state)
	})
	e.GET("/gamews", func(c echo.Context) error {
		return handlers.HandleWS(c, state)
	})

	e.POST("/login", func(c echo.Context) error {
		return handlers.LoginAnonymous(c, state)
	})
	e.POST("/loginanon", func(c echo.Context) error {
		return handlers.LoginAnonymous(c, state)
	})
	e.POST("/create", func(c echo.Context) error {
		return handlers.CreateGame(c, state)
	})

	//-------------------------------------

	e.Logger.Info(e.Start(":8080"))
}
