package status

import "github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status/code"

var (
	ErrBadRequest   = Error(code.BadRequest, "bad request")
	ErrConflict     = Error(code.BadRequest, "conflict")
	ErrUnauthorized = Error(code.Unauthorized, "unauthorized")
	ErrInternal     = Error(code.Internal, "internal server error")
)
