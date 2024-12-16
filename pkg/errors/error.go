package errors

type InvalidKeyError struct{}

func (e *InvalidKeyError) Error() string {
	return "ERROR: InvalidAPIKey"
}

func (e *InvalidKeyError) Is(target error) bool {
	_, ok := target.(*InvalidKeyError)
	return ok
}

var ErrInvalidKey = &InvalidKeyError{}
