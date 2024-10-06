package main

import (
	"time"

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

}
