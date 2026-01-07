package objects

import (
	"context"
	"refl/ast"
	"refl/runtime"
	"refl/runtime/eventloop"
)

type Evaluator interface {
	Context() context.Context
	EvalBlock(block *ast.BlockStatement, env *runtime.Environment) (runtime.Object, error)
	FireEvent(event string, args []runtime.Object)
	EnqueueTask(task eventloop.Task)
}
