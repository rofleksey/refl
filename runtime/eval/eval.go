package eval

import (
	"context"
	"refl/ast"
	"refl/runtime"
	"refl/runtime/eventloop"
	"refl/runtime/objects"
)

func New(ctx context.Context, program *ast.Program, env *runtime.Environment, opts ...Option) *Evaluator {
	var options Options

	for _, opt := range opts {
		opt.Apply(&options)
	}

	ctx = context.WithValue(ctx, "options", options)

	evaluator := &Evaluator{
		program: program,
		env:     env,
	}

	ctx = context.WithValue(ctx, "evaluator", evaluator)

	env.Define("math", createMathObject())
	env.Define("strings", createStringObject())
	env.Define("errors", createErrorsObject())
	env.Define("io", createIoObject())
	env.Define("time", createTimeObject())
	if !options.disableEvents {
		env.Define("events", createEventsObject())
	}

	defEnvBuiltinFunc("type", env, builtinTypeFunc)
	defEnvBuiltinFunc("str", env, builtinStrFunc)
	defEnvBuiltinFunc("number", env, builtinNumberFunc)
	defEnvBuiltinFunc("len", env, builtinLenFunc)
	defEnvBuiltinFunc("range", env, builtinRangeFunc)
	defEnvBuiltinFunc("clone", env, builtinCloneFunc)

	var eventLoop *eventloop.EventLoop

	if !options.disableEvents {
		eventLoop = eventloop.New(ctx)
		evaluator.eventLoop = eventLoop
		ctx = context.WithValue(ctx, "event_loop", eventLoop)
		defEnvBuiltinFunc("eval", env, builtinEvalFunc)
	}

	env.Define("$", &globalRefObject{env: env})

	evaluator.ctx = ctx

	return evaluator
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
