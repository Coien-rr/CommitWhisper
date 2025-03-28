package errors

type InvalidKeyError struct {
	msg string
}

func NewInvalidKeyError(errMsg string) *InvalidKeyError {
	return &InvalidKeyError{
		msg: removeRequestID(errMsg),
	}
}

func (e *InvalidKeyError) Error() string {
	return e.msg
}

func (e *InvalidKeyError) Is(target error) bool {
	_, ok := target.(*InvalidKeyError)
	return ok
}

var ErrInvalidKey = &InvalidKeyError{}

type TooManyReqError struct {
	msg string
}

func NewTooManyReqError(errMsg string) *TooManyReqError {
	return &TooManyReqError{
		msg: errMsg,
	}
}

func (e *TooManyReqError) Error() string {
	return e.msg
}

func (e *TooManyReqError) Is(target error) bool {
	_, ok := target.(*TooManyReqError)
	return ok
}

type NotFoundError struct {
	msg string
}

func NewNotFoundError(errMsg string) *NotFoundError {
	return &NotFoundError{
		msg: removeRequestID(errMsg),
	}
}

func (e *NotFoundError) Error() string {
	return e.msg
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}
