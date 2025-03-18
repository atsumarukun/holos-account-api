package handler

import (
	"log"
	"net/http"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
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

func (h *accountHandler) Create(c *gin.Context) {
	var req schema.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		errors.Handle(c, status.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()

	account, err := h.accountUC.Create(ctx, req.Name, req.Password, req.ConfirmPassword)
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToAccountResponse(account))
}
