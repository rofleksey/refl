package eval

import (
	"context"
	"fmt"
	"refl/ast"
	"refl/runtime"
	"refl/runtime/eventloop"
	"refl/runtime/objects"
)

type Evaluator struct {
	ctx       context.Context
	program   *ast.Program
	env       *runtime.Environment
	eventLoop *eventloop.EventLoop
}

func (e *Evaluator) Context() context.Context {
	return e.ctx
}

func (e *Evaluator) FireEvent(event string, args []runtime.Object) {
	if e.eventLoop != nil {
		e.eventLoop.Fire(event, args)
	}
}

func (e *Evaluator) Run() (runtime.Object, error) {
	result, err := e.evalProgram(e.program, e.env)
	if err != nil {
		return result, err
	}

	if e.eventLoop != nil {
		e.eventLoop.Start()
		e.eventLoop.Wait()

		r := e.eventLoop.LastPanic()
		if r != nil {
			return nil, runtime.NewPanic(fmt.Sprintf("Event loop panic: %v", r), 0, 0)
		}
	}

	return result, nil
}

func (e *Evaluator) evalGeneric(node ast.Node, env *runtime.Environment) (runtime.Object, error) {
	switch n := node.(type) {
	case *ast.Program:
		return e.evalProgram(n, env)
	case *ast.BlockStatement:
		return e.EvalBlock(n, env)
	case *ast.VarDeclaration:
		return e.evalVarDeclaration(n, env)
	case *ast.ExpressionStatement:
		return e.evalExpressionStatement(n, env)
	case *ast.IfStatement:
		return e.evalIfStatement(n, env)
	case *ast.WhileStatement:
		return e.evalWhileStatement(n, env)
	case *ast.ForStatement:
		return e.evalForStatement(n, env)
	case *ast.ReturnStatement:
		return e.evalReturnStatement(n, env)
	case *ast.BreakStatement:
		return e.evalBreakStatement()
	case *ast.ContinueStatement:
		return e.evalContinueStatement()
	case *ast.Identifier:
		return e.evalIdentifier(n, env)
	case *ast.NumberLiteral:
		return e.evalNumberLiteral(n)
	case *ast.StringLiteral:
		return e.evalStringLiteral(n)
	case *ast.RawStringLiteral:
		return e.evalRawStringLiteral(n)
	case *ast.NilLiteral:
		return e.evalNilLiteral()
	case *ast.ObjectLiteral:
		return e.evalObjectLiteral(n, env)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(n, env)
	case *ast.FunctionLiteral:
		return e.evalFunctionLiteral(n, env)
	case *ast.MemberDot:
		return e.evalMemberDot(n, env)
	case *ast.MemberBracket:
		return e.evalMemberBracket(n, env)
	case *ast.FunctionCall:
		return e.evalFunctionCall(n, env)
	case *ast.MethodCall:
		return e.evalMethodCall(n, env)
	case *ast.UnaryExpression:
		return e.evalUnaryExpression(n, env)
	case *ast.BinaryExpression:
		return e.evalBinaryExpression(n, env)
	case *ast.Assignment:
		return e.evalAssignment(n, env)
	default:
		return nil, runtime.NewPanic(fmt.Sprintf("unknown node type: %T", node), 0, 0)
	}
}

func (e *Evaluator) evalProgram(program *ast.Program, env *runtime.Environment) (runtime.Object, error) {
	var result runtime.Object = objects.NilInstance

	for _, stmt := range program.Statements {
		val, err := e.evalGeneric(stmt, env)
		if err != nil {
			return nil, err
		}
		result = val
	}

	return result, nil
}

