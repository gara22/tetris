package entities

import (
	"fmt"
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
	Tiles [4]Tile
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

// func GenerateRandomShape() Shape {
// 	shapes := []string{"I", "O", "T", "L1", "L2", "Z1", "Z2"}
// 	rand := rand.Intn(len(shapes))
// 	fmt.Println(rand)
// 	shape := NewShape(shapes[rand])
// 	return shape
// }

func NewShape(kind string, color string) Shape {
	shape := Shape{}
	trueval := true
	switch kind {
	case "I":
		// X R X X
		// 0 0 0 0
		shape.Tiles = [4]Tile{
			{Row: 1, Column: 3, Display: "I", Color: color}, {Row: 1, Column: 4, Display: "I", IsFixed: &trueval, Color: color}, {Row: 1, Column: 5, Display: "I", Color: color}, {Row: 1, Column: 6, Display: "I", Color: color},
		}
		// case "O":
		// 	// X X 0 0
		// 	// X X 0 0
		// 	shape.Tiles = [4]Tile{
		// 		{Tile{X: 0, Y: 0}, Tile{X: 0, Y: 1}, nil, nil},
		// 		{Tile{X: 1, Y: 0}, Tile{X: 1, Y: 1}, nil, nil},
		// 	}
		// case "T":
		// 	// X R X 0
		// 	// 0 X 0 0
		// 	shape.Tiles = [4]Tile{
		// 		{Tile{X: 0, Y: 0}, Tile{X: 0, Y: 1}, Tile{X: 0, Y: 2}, nil},
		// 		{nil, Tile{X: 1, Y: 1}, nil},
		// 	}
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

func (s Shape) Move(direction string) Shape {
	for i := 0; i < len(s.Tiles); i++ {
		if direction == "left" {
			s.Tiles[i].Column -= 1
		} else if direction == "right" {
			s.Tiles[i].Column += 1
		} else if direction == "down" {
			s.Tiles[i].Row += 1
		}
	}
	return s
}

func (s Shape) GetFixed() Tile {
	for i := 0; i < 4; i++ {
		if *s.Tiles[i].IsFixed {
			return s.Tiles[i]
		}
	}
	return Tile{}
}

func (s *Shape) MostOuterTiles(direction string) []Tile {
	mostOuters := []Tile{s.Tiles[0]}
	for i := 1; i < 4; i++ {
		switch direction {
		case "left":
			if s.Tiles[i].Column < mostOuters[0].Column {
				mostOuters = []Tile{s.Tiles[i]}
			}
			if s.Tiles[i].Column == mostOuters[0].Column {
				mostOuters = append(mostOuters, s.Tiles[i])
			}
		case "right":

			if s.Tiles[i].Column > mostOuters[0].Column {
				mostOuters = []Tile{s.Tiles[i]}
			}
			if s.Tiles[i].Column == mostOuters[0].Column {
				mostOuters = append(mostOuters, s.Tiles[i])
			}
		case "down":
			if s.Tiles[i].Row > mostOuters[0].Row {
				mostOuters = []Tile{s.Tiles[i]}
			}
			if s.Tiles[i].Row == mostOuters[0].Row {
				mostOuters = append(mostOuters, s.Tiles[i])
			}
		}
	}
	return mostOuters
}

func (s *Shape) Block() {
	for i := 0; i < 4; i++ {
		s.Tiles[i].Blocked = true
	}
}

// TODO: Implement this
func (s Shape) Rotate() {
	fixed := s.GetFixed()
	for i := 0; i < 4; i++ {
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
	for _, t := range s.MostOuterTiles(direction) {
		if grid.Tiles[t.GetCoordinates()].Blocked {
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
	emptyTiles := make(map[string]Tile)

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
			emptyTiles[key] = tile
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
			g.Tiles[GetCoordinates(i, j)] = NewTile(i, j)
		}
	}
}

func (g *Grid) RenderShapes(shapes []Shape) {
	// fmt.Println(g.OldShapes)

	// for _, oldShape := range g.OldShapes {
	// 	for _, oldTile := range oldShape.Tiles {
	// 		g.Tiles[oldTile.GetCoordinates()] = NewTile(oldTile.Row, oldTile.Column)
	// 	}
	// }

	g.Reset()

	for _, newShape := range shapes {
		for _, newTile := range newShape.Tiles {
			g.Tiles[newTile.GetCoordinates()] = newTile
		}
	}

	// if len(g.OldShapes) == len(shapes) {
	// 	copy(g.OldShapes, shapes)
	// } else {
	// 	g.OldShapes = append(g.OldShapes, shapes...)
	// }

	// fmt.Println(g.OldShapes)
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
	case "I":
		return color
	case "0":
		return Black
	case "X":
		return White
	}
	return ""
}
