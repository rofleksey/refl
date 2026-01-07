package objects

import (
	"context"
	"refl/ast"
	"refl/runtime"
)

type Evaluator interface {
	Context() context.Context
	EvalBlock(block *ast.BlockStatement, env *runtime.Environment) (runtime.Object, error)
}
