package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	app_service "github.com/gara22/tetris/app-service"
	handler "github.com/gara22/tetris/http"
	"github.com/gara22/tetris/repository"
	socket "github.com/gara22/tetris/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: get path from env
	repo := repository.New("./games.json")
	appService := app_service.NewAppService(&repo)
	// tetrisGame := game.NewTetrisGame(*hub)
	// tetrisGame.StartGame()

	// allow cors for localhost:3000

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// TODO: get this from env
		// AllowOrigins:     []string{"http://husi.lol", "http://goblin.rest"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			spew.Dump(origin)
			return origin == "https://husi.lol" || origin == "https://goblin.rest" || origin == "http://localhost:5173"
		},
	}))

	handler := handler.NewHTTPHandler(*appService)
	router.POST("/new-game", handler.NewTetrisGame)
	router.GET("/hello", func(c *gin.Context) {
		spew.Dump("hello")
		c.JSON(200, gin.H{
			"message": "hello",
		})
	})
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
	router.Run(":3200")

}
