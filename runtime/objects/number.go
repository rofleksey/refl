package objects

import (
	"math"
	"refl/runtime"
	"strconv"
)

type Number struct {
	Value float64
}

func NewNumber(value float64) *Number {
	return &Number{Value: value}
}

func (n *Number) Type() runtime.ObjectType { return runtime.NumberType }
func (n *Number) String() string           { return strconv.FormatFloat(n.Value, 'f', -1, 64) }
func (n *Number) Truthy() bool             { return n.Value != 0 }
func (n *Number) Equal(other runtime.Object) bool {
	if other.Type() != runtime.NumberType {
		return false
	}
	return n.Value == other.(*Number).Value
}
func (n *Number) Clone() runtime.Object { return NewNumber(n.Value) }

func (n *Number) Add(other runtime.Object) (runtime.Object, *runtime.Error) {
	switch o := other.(type) {
	case *Number:
		return NewNumber(n.Value + o.Value), nil
	case *String:
		return NewString(n.String() + o.Value), nil
	default:
		return nil, runtime.NewError("cannot add number to "+string(other.Type()), 0, 0)
	}
}

func (n *Number) Sub(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot subtract non-number from number", 0, 0)
	}
	return NewNumber(n.Value - other.(*Number).Value), nil
}

func (n *Number) Mul(other runtime.Object) (runtime.Object, *runtime.Error) {
	switch o := other.(type) {
	case *Number:
		return NewNumber(n.Value * o.Value), nil
	case *String:
		result := ""
		count := int(n.Value)
		for i := 0; i < count; i++ {
			result += o.Value
		}
		return NewString(result), nil
	default:
		return nil, runtime.NewError("cannot multiply number by "+string(other.Type()), 0, 0)
	}
}

func (n *Number) Div(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot divide number by non-number", 0, 0)
	}
	divisor := other.(*Number).Value
	if divisor == 0 {
		return nil, runtime.NewError("division by zero", 0, 0)
	}
	return NewNumber(n.Value / divisor), nil
}

func (n *Number) Mod(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot modulo number by non-number", 0, 0)
	}
	mod := other.(*Number).Value
	if mod == 0 {
		return nil, runtime.NewError("modulo by zero", 0, 0)
	}
	return NewNumber(math.Mod(n.Value, mod)), nil
}

func (n *Number) Negate() (runtime.Object, *runtime.Error) {
	return NewNumber(-n.Value), nil
}

func (n *Number) LessThan(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot compare number with non-number", 0, 0)
	}
	return NewBoolean(n.Value < other.(*Number).Value), nil
}

func (n *Number) GreaterThan(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot compare number with non-number", 0, 0)
	}
	return NewBoolean(n.Value > other.(*Number).Value), nil
}

func (n *Number) LessThanEqual(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot compare number with non-number", 0, 0)
	}
	return NewBoolean(n.Value <= other.(*Number).Value), nil
}

func (n *Number) GreaterThanEqual(other runtime.Object) (runtime.Object, *runtime.Error) {
	if other.Type() != runtime.NumberType {
		return nil, runtime.NewError("cannot compare number with non-number", 0, 0)
	}
	return NewBoolean(n.Value >= other.(*Number).Value), nil
}

func (n *Number) Not() runtime.Object {
	return NewBoolean(!n.Truthy())
}

func (n *Number) HashKey() runtime.HashKey {
	return runtime.HashKey("num_" + strconv.FormatFloat(n.Value, 'f', -1, 64))
}

func NewBoolean(value bool) *Number {
	if value {
		return &Number{Value: 1}
	}

	return &Number{Value: 0}
}
