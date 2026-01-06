package objects

import (
	"refl/runtime"
)

type BreakSignal struct{}

func (b *BreakSignal) Type() runtime.ObjectType { return runtime.BreakSignal }
func (b *BreakSignal) String() string           { return "break" }
func (b *BreakSignal) Truthy() bool             { return false }
func (b *BreakSignal) Equal(other runtime.Object) bool {
	_, ok := other.(*BreakSignal)
	return ok
}
func (b *BreakSignal) Clone() runtime.Object { return b }

func (b *BreakSignal) Not() runtime.Object {
	return NewBoolean(!b.Truthy())
}

func (b *BreakSignal) HashKey() runtime.HashKey {
	return "break"
}

type ContinueSignal struct{}

func (c *ContinueSignal) Type() runtime.ObjectType { return runtime.ContinueSignal }
func (c *ContinueSignal) String() string           { return "continue" }
func (c *ContinueSignal) Truthy() bool             { return false }
func (c *ContinueSignal) Equal(other runtime.Object) bool {
	_, ok := other.(*ContinueSignal)
	return ok
}
func (c *ContinueSignal) Clone() runtime.Object { return c }

func (c *ContinueSignal) Not() runtime.Object {
	return NewBoolean(!c.Truthy())
}

func (c *ContinueSignal) HashKey() runtime.HashKey {
	return "continue"
}

type ReturnSignal struct {
	Value runtime.Object
}

func (r *ReturnSignal) Type() runtime.ObjectType { return runtime.ReturnSignal }
func (r *ReturnSignal) String() string           { return "return" }
func (r *ReturnSignal) Truthy() bool             { return false }
func (r *ReturnSignal) Equal(other runtime.Object) bool {
	_, ok := other.(*ReturnSignal)
	return ok
}
func (r *ReturnSignal) Clone() runtime.Object { return &ReturnSignal{Value: r.Value.Clone()} }

func (r *ReturnSignal) Not() runtime.Object {
	return NewBoolean(!r.Truthy())
}

func (r *ReturnSignal) HashKey() runtime.HashKey {
	return "return"
}
