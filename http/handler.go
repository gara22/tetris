package handler

import (
	"net/http"

	"github.com/gara22/tetris/game"
	"github.com/gara22/tetris/messages"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc game.Game
}

func NewHTTPHandler(GameService game.Game) *HTTPHandler {
	return &HTTPHandler{
		svc: GameService,
	}
}

func (h *HTTPHandler) Move(ctx *gin.Context) {
	var message messages.MoveMessage
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})

		return
	}

	moveParams := game.MoveParams{
		Direction: message.Direction,
	}

	game, err := h.svc.Move(moveParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// gridBytes, err := json.Marshal(game.Grid)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusCreated, game.Grid)
}

func (h *HTTPHandler) GetState(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.svc.GetState())
}
