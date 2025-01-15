package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"projekt-paw/data"
	"projekt-paw/handlers"
	"projekt-paw/views"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var schema []byte
	if schema, err = os.ReadFile("env/schema.sql"); err != nil {
		fmt.Println("error: could not find schema.sql")
		return
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		fmt.Println("error on creating database:", err)
		return
	}

	state := &data.State{Games: make([]*data.Game, 0), Db: db, SessionId: 0}

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
	e.GET("/gamelist", func(c echo.Context) error {
		return handlers.GetList(c, state)
	})

	e.POST("/login", func(c echo.Context) error {
		return handlers.Login(c, state)
	})
	e.POST("/loginanon", func(c echo.Context) error {
		return handlers.LoginAnonymous(c, state)
	})
	e.POST("/register", func(c echo.Context) error {
		return handlers.Register(c, state)
	})
	e.POST("/create", func(c echo.Context) error {
		return handlers.CreateGame(c, state)
	})

	//-------------------------------------

	e.Logger.Info(e.Start(":8080"))
}
