package errno

type Errno struct {
	Code    int
	Message string
}

func (e *Errno) Error() string {
	return e.Message
}

func (e *Errno) SetMessage(msg string) *Errno {
	e.Message = msg
	return e
}
