package middleware

import (
	"log"
	"strings"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
)

type AuthenticationMiddleware interface {
	Authenticate(*gin.Context)
}

type authenticationMiddleware struct {
	sessionUC usecase.SessionUsecase
}

func NewAuthenticationMiddleware(sessionUC usecase.SessionUsecase) AuthenticationMiddleware {
	return &authenticationMiddleware{
		sessionUC: sessionUC,
	}
}

func (m *authenticationMiddleware) Authenticate(c *gin.Context) {
	sessionToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(sessionToken) != 2 || sessionToken[0] != "Session" {
		err := status.ErrUnauthorized
		log.Println(err)
		errors.Handle(c, err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	account, err := m.sessionUC.Authenticate(ctx, sessionToken[1])
	if err != nil {
		log.Println(err)
		errors.Handle(c, err)
		c.Abort()
		return
	}
	if account == nil {
		err := status.ErrUnauthorized
		log.Println(err)
		errors.Handle(c, err)
		c.Abort()
		return
	}

	c.Set("accountID", account.ID)
	c.Next()
}
