package eval

import (
	"context"
	"iter"
	"refl/parser"
	"refl/runtime"
	"refl/runtime/objects"
)

type builtinFunction struct {
	ctx  context.Context
	name string
	fn   func(context.Context, []runtime.Object) (runtime.Object, error)
}

func (f *builtinFunction) Type() runtime.ObjectType { return runtime.FunctionType }
func (f *builtinFunction) String() string           { return "function" }
func (f *builtinFunction) Truthy() bool             { return true }
func (f *builtinFunction) Equal(other runtime.Object) bool {
	return f == other
}
func (f *builtinFunction) Clone() runtime.Object { return f }

func (f *builtinFunction) Call(args []runtime.Object) (runtime.Object, error) {
	return f.fn(f.ctx, args)
}
func (f *builtinFunction) Not() runtime.Object {
	return objects.NewBoolean(!f.Truthy())
}
func (f *builtinFunction) HashKey() runtime.HashKey {
	return runtime.HashKey("builtin_" + f.name)
}

func builtinTypeFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("type() expects exactly 1 argument", 0, 0)
	}
	return objects.NewString(string(args[0].Type())), nil
}

func builtinStrFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("str() expects exactly 1 argument", 0, 0)
	}
	return objects.NewString(args[0].String()), nil
}

func builtinNumberFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("number() expects exactly 1 argument", 0, 0)
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
		return nil, runtime.NewPanic("len() expects exactly 1 argument", 0, 0)
	}

	indexable, ok := args[0].(runtime.Indexable)
	if !ok {
		return nil, runtime.NewPanic("len() can only be called on indexable objects", 0, 0)
	}

	return objects.NewNumber(float64(indexable.Length())), nil
}

func builtinCloneFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("clone() expects exactly 1 argument", 0, 0)
	}

	obj := args[0]

	return obj.Clone(), nil
}

func builtinEvalFunc(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 1 {
		return nil, runtime.NewPanic("refl() expects at least 1 argument", 0, 0)
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

func builtinPanicFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	msg := "exit called"
	if len(args) > 0 {
		msg = args[0].String()
	}
	return nil, runtime.NewPanic(msg, 0, 0)
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
