package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/builder"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
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

func (h *sessionHandler) Login(c *gin.Context) {
	var req schema.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		errors.Handle(c, status.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()

	session, err := h.sessionUC.Login(ctx, req.AccountName, req.Password)
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, builder.ToSessionResponse(session))
}
