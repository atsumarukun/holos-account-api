package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/parameter"
	hdlerr "github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-api-pkg/errors"
)

type SessionHandler interface {
	Login(*gin.Context)
	Logout(*gin.Context)
	Authorize(*gin.Context)
}

type sessionHandler struct {
	sessionUC usecase.SessionUsecase
}

func NewSessionHandler(sessionUC usecase.SessionUsecase) SessionHandler {
	return &sessionHandler{
		sessionUC: sessionUC,
	}
}

func (h *sessionHandler) Login(c *gin.Context) {
	var req schema.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, "failed to login"))
		return
	}

	ctx := c.Request.Context()

	session, err := h.sessionUC.Login(ctx, req.AccountName, req.Password)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToSessionResponse(session))
}

func (h *sessionHandler) Logout(c *gin.Context) {
	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to logout"))
		return
	}

	ctx := c.Request.Context()

	if err := h.sessionUC.Logout(ctx, accountID); err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) Authorize(c *gin.Context) {
	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to authorize"))
		return
	}

	ctx := c.Request.Context()

	account, err := h.sessionUC.Authorize(ctx, accountID)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToAauthorizationResponse(account))
}
