package objects

import (
	"fmt"
	"refl/runtime"
)

type UserError struct {
	ID   string
	text string
}

func NewError(text string) *UserError {
	result := &UserError{text: text}
	result.ID = fmt.Sprintf("%p", result)

	return result
}

func (e *UserError) Type() runtime.ObjectType { return runtime.ErrorType }
func (e *UserError) String() string           { return e.text }
func (e *UserError) Truthy() bool             { return true }
func (e *UserError) Equal(other runtime.Object) bool {
	return e == other
}
func (e *UserError) Clone() runtime.Object { return NewError(e.text) }

func (e *UserError) Add(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support addition", 0, 0)
}

func (e *UserError) Sub(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support subtraction", 0, 0)
}

func (e *UserError) Mul(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support multiplication", 0, 0)
}

func (e *UserError) Div(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support division", 0, 0)
}

func (e *UserError) Mod(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support modulo", 0, 0)
}

func (e *UserError) Negate() (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support negation", 0, 0)
}

func (e *UserError) LessThan(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support comparison", 0, 0)
}

func (e *UserError) GreaterThan(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support comparison", 0, 0)
}

func (e *UserError) LessThanEqual(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support comparison", 0, 0)
}

func (e *UserError) GreaterThanEqual(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("errors do not support comparison", 0, 0)
}

func (e *UserError) Not() runtime.Object {
	return &Number{Value: 0}
}

func (e *UserError) HashKey() runtime.HashKey {
	return runtime.HashKey("err_" + e.ID)
}
