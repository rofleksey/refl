package eval

import (
	"context"
	"fmt"
	"refl/ast"
	"refl/runtime"
	eventloop "refl/runtime/eventloop"
	"refl/runtime/objects"
)

func Eval(ctx context.Context, program *ast.Program, env *runtime.Environment) (runtime.Object, error) {
	env.Define("math", createMathObject())
	env.Define("strings", createStringObject())
	env.Define("errors", createErrorsObject())
	env.Define("io", createIoObject())
	env.Define("time", createTimeObject())
	env.Define("events", createEventsObject())

	defEnvBuiltinFunc("type", env, builtinTypeFunc)
	defEnvBuiltinFunc("str", env, builtinStrFunc)
	defEnvBuiltinFunc("number", env, builtinNumberFunc)

	defEnvBuiltinFunc("len", env, builtinLenFunc)

	defEnvBuiltinFunc("range", env, builtinRangeFunc)

	defEnvBuiltinFunc("clone", env, builtinCloneFunc)
	defEnvBuiltinFunc("eval", env, builtinEvalFunc)

	env.Define("$", &globalRefObject{env: env})

	eventLoop := eventloop.New(ctx)
	ctx = context.WithValue(ctx, "event_loop", eventLoop)

	evaluator := &Evaluator{
		ctx: ctx,
	}

	result, err := evaluator.evalProgram(program, env)
	if err != nil {
		return result, err
	}

	eventLoop.Start()
	eventLoop.Wait()

	r := eventLoop.LastPanic()
	if r != nil {
		return nil, runtime.NewPanic(fmt.Sprintf("Event loop panic: %v", r), 0, 0)
	}

	return result, nil
}

func defLiteralBuiltinFunc(
	name string,
	obj *objects.ReflObject,
	fn func(_ context.Context, args []runtime.Object) (runtime.Object, error),
) {
	obj.SetLiteral(name, objects.NewWrapperFunction(fn))
}

func defEnvBuiltinFunc(
	name string,
	env *runtime.Environment,
	fn func(_ context.Context, args []runtime.Object) (runtime.Object, error),
) {
	env.Define(name, objects.NewWrapperFunction(fn))
}
