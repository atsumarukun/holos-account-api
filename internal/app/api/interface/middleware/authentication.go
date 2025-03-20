package middleware

import (
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

func (m *authenticationMiddleware) Authenticate(c *gin.Context) {}
