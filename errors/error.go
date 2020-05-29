package errors

import (
	"fmt"
)

type Error interface {
	Error() string

}

type gError struct {
	msg 		string
}

func (g* gError)Error() string {
	return g.msg
}

func New(format string, a ...interface{})Error {

	return fmt.Errorf(format, a...)
}