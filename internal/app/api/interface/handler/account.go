package handler

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	Create(*gin.Context)
}

type accountHandler struct {
	accountUC usecase.AccountUsecase
}

func NewAccountHandler(accountUC usecase.AccountUsecase) AccountHandler {
	return &accountHandler{
		accountUC: accountUC,
	}
}

func (h *accountHandler) Create(c *gin.Context) {}
