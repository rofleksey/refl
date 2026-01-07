package objects

import (
	"fmt"
	"iter"
	"refl/runtime"
)

type Iterator struct {
	ID   string
	Iter iter.Seq2[runtime.Object, runtime.Object]
}

func NewIterator(iter iter.Seq2[runtime.Object, runtime.Object]) *Iterator {
	result := &Iterator{
		Iter: iter,
	}

	result.ID = fmt.Sprintf("%p", result)

	return result
}

func (it *Iterator) Type() runtime.ObjectType { return runtime.IteratorType }
func (it *Iterator) String() string           { return "iterator" }
func (it *Iterator) Truthy() bool             { return true }
func (it *Iterator) Equal(other runtime.Object) bool {
	return it == other
}
func (it *Iterator) Clone() runtime.Object { return it }
func (it *Iterator) Iterator() iter.Seq2[runtime.Object, runtime.Object] {
	return it.Iter
}

func (it *Iterator) Not() runtime.Object {
	return NewBoolean(!it.Truthy())
}

func (it *Iterator) HashKey() runtime.HashKey {
	return runtime.HashKey("iter_" + it.ID)
}
