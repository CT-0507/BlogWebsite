package contracts

type ErrDuplicate struct {
	Code    int
	Message string
}

func (e *ErrDuplicate) Error() string {
	return e.Message
}
