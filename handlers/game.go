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

const delim string = "»"

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

func GetList(ctx echo.Context, state *data.State) error {
	return RenderView(200, ctx, views.RenderList(state.Games))
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
		Conn:      conn,
		Character: -1,
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
		fmt.Println("Message received:", mesg)
		switch mesg[0] {
		// TODO: change delimiter from . to a rarer character
		case 'd':
			/* INFO: d - user sends data upon opening connection
			   format is d.name.id
			*/
			tokens := strings.Split(mesg, delim)
			player.Name = tokens[1]
			id, err := strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				game.Broadcast(fmt.Sprintf("x%sinvalid data - game closed", delim))
			}
			player.Id = id
			err = game.LogBroadcast(fmt.Sprintf("Joined: %s", player.Name))

			if game.Player2 == nil {
				err = game.SendTo(game.Player1, "w")
				err = game.LogTo(game.Player1, "Waiting for other player...")
			} else {
				err = game.LogTo(player, fmt.Sprintf("Your opponent is %s", game.Player1.Name))
				err = game.Broadcast("s")
				err = game.LogBroadcast("Select your flag")
				game.Started = true
			}
			if err != nil {
				break
			}
		case 's':
			/* INFO: s - user selected the character
			   format is s.number
			*/
			tokens := strings.Split(mesg, delim)
			c, err := strconv.ParseInt(tokens[1], 10, 8)
			if err != nil {
				break
			}
			player.Character = int8(c)
			game.Mutex.Lock()
			if game.Player1.Character >= 0 && game.Player2.Character >= 0 {
				game.Begin()
			}
			game.Mutex.Unlock()

		case 'a', 'r':
			/* INFO: a - user asks
			   format is a.question
			*/

			/* INFO: r - user replies
			   format is r.y if yes or r.n if no
			*/
			var recv *data.Player
			if player == game.Player1 {
				recv = game.Player2
			} else {
				recv = game.Player1
			}
			game.SendTo(recv, mesg)
			if mesg[0] == 'r' {
				a := "Yes"
				if []rune(mesg)[2] == 'n' {
					a = "No"
				}
				game.LogBroadcast(fmt.Sprintf("Answer: %s", a))
				game.SwitchTurn()
			} else {
				a := strings.Split(mesg, delim)[1]
				game.LogBroadcast(fmt.Sprintf("%s asks: %s", player.Name, a))
			}
		case 'g':
			/* INFO: user guesses
			   format is g.guess
			*/
			tokens := strings.Split(mesg, delim)
			c, err := strconv.ParseInt(tokens[1], 10, 8)
			if err != nil {
				break
			}
			var recv *data.Player
			if player == game.Player1 {
				recv = game.Player2
			} else {
				recv = game.Player1
			}
			if recv.Character == int8(c) {
				game.SendTo(player, fmt.Sprintf("r%sy", delim))
				game.LogBroadcast(fmt.Sprintf("Correct! %s wins!", player.Name))
				game.Broadcast("q")
			} else {
				game.SendTo(player, fmt.Sprintf("r%sn", delim))
				game.LogBroadcast("Wrong guess!")
				game.SwitchTurn()
			}
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
