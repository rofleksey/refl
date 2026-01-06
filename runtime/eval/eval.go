package eval

import (
	"fmt"
	"refl/ast"
	"refl/runtime"
	"refl/runtime/objects"
)

func evalGeneric(node ast.Node, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.BlockStatement:
		return evalBlock(n, env)
	case *ast.VarDeclaration:
		return evalVarDeclaration(n, env)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(n, env)
	case *ast.IfStatement:
		return evalIfStatement(n, env)
	case *ast.WhileStatement:
		return evalWhileStatement(n, env)
	case *ast.ForStatement:
		return evalForStatement(n, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(n, env)
	case *ast.BreakStatement:
		return evalBreakStatement()
	case *ast.ContinueStatement:
		return evalContinueStatement()
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.NumberLiteral:
		return evalNumberLiteral(n)
	case *ast.StringLiteral:
		return evalStringLiteral(n)
	case *ast.RawStringLiteral:
		return evalRawStringLiteral(n)
	case *ast.NilLiteral:
		return evalNilLiteral()
	case *ast.ObjectLiteral:
		return evalObjectLiteral(n, env)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(n, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(n, env)
	case *ast.MemberDot:
		return evalMemberDot(n, env)
	case *ast.MemberBracket:
		return evalMemberBracket(n, env)
	case *ast.FunctionCall:
		return evalFunctionCall(n, env)
	case *ast.MethodCall:
		return evalMethodCall(n, env)
	case *ast.UnaryExpression:
		return evalUnaryExpression(n, env)
	case *ast.BinaryExpression:
		return evalBinaryExpression(n, env)
	case *ast.Assignment:
		return evalAssignment(n, env)
	default:
		return nil, runtime.NewError(fmt.Sprintf("unknown node type: %T", node), 0, 0)
	}
}

func Eval(program *ast.Program, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	env.Define("type", &builtinFunction{
		name: "type",
		fn:   builtinTypeFunc,
	})
	env.Define("str", &builtinFunction{
		name: "str",
		fn:   builtinStrFunc,
	})
	env.Define("number", &builtinFunction{
		name: "number",
		fn:   builtinNumberFunc,
	})
	env.Define("len", &builtinFunction{
		name: "len",
		fn:   builtinLenFunc,
	})
	env.Define("clone", &builtinFunction{
		name: "clone",
		fn:   builtinCloneFunc,
	})
	env.Define("exit", &builtinFunction{
		name: "exit",
		fn:   builtinExitFunc,
	})

	if _, exists := env.Get("$"); !exists {
		env.Define("$", &globalRefObject{env: env})
	}

	return evalProgram(program, env)
}

func evalProgram(program *ast.Program, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	var result runtime.Object = objects.NilInstance

	for _, stmt := range program.Statements {
		val, err := evalGeneric(stmt, env)
		if err != nil {
			return nil, err
		}
		result = val
	}

	return result, nil
}

func evalBlock(block *ast.BlockStatement, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	var result runtime.Object = objects.NilInstance
	blockEnv := runtime.NewEnvironment(env)

	for _, stmt := range block.Statements {
		val, err := evalGeneric(stmt, blockEnv)
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

func evalVarDeclaration(vd *ast.VarDeclaration, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	var value runtime.Object = objects.NilInstance

	if vd.Value != nil {
		val, err := evalGeneric(vd.Value, env)
		if err != nil {
			return nil, err
		}
		value = val
	}

	env.Define(vd.Name, value)
	return value, nil
}

func evalExpressionStatement(es *ast.ExpressionStatement, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	return evalGeneric(es.Expression, env)
}

func evalIfStatement(is *ast.IfStatement, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	cond, err := evalGeneric(is.Condition, env)
	if err != nil {
		return nil, err
	}

	if cond.Truthy() {
		result, err := evalBlock(is.Then, env)
		if err != nil {
			return nil, err
		}

		if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
			return ret, nil
		}

		return result, nil
	}

	for _, elif := range is.Elif {
		cond, err := evalGeneric(elif.Condition, env)
		if err != nil {
			return nil, err
		}

		if cond.Truthy() {
			result, err := evalBlock(elif.Body, env)
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
		result, err := evalBlock(is.Else, env)
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

func evalWhileStatement(ws *ast.WhileStatement, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	var result runtime.Object = objects.NilInstance

	for {
		cond, err := evalGeneric(ws.Condition, env)
		if err != nil {
			return nil, err
		}

		if !cond.Truthy() {
			break
		}

		val, err := evalBlock(ws.Body, env)
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

func evalForStatement(fs *ast.ForStatement, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	obj, err := evalGeneric(fs.Object, env)
	if err != nil {
		return nil, err
	}

	iterable, ok := obj.(runtime.Iterable)
	if !ok {
		return nil, runtime.NewError("cannot iterate over non-iterable object", fs.Pos.Line, fs.Pos.Column)
	}

	var result runtime.Object = objects.NilInstance

	for key, value := range iterable.Iterator() {
		forEnv := runtime.NewEnvironment(env)
		forEnv.Define(fs.Key, key)
		if fs.Value != "" {
			forEnv.Define(fs.Value, value)
		}

		val, err := evalBlock(fs.Body, forEnv)
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

func evalReturnStatement(rs *ast.ReturnStatement, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	var value runtime.Object = objects.NilInstance

	if rs.Value != nil {
		val, err := evalGeneric(rs.Value, env)
		if err != nil {
			return nil, err
		}
		value = val
	}

	return &objects.ReturnSignal{Value: value}, nil
}

func evalBreakStatement() (runtime.Object, *runtime.Error) {
	return &objects.BreakSignal{}, nil
}

func evalContinueStatement() (runtime.Object, *runtime.Error) {
	return &objects.ContinueSignal{}, nil
}

func evalIdentifier(id *ast.Identifier, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	val, ok := env.Get(id.Name)
	if !ok {
		return objects.NilInstance, nil
	}
	return val, nil
}

func evalNumberLiteral(nl *ast.NumberLiteral) (runtime.Object, *runtime.Error) {
	return objects.NewNumber(nl.Value), nil
}

func evalStringLiteral(sl *ast.StringLiteral) (runtime.Object, *runtime.Error) {
	return objects.NewString(sl.Value), nil
}

func evalRawStringLiteral(rsl *ast.RawStringLiteral) (runtime.Object, *runtime.Error) {
	return objects.NewString(rsl.Value), nil
}

func evalNilLiteral() (runtime.Object, *runtime.Error) {
	return objects.NilInstance, nil
}

func evalObjectLiteral(ol *ast.ObjectLiteral, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	obj := objects.NewObject()

	for key, expr := range ol.Properties {
		val, err := evalGeneric(expr, env)
		if err != nil {
			return nil, err
		}

		_ = obj.Set(objects.NewString(key), val)
	}

	return obj, nil
}

func evalArrayLiteral(al *ast.ArrayLiteral, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	obj := objects.NewObject()

	for i, expr := range al.Elements {
		val, err := evalGeneric(expr, env)
		if err != nil {
			return nil, err
		}

		_ = obj.Set(objects.NewNumber(float64(i)), val)
	}

	return obj, nil
}

func evalFunctionLiteral(fl *ast.FunctionLiteral, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	return objects.NewFunction(fl.Parameters, fl.Body, runtime.NewEnvironment(env), evalBlock), nil
}

func evalMemberDot(md *ast.MemberDot, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	obj, err := evalGeneric(md.Object, env)
	if err != nil {
		return nil, err
	}

	indexable, ok := obj.(runtime.Indexable)
	if !ok {
		return nil, runtime.NewError("cannot access member of non-indexable object", md.Pos.Line, md.Pos.Column)
	}

	key := objects.NewString(md.Member)
	return indexable.Get(key)
}

func evalMemberBracket(mb *ast.MemberBracket, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	obj, err := evalGeneric(mb.Object, env)
	if err != nil {
		return nil, err
	}

	key, err := evalGeneric(mb.Member, env)
	if err != nil {
		return nil, err
	}

	indexable, ok := obj.(runtime.Indexable)
	if !ok {
		return nil, runtime.NewError("cannot access member of non-indexable object", mb.Pos.Line, mb.Pos.Column)
	}

	return indexable.Get(key)
}

func evalFunctionCall(fc *ast.FunctionCall, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	// Evaluate function expression
	fnExpr, err := evalGeneric(fc.Function, env)
	if err != nil {
		return nil, err
	}

	// Check if it's callable
	callable, ok := fnExpr.(runtime.Callable)
	if !ok {
		return nil, runtime.NewError("attempt to call non-function", fc.Pos.Line, fc.Pos.Column)
	}

	// Evaluate arguments
	var args []runtime.Object
	for _, argExpr := range fc.Arguments {
		arg, err := evalGeneric(argExpr, env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	// Call the function
	result, err := callable.Call(args)
	if err != nil {
		return nil, err
	}

	// Unwrap return signals
	if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
		return ret.Value, nil
	}

	return result, nil
}

func evalMethodCall(mc *ast.MethodCall, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	// Evaluate object
	obj, err := evalGeneric(mc.Object, env)
	if err != nil {
		return nil, err
	}

	// Get method from object
	indexable, ok := obj.(runtime.Indexable)
	if !ok {
		return nil, runtime.NewError("cannot access method of non-indexable object", mc.Pos.Line, mc.Pos.Column)
	}

	method, err := indexable.Get(objects.NewString(mc.Method))
	if err != nil {
		return nil, err
	}

	// Check if method is callable
	callable, ok := method.(runtime.Callable)
	if !ok {
		return nil, runtime.NewError("attempt to call non-function method", mc.Pos.Line, mc.Pos.Column)
	}

	// Evaluate arguments
	var args []runtime.Object
	for _, argExpr := range mc.Arguments {
		arg, err := evalGeneric(argExpr, env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	// Call method with object as first argument
	allArgs := append([]runtime.Object{obj}, args...)
	result, err := callable.Call(allArgs)
	if err != nil {
		return nil, err
	}

	// Unwrap return signals
	if ret, isReturn := result.(*objects.ReturnSignal); isReturn {
		return ret.Value, nil
	}

	return result, nil
}

func evalUnaryExpression(ue *ast.UnaryExpression, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	right, err := evalGeneric(ue.Right, env)
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
		return nil, runtime.NewError("cannot negate non-number", ue.Pos.Line, ue.Pos.Column)
	default:
		return nil, runtime.NewError(fmt.Sprintf("unknown operator: %s", ue.Operator), ue.Pos.Line, ue.Pos.Column)
	}
}

func evalBinaryExpression(be *ast.BinaryExpression, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	left, err := evalGeneric(be.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := evalGeneric(be.Right, env)
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
		if !left.Truthy() {
			return left, nil
		}
		return right, nil
	case "||":
		if left.Truthy() {
			return left, nil
		}
		return right, nil
	default:
		return nil, runtime.NewError(fmt.Sprintf("unknown operator: %s", be.Operator), be.Pos.Line, be.Pos.Column)
	}

	return nil, runtime.NewError(fmt.Sprintf("cannot apply operator %s to types %s and %s",
		be.Operator, left.Type(), right.Type()), be.Pos.Line, be.Pos.Column)
}

func evalAssignment(a *ast.Assignment, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	right, err := evalGeneric(a.Right, env)
	if err != nil {
		return nil, err
	}

	switch left := a.Left.(type) {
	case *ast.Identifier:
		env.Set(left.Name, right)
		return right, nil
	case *ast.MemberDot:
		obj, err := evalGeneric(left.Object, env)
		if err != nil {
			return nil, err
		}

		indexable, ok := obj.(runtime.Indexable)
		if !ok {
			return nil, runtime.NewError("cannot assign to member of non-indexable object", left.Pos.Line, left.Pos.Column)
		}

		key := objects.NewString(left.Member)
		if err := indexable.Set(key, right); err != nil {
			return nil, err
		}
		return right, nil
	case *ast.MemberBracket:
		obj, err := evalGeneric(left.Object, env)
		if err != nil {
			return nil, err
		}

		key, err := evalGeneric(left.Member, env)
		if err != nil {
			return nil, err
		}

		indexable, ok := obj.(runtime.Indexable)
		if !ok {
			return nil, runtime.NewError("cannot assign to member of non-indexable object", left.Pos.Line, left.Pos.Column)
		}

		if err := indexable.Set(key, right); err != nil {
			return nil, err
		}
		return right, nil
	default:
		return nil, runtime.NewError("invalid assignment target", a.Pos.Line, a.Pos.Column)
	}
}
