package code

type StatusCode string

var (
	BadRequest   StatusCode = "BAD_REQUEST"
	Conflict     StatusCode = "CONFLICT"
	Unauthorized StatusCode = "UNAUTHORIZED"
	Internal     StatusCode = "INTERNAL"
)
