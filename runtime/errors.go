package runtime

import "fmt"

type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Runtime error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func NewError(msg string, line, column int) *Error {
	return &Error{Message: msg, Line: line, Column: column}
}
