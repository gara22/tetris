package messages

import "github.com/gara22/tetris/entities"

type MoveMessage struct {
	Direction string `json:"direction"`
}

type GameStateMessage struct {
	Grid         entities.Grid `json:"grid"`
	Level        int           `json:"level"`
	LinesCleared int           `json:"linesCleared"`
	Score        int           `json:"score"`
	IsGameOver   bool          `json:"isGameOver"`
}
