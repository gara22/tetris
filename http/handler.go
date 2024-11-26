package handler

import (
	"net/http"

	app_service "github.com/gara22/tetris/app-service"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc app_service.AppService
}

func NewHTTPHandler(appService app_service.AppService) *HTTPHandler {
	return &HTTPHandler{
		svc: appService,
	}
}

// func (h *HTTPHandler) Move(ctx *gin.Context) {
// 	var message messages.MoveMessage
// 	if err := ctx.ShouldBindJSON(&message); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"Error": err,
// 		})

// 		return
// 	}

// 	moveParams := game.MoveParams{
// 		Direction: message.Direction,
// 	}

// 	game, err := h.svc.Move(moveParams)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 		return
// 	}

// 	// gridBytes, err := json.Marshal(game.Grid)
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusBadRequest, gin.H{
// 	// 		"error": err,
// 	// 	})
// 	// 	return
// 	// }

// 	ctx.JSON(http.StatusCreated, game.Grid)
// }

// func (h *HTTPHandler) GetState(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, h.svc.GetState())
// }

func (h *HTTPHandler) NewTetrisGame(ctx *gin.Context) {
	gameId, err := h.svc.NewGame()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Game started",
		"gameId":  gameId,
	})
}

// func (h *HTTPHandler) JoinGame(ctx *gin.Context) {
// 	id := ctx.Query("id")
// 	if id == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "game id is required",
// 		})
// 		return
// 	}

// 	err := h.svc.JoinGame(id)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Game joined",
// 	})
// }
