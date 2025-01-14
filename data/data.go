package data

import (
	"database/sql"
	"fmt"
	"sync"

	ws "github.com/gorilla/websocket"
)

type Player struct {
	Name      string
	Id        int64
	Conn      *ws.Conn
	Character int8
}

type UserData struct {
	Name      string `json:"username"`
	SessionId int64  `json:"session_id"`
	UserId    int64  `json:"user_id"`
}

type Game struct {
	Id      string
	Player1 *Player
	Player2 *Player
	Turn    bool
	Started bool
	Mutex   sync.Mutex
}

func (g *Game) Ready() bool {
	return g.Player1 != nil && g.Player2 != nil
}

func (g *Game) SendTo(player *Player, msg string) error {
	if player == nil {
		return nil
	}
	return player.Conn.WriteMessage(ws.TextMessage, []byte(msg))
}

func (g *Game) Broadcast(msg string) error {
	err := g.SendTo(g.Player1, msg)
	if err != nil {
		return err
	}
	return g.SendTo(g.Player2, msg)
}

func (g *Game) LogTo(player *Player, msg string) error {
	return g.SendTo(player, fmt.Sprintf("l»%s", msg))
}

func (g *Game) LogBroadcast(msg string) error {
	return g.Broadcast(fmt.Sprintf("l»%s", msg))
}

func (g *Game) Begin() {
	g.Started = true
	g.Turn = false
	g.Broadcast("o")
	g.LogBroadcast("Game starts!")
	g.SwitchTurn()
}

func (g *Game) SwitchTurn() {
	g.Turn = !g.Turn
	if g.Turn {
		g.SendTo(g.Player1, "t")
		g.LogTo(g.Player1, "Your turn")
	} else {
		g.SendTo(g.Player2, "t")
		g.LogTo(g.Player2, "Your turn")
	}
}

type State struct {
	Games     []*Game
	Db        *sql.DB
	SessionId int64
}

func (s *State) DeleteGame(game *Game) {
	numGames := len(s.Games)
	for i, g := range s.Games {
		if g == game {
			s.Games[i] = s.Games[numGames-1]
			s.Games[numGames-1] = nil
			s.Games = s.Games[:numGames-1]
		}
	}
}
