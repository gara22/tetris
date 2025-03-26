package handler

import (
	"net/http"

	app_service "github.com/gara22/tetris/app-service"
	"github.com/gara22/tetris/messages"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HTTPHandler struct {
	svc app_service.AppService
}

func NewHTTPHandler(appService app_service.AppService) *HTTPHandler {
	return &HTTPHandler{
		svc: appService,
	}
}

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

func (h *HTTPHandler) AddScore(ctx *gin.Context) {
	var message messages.SavePlayerScoreMessage
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	validator := validator.New()
	if err := validator.Struct(message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid request",
		})
		return
	}

	err := h.svc.AddScore(message.GameId, message.PlayerName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Score added",
	})
}
