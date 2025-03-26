package repository

import "github.com/gara22/tetris/game"

type Repository interface {
	SaveGame(game.PersistedGame) error
	GetGames() ([]game.PersistedGame, error)
	GetByID(id string) (game.PersistedGame, error)
}
