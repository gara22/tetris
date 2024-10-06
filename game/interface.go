package game

import "github.com/gara22/tetris/entities"

type Game interface {
	StartGame()
	Move(params MoveParams) (TetrisGame, error)
	GetState() entities.Grid
}

type MoveParams struct {
	Direction string
}
