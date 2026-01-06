package eval

import (
	"context"
	"fmt"
	"refl/runtime"
	"refl/runtime/objects"
)

func builtinIOPrintFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}

	return objects.NilInstance, nil
}

func builtinIOPrintlnFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}

	fmt.Println()

	return objects.NilInstance, nil
}

func builtinIOPrintfFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
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

	fmt.Print(result)

	return objects.NilInstance, nil
}

func createIoObject(ctx context.Context) runtime.Object {
	ioObj := objects.NewObject()

	defLiteralBuiltinFunc(ctx, "print", ioObj, builtinIOPrintFunc)
	defLiteralBuiltinFunc(ctx, "println", ioObj, builtinIOPrintlnFunc)
	defLiteralBuiltinFunc(ctx, "printf", ioObj, builtinIOPrintfFunc)

	return ioObj
}
