package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gara22/tetris/game"
	handler "github.com/gara22/tetris/http"
	socket "github.com/gara22/tetris/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	hub := socket.NewHub()
	go hub.Run()
	tetrisGame := game.NewTetrisGame(*hub)
	tetrisGame.StartGame()

	tetrisGame.Grid.RenderShapes(tetrisGame.Shapes)
	// tetrisGame.Grid.Print()

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

	// allow cors for localhost:3000

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},                   // Allow localhost for front-end development (adjust port if needed)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := handler.NewHTTPHandler(&tetrisGame)
	router.POST("/move", handler.Move)
	router.GET("/state", handler.GetState)
	router.GET("/ws", func(c *gin.Context) {
		socket.ServeWs(hub, c.Writer, c.Request)
	})
	router.Run(":8080")

	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		close(done)
	case <-done:
	}

}
