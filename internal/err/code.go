package fwerr

type Code string

// System errors.
const (
	CodeInternalError Code = "INTERNAL_ERROR"
	CodeAuthError     Code = "AUTH_ERROR"
	CodeForbidden     Code = "FORBIDDEN"
	CodeDBError       Code = "DB_ERROR"
	CodeBadRequest    Code = "BAD_REQUEST"
	CodeNotFound      Code = "NOT_FOUND"
)
