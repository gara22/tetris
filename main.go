package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/gara22/tetris/game"
)

func main() {

	tetrisGame := game.NewTetrisGame()
	tetrisGame.StartGame()

	tetrisGame.Grid.RenderShapes(tetrisGame.Shapes)
	tetrisGame.Grid.Print()

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
	// // Main goroutine to handle shape movement
	go func() {
		for {
			// fmt.Println("shapes", len(shapes))
			select {
			case <-done:
				return
			case input := <-userInput:
				tetrisGame.Move(game.MoveParams{Direction: input})

			}
		}
	}()

	// add a http listener for user input
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, World!")

		// userInput <- r.RequestURI
		// //get direction from uri param
		direction := r.URL.Query().Get("direction")
		fmt.Println("direction", direction)
		userInput <- direction
		w.Write([]byte(direction))
		// w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		close(done)
	case <-done:
	}

}
