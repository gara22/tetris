package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/gara22/tetris/entities"
)

const (
	HEIGHT = 21
	WIDTH  = 11
)

func main() {
	grid := entities.NewGrid(WIDTH, HEIGHT)

	// // entities.GenerateRandomShape()

	iShape := entities.NewShape("I")

	shapes := []entities.Shape{iShape}
	grid.RenderShapes(shapes)
	grid.Print()

	activeIndex := 0

	done := make(chan bool)
	userInput := make(chan string)

	// Initialize the keyboard listener
	if err := keyboard.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer keyboard.Close()

	// Capture user input in a separate goroutine
	go func() {
		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Println(err)
				return
			}
			switch key {
			case keyboard.KeyArrowLeft:
				userInput <- "left"
			case keyboard.KeyArrowRight:
				userInput <- "right"
			case keyboard.KeyArrowDown:
				userInput <- "down"
			case keyboard.KeyEsc:
				done <- true
				return
			default:
				fmt.Println("Invalid input")
			}
		}
	}()
	// Main goroutine to handle shape movement
	go func() {
		for {
			fmt.Println("shapes", len(shapes))
			select {
			case <-done:
				return
			case input := <-userInput:
				var newShape entities.Shape
				switch input {
				case "left":
					newShape = shapes[activeIndex].Move("left")
				case "right":

					newShape = shapes[activeIndex].Move("right")
				case "down":
					newShape = shapes[activeIndex].Move("down")
					if newShape.IsColliding(grid) {
						fmt.Println("Shape is stuck")
						shapes[activeIndex].Block()
						activeIndex++
						// TODO: generate a random shape here and append it to the shapes slice
						shapes = append(shapes, entities.NewShape("I"))
						continue
					}
				default:
					fmt.Println("Invalid input")
				}
				// check if new shape is colliding with the grid
				if newShape.IsColliding(grid) {
					fmt.Println("Shape is colliding")
					continue
				}

				shapes[activeIndex] = newShape

				grid.RenderShapes(shapes)
				grid.Print()
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		close(done)
	case <-done:
	}

}
