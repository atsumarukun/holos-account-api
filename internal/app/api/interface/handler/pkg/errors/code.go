package errors

import (
	"net/http"

	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status/code"
)

var codes = map[code.StatusCode]int{
	code.BadRequest:   http.StatusBadRequest,
	code.Conflict:     http.StatusConflict,
	code.Unauthorized: http.StatusUnauthorized,
	code.Internal:     http.StatusInternalServerError,
}
