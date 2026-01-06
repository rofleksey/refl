package objects

import (
	"fmt"
	"iter"
	"refl/runtime"
	"sort"
)

type ReflObject struct {
	id          string
	numFields   map[float64]runtime.Object
	otherFields map[runtime.HashKey]otherFieldCarriage
}

type otherFieldCarriage struct {
	Key   runtime.Object
	Value runtime.Object
}

func NewObject() *ReflObject {
	result := &ReflObject{
		numFields:   make(map[float64]runtime.Object),
		otherFields: make(map[runtime.HashKey]otherFieldCarriage),
	}

	result.id = fmt.Sprintf("%p", result)

	return result
}

func (o *ReflObject) Type() runtime.ObjectType { return runtime.ObjectType_ }
func (o *ReflObject) String() string           { return "object" }
func (o *ReflObject) Truthy() bool             { return true }
func (o *ReflObject) Equal(other runtime.Object) bool {
	return o == other
}
func (o *ReflObject) Clone() runtime.Object {
	cloned := NewObject()
	for key, value := range o.numFields {
		cloned.numFields[key] = value.Clone()
	}
	for _, carriage := range o.otherFields {
		newKey := carriage.Key.Clone()
		newValue := carriage.Value.Clone()

		cloned.otherFields[newKey.HashKey()] = otherFieldCarriage{
			Key:   newKey,
			Value: newValue,
		}
	}
	return cloned
}

func (o *ReflObject) Add(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support addition", 0, 0)
}

func (o *ReflObject) Sub(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support subtraction", 0, 0)
}

func (o *ReflObject) Mul(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support multiplication", 0, 0)
}

func (o *ReflObject) Div(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support division", 0, 0)
}

func (o *ReflObject) Mod(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support modulo", 0, 0)
}

func (o *ReflObject) Negate() (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support negation", 0, 0)
}

func (o *ReflObject) LessThan(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support comparison", 0, 0)
}

func (o *ReflObject) GreaterThan(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support comparison", 0, 0)
}

func (o *ReflObject) LessThanEqual(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support comparison", 0, 0)
}

func (o *ReflObject) GreaterThanEqual(other runtime.Object) (runtime.Object, *runtime.Panic) {
	return nil, runtime.NewPanic("objects do not support comparison", 0, 0)
}

func (o *ReflObject) Not() runtime.Object {
	return NewBoolean(!o.Truthy())
}

func (o *ReflObject) Get(key runtime.Object) (runtime.Object, *runtime.Panic) {
	if key.Type() == runtime.NumberType {
		numKey, _ := key.(*Number)

		val, exists := o.numFields[numKey.Value]
		if !exists {
			return NilInstance, nil
		}
		return val, nil
	}

	keyStr := key.HashKey()
	carriage, exists := o.otherFields[keyStr]
	if !exists {
		return NilInstance, nil
	}

	return carriage.Value, nil
}

func (o *ReflObject) Set(key, value runtime.Object) *runtime.Panic {
	if key.Type() == runtime.NumberType {
		numKey, _ := key.(*Number)

		o.numFields[numKey.Value] = value

		return nil
	}

	keyStr := key.HashKey()

	o.otherFields[keyStr] = otherFieldCarriage{
		Key:   key,
		Value: value,
	}

	return nil
}

func (o *ReflObject) SetLiteral(key string, value runtime.Object) {
	_ = o.Set(NewString(key), value)
}

func (o *ReflObject) Length() int {
	return len(o.numFields) + len(o.otherFields)
}

func (o *ReflObject) Iterator() iter.Seq2[runtime.Object, runtime.Object] {
	return func(yield func(runtime.Object, runtime.Object) bool) {
		if len(o.numFields) > 0 {
			keys := make([]float64, 0, len(o.numFields))
			for key := range o.numFields {
				keys = append(keys, key)
			}
			sort.Float64s(keys)

			for _, key := range keys {
				keyObj := NewNumber(key)
				value := o.numFields[key]
				if !yield(keyObj, value) {
					return
				}
			}
		}

		for _, carriage := range o.otherFields {
			if !yield(carriage.Key, carriage.Value) {
				return
			}
		}
	}
}

func (o *ReflObject) HashKey() runtime.HashKey {
	return runtime.HashKey("obj_" + o.id)
}
