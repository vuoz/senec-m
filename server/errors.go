package server

type HandlerError struct {
	reason string
}

func (e *HandlerError) Error() string {
	return e.reason
}
func NewHandlerError(reason string) error {
	return &HandlerError{
		reason: reason,
	}
}
