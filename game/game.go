package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gara22/tetris/entities"
	"github.com/gara22/tetris/messages"
	socket "github.com/gara22/tetris/websocket"
)

type TetrisGame struct {
	Grid        entities.Grid `json:"grid"`
	ActiveShape entities.Shape
	Ticker      *time.Ticker
	Hub         socket.Hub
}

var (
	// Make sure that TurnaroundManager implements the AppService interface
	_ Game = (*TetrisGame)(nil)
)

const (
	HEIGHT = 21
	WIDTH  = 11
)

func NewTetrisGame(hub socket.Hub) TetrisGame {
	return TetrisGame{
		Grid:        entities.NewGrid(WIDTH, HEIGHT),
		ActiveShape: entities.Shape{},
		Ticker:      nil,
		Hub:         hub,
	}
}

func (t *TetrisGame) StartGame() {
	// TODO: generate a random shape here and append it to the shapes slice
	shapeKind := entities.GenerateRandomShape()
	fmt.Println("start game")
	t.ActiveShape = entities.NewShape(shapeKind)
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
				}
				t = &game

				// game.Move()

				bytes, err := json.Marshal(game.Grid)
				if err != nil {
					fmt.Println("Error marshalling grid")
				}
				t.Hub.PublishMessage(bytes)
			}
		}
	}()
}

func (t TetrisGame) Move(params MoveParams) (TetrisGame, error) {
	spew.Dump(t.ActiveShape)
	newShape, err := t.calculateNewShape(params, t.ActiveShape)
	// spew.Dump(t.ActiveShape)
	if err != nil && err.Error() == "Shape is stuck" {
		fmt.Println("Shape is stuck")
		newShape = t.ActiveShape.Block()
		t.Grid.RenderShape(newShape)

		shapeKind := entities.GenerateRandomShape()
		t.ActiveShape = entities.NewShape(shapeKind)
		// TODO: do we need to render here?
		t.checkForFullRows()
		t.Grid.RenderShape(t.ActiveShape)
		spew.Dump(t.ActiveShape)
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
	switch params.Direction {
	case "left":
		return activeShape.Move("left", t.Grid)
	case "right":
		return activeShape.Move("right", t.Grid)
	case "down":
		newShape, err := activeShape.Move("down", t.Grid)

		if err != nil {
			return entities.Shape{}, fmt.Errorf("Shape is stuck")

		}
		return newShape, err

		// if newShape.IsColliding(t.Grid, params.Direction) {

		// 	return entities.Shape{}, true, t, nil
		// }
	default:
		return entities.Shape{}, fmt.Errorf("Invalid input")
	}
}

// TODO: double check later if it works if multiple rows are full
func (t *TetrisGame) checkForFullRows() {
	for row := 0; row < t.Grid.Height-1; row++ {
		if t.Grid.IsRowFull(row) {
			t.Grid.ClearRow(row)
		}
	}
}

func (t *TetrisGame) StopGame() {
	t.Ticker.Stop()
}

func (t *TetrisGame) StartTicker() {
	t.Ticker = time.NewTicker(1 * time.Second)
}

func (t *TetrisGame) StopTicker() {
	t.Ticker.Stop()
	t.Ticker = nil
}

func (t TetrisGame) GetState() entities.Grid {
	return t.Grid
}

// func (t *TetrisGame)isGameOver() bool {

// }

func deepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
