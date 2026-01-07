package eval

import (
	"context"
	"iter"
	"refl/parser"
	"refl/runtime"
	"refl/runtime/objects"
)

func builtinTypeFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return objects.NewError("type() expects exactly 1 argument"), nil
	}
	return objects.NewString(string(args[0].Type())), nil
}

func builtinStrFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return objects.NewError("str() expects exactly 1 argument"), nil
	}
	return objects.NewString(args[0].String()), nil
}

func builtinNumberFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return objects.NewError("number() expects exactly 1 argument"), nil
	}

	switch arg := args[0].(type) {
	case *objects.Number:
		return arg, nil
	case *objects.String:
		return arg.ToNumber()
	default:
		return objects.NilInstance, nil
	}
}

func builtinLenFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return objects.NewError("len() expects exactly 1 argument"), nil
	}

	indexable, ok := args[0].(runtime.Indexable)
	if !ok {
		return objects.NewError("len() can only be called on indexable objects"), nil
	}

	return objects.NewNumber(float64(indexable.Length())), nil
}

func builtinCloneFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return objects.NewError("clone() expects exactly 1 argument"), nil
	}

	obj := args[0]

	return obj.Clone(), nil
}

func builtinEvalFunc(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 1 {
		return objects.NewError("refl() expects at least 1 argument"), nil
	}

	code := args[0].String()

	p := parser.New()

	program, err := p.Parse(code)
	if err != nil {
		return objects.NewError(err.Error()), nil
	}

	result, err := Eval(ctx, program, runtime.NewEnvironment(nil))
	if err != nil {
		return objects.NewError(err.Error()), nil
	}

	return result, nil
}

func builtinRangeFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 2 || len(args) > 3 {
		return objects.NewError("range() expects 2 or 3 arguments"), nil
	}

	if args[0].Type() != runtime.NumberType || args[1].Type() != runtime.NumberType {
		return objects.NewError("range() can only be called on numbers"), nil
	}

	start := args[0].(*objects.Number).Value
	finish := args[1].(*objects.Number).Value
	step := 1.0

	if len(args) == 3 {
		if args[2].Type() != runtime.NumberType {
			return objects.NewError("range() can only be called on numbers"), nil
		}

		step = args[2].(*objects.Number).Value
		if step == 0 {
			return objects.NewError("step is zero"), nil
		}
	}

	return objects.NewIterator(func(yield func(runtime.Object, runtime.Object) bool) {
		i := 0.0

		if step > 0 {
			for num := start; num < finish; num += step {
				if !yield(objects.NewNumber(i), objects.NewNumber(num)) {
					break
				}

				i++
			}
		} else {
			for num := start; num > finish; num += step {
				if !yield(objects.NewNumber(i), objects.NewNumber(num)) {
					break
				}

				i++
			}
		}

	}), nil
}

type globalRefObject struct {
	env *runtime.Environment
}

func (g *globalRefObject) Type() runtime.ObjectType { return runtime.ObjectType_ }
func (g *globalRefObject) String() string           { return "$" }
func (g *globalRefObject) Truthy() bool             { return true }
func (g *globalRefObject) Equal(other runtime.Object) bool {
	return g == other
}
func (g *globalRefObject) Clone() runtime.Object { return g }

func (g *globalRefObject) Get(key runtime.Object) (runtime.Object, error) {
	keyStr := key.String()
	if val, ok := g.env.Get(keyStr); ok {
		return val, nil
	}
	return objects.NilInstance, nil
}

func (g *globalRefObject) Set(key, value runtime.Object) error {
	return runtime.NewPanic("cannot modify $ object directly", 0, 0)
}

func (g *globalRefObject) Length() int {
	return 0
}

func (g *globalRefObject) Not() runtime.Object {
	return objects.NewBoolean(!g.Truthy())
}

func (g *globalRefObject) HashKey() runtime.HashKey {
	return "$"
}

func (g *globalRefObject) Iterator() iter.Seq2[runtime.Object, runtime.Object] {
	return func(yield func(runtime.Object, runtime.Object) bool) {
		seq := g.env.GlobalsIterator()

		seq(func(name string, value runtime.Object) bool {
			key := objects.NewString(name)

			if !yield(key, value) {
				return false
			}
			return true
		})
	}
}
