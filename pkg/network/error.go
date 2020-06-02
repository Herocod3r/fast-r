package network

type Error struct {
	InternalError error
}

func (e *Error) Error() string {
	return e.InternalError.Error()
}
