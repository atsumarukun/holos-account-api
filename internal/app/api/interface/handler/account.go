package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/parameter"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
)

type AccountHandler interface {
	Create(*gin.Context)
	UpdateName(*gin.Context)
	UpdatePassword(*gin.Context)
	Delete(*gin.Context)
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

func (h *accountHandler) UpdateName(c *gin.Context) {
	var req schema.UpdateAccountNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		errors.Handle(c, status.ErrBadRequest)
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	account, err := h.accountUC.UpdateName(ctx, accountID, req.Name)
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToAccountResponse(account))
}

func (h *accountHandler) UpdatePassword(c *gin.Context) {
	var req schema.UpdateAccountPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		errors.Handle(c, status.ErrBadRequest)
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	account, err := h.accountUC.UpdatePassword(ctx, accountID, req.Password, req.ConfirmPassword)
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToAccountResponse(account))
}

func (h *accountHandler) Delete(c *gin.Context) {
	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	ctx := c.Request.Context()

	if err := h.accountUC.Delete(ctx, accountID); err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
