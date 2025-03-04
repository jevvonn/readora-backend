package errorpkg

type ErrorResp struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func NewError(message string, code int) *ErrorResp {
	return &ErrorResp{
		Message:    message,
		StatusCode: code,
	}
}

func (e *ErrorResp) Error() string {
	return e.Message
}

func (e *ErrorResp) WithCustomMessage(msg string) *ErrorResp {
	e.Message = msg
	return e
}
