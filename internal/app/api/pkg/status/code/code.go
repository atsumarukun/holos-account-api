package code

type StatusCode string

var (
	BadRequest   StatusCode = "BAD_REQUEST"
	Unauthorized StatusCode = "UNAUTHORIZED"
	Internal     StatusCode = "INTERNAL"
)
