package handler

import (
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/parameter"
	hdlerr "github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
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
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, "failed to create account"))
		return
	}

	ctx := c.Request.Context()

	account, err := h.accountUC.Create(ctx, req.Name, req.Password, req.ConfirmPassword)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToAccountResponse(account))
}

func (h *accountHandler) UpdateName(c *gin.Context) {
	var req schema.UpdateAccountNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, "failed to update account name"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to update account name"))
		return
	}

	ctx := c.Request.Context()

	account, err := h.accountUC.UpdateName(ctx, accountID, req.Password, req.Name)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToAccountResponse(account))
}

func (h *accountHandler) UpdatePassword(c *gin.Context) {
	var req schema.UpdateAccountPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, "failed to update account password"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to update account password"))
		return
	}

	ctx := c.Request.Context()

	account, err := h.accountUC.UpdatePassword(ctx, accountID, req.Password, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToAccountResponse(account))
}

func (h *accountHandler) Delete(c *gin.Context) {
	var req schema.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, "failed to delete account"))
		return
	}

	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to delete account"))
		return
	}

	ctx := c.Request.Context()

	if err := h.accountUC.Delete(ctx, accountID, req.Password); err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
