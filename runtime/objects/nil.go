package objects

import (
	"refl/runtime"
)

type Nil struct{}

var NilInstance = &Nil{}

func (n *Nil) Type() runtime.ObjectType { return runtime.NilType }
func (n *Nil) String() string           { return "nil" }
func (n *Nil) Truthy() bool             { return false }
func (n *Nil) Equal(other runtime.Object) bool {
	_, ok := other.(*Nil)
	return ok
}
func (n *Nil) Clone() runtime.Object { return n }

func (n *Nil) Add(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support addition", 0, 0)
}

func (n *Nil) Sub(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support subtraction", 0, 0)
}

func (n *Nil) Mul(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support multiplication", 0, 0)
}

func (n *Nil) Div(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support division", 0, 0)
}

func (n *Nil) Mod(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support modulo", 0, 0)
}

func (n *Nil) Negate() (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support negation", 0, 0)
}

func (n *Nil) LessThan(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support comparison", 0, 0)
}

func (n *Nil) GreaterThan(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support comparison", 0, 0)
}

func (n *Nil) LessThanEqual(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support comparison", 0, 0)
}

func (n *Nil) GreaterThanEqual(other runtime.Object) (runtime.Object, error) {
	return nil, runtime.NewPanic("nil does not support comparison", 0, 0)
}

func (n *Nil) Not() runtime.Object {
	return &Number{Value: 1}
}

func (n *Nil) HashKey() runtime.HashKey {
	return "nil"
}
