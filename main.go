package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	app_service "github.com/gara22/tetris/app-service"
	handler "github.com/gara22/tetris/http"
	socket "github.com/gara22/tetris/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	appService := app_service.NewAppService()
	// tetrisGame := game.NewTetrisGame(*hub)
	// tetrisGame.StartGame()

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

	handler := handler.NewHTTPHandler(*appService)
	router.POST("/new-game", handler.NewTetrisGame)
	router.GET("/ws", func(c *gin.Context) {
		// TODO: validate game id
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": "game id is required",
			})
			return
		}

		game := appService.Games[id]

		if game == nil {
			spew.Dump("game not found")
			c.JSON(400, gin.H{
				"error": "game not found",
			})
			return
		}

		socket.ServeWs(&appService.Games[id].Hub, c.Writer, c.Request)
		game.PublishGameState()

	})
	router.Run("localhost:8080")

}
