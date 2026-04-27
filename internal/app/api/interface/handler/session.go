package handler

import (
	stderr "errors"
	"net/http"
	"strings"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	hdlerr "github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/parameter"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
)

var ErrInvalidToken = stderr.New("invalid token")

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
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeBadRequest, "failed to create session"))
		return
	}

	ctx := c.Request.Context()

	session, err := h.sessionUC.Create(ctx, req.AccountName, req.Password)
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusCreated, builder.ToSessionResponse(session))
}

func (h *sessionHandler) Delete(c *gin.Context) {
	accountID, err := parameter.GetContextParameter[uuid.UUID](c, "accountID")
	if err != nil {
		hdlerr.Handle(c, errors.Wrap(err, errors.CodeUnauthenticated, "failed to delete session"))
		return
	}

	ctx := c.Request.Context()

	if err := h.sessionUC.Delete(ctx, accountID); err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *sessionHandler) Verify(c *gin.Context) {
	sessionToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(sessionToken) != 2 || sessionToken[0] != "Session" {
		err := errors.Wrap(ErrInvalidToken, errors.CodeUnauthenticated, "failed to verify")
		hdlerr.Handle(c, err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	account, err := h.sessionUC.Verify(ctx, sessionToken[1])
	if err != nil {
		hdlerr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToVerifiedSessionResponse(account))
}
