package data

import (
	"database/sql"
	ws "github.com/gorilla/websocket"
	"sync"
)

type Player struct {
	Name string
	Id   int64
	Conn *ws.Conn
}

type Game struct {
	Id               string
	Player1          *Player
	Player2          *Player
	Player1Character int8
	Player2Character int8
	Turn             bool
	Started          bool
	Mutex            sync.Mutex
}

func (g *Game) Ready() bool {
	return g.Player1 != nil && g.Player2 != nil
}

func (g *Game) SendTo(player *Player, msg string) error {
	return player.Conn.WriteMessage(ws.TextMessage, []byte(msg))
}

func (g *Game) Broadcast(msg string) error {
	err := g.SendTo(g.Player1, msg)
	if err != nil {
		return err
	}
	return g.SendTo(g.Player2, msg)
}

type State struct {
	Games []*Game
	Db    *sql.DB
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
