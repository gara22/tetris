package game

import (
	"fmt"
	"time"

	"github.com/gara22/tetris/entities"
)

type TetrisGame struct {
	Grid         entities.Grid
	Shapes       []entities.Shape
	ActiveIndex  int
	Ticker       *time.Ticker
	ColorCounter int
}

var (
	// Make sure that TurnaroundManager implements the AppService interface
	_ Game = (*TetrisGame)(nil)
)

const (
	HEIGHT = 21
	WIDTH  = 11
)

func NewTetrisGame() TetrisGame {
	return TetrisGame{
		Grid:         entities.NewGrid(WIDTH, HEIGHT),
		Shapes:       []entities.Shape{},
		ActiveIndex:  0,
		Ticker:       nil,
		ColorCounter: 0,
	}
}

func (t *TetrisGame) StartGame() TetrisGame {
	// TODO: generate a random shape here and append it to the shapes slice
	t.Shapes = []entities.Shape{entities.NewShape("I", TileColors[t.ColorCounter])}
	// setup a ticker to move the shape down every second
	go func() {
		t.Ticker = time.NewTicker(1 * time.Second)
		for {
			select {
			case <-t.Ticker.C:
				t.Move(MoveParams{Direction: "down"})

			}
		}
	}()

	return *t
}

func (t *TetrisGame) Move(params MoveParams) {
	var newShape entities.Shape
	switch params.Direction {
	case "left":
		newShape = t.Shapes[t.ActiveIndex].Move("left")
	case "right":
		newShape = t.Shapes[t.ActiveIndex].Move("right")
	case "down":
		newShape = t.Shapes[t.ActiveIndex].Move("down")
		if newShape.IsColliding(t.Grid, params.Direction) {
			fmt.Println("Shape is stuck")
			t.Shapes[t.ActiveIndex].Block()
			t.ActiveIndex++
			// TODO: generate a random shape here and append it to the shapes slice
			t.Shapes = append(t.Shapes, entities.NewShape("I", t.NextColor()))
			return
		}
	default:
		fmt.Println("Invalid input")
	}
	if newShape.IsColliding(t.Grid, params.Direction) {
		fmt.Println("Shape is colliding")
		return
	}

	t.Shapes[t.ActiveIndex] = newShape

	t.Grid.RenderShapes(t.Shapes)
	t.Grid.Print()
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

var TileColors = []string{entities.Cyan, entities.Green, entities.Blue, entities.Red, entities.Yellow, entities.Magenta, entities.Orange}

// func (t *TetrisGame)isGameOver() bool {

// }
