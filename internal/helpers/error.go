package helpers

import "errors"

type Error struct {
	Code  int
	Error error
}

func ErrorBuilder(code int, errString string) *Error {
	return &Error{
		Code:  code,
		Error: errors.New(errString),
	}
}

func (err *Error) String() string {
	return err.Error.Error()
}
