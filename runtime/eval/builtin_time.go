package eval

import (
	"context"
	"fmt"
	"refl/runtime"
	"refl/runtime/objects"
	"time"
)

func builtinTimeParseFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("time.parse() expects exactly 1 argument", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("time.parse() argument must be a string", 0, 0)
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, str.Value); err == nil {
			return objects.NewNumber(float64(t.UnixNano()) / 1e6), nil
		}
	}

	return nil, runtime.NewPanic(fmt.Sprintf("time.parse() could not parse string: %s", str.Value), 0, 0)
}

func builtinTimeFormatFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("time.format() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("time.format() argument must be a number", 0, 0)
	}

	t := time.UnixMilli(int64(num.Value))
	formatted := t.UTC().Format(time.RFC3339)

	return objects.NewString(formatted), nil
}

func builtinTimeNowFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 0 {
		return nil, runtime.NewPanic("time.now() expects no arguments", 0, 0)
	}

	now := time.Now()
	millis := now.UnixMilli()

	return objects.NewNumber(float64(millis)), nil
}

func builtinTimeSleepFunc(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("time.sleep() expects exactly 1 argument", 0, 0)
	}

	num, ok := args[0].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("time.sleep() argument must be a number", 0, 0)
	}

	if num.Value < 0 {
		return nil, runtime.NewPanic("time.sleep() argument must be non-negative", 0, 0)
	}

	duration := time.Duration(num.Value * float64(time.Millisecond))

	select {
	case <-time.After(duration):
		return objects.NilInstance, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func createTimeObject() runtime.Object {
	obj := objects.NewObject()

	defLiteralBuiltinFunc("parse", obj, builtinTimeParseFunc)
	defLiteralBuiltinFunc("format", obj, builtinTimeFormatFunc)
	defLiteralBuiltinFunc("now", obj, builtinTimeNowFunc)
	defLiteralBuiltinFunc("sleep", obj, builtinTimeSleepFunc)

	return obj
}