func (e *Evaluator) EvalBlock(block *ast.BlockStatement, env *runtime.Environment) (runtime.Object, error) {
	var result runtime.Object = objects.NilInstance
	blockEnv := runtime.NewEnvironment(env)

	for _, stmt := range block.Statements {
		val, err := e.evalGeneric(stmt, blockEnv)
		if err != nil {
			return nil, err
		}

		switch val.(type) {
		case *objects.Nil:
			// Continue execution
		default:
			result = val
		}

		// Check for control flow signals
		if _, isBreak := val.(*objects.BreakSignal); isBreak {
			return &objects.BreakSignal{}, nil
		}
		if _, isContinue := val.(*objects.ContinueSignal); isContinue {
			return &objects.ContinueSignal{}, nil
		}
		if ret, isReturn := val.(*objects.ReturnSignal); isReturn {
			return ret, nil
		}
	}

	return result, nil
}

func (e *Evaluator) evalVarDeclaration(vd *ast.VarDeclaration, env *runtime.Environment) (runtime.Object, error) {
	var value runtime.Object = objects.NilInstance

	if vd.Value != nil {
		val, err := e.evalGeneric(vd.Value, env)
		if err != nil {
			return nil, err
		}
		value = val
	}

	env.Define(vd.Name, value)
	return value, nil
}

func (e *Evaluator) evalExpressionStatement(es *ast.ExpressionStatement, env *runtime.Environment) (runtime.Object, error) {
	return e.evalGeneric(es.Expression, env)
}

func (e *Evaluator) evalIfStatement(is *ast.IfStatement, env *runtime.Environment) (runtime.Object, error) {
	cond, err := e.evalGeneric(is.Condition, env)
	if err != nil {
		return nil, err
	}

	if cond.Truthy() {
		result, err := e.EvalBlock(is.Then, env)
		if err != nil {
			return nil, err
		}

		if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
			return ret, nil
		}

		return result, nil
	}

	for _, elif := range is.Elif {
		cond, err := e.evalGeneric(elif.Condition, env)
		if err != nil {
			return nil, err
		}

		if cond.Truthy() {
			result, err := e.EvalBlock(elif.Body, env)
			if err != nil {
				return nil, err
			}

			if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
				return ret, nil
			}

			return result, nil
		}
	}

	if is.Else != nil {
		result, err := e.EvalBlock(is.Else, env)
		if err != nil {
			return nil, err
		}

		if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
			return ret, nil
		}

		return result, nil
	}

	return objects.NilInstance, nil
}

func (e *Evaluator) evalWhileStatement(ws *ast.WhileStatement, env *runtime.Environment) (runtime.Object, error) {
	var result runtime.Object = objects.NilInstance

	for {
		select {
		case <-e.ctx.Done():
			return nil, runtime.NewPanic("context cancelled", 0, 0)
		default:
		}

		cond, err := e.evalGeneric(ws.Condition, env)
		if err != nil {
			return nil, err
		}

		if !cond.Truthy() {
			break
		}

		val, err := e.EvalBlock(ws.Body, env)
		if err != nil {
			return nil, err
		}

		switch v := val.(type) {
		case *objects.BreakSignal:
			return objects.NilInstance, nil
		case *objects.ContinueSignal:
			continue
		case *objects.ReturnSignal:
			return v, nil
		default:
			result = v
		}
	}

	return result, nil
}

func (e *Evaluator) evalForStatement(fs *ast.ForStatement, env *runtime.Environment) (funcRes runtime.Object, funcErr error) {
	defer func() {
		if r := recover(); r != nil {
			funcErr = runtime.NewPanic(fmt.Sprintf("iteration error: %v", r), fs.Pos.Line, fs.Pos.Column)
		}
	}()

	obj, err := e.evalGeneric(fs.Object, env)
	if err != nil {
		return nil, err
	}

	iterable, ok := obj.(runtime.Iterable)
	if !ok {
		return nil, runtime.NewPanic("cannot iterate over non-iterable object", fs.Pos.Line, fs.Pos.Column)
	}

	var result runtime.Object = objects.NilInstance

	for key, value := range iterable.Iterator() {
		select {
		case <-e.ctx.Done():
			return nil, runtime.NewPanic("context cancelled", 0, 0)
		default:
		}

		forEnv := runtime.NewEnvironment(env)
		forEnv.Define(fs.Key, key)
		if fs.Value != "" {
			forEnv.Define(fs.Value, value)
		}

		val, err := e.EvalBlock(fs.Body, forEnv)
		if err != nil {
			return nil, err
		}

		switch v := val.(type) {
		case *objects.BreakSignal:
			return objects.NilInstance, nil
		case *objects.ContinueSignal:
			continue
		case *objects.ReturnSignal:
			return v, nil
		default:
			result = v
		}
	}

	return result, nil
}

