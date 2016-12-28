package kate

type ErrorInfo interface {
	error
	Code() int
}

type errSimple struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

func NewError(code int, message string) ErrorInfo {
	return &errSimple{code, message}
}

func (e *errSimple) Code() int {
	return e.ErrCode
}

func (e *errSimple) Error() string {
	return e.ErrMessage
}

var (
	ErrBadRequest       = NewError(400, "bad request")
	ErrUnauthorized     = NewError(401, "unauthorized")
	ErrNotFound         = NewError(404, "not found")
	ErrMethodNotAllowed = NewError(405, "method not allowed")
	ErrRequestExpired   = NewError(406, "request expired")
	ErrTimeout          = NewError(408, "time out")

	ErrServerInternal     = NewError(500, "server internal error")
	ErrNotImplemented     = NewError(501, "not implemented")
	ErrServiceUnavailable = NewError(503, "service unavailable")
)
