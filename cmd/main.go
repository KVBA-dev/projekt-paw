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

	state := &data.State{Games: make([]data.Game, 0), Ctx: nil}

	e := echo.New()
	e.Static("/css", "css")
	e.Static("/img", "img")
	e.Static("/js", "js")
	//define your routes and endpoints here
	e.GET("/", func(c echo.Context) error {
		return handlers.RenderView(200, c, views.Index())
	})

	e.POST("/login", func(c echo.Context) error {
		state.Ctx = c
		return handlers.Login(c, state)
	})
	e.POST("loginanon", func(c echo.Context) error {
		state.Ctx = c
		return handlers.Login(c, state)
	})

	e.POST("/create/:name", func(c echo.Context) error {
		state.Ctx = c
		return handlers.CreateGame(state)
	})

	e.GET("/list", func(c echo.Context) error { return handlers.ShowList(state, 1) })
	//-------------------------------------

	e.Logger.Fatal(e.Start(":8080"))
}