func (e *Evaluator) evalReturnStatement(rs *ast.ReturnStatement, env *runtime.Environment) (runtime.Object, error) {
	var value runtime.Object = objects.NilInstance

	if rs.Value != nil {
		val, err := e.evalGeneric(rs.Value, env)
		if err != nil {
			return nil, err
		}
		value = val
	}

	return &objects.ReturnSignal{Value: value}, nil
}

func (e *Evaluator) evalBreakStatement() (runtime.Object, error) {
	return &objects.BreakSignal{}, nil
}

func (e *Evaluator) evalContinueStatement() (runtime.Object, error) {
	return &objects.ContinueSignal{}, nil
}

func (e *Evaluator) evalIdentifier(id *ast.Identifier, env *runtime.Environment) (runtime.Object, error) {
	val, ok := env.Get(id.Name)
	if !ok {
		return objects.NilInstance, nil
	}
	return val, nil
}

func (e *Evaluator) evalNumberLiteral(nl *ast.NumberLiteral) (runtime.Object, error) {
	return objects.NewNumber(nl.Value), nil
}

func (e *Evaluator) evalStringLiteral(sl *ast.StringLiteral) (runtime.Object, error) {
	return objects.NewString(sl.Value), nil
}

func (e *Evaluator) evalRawStringLiteral(rsl *ast.RawStringLiteral) (runtime.Object, error) {
	return objects.NewString(rsl.Value), nil
}

func (e *Evaluator) evalNilLiteral() (runtime.Object, error) {
	return objects.NilInstance, nil
}

func (e *Evaluator) evalObjectLiteral(ol *ast.ObjectLiteral, env *runtime.Environment) (runtime.Object, error) {
	obj := objects.NewObject()

	for key, expr := range ol.Properties {
		val, err := e.evalGeneric(expr, env)
		if err != nil {
			return nil, err
		}

		_ = obj.Set(objects.NewString(key), val)
	}

	return obj, nil
}

func (e *Evaluator) evalArrayLiteral(al *ast.ArrayLiteral, env *runtime.Environment) (runtime.Object, error) {
	obj := objects.NewObject()

	for i, expr := range al.Elements {
		val, err := e.evalGeneric(expr, env)
		if err != nil {
			return nil, err
		}

		_ = obj.Set(objects.NewNumber(float64(i)), val)
	}

	return obj, nil
}

func (e *Evaluator) evalFunctionLiteral(fl *ast.FunctionLiteral, env *runtime.Environment) (runtime.Object, error) {
	return objects.NewFunction(fl.Parameters, fl.Body, runtime.NewEnvironment(env)), nil
}

func (e *Evaluator) evalMemberDot(md *ast.MemberDot, env *runtime.Environment) (runtime.Object, error) {
	obj, err := e.evalGeneric(md.Object, env)
	if err != nil {
		return nil, err
	}

	indexable, ok := obj.(runtime.Indexable)
	if !ok {
		return nil, runtime.NewPanic("cannot access member of non-indexable object", md.Pos.Line, md.Pos.Column)
	}

	key := objects.NewString(md.Member)
	return indexable.Get(key)
}

