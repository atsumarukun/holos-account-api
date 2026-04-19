package middleware

import (
	stderr "errors"
	"strings"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"

	hdlerr "github.com/atsumarukun/holos-account-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
)

var (
	ErrInvalidToken    = stderr.New("invalid token")
	ErrAccountNotFound = stderr.New("account not found")
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
		err := errors.Wrap(ErrInvalidToken, errors.CodeUnauthenticated, "failed to authenticate")
		hdlerr.Handle(c, err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	account, err := m.sessionUC.Authenticate(ctx, sessionToken[1])
	if err != nil {
		hdlerr.Handle(c, err)
		c.Abort()
		return
	}
	if account == nil {
		err := errors.Wrap(ErrAccountNotFound, errors.CodeUnauthenticated, "failed to authenticate")
		hdlerr.Handle(c, err)
		c.Abort()
		return
	}

	c.Set("accountID", account.ID)
	c.Next()
}
