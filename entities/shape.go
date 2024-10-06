package entities

import (
	"fmt"

	"golang.org/x/exp/rand"
)

type Tile struct {
	Row        int `json:"column"`
	Column     int `json:"row"`
	ColorValue int
	Color      string `json:"color"`
	Display    string `json:"display"`
	Blocked    bool
	IsFixed    *bool
}

type Shape struct {
	Tiles []Tile
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

func GenerateRandomShape() string {
	shapes := []string{"I", "O", "T"} //"L1", "L2", "Z1", "Z2"

	rand := rand.Intn(len(shapes))
	return shapes[rand]
}

func NewShape(kind string, color string) Shape {
	shape := Shape{}
	trueval := true
	switch kind {
	case "I":
		// X R X X
		// 0 0 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "I", Color: color}, {Row: 1, Column: 4, Display: "I", IsFixed: &trueval, Color: color}, {Row: 1, Column: 5, Display: "I", Color: color}, {Row: 1, Column: 6, Display: "I", Color: color},
		}
	case "O":
		// X X 0 0
		// X X 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "O", Color: color}, {Row: 1, Column: 4, Display: "O", Color: color}, {Row: 2, Column: 3, Display: "O", Color: color}, {Row: 2, Column: 4, Display: "O", Color: color},
		}
	case "T":
		// X R X 0
		// 0 X 0 0
		shape.Tiles = []Tile{
			{Row: 1, Column: 3, Display: "T", Color: color}, {Row: 1, Column: 4, Display: "T", IsFixed: &trueval, Color: color}, {Row: 1, Column: 5, Display: "T", Color: color}, {Row: 2, Column: 4, Display: "T", Color: color},
		}
		// case "L1":
		// 	// X R X 0
		// 	// X 0 0 0
		// 	shape.Tiles = [4]Tile{
		// 		{Tile{X: 0, Y: 0}, Tile{X: 0, Y: 1}, Tile{X: 0, Y: 3}, nil},
		// 		{Tile{X: 1, Y: 0}, nil, nil, nil},
		// 	}
		// case "L2":
		// 	// X 0 0 0
		// 	// X R X 0
		// 	shape.Tiles = [4]Tile{
		// 		{Tile{X: 0, Y: 0}, nil, nil, nil},
		// 		{Tile{X: 1, Y: 0}, Tile{X: 1, Y: 1}, Tile{X: 1, Y: 3}, nil},
		// 	}
		// case "Z1":
		// 	// 0 X X 0
		// 	// X R 0 0
		// 	shape.Tiles = [4]Tile{
		// 		{nil, Tile{X: 0, Y: 1}, Tile{X: 0, Y: 2}, nil},
		// 		{Tile{X: 1, Y: 0}, Tile{X: 1, Y: 1}, nil, nil},
		// 	}
		// case "Z2":
		// 	// X X 0 0
		// 	// 0 R X 0
		// 	shape.Tiles = [4]Tile{
		// 		{Tile{X: 0, Y: 0}, Tile{X: 0, Y: 1}, nil, nil},
		// 		{nil, Tile{X: 1, Y: 1}, Tile{X: 1, Y: 2}, nil},
		// 	}
	}
	return shape
}

func (s Shape) Move(direction string, grid Grid) (Shape, error) {
	tempTiles := make([]Tile, len(s.Tiles))
	for i := 0; i < len(s.Tiles); i++ {
		newTile := s.Tiles[i]
		if direction == "left" {
			newTile.Column -= 1
			if grid.Tiles[newTile.GetCoordinates()].Blocked {
				return s, fmt.Errorf("Shape is colliding")
			}
			tempTiles[i] = newTile
		} else if direction == "right" {
			newTile.Column += 1
			if grid.Tiles[newTile.GetCoordinates()].Blocked {
				return s, fmt.Errorf("Shape is colliding")
			}
			tempTiles[i] = newTile
		} else if direction == "down" {
			newTile.Row += 1
			if grid.Tiles[newTile.GetCoordinates()].Blocked {
				return s, fmt.Errorf("Shape is colliding")
			}
			tempTiles[i] = newTile
		}
	}
	for i := 0; i < len(s.Tiles); i++ {
		s.Tiles[i] = tempTiles[i]
	}

	return s, nil
}

func (s Shape) GetFixed() Tile {
	for i := 0; i < len(s.Tiles); i++ {
		if *s.Tiles[i].IsFixed {
			return s.Tiles[i]
		}
	}
	return Tile{}
}

func (s Shape) Block() Shape {
	fmt.Println(len(s.Tiles))
	newShape := s
	for i := 0; i < len(s.Tiles); i++ {
		newShape.Tiles[i].Blocked = true
	}
	return newShape
}

// TODO: Implement this
func (s Shape) Rotate() {
	fixed := s.GetFixed()
	for i := 0; i < len(s.Tiles); i++ {
		if *s.Tiles[i].IsFixed {
			continue
		}
		if s.Tiles[i].Row == fixed.Row {
			s.Tiles[i].Row = fixed.Row
			s.Tiles[i].Column = fixed.Column
		} else {
			s.Tiles[i].Row = fixed.Row - i
			s.Tiles[i].Column = fixed.Column
		}
	}
}

func (s Shape) IsColliding(grid Grid, direction string) bool {

	// Check if the tile is out of grid bounds
	for i := 0; i < len(s.Tiles); i++ {
		if grid.Tiles[s.Tiles[i].GetCoordinates()].Blocked {
			return true
		}
	}
	return false
}

// 11 * 21 default
type Grid struct {
	Width, Height int
	Tiles         map[string]Tile `json:"tiles" mapstructure:"tiles"`
	ColorCounter  int             `json:"colorCounter"`
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
			newTile := g.Tiles[GetCoordinates(i, j)]
			newTile.Row += 1
			g.Tiles[GetCoordinates(i, j)] = newTile
		}
	}
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
