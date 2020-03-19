package domain

import "github.com/sirupsen/logrus"

type Operation string

type Error struct {
	Err error
	Code Code
	Operation Operation
	Severity logrus.Level
}

func NewError(operation Operation, Code Code, severity logrus.Level, err error) *Error {
	return &Error{
		Err:       err,
		Code:      Code,
		Severity:  severity,
		Operation: operation,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Stack() []Operation {
	stack := []Operation{e.Operation}
	wrappedError, ok := e.Err.(*Error)
	if !ok {
		return stack
	}
	stack = append(stack, wrappedError.Stack()...)
	return stack
}
