package parameter

import (
	"github.com/gin-gonic/gin"

	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
)

func GetContextParameter[T any](c *gin.Context, name string) (T, error) {
	var zero T

	param, exists := c.Get(name)
	if !exists {
		return zero, status.ErrInternal
	}

	v, ok := param.(T)
	if !ok {
		return zero, status.ErrInternal
	}

	return v, nil
}
