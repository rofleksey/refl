package objects

import (
	"context"
	"fmt"
	"refl/runtime"
)

type WrapperFunction struct {
	id string
	fn func(context.Context, []runtime.Object) (runtime.Object, error)
}

func NewWrapperFunction(fn func(context.Context, []runtime.Object) (runtime.Object, error)) *WrapperFunction {
	result := &WrapperFunction{
		fn: fn,
	}

	result.id = fmt.Sprintf("%p", result)

	return result
}

func (f *WrapperFunction) Type() runtime.ObjectType { return runtime.FunctionType }
func (f *WrapperFunction) String() string           { return "function" }
func (f *WrapperFunction) Truthy() bool             { return true }
func (f *WrapperFunction) Equal(other runtime.Object) bool {
	return f == other
}
func (f *WrapperFunction) Clone() runtime.Object { return f }

func (f *WrapperFunction) Call(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
	return f.fn(ctx, args)
}
func (f *WrapperFunction) Not() runtime.Object {
	return NewBoolean(!f.Truthy())
}

func (f *WrapperFunction) HashKey() runtime.HashKey {
	return runtime.HashKey("builtin_" + f.id)
}
