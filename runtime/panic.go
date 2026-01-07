package runtime

import "fmt"

type Panic struct {
	Message string
	Line    int
	Column  int
}

func (e *Panic) Error() string {
	if e.Line == 0 && e.Column == 0 {
		return e.Message
	}

	return fmt.Sprintf("Panic at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

func NewPanic(msg string, line, column int) *Panic {
	return &Panic{Message: msg, Line: line, Column: column}
}
