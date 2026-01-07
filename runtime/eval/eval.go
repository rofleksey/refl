package eval

import (
	"context"
	"refl/ast"
	"refl/runtime"
	"refl/runtime/objects"
)

func Eval(ctx context.Context, program *ast.Program, env *runtime.Environment) (runtime.Object, error) {
	env.Define("math", createMathObject())
	env.Define("strings", createStringObject())
	env.Define("io", createIoObject())

	defEnvBuiltinFunc("panic", env, builtinPanicFunc)

	defEnvBuiltinFunc("newerr", env, builtinNewErrFunc)
	defEnvBuiltinFunc("errfmt", env, builtinErrFmtFunc)
	defEnvBuiltinFunc("iserr", env, builtinIsErrFunc)

	defEnvBuiltinFunc("type", env, builtinTypeFunc)
	defEnvBuiltinFunc("str", env, builtinStrFunc)
	defEnvBuiltinFunc("number", env, builtinNumberFunc)

	defEnvBuiltinFunc("len", env, builtinLenFunc)

	defEnvBuiltinFunc("range", env, builtinRangeFunc)

	defEnvBuiltinFunc("clone", env, builtinCloneFunc)
	defEnvBuiltinFunc("eval", env, builtinEvalFunc)

	env.Define("$", &globalRefObject{env: env})

	evaluator := &Evaluator{
		ctx: ctx,
	}

	return evaluator.evalProgram(program, env)
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
