package code

type StatusCode string

var (
	BadRequest           StatusCode = "BAD_REQUEST"
	Conflict             StatusCode = "CONFLICT"
	UnprocessableContent StatusCode = "UNPROCESSABLE_CONTENT"
	Unauthorized         StatusCode = "UNAUTHORIZED"
	Internal             StatusCode = "INTERNAL"
)
