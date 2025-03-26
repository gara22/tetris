package app_service

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gara22/tetris/game"
	"github.com/gara22/tetris/repository"
	socket "github.com/gara22/tetris/websocket"
	"github.com/google/uuid"
)

type AppService struct {
	Games      map[string]*game.TetrisGame
	Repository repository.Repository
}

func NewAppService(repository repository.Repository) *AppService {
	return &AppService{
		Games:      make(map[string]*game.TetrisGame),
		Repository: repository,
	}
}

func (a *AppService) NewGame() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	hub := socket.NewHub()
	hub.ID = id.String()

	go hub.Run()

	// TODO: defer hub.Close()

	gameOverChannel := make(chan bool)

	newGame := game.NewTetrisGame(hub, gameOverChannel)
	a.Games[id.String()] = &newGame

	newGame.StartGame()

	go func() {
		<-gameOverChannel
		err := a.Repository.SaveGame(a.Games[id.String()].ToPersistedGame())
		if err != nil {
			spew.Dump(err)
		}
		delete(a.Games, id.String())
	}()

	return id.String(), nil
}

func (a *AppService) AddScore(gameId string, playerName string) error {
	retrievedGame, err := a.Repository.GetByID(gameId)
	if err != nil {
		return err
	}

	if retrievedGame.Player != "" {
		return fmt.Errorf("player already set for this game")
	}
	retrievedGame.Player = playerName

	err = a.Repository.SaveGame(retrievedGame)
	if err != nil {
		return err
	}

	highScore := game.HighScore{
		Player: playerName,
		Score:  retrievedGame.Score,
	}

	err = a.Repository.SaveHighScore(highScore)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppService) GetHighScores() ([]game.HighScore, error) {
	return a.Repository.GetHighScores()
}
