package entities

import (
	"fmt"
	"slices"
	"time"

	"golang.org/x/exp/rand"
)

type Tile struct {
	Row        int `json:"column"`
	Column     int `json:"row"`
	ColorValue int
	Color      string `json:"color"`
	Display    string `json:"display"`
	Blocked    bool
	IsFixed    bool
}

type Shape struct {
	Tiles    []Tile
	kind     string
	Rotation int
}

func NewTile(x, y int) Tile {
	// TODO: Add color
	return Tile{Row: x, Column: y, ColorValue: 0, Display: "0", Blocked: false}
}

func (t Tile) GetCoordinates() string {
	return fmt.Sprintf("%d-%d", t.Row, t.Column)
}

func (t *Tile) Reset() {
	t = &Tile{Row: t.Row, Column: t.Column, ColorValue: 0, Display: "0"}
}

func (t Tile) Print() {
	color := getTileColor(t.Display, t.Color)
	// fmt.Printf("%s", color)
	fmt.Printf("%s%s%s", color, t.Display, Reset)
}

func GenerateRandomShape() Shape {
	shapes := []string{"L1", "L2", "Z1", "Z2", "I", "O", "T"}

	rand.Seed(uint64(time.Now().UnixNano()))
	rand := rand.Intn(len(shapes))
	return NewShape(shapes[rand])
}

func NewShape(kind string) Shape {
	shape := Shape{}
	switch kind {
	case "I":
		// X R X X
		// 0 0 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "I", Color: Red},
			{Row: 1, Column: 4, Display: "I", IsFixed: true, Color: Red},
			{Row: 1, Column: 5, Display: "I", Color: Red},
			{Row: 1, Column: 6, Display: "I", Color: Red},
		}
	case "O":
		// X X 0 0
		// X X 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "O", Color: Blue},
			{Row: 1, Column: 4, Display: "O", Color: Blue},
			{Row: 2, Column: 3, Display: "O", Color: Blue},
			{Row: 2, Column: 4, Display: "O", Color: Blue},
		}
	case "T":
		// X R X 0
		// 0 X 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "T0", Color: Yellow},
			{Row: 1, Column: 4, Display: "T1", IsFixed: true, Color: Yellow},
			{Row: 1, Column: 5, Display: "T2", Color: Yellow},
			{Row: 2, Column: 4, Display: "T3", Color: Yellow},
		}
	case "L1":
		// X R X 0
		// X 0 0 0

		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "L10", Color: Cyan},
			{Row: 1, Column: 4, Display: "L11", IsFixed: true, Color: Cyan},
			{Row: 1, Column: 5, Display: "L12", Color: Cyan},
			{Row: 2, Column: 3, Display: "L13", Color: Cyan},
		}

	case "L2":
		// X 0 0 0
		// X R X 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "L20", Color: Cyan},
			{Row: 1, Column: 4, Display: "L21", IsFixed: true, Color: Cyan},
			{Row: 1, Column: 5, Display: "L22", Color: Cyan},
			{Row: 2, Column: 5, Display: "L23", Color: Cyan},
		}
	case "Z1":
		// 0 X X 0
		// X R 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "Z10", Color: Magenta},
			{Row: 1, Column: 4, Display: "Z11", IsFixed: true, Color: Magenta},
			{Row: 2, Column: 4, Display: "Z12", Color: Magenta},
			{Row: 2, Column: 5, Display: "Z13", Color: Magenta},
		}
	case "Z2":
		// X X 0 0
		// 0 R X X
		shape.Tiles = []Tile{
			{Row: 1, Column: 4, Display: "Z20", Color: Magenta},
			{Row: 1, Column: 5, Display: "Z21", IsFixed: true, Color: Magenta},
			{Row: 2, Column: 3, Display: "Z22", Color: Magenta},
			{Row: 2, Column: 4, Display: "Z23", Color: Magenta},
		}

	}
	shape.kind = kind
	return shape
}

func (s Shape) Move(direction string, grid Grid) Shape {
	newShape := s.Clone()
	for i := 0; i < len(newShape.Tiles); i++ {

		if direction == "left" {
			newShape.Tiles[i].Column -= 1
		} else if direction == "right" {
			newShape.Tiles[i].Column += 1
		} else if direction == "down" {
			newShape.Tiles[i].Row += 1
		}
	}

	return newShape
}

func (s Shape) Block() Shape {
	fmt.Println(len(s.Tiles))
	newShape := s
	for i := 0; i < len(s.Tiles); i++ {
		newShape.Tiles[i].Blocked = true
	}
	return newShape
}

