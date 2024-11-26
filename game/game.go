package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gara22/tetris/entities"
	"github.com/gara22/tetris/messages"
	socket "github.com/gara22/tetris/websocket"
)

type TetrisGame struct {
	Grid            entities.Grid
	ActiveShape     entities.Shape
	Ticker          *time.Ticker
	Hub             socket.Hub
	IsGameOver      bool
	GameOverChannel chan bool
}

var (
	// Make sure that TetrisGame implements the Game interface
	_ Game = (*TetrisGame)(nil)
)

const (
	HEIGHT = 21
	WIDTH  = 11
)

func NewTetrisGame(hub *socket.Hub, gameOverChannel chan bool) TetrisGame {
	return TetrisGame{
		Grid:            entities.NewGrid(WIDTH, HEIGHT),
		ActiveShape:     entities.Shape{},
		Ticker:          nil,
		Hub:             *hub,
		GameOverChannel: gameOverChannel,
	}
}

func (t *TetrisGame) StartGame() {
	fmt.Println("start game")
	t.ActiveShape = entities.GenerateRandomShape()
	// publish initial game state
	err := t.PublishGameState()
	if err != nil {
		fmt.Println("Error publishing game state")
	}
	// setup a ticker to move the shape down every second
	msg := t.Hub.ReadMessage()

	go func() {
		t.Ticker = time.NewTicker(1 * time.Second)
		for {
			select {
			case <-t.Ticker.C:
				// t.Move(MoveParams{Direction: "down"})

			case m := <-msg:
				var message messages.MoveMessage
				json.Unmarshal(m, &message)
				// spew.Dump(message)

				game, err := t.Move(MoveParams{Direction: message.Direction})
				if err != nil {
					fmt.Println("Error moving shape: %s", err)
					if err.Error() == "Game over" {
						t.EndGame()
						game.IsGameOver = true
					}
				}
				t = &game

				// game.Move()

				err = t.PublishGameState()
				if err != nil {
					fmt.Println("Error publishing game state")
				}
			}
		}
	}()
}

func (t TetrisGame) Move(params MoveParams) (TetrisGame, error) {
	newShape, err := t.calculateNewShape(params, t.ActiveShape)
	// spew.Dump(t.ActiveShape)
	if err != nil && err.Error() == "Shape is stuck" {
		fmt.Println("Shape is stuck")
		newShape = t.ActiveShape.Block()
		t.Grid.RenderShape(newShape)

		t.ActiveShape = entities.GenerateRandomShape()
		// TODO: do we need to render here?
		t.checkForFullRows()
		if t.isGameOver() {
			return t, errors.New("Game over")
		}
		t.Grid.RenderShape(t.ActiveShape)
		t.Grid.Print()
		return t, nil
		// spew.Dump(t.Grid)
	}
	if err != nil {
		return t, err
	}

	t.checkForFullRows()
	t.Grid.RenderShape(newShape)
	t.ActiveShape = newShape
	fmt.Println("assigning new shape")

	t.Grid.Print()
	return t, nil
}

func (t TetrisGame) calculateNewShape(params MoveParams, activeShape entities.Shape) (entities.Shape, error) {
	var newShape entities.Shape
	switch params.Direction {
	case "left":
		newShape = activeShape.Move("left", t.Grid)
	case "right":
		newShape = activeShape.Move("right", t.Grid)
	case "down":
		newShape := activeShape.Move("down", t.Grid)

		if t.isShapeColliding(newShape) {
			return entities.Shape{}, fmt.Errorf("Shape is stuck")
		}
		return newShape, nil
	case "up":
		newShape = activeShape.Rotate()
	default:
		return entities.Shape{}, fmt.Errorf("Invalid input")
	}
	if t.isShapeColliding(newShape) {
		return entities.Shape{}, fmt.Errorf("Cannot move shape")
	}
	return newShape, nil
}

func (t TetrisGame) isShapeColliding(newShape entities.Shape) bool {
	for _, tile := range newShape.Tiles {
		if t.Grid.Tiles[tile.GetCoordinates()].Blocked {
			return true
		}
	}
	return false
}

// TODO: double check later if it works if multiple rows are full
func (t *TetrisGame) checkForFullRows() {
	for row := 0; row < t.Grid.Height-1; row++ {
		if t.Grid.IsRowFull(row) {
			t.Grid.ClearRow(row)
		}
	}
}

func (t *TetrisGame) PublishGameState() error {
	gameStateMessage := messages.GameStateMessage{
		Grid:       t.GetState(),
		IsGameOver: t.IsGameOver,
	}

	bytes, err := json.Marshal(gameStateMessage)
	if err != nil {
		fmt.Println("Error marshalling grid")
		return err
	}
	t.Hub.PublishMessage(bytes)
	return nil
}

func (t *TetrisGame) EndGame() {
	t.Ticker.Stop()
	t.GameOverChannel <- true
}

func (t *TetrisGame) StartTicker() {
	t.Ticker = time.NewTicker(1 * time.Second)
}

func (t TetrisGame) GetState() entities.Grid {
	return t.Grid
}

func (t *TetrisGame) isGameOver() bool {
	for _, tile := range t.ActiveShape.Tiles {
		if t.Grid.Tiles[tile.GetCoordinates()].Blocked {
			return true
		}
	}
	return false
}
