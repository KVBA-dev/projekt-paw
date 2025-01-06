package handlers

import (
	"fmt"
	"net/http"
	"projekt-paw/data"
	"projekt-paw/views"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ShowList(ctx echo.Context, state *data.State, page uint) error {
	offset := (page - 1) * 10
	shownGames := state.Games[offset:]
	if len(shownGames) >= 10 {
		shownGames = shownGames[:10]
	}
	return RenderView(200, ctx, views.List("", shownGames))
}

func CreateGame(ctx echo.Context, state *data.State) error {
	game := &data.Game{
		Id:      uuid.New().String(),
		Started: false,
		Player1: nil,
		Player2: nil,
	}
	state.Games = append(state.Games, game)
	ctx.Logger().Print("game created - id:", game.Id)
	return ctx.String(200, game.Id)
}

func GetGamePage(ctx echo.Context, state *data.State) error {
	gameId := ctx.QueryParam("id")

	var game *data.Game = nil

	for _, g := range state.Games {
		if g.Id == gameId {
			game = g
			break
		}
	}

	if game == nil {
		return ctx.String(404, "a game with given id was not found")
	}

	return RenderView(200, ctx, views.Game(game))
}

func HandleWS(ctx echo.Context, state *data.State) error {
	gameid := ctx.QueryParam("game")
	fmt.Println(gameid)
	var game *data.Game = nil
	for _, g := range state.Games {
		if g.Id == gameid {
			game = g
			break
		}
	}

	if game == nil {
		return ctx.NoContent(404)
	}

	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	player := &data.Player{
		Conn: conn,
	}

	game.Mutex.Lock()
	if game.Player1 == nil {
		game.Player1 = player
		game.Mutex.Unlock()
	} else if game.Player2 == nil {
		game.Player2 = player
		game.Mutex.Unlock()
	} else {
		game.Mutex.Unlock()
		conn.Close()
		return ctx.NoContent(422)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		mesg := string(msg[:])
		switch mesg[0] {
		// TODO: change delimiter from . to a rarer character
		case 'd':
			/* INFO: user sends data upon opening connection
			   format is d.name.id
			*/
			tokens := strings.Split(mesg, ".")
			player.Name = tokens[1]
			id, err := strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				break
			}
			player.Id = id

			if game.Player2 == nil {
				err = game.SendTo(game.Player1, "w")
			} else {
				err = game.Broadcast("s")
			}
			if err != nil {
				break
			}
		case 'a':
			/* INFO: user asks another one
			   format is s.question
			*/
		case 'r':
			/* INFO: user replies
			   format is r.y if yes or r.n if no
			*/
		case 'g':
			/* INFO: user guesses
			   format is g.guess
			*/
		case 'e':
			/* INFO: user responds to a guess
			   format is alternative to r: e.y or e.n
			   if e.y is sent, the game ends
			*/
		default:
		}
	}
	if game.Player1 != nil {
		game.Player1.Conn.Close()
	}
	if game.Player2 != nil {
		game.Player2.Conn.Close()
	}

	fmt.Println("client disconnected - game terminated")
	state.DeleteGame(game)
	return nil
}
