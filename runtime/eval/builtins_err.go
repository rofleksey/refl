package eval

import (
	"context"
	"refl/runtime"
	"refl/runtime/objects"
)

func builtinErrFmtFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 1 {
		return nil, runtime.NewPanic("newerr() expects at least 1 argument", 0, 0)
	}

	patternObj, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("newerr() first argument must be a string pattern", 0, 0)
	}

	pattern := patternObj.Value
	result := ""
	argIndex := 0
	patternIdx := 0

	for patternIdx < len(pattern) {
		if pattern[patternIdx] == '$' {
			if argIndex+1 < len(args) {
				result += args[argIndex+1].String()
				argIndex++
			} else {
				result += "$MISSING"
			}
			patternIdx++
		} else {
			result += string(pattern[patternIdx])
			patternIdx++
		}
	}

	for i := argIndex + 1; i < len(args); i++ {
		result += " $ERROR{" + args[i].String() + "}"
	}

	return objects.NewError(result), nil
}

func builtinNewErrFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("newerr() expects exactly 1 argument", 0, 0)
	}

	return objects.NewError(args[0].String()), nil
}

func builtinIsErrFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("iserr() expects exactly 1 argument", 0, 0)
	}

	_, ok := args[0].(*objects.UserError)

	return objects.NewBoolean(ok), nil
}

func builtinPanicFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	msg := "< no message >"
	if len(args) > 0 {
		msg = args[0].String()
	}
	return nil, runtime.NewPanic(msg, 0, 0)
}