func (s Shape) Clone() Shape {
	newShape := Shape{}
	newShape.Tiles = slices.Clone(s.Tiles)
	newShape.kind = s.kind
	newShape.Rotation = s.Rotation
	return newShape
}

func (s Shape) Rotate() Shape {
	fmt.Println("Rotating")
	newShape := Shape{}
	newShape.Tiles = slices.Clone(s.Tiles)
	newShape.kind = s.kind
	newShape.Rotation = s.Rotation

	switch s.kind {
	case "T":
		switch s.Rotation {
		case 0:
			// X R X 0
			// 0 X 0 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column -= 1
			newShape.Tiles[3].Row -= 1
			newShape.Tiles[3].Column -= 1
			newShape.Rotation = 1
		case 1:
			// 0 X 0
			// 0 R X
			// 0 X 0
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column -= 1
			newShape.Tiles[3].Row -= 1
			newShape.Tiles[3].Column += 1
			newShape.Rotation = 2

		case 2:
			// 0 X 0
			// X R X
			// 0 X 0
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column += 1
			newShape.Tiles[3].Row += 1
			newShape.Tiles[3].Column += 1
			newShape.Rotation = 3
		case 3:
			// X R X 0
			// 0 X 0 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column += 1
			newShape.Tiles[3].Row += 1
			newShape.Tiles[3].Column -= 1
			newShape.Rotation = 0
		}
	case "I":
		switch s.Rotation {
		case 0:
			// X
			// R
			// X
			// X
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column -= 1
			newShape.Tiles[3].Row += 2
			newShape.Tiles[3].Column -= 2
			newShape.Rotation = 1
		case 1:
			// X X X R
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column += 1
			newShape.Tiles[3].Row -= 2
			newShape.Tiles[3].Column += 2
			newShape.Rotation = 0
		}
	case "L1":
		switch s.Rotation {
		case 0:
			// X R X 0
			// X 0 0 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column -= 1
			newShape.Tiles[3].Row -= 2
			// newShape.Tiles[3].Column -= 2
			newShape.Rotation = 1
		case 1:
			// X X
			// 0 R
			// 0 X
			// 0 X
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column += 1
			// newShape.Tiles[3].Row -= 2
			newShape.Tiles[3].Column += 2
			newShape.Rotation = 2
		case 2:
			// 0 X 0 0
			// X R X 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column -= 1
			newShape.Tiles[3].Row += 2
			// newShape.Tiles[3].Column += 2
			newShape.Rotation = 3
		case 3:
			// X 0
			// X R
			// X 0
			// 0 0
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column += 1
			// newShape.Tiles[3].Row += 2
			newShape.Tiles[3].Column -= 2
			newShape.Rotation = 0
		}
	case "L2":
		switch s.Rotation {
		case 0:
			// X 0 0 0
			// X R X 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column -= 1
			// newShape.Tiles[3].Row += 2
			newShape.Tiles[3].Column -= 2
			newShape.Rotation = 1
		case 1:
			// 0 X
			// X R
			// 0 X
			// 0 0
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column += 1
			newShape.Tiles[3].Row -= 2
			// newShape.Tiles[3].Column += 2
			newShape.Rotation = 2
		case 2:
			// 0 0 X
			// X R X
			// 0 X 0
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column -= 1
			// newShape.Tiles[3].Row -= 2
			newShape.Tiles[3].Column += 2
			newShape.Rotation = 3

		case 3:
			// X 0
			// X R
			// 0 X
			// 0 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column += 1

			newShape.Tiles[3].Row += 2
			// newShape.Tiles[3].Column -= 2
			newShape.Rotation = 0

		}
	case "Z1":
		switch s.Rotation {
		case 0:
			// 0 X X 0
			// X R 0 0
			// newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 2
			newShape.Tiles[1].Row += 1
			newShape.Tiles[1].Column += 1
			newShape.Tiles[3].Row += 1
			newShape.Tiles[3].Column -= 1
			newShape.Rotation = 1
		case 1:
			// X 0
			// R X
			// 0 0
			// 0 0
			newShape.Tiles[0].Row += 2
			newShape.Tiles[1].Row += 1
			newShape.Tiles[1].Column -= 1
			newShape.Tiles[3].Row -= 1
			newShape.Tiles[3].Column -= 1
			newShape.Rotation = 2
		case 2:
			// 0 X X 0
			// X R 0 0
			// newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column -= 2
			newShape.Tiles[1].Row -= 1
			newShape.Tiles[1].Column -= 1
			newShape.Tiles[3].Row -= 1
			newShape.Tiles[3].Column += 1
			newShape.Rotation = 3
		case 3:
			// X 0
			// R X
			// 0 0
			// 0 0
			newShape.Tiles[0].Row -= 2
			newShape.Tiles[1].Row -= 1
			newShape.Tiles[1].Column += 1
			newShape.Tiles[3].Row += 1
			newShape.Tiles[3].Column += 1
			newShape.Rotation = 0
		}
	case "Z2":
		switch s.Rotation {
		case 0:
			// X X 0 0
			// 0 R X X
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[1].Row += 2
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column += 1
			newShape.Rotation = 1
		case 1:
			// 0 X
			// X R
			// X 0
			// 0 0
			newShape.Tiles[0].Row += 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[1].Column -= 2
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column += 1
			newShape.Rotation = 2

		case 2:
			// 0 0 X
			// X R X
			// X 0 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column -= 1
			newShape.Tiles[1].Row -= 2
			newShape.Tiles[2].Row += 1
			newShape.Tiles[2].Column -= 1
			newShape.Rotation = 3
		case 3:
			// X 0
			// R X
			// 0 X
			// 0 0
			newShape.Tiles[0].Row -= 1
			newShape.Tiles[0].Column += 1
			newShape.Tiles[1].Column += 2
			newShape.Tiles[2].Row -= 1
			newShape.Tiles[2].Column -= 1
			newShape.Rotation = 0
		}

	case "O":
		return s

	}
	return newShape
}

