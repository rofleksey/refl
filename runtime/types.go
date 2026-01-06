package runtime

import "iter"

type ObjectType string

const (
	NilType      ObjectType = "nil"
	NumberType   ObjectType = "number"
	StringType   ObjectType = "string"
	ObjectType_  ObjectType = "object"
	FunctionType ObjectType = "function"
	ErrorType    ObjectType = "error"

	BreakSignal    ObjectType = "break"
	ContinueSignal ObjectType = "continue"
	ReturnSignal   ObjectType = "return"
)

type HashKey string

type Object interface {
	Type() ObjectType
	String() string
	Truthy() bool
	Not() Object
	Equal(other Object) bool
	Clone() Object
	HashKey() HashKey
}

type Indexable interface {
	Get(key Object) (Object, error)
	Set(key, value Object) error
	Length() int
}

type Iterable interface {
	Iterator() iter.Seq2[Object, Object]
}

type Callable interface {
	Call(args []Object) (Object, error)
}