func (e *Evaluator) evalMemberBracket(mb *ast.MemberBracket, env *runtime.Environment) (runtime.Object, error) {
	obj, err := e.evalGeneric(mb.Object, env)
	if err != nil {
		return nil, err
	}

	key, err := e.evalGeneric(mb.Member, env)
	if err != nil {
		return nil, err
	}

	indexable, ok := obj.(runtime.Indexable)
	if !ok {
		return nil, runtime.NewPanic("cannot access member of non-indexable object", mb.Pos.Line, mb.Pos.Column)
	}

	return indexable.Get(key)
}

func (e *Evaluator) evalFunctionCall(fc *ast.FunctionCall, env *runtime.Environment) (runtime.Object, error) {
	// Evaluate function expression
	fnExpr, err := e.evalGeneric(fc.Function, env)
	if err != nil {
		return nil, err
	}

	if fnExpr == nil {
		return nil, runtime.NewPanic("nil is not callable", fc.Pos.Line, fc.Pos.Column)
	}

	// Check if it's callable
	callable, ok := fnExpr.(runtime.Callable)
	if !ok {
		return nil, runtime.NewPanic("attempt to call non-function", fc.Pos.Line, fc.Pos.Column)
	}

	// Evaluate arguments
	var args []runtime.Object
	for _, argExpr := range fc.Arguments {
		arg, err := e.evalGeneric(argExpr, env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	// Call the function
	result, err := callable.Call(e.ctx, args)
	if err != nil {
		return nil, err
	}

	// Unwrap return signals
	if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
		return ret.Value, nil
	}

	return result, nil
}

func (e *Evaluator) evalMethodCall(mc *ast.MethodCall, env *runtime.Environment) (runtime.Object, error) {
	// Evaluate object
	obj, err := e.evalGeneric(mc.Object, env)
	if err != nil {
		return nil, err
	}

	// Get method from object
	indexable, ok := obj.(runtime.Indexable)
	if !ok {
		return nil, runtime.NewPanic("cannot access method of non-indexable object", mc.Pos.Line, mc.Pos.Column)
	}

	method, err := indexable.Get(objects.NewString(mc.Method))
	if err != nil {
		return nil, err
	}

	if method == nil {
		return nil, runtime.NewPanic("method '"+mc.Method+"' not found", mc.Pos.Line, mc.Pos.Column)
	}

	// Check if method is callable
	callable, ok := method.(runtime.Callable)
	if !ok {
		return nil, runtime.NewPanic("attempt to call non-function method", mc.Pos.Line, mc.Pos.Column)
	}

	// Evaluate arguments
	var args []runtime.Object
	for _, argExpr := range mc.Arguments {
		arg, err := e.evalGeneric(argExpr, env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	// Call method with object as first argument
	allArgs := append([]runtime.Object{obj}, args...)
	result, err := callable.Call(e.ctx, allArgs)
	if err != nil {
		return nil, err
	}

	// Unwrap return signals
	if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
		return ret.Value, nil
	}

	return result, nil
}

func (e *Evaluator) evalUnaryExpression(ue *ast.UnaryExpression, env *runtime.Environment) (runtime.Object, error) {
	right, err := e.evalGeneric(ue.Right, env)
	if err != nil {
		return nil, err
	}

	switch ue.Operator {
	case "!":
		return right.Not(), nil
	case "-":
		// Check if it's a number
		if num, ok := right.(*objects.Number); ok {
			return num.Negate()
		}
		return nil, runtime.NewPanic("cannot negate non-number", ue.Pos.Line, ue.Pos.Column)
	default:
		return nil, runtime.NewPanic(fmt.Sprintf("unknown operator: %s", ue.Operator), ue.Pos.Line, ue.Pos.Column)
	}
}

func (e *Evaluator) evalBinaryExpression(be *ast.BinaryExpression, env *runtime.Environment) (runtime.Object, error) {
	left, err := e.evalGeneric(be.Left, env)
	if err != nil {
		return nil, err
	}

	if be.Operator == "||" && left.Truthy() {
		return left, nil
	}

	if be.Operator == "&&" && !left.Truthy() {
		return left, nil
	}

	right, err := e.evalGeneric(be.Right, env)
	if err != nil {
		return nil, err
	}

	switch be.Operator {
	case "+":
		if num, ok := left.(*objects.Number); ok {
			return num.Add(right)
		}
		if str, ok := left.(*objects.String); ok {
			return str.Add(right)
		}
	case "-":
		if num, ok := left.(*objects.Number); ok {
			return num.Sub(right)
		}
	case "*":
		if num, ok := left.(*objects.Number); ok {
			return num.Mul(right)
		}
		if str, ok := left.(*objects.String); ok {
			return str.Mul(right)
		}
	case "/":
		if num, ok := left.(*objects.Number); ok {
			return num.Div(right)
		}
	case "%":
		if num, ok := left.(*objects.Number); ok {
			return num.Mod(right)
		}
	case "<":
		if num, ok := left.(*objects.Number); ok {
			return num.LessThan(right)
		}
		if str, ok := left.(*objects.String); ok {
			return str.LessThan(right)
		}
	case ">":
		if num, ok := left.(*objects.Number); ok {
			return num.GreaterThan(right)
		}
		if str, ok := left.(*objects.String); ok {
			return str.GreaterThan(right)
		}
	case "<=":
		if num, ok := left.(*objects.Number); ok {
			return num.LessThanEqual(right)
		}
		if str, ok := left.(*objects.String); ok {
			return str.LessThanEqual(right)
		}
	case ">=":
		if num, ok := left.(*objects.Number); ok {
			return num.GreaterThanEqual(right)
		}
		if str, ok := left.(*objects.String); ok {
			return str.GreaterThanEqual(right)
		}
	case "==":
		return objects.NewBoolean(left.Equal(right)), nil
	case "!=":
		return objects.NewBoolean(!left.Equal(right)), nil
	case "&&":
		return right, nil
	case "||":
		return right, nil
	default:
		return nil, runtime.NewPanic(fmt.Sprintf("unknown operator: %s", be.Operator), be.Pos.Line, be.Pos.Column)
	}

	return nil, runtime.NewPanic(fmt.Sprintf("cannot apply operator %s to types %s and %s",
		be.Operator, left.Type(), right.Type()), be.Pos.Line, be.Pos.Column)
}

func (e *Evaluator) evalAssignment(a *ast.Assignment, env *runtime.Environment) (runtime.Object, error) {
	right, err := e.evalGeneric(a.Right, env)
	if err != nil {
		return nil, err
	}

	switch left := a.Left.(type) {
	case *ast.Identifier:
		env.Set(left.Name, right)
		return right, nil
	case *ast.MemberDot:
		obj, err := e.evalGeneric(left.Object, env)
		if err != nil {
			return nil, err
		}

		indexable, ok := obj.(runtime.Indexable)
		if !ok {
			return nil, runtime.NewPanic("cannot assign to member of non-indexable object", left.Pos.Line, left.Pos.Column)
		}

		key := objects.NewString(left.Member)
		if err := indexable.Set(key, right); err != nil {
			return nil, err
		}
		return right, nil
	case *ast.MemberBracket:
		obj, err := e.evalGeneric(left.Object, env)
		if err != nil {
			return nil, err
		}

		key, err := e.evalGeneric(left.Member, env)
		if err != nil {
			return nil, err
		}

		indexable, ok := obj.(runtime.Indexable)
		if !ok {
			return nil, runtime.NewPanic("cannot assign to member of non-indexable object", left.Pos.Line, left.Pos.Column)
		}

		if err := indexable.Set(key, right); err != nil {
			return nil, err
		}
		return right, nil
	default:
		return nil, runtime.NewPanic("invalid assignment target", a.Pos.Line, a.Pos.Column)
	}
}
