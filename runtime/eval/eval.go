package eval

import (
	"context"
	"refl/ast"
	"refl/runtime"
	"refl/runtime/objects"
)

func Eval(ctx context.Context, program *ast.Program, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	runtimeObj := objects.NewObject()
	env.Define("runtime", runtimeObj)
	defLiteralBuiltinFunc(ctx, "exit", runtimeObj, builtinExitFunc)

	defEnvBuiltinFunc(ctx, "type", env, builtinTypeFunc)
	defEnvBuiltinFunc(ctx, "str", env, builtinStrFunc)
	defEnvBuiltinFunc(ctx, "number", env, builtinNumberFunc)
	defEnvBuiltinFunc(ctx, "len", env, builtinLenFunc)
	defEnvBuiltinFunc(ctx, "clone", env, builtinCloneFunc)

	env.Define("$", &globalRefObject{env: env})

	evaluator := &Evaluator{
		ctx: ctx,
	}

	return evaluator.evalProgram(program, env)
}

func defLiteralBuiltinFunc(
	ctx context.Context,
	name string,
	obj *objects.ReflObject,
	fn func(_ context.Context, args []runtime.Object) (runtime.Object, *runtime.Error),
) {
	obj.SetLiteral(name, &builtinFunction{
		ctx:  ctx,
		name: name,
		fn:   fn,
	})
}

func defEnvBuiltinFunc(
	ctx context.Context,
	name string,
	env *runtime.Environment,
	fn func(_ context.Context, args []runtime.Object) (runtime.Object, *runtime.Error),
) {
	env.Define(name, &builtinFunction{
		ctx:  ctx,
		name: name,
		fn:   fn,
	})
}
