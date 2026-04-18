package parameter

import (
	stderr "errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrParameterMissing     = stderr.New("parameter is missing")
	ErrInvalidParameterType = stderr.New("invalid parameter type")
)

func GetContextParameter[T any](c *gin.Context, name string) (T, error) {
	var zero T

	param, exists := c.Get(name)
	if !exists {
		return zero, ErrParameterMissing
	}

	v, ok := param.(T)
	if !ok {
		return zero, ErrInvalidParameterType
	}

	return v, nil
}
