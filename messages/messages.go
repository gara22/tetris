package messages

import "github.com/gara22/tetris/entities"

type MoveMessage struct {
	Type      string `json:"type"`
	Direction string `json:"direction"`
}

type GameStateMessage struct {
	Grid         entities.Grid `json:"grid"`
	Level        int           `json:"level"`
	LinesCleared int           `json:"linesCleared"`
	Score        int           `json:"score"`
	IsGameOver   bool          `json:"isGameOver"`
	IsGamePaused bool          `json:"isGamePaused"`
	NextShape    string        `json:"nextShape"`
}

type IncomingMessage struct {
	Type string `json:"type"`
}
