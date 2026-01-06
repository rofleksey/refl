package runtime

import "fmt"

type Panic struct {
	Message string
	Line    int
	Column  int
}

func (e *Panic) Error() string {
	return fmt.Sprintf("Panic at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func NewPanic(msg string, line, column int) *Panic {
	return &Panic{Message: msg, Line: line, Column: column}
}
