package errors

import "github.com/go-playground/validator/v10"

var (
	// general
	// InputError ...
	InputError = NewCodeError(1, "Input error")
	// AuthError ...
	AuthError = NewCodeError(2, "Auth error")
	// DBError ...
	DBError = NewCodeError(3, "DB error")
	// NotFound ...
	NotFound = NewCodeError(4, "Not found")
	// UnknownError ...
	UnknownError = NewCodeError(100, "Unknown error")
)

// IsCodeErrorEqual check if code of err is equal to codeErr's
func IsCodeErrorEqual(err error, codeErr *CodeError) bool {
	err = Cause(err)
	if e, ok := err.(*CodeError); ok {
		if e.Code() == codeErr.Code() {
			return true
		}
	}

	return false
}

// Err2Code convert err to CodeError
func Err2Code(err error) *CodeError {
	err = Cause(err)
	if _, ok := err.(validator.ValidationErrors); ok {
		return InputError
	}
	if codeErr, ok := err.(*CodeError); ok {
		return codeErr
	}
	return UnknownError
}
