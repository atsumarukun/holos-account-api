package status

import "github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status/code"

var (
	ErrBadRequest           = Error(code.BadRequest, "bad request")
	ErrConflict             = Error(code.Conflict, "conflict")
	ErrUnprocessableContent = Error(code.UnprocessableContent, "unprocessable content")
	ErrUnauthorized         = Error(code.Unauthorized, "unauthorized")
	ErrInternal             = Error(code.Internal, "internal server error")
)
