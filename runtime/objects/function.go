package objects

import (
	"context"
	"fmt"
	"refl/ast"
	"refl/runtime"
)

type Function struct {
	ID         string
	Parameters []string
	Body       *ast.BlockStatement
	Env        *runtime.Environment
}

func NewFunction(
	params []string,
	body *ast.BlockStatement,
	env *runtime.Environment,
) *Function {
	result := &Function{
		Parameters: params,
		Body:       body,
		Env:        runtime.NewEnvironment(env),
	}

	result.ID = fmt.Sprintf("%p", result)

	return result
}

func (f *Function) Type() runtime.ObjectType { return runtime.FunctionType }
func (f *Function) String() string           { return "function" }
func (f *Function) Truthy() bool             { return true }
func (f *Function) Equal(other runtime.Object) bool {
	_, ok := other.(*Function)
	return ok && f == other
}
func (f *Function) Clone() runtime.Object { return f }

func (f *Function) Call(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
	evaluator, ok := ctx.Value("evaluator").(Evaluator)
	if !ok {
		return nil, runtime.NewPanic("evaluator not found in context", 0, 0)
	}

	funcEnv := runtime.NewEnvironment(f.Env)

	for i, param := range f.Parameters {
		if i < len(args) {
			funcEnv.Define(param, args[i])
		} else {
			funcEnv.Define(param, NilInstance)
		}
	}

	argsObj := NewObject()
	for i, arg := range args {
		_ = argsObj.Set(NewNumber(float64(i)), arg)
	}
	funcEnv.Define("args", argsObj)

	return evaluator.EvalBlock(f.Body, funcEnv)
}

func (f *Function) Not() runtime.Object {
	return NewBoolean(!f.Truthy())
}

func (f *Function) HashKey() runtime.HashKey {
	return runtime.HashKey("fun_" + f.ID)
}
