package game

type Game interface {
	StartGame() TetrisGame
	Move(params MoveParams)
}

type MoveParams struct {
	Direction string
}
