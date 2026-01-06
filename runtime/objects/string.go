package objects

import (
	"iter"
	"refl/runtime"
	"strconv"
)

type String struct {
	Value string
}

func NewString(value string) *String {
	return &String{Value: value}
}

func (s *String) Type() runtime.ObjectType { return runtime.StringType }
func (s *String) String() string           { return s.Value }
func (s *String) Truthy() bool             { return len(s.Value) > 0 }
func (s *String) Equal(other runtime.Object) bool {
	if other.Type() != runtime.StringType {
		return false
	}
	return s.Value == other.(*String).Value
}
func (s *String) Clone() runtime.Object { return NewString(s.Value) }

func (s *String) Add(other runtime.Object) (runtime.Object, *runtime.Error) {
	return NewString(s.Value + other.String()), nil
}

func (s *String) Sub(other runtime.Object) (runtime.Object, *runtime.Error) {
	return nil, runtime.NewError("strings do not support subtraction", 0, 0)
}

func (s *String) Mul(other runtime.Object) (runtime.Object, *runtime.Error) {
	switch o := other.(type) {
	case *Number:
		result := ""
		count := int(o.Value)
		for i := 0; i < count; i++ {
			result += s.Value
		}
		return NewString(result), nil
	default:
		return nil, runtime.NewError("cannot multiply string by non-number", 0, 0)
	}
}

func (s *String) Div(other runtime.Object) (runtime.Object, *runtime.Error) {
	return nil, runtime.NewError("strings do not support division", 0, 0)
}

func (s *String) Mod(other runtime.Object) (runtime.Object, *runtime.Error) {
	return nil, runtime.NewError("strings do not support modulo", 0, 0)
}

func (s *String) Negate() (runtime.Object, *runtime.Error) {
	return nil, runtime.NewError("strings do not support negation", 0, 0)
}

func (s *String) LessThan(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.StringType {
		return nil, runtime.NewError("cannot compare string with non-string", 0, 0)
	}
	return NewBoolean(s.Value < other.(*String).Value), nil
}

func (s *String) GreaterThan(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.StringType {
		return nil, runtime.NewError("cannot compare string with non-string", 0, 0)
	}
	return NewBoolean(s.Value > other.(*String).Value), nil
}

func (s *String) LessThanEqual(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.StringType {
		return nil, runtime.NewError("cannot compare string with non-string", 0, 0)
	}
	return NewBoolean(s.Value <= other.(*String).Value), nil
}

func (s *String) GreaterThanEqual(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.StringType {
		return nil, runtime.NewError("cannot compare string with non-string", 0, 0)
	}
	return NewBoolean(s.Value >= other.(*String).Value), nil
}

func (s *String) Get(key runtime.Object) (runtime.Object, *runtime.Error) {
	if key.Type() != runtime.NumberType {
		return nil, runtime.NewError("string index must be a number", 0, 0)
	}

	idx := int(key.(*Number).Value)
	if idx < 0 || idx >= len(s.Value) {
		return nil, runtime.NewError("string index out of bounds", 0, 0)
	}

	return NewString(string(s.Value[idx])), nil
}

func (s *String) Set(key, value runtime.Object) *runtime.Error {
	return runtime.NewError("strings are immutable", 0, 0)
}

func (s *String) Length() int {
	return len(s.Value)
}

func (s *String) ToNumber() (runtime.Object, *runtime.Error) {
	val, err := strconv.ParseFloat(s.Value, 64)
	if err != nil {
		return NilInstance, nil
	}
	return NewNumber(val), nil
}

func (s *String) Not() runtime.Object {
	return NewBoolean(!s.Truthy())
}

func (s *String) HashKey() runtime.HashKey {
	return runtime.HashKey("str_" + s.Value)
}

func (s *String) Iterator() iter.Seq2[runtime.Object, runtime.Object] {
	return func(yield func(runtime.Object, runtime.Object) bool) {
		for i, ch := range s.Value {
			key := NewNumber(float64(i))
			value := NewString(string(ch))
			if !yield(key, value) {
				return
			}
		}
	}
}
