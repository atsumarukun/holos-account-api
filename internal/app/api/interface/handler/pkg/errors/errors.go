package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
)

func Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if v, ok := err.(*status.Status); ok {
		c.JSON(codes[v.Code()], map[string]string{"message": v.Message()})
	} else {
		c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}
}
