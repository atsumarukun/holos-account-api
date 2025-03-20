package handler

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
)

type SessionHandler interface {
	Login(*gin.Context)
}

type sessionHandler struct {
	sessionUC usecase.SessionUsecase
}

func NewSessionHandler(sessionUC usecase.SessionUsecase) SessionHandler {
	return &sessionHandler{
		sessionUC: sessionUC,
	}
}

func (h *sessionHandler) Login(c *gin.Context) {}
