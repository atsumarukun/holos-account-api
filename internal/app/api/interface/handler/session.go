package handler

import (
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	hdlerr "github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/parameter"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
)

type SessionHandler interface {
	Create(*gin.Context)
	Delete(*gin.Context)
	Verify(*gin.Context)
}

type sessionHandler struct {
	sessionUC usecase.SessionUsecase
}

func NewSessionHandler(sessionUC usecase.SessionUsecase) SessionHandler {
	return &sessionHandler{
		sessionUC: sessionUC,
	}
}

func (h *sessionHandler) Create(c *gin.Context) {
	var req schema.CreateSessionRequest
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

	c.JSON(http.StatusCreated, builder.ToSessionResponse(session))
}

func (h *sessionHandler) Delete(c *gin.Context) {
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

func (h *sessionHandler) Verify(c *gin.Context) {
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
