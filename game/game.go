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
	Grid             entities.Grid `json:"grid"`
	Shapes           []entities.Shape
	ActiveShapeIndex int
	Ticker           *time.Ticker
	ColorCounter     int
	Hub              socket.Hub
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
		Grid:             entities.NewGrid(WIDTH, HEIGHT),
		Shapes:           []entities.Shape{},
		ActiveShapeIndex: 0,
		Ticker:           nil,
		ColorCounter:     0,
		Hub:              hub,
	}
}

func (t *TetrisGame) StartGame() TetrisGame {
	// TODO: generate a random shape here and append it to the shapes slice
	t.Shapes = []entities.Shape{entities.NewShape("I", TileColors[t.ColorCounter])}
	// setup a ticker to move the shape down every second
	msg := t.Hub.ReadMessage()

	go func() {
		t.Ticker = time.NewTicker(1 * time.Second)
		for {
			select {
			case <-t.Ticker.C:
				t.Move(MoveParams{Direction: "down"})

			case m := <-msg:
				var msg messages.MoveMessage
				json.Unmarshal(m, &msg)
				spew.Dump(msg)
				t.Move(MoveParams{Direction: msg.Direction})
			}
		}
	}()

	return *t
}

func (t *TetrisGame) Move(params MoveParams) (TetrisGame, error) {
	var newShape entities.Shape
	switch params.Direction {
	case "left":
		newShape = t.Shapes[t.ActiveShapeIndex].Move("left")
	case "right":
		newShape = t.Shapes[t.ActiveShapeIndex].Move("right")
	case "down":
		newShape = t.Shapes[t.ActiveShapeIndex].Move("down")
		if newShape.IsColliding(t.Grid, params.Direction) {
			fmt.Println("Shape is stuck")
			t.Shapes[t.ActiveShapeIndex].Block()
			t.ActiveShapeIndex++
			// TODO: generate a random shape here and append it to the shapes slice
			t.Shapes = append(t.Shapes, entities.NewShape("I", t.NextColor()))
			return *t, nil
		}
	default:
		fmt.Println("Invalid input")
	}
	if newShape.IsColliding(t.Grid, params.Direction) {
		fmt.Println("Shape is colliding")
		return *t, nil
	}

	t.Shapes[t.ActiveShapeIndex] = newShape

	t.Grid.RenderShapes(t.Shapes)
	bytes, err := json.Marshal(t.Grid)
	if err != nil {
		fmt.Println("Error marshalling grid")
		return *t, err
	}
	t.Hub.PublishMessage(bytes)
	t.Grid.Print()
	return *t, nil
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

func (t *TetrisGame) NextColor() string {
	if t.ColorCounter == len(TileColors)-1 {
		t.ColorCounter = 0
		return TileColors[t.ColorCounter]
	}
	t.ColorCounter++
	return TileColors[t.ColorCounter]
}

func (t TetrisGame) GetState() entities.Grid {
	return t.Grid
}

var TileColors = []string{entities.Cyan, entities.Green, entities.Blue, entities.Red, entities.Yellow, entities.Magenta, entities.Orange}

// func (t *TetrisGame)isGameOver() bool {

// }