// 11 * 21 default
type Grid struct {
	Width, Height int
	Tiles         map[string]Tile `json:"tiles" mapstructure:"tiles"`
}

func NewGrid(width, height int) Grid {
	fmt.Println("newGrd")
	tiles := make(map[string]Tile)

	grid := Grid{Width: width, Height: height, Tiles: tiles}

	for i := 0; i < grid.Height; i++ {
		for j := 0; j < grid.Width; j++ {
			key := GetCoordinates(i, j)
			var tile Tile

			tile = NewTile(i, j)
			// edges
			if j == 0 || i == 0 || j == grid.Width-1 || i == grid.Height-1 {
				tile = Tile{Row: i, Column: j, ColorValue: 0, Display: "X", Blocked: true}
			}
			tiles[key] = tile
		}
	}

	return grid
}

func (g Grid) Print() {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			// fmt.Printf(" %d-%d",  g.Tiles[j][i].X, g.Tiles[j][i].Y)
			// fmt.Printf(" %s", g.Tiles[GetCoordinates(i, j)].Display)
			g.Tiles[GetCoordinates(i, j)].Print()
			//
		}
		fmt.Printf("\n")
	}

	fmt.Printf("------------\n")

}

func (g *Grid) Reset() {
	for i := 0; i < g.Height-1; i++ {
		for j := 1; j < g.Width-1; j++ {
			if g.Tiles[GetCoordinates(i, j)].Blocked {
				continue
			}
			g.Tiles[GetCoordinates(i, j)] = NewTile(i, j)
		}
	}
}

func (g *Grid) RenderShape(shape Shape) {

	fmt.Println("Rendering shape")

	// spew.Dump(shape)

	g.Reset()

	for _, tile := range shape.Tiles {
		g.Tiles[tile.GetCoordinates()] = tile
	}

}

func (g *Grid) ClearRow(row int) {
	for i := 1; i < g.Width-1; i++ {
		g.Tiles[GetCoordinates(row, i)] = NewTile(row, i)
	}
	// shift all rows above down
	for i := row; i > 1; i-- {
		for j := 1; j < g.Width-1; j++ {
			newTile := g.Tiles[GetCoordinates(i-1, j)]
			newTile.Row = i
			g.Tiles[GetCoordinates(i, j)] = newTile
		}
	}
}

func (g Grid) IsRowFull(row int) bool {
	for i := 1; i < g.Width-1; i++ {
		if !g.Tiles[GetCoordinates(row, i)].Blocked {
			return false
		}
	}
	return true
}

func GetCoordinates(x, y int) string {
	return fmt.Sprintf("%d-%d", x, y)
}

var Cyan = "\033[36m"
var Green = "\033[32m"
var Blue = "\033[34m"
var Red = "\033[31m"
var Yellow = "\033[33m"
var Magenta = "\033[35m"
var Orange = "\033[91m"
var Gray = "\033[37m"
var White = "\033[97m"
var Black = "\033[30m"
var Reset = "\033[0m"

func getTileColor(tile string, color string) string {
	switch tile {
	case "0":
		return Black
	case "X":
		return White
	default:
		return color
	}
}
