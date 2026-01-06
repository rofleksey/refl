package eval

import (
	"context"
	"math"
	"math/rand"
	"refl/runtime"
	"refl/runtime/objects"
)

func builtinMathAbsFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("math.abs() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("math.abs() argument must be a number", 0, 0)
	}

	return objects.NewNumber(math.Abs(num.Value)), nil
}

func builtinMathFloorFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("math.floor() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("math.floor() argument must be a number", 0, 0)
	}

	return objects.NewNumber(math.Floor(num.Value)), nil
}

func builtinMathCeilFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("math.ceil() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("math.ceil() argument must be a number", 0, 0)
	}

	return objects.NewNumber(math.Ceil(num.Value)), nil
}

func builtinMathRoundFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("math.round() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("math.round() argument must be a number", 0, 0)
	}

	return objects.NewNumber(math.Round(num.Value)), nil
}

func builtinMathSqrtFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("math.sqrt() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("math.sqrt() argument must be a number", 0, 0)
	}

	if num.Value < 0 {
		return nil, runtime.NewPanic("math.sqrt() argument must be non-negative", 0, 0)
	}

	return objects.NewNumber(math.Sqrt(num.Value)), nil
}

func builtinMathPowFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("math.pow() expects exactly 2 arguments", 0, 0)
	}

	base, ok1 := args[0].(*objects.Number)
	exp, ok2 := args[1].(*objects.Number)

	if !ok1 || !ok2 {
		return nil, runtime.NewPanic("math.pow() arguments must be numbers", 0, 0)
	}

	return objects.NewNumber(math.Pow(base.Value, exp.Value)), nil
}

func builtinMathMaxFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 1 {
		return nil, runtime.NewPanic("math.max() expects at least 1 argument", 0, 0)
	}

	max := math.Inf(-1)
	for _, arg := range args {
		num, ok := arg.(*objects.Number)
		if !ok {
			return nil, runtime.NewPanic("math.max() arguments must be numbers", 0, 0)
		}
		max = math.Max(max, num.Value)
	}

	return objects.NewNumber(max), nil
}

func builtinMathMinFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 1 {
		return nil, runtime.NewPanic("math.min() expects at least 1 argument", 0, 0)
	}

	min := math.Inf(1)
	for _, arg := range args {
		num, ok := arg.(*objects.Number)
		if !ok {
			return nil, runtime.NewPanic("math.min() arguments must be numbers", 0, 0)
		}
		min = math.Min(min, num.Value)
	}

	return objects.NewNumber(min), nil
}

func builtinMathRandomFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) > 0 {
		return nil, runtime.NewPanic("math.random() expects no arguments", 0, 0)
	}

	return objects.NewNumber(rand.Float64()), nil
}

func createMathObject(ctx context.Context) runtime.Object {
	mathObj := objects.NewObject()

	defLiteralBuiltinFunc(ctx, "abs", mathObj, builtinMathAbsFunc)
	defLiteralBuiltinFunc(ctx, "floor", mathObj, builtinMathFloorFunc)
	defLiteralBuiltinFunc(ctx, "ceil", mathObj, builtinMathCeilFunc)
	defLiteralBuiltinFunc(ctx, "round", mathObj, builtinMathRoundFunc)
	defLiteralBuiltinFunc(ctx, "sqrt", mathObj, builtinMathSqrtFunc)
	defLiteralBuiltinFunc(ctx, "pow", mathObj, builtinMathPowFunc)
	defLiteralBuiltinFunc(ctx, "max", mathObj, builtinMathMaxFunc)
	defLiteralBuiltinFunc(ctx, "min", mathObj, builtinMathMinFunc)
	defLiteralBuiltinFunc(ctx, "random", mathObj, builtinMathRandomFunc)

	mathObj.SetLiteral("PI", objects.NewNumber(math.Pi))
	mathObj.SetLiteral("E", objects.NewNumber(math.E))
	mathObj.SetLiteral("INF", objects.NewNumber(math.Inf(1)))
	mathObj.SetLiteral("NEG_INF", objects.NewNumber(math.Inf(-1)))
	mathObj.SetLiteral("NAN", objects.NewNumber(math.NaN()))

	return mathObj
}
