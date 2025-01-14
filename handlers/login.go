package handlers

import (
	"database/sql"
	"fmt"
	"projekt-paw/data"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	// "projekt-paw/views"
)

func LoginAnonymous(ctx echo.Context, state *data.State) error {
	userName := ctx.QueryParam("username")
	ctx.Logger().Print(fmt.Sprintf("user %s logged in anonymously", userName))
	reply := &data.UserData{
		UserId:    -1,
		SessionId: state.SessionId,
		Name:      userName,
	}
	state.SessionId += 1
	return ctx.JSON(200, reply)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx echo.Context, state *data.State) error {
	var req LoginRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(422)
	}

	if req.Password == "" || req.Username == "" {
		return ctx.NoContent(422)
	}

	var username, storedHash string
	var userId int64
	err = state.Db.QueryRow("SELECT id, name, passhash FROM users WHERE name = ?", req.Username).Scan(&userId, &username, &storedHash)
	if err == sql.ErrNoRows {
		return ctx.NoContent(403)
	} else if err != nil {
		return ctx.NoContent(500)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
	if err != nil {
		return ctx.NoContent(403)
	}

	reply := &data.UserData{
		Name:      req.Username,
		SessionId: state.SessionId,
		UserId:    userId,
	}
	state.SessionId += 1
	return ctx.JSON(200, reply)
}

func Register(ctx echo.Context, state *data.State) error {
	var req LoginRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.NoContent(422)
	}

	if req.Password == "" || req.Username == "" {
		return ctx.NoContent(422)
	}

	var exists bool
	err = state.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)", req.Username).Scan(&exists)
	if err != nil {
		return ctx.NoContent(500)
	}
	if exists {
		return ctx.NoContent(409)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.NoContent(500)
	}

	result, err := state.Db.Exec("INSERT INTO users (name, passhash) VALUES (?, ?)", req.Username, hash)
	if err != nil {
		return ctx.NoContent(500)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return ctx.NoContent(500)
	}

	reply := &data.UserData{
		UserId:    id,
		SessionId: state.SessionId,
		Name:      req.Username,
	}
	state.SessionId += 1

	return ctx.JSON(200, reply)
}
