package objects

import (
	"fmt"
	"refl/ast"
	"refl/runtime"
)

type Function struct {
	ID         string
	Parameters []string
	Body       *ast.BlockStatement
	Env        *runtime.Environment

	evaluator func(block *ast.BlockStatement, env *runtime.Environment) (runtime.Object, error)
}

func NewFunction(
	params []string,
	body *ast.BlockStatement,
	env *runtime.Environment,
	evaluator func(block *ast.BlockStatement, env *runtime.Environment) (runtime.Object, error),
) *Function {
	result := &Function{
		Parameters: params,
		Body:       body,
		Env:        runtime.NewEnvironment(env),
		evaluator:  evaluator,
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

func (f *Function) Call(args []runtime.Object) (runtime.Object, error) {
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

	return f.evaluator(f.Body, funcEnv)
}

func (f *Function) Not() runtime.Object {
	return NewBoolean(!f.Truthy())
}

func (f *Function) HashKey() runtime.HashKey {
	return runtime.HashKey("fun_" + f.ID)
}
