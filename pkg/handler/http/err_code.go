package http

type ErrCode int

const (
	ErrCodeInvalidInput ErrCode = 1
	ErrCodeInternal     ErrCode = 2
)
