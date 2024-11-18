package app_service

import (
	"github.com/gara22/tetris/game"
	socket "github.com/gara22/tetris/websocket"
	"github.com/google/uuid"
)

type AppService struct {
	Games map[string]*game.TetrisGame
}

func NewAppService() *AppService {
	return &AppService{
		Games: make(map[string]*game.TetrisGame),
	}
}

func (a *AppService) NewGame() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	hub := socket.NewHub()

	go hub.Run()

	// TODO: defer hub.Close()

	gameOverChannel := make(chan bool)

	newGame := game.NewTetrisGame(hub, gameOverChannel)
	a.Games[id.String()] = &newGame

	newGame.StartGame()

	go func() {
		<-gameOverChannel
		delete(a.Games, id.String())
	}()

	return id.String(), nil
}
