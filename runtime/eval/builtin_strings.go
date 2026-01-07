package eval

import (
	"context"
	"refl/runtime"
	"refl/runtime/objects"
	"strings"
)

func builtinStringUpperFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("string.upper() expects exactly 1 argument", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.upper() argument must be a string", 0, 0)
	}

	return objects.NewString(strings.ToUpper(str.Value)), nil
}

func builtinStringLowerFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("string.lower() expects exactly 1 argument", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.lower() argument must be a string", 0, 0)
	}

	return objects.NewString(strings.ToLower(str.Value)), nil
}

func builtinStringTrimSpaceFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 1 {
		return nil, runtime.NewPanic("string.trimspace() expects exactly 1 argument", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.trimspace() argument must be a string", 0, 0)
	}

	return objects.NewString(strings.TrimSpace(str.Value)), nil
}

func builtinStringSplitFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.split() expects exactly 2 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.split() first argument must be a string", 0, 0)
	}

	sep, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.split() second argument must be a string", 0, 0)
	}

	parts := strings.Split(str.Value, sep.Value)
	result := objects.NewObject()

	for i, part := range parts {
		_ = result.Set(objects.NewNumber(float64(i)), objects.NewString(part))
	}

	return result, nil
}

func builtinStringJoinFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.join() expects exactly 2 arguments", 0, 0)
	}

	sep, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.join() first argument must be a string", 0, 0)
	}

	arr, ok := args[1].(*objects.ReflObject)
	if !ok {
		return nil, runtime.NewPanic("string.join() second argument must be an object", 0, 0)
	}

	var stringParts []string
	for key, value := range arr.Iterator() {
		if key.Type() == runtime.NumberType {
			stringParts = append(stringParts, value.String())
		}
	}

	return objects.NewString(strings.Join(stringParts, sep.Value)), nil
}

func builtinStringContainsFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.contains() expects exactly 2 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.contains() first argument must be a string", 0, 0)
	}

	substr, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.contains() second argument must be a string", 0, 0)
	}

	return objects.NewBoolean(strings.Contains(str.Value, substr.Value)), nil
}

func builtinStringHasPrefixFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.hasprefix() expects exactly 2 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.hasprefix() first argument must be a string", 0, 0)
	}

	prefix, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.hasprefix() second argument must be a string", 0, 0)
	}

	return objects.NewBoolean(strings.HasPrefix(str.Value, prefix.Value)), nil
}

func builtinStringHasSuffixFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.hassuffix() expects exactly 2 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.hassuffix() first argument must be a string", 0, 0)
	}

	suffix, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.hassuffix() second argument must be a string", 0, 0)
	}

	return objects.NewBoolean(strings.HasSuffix(str.Value, suffix.Value)), nil
}

func builtinStringReplaceFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 3 && len(args) != 4 {
		return nil, runtime.NewPanic("string.replace() expects 3 or 4 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.replace() first argument must be a string", 0, 0)
	}

	oldStr, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.replace() second argument must be a string", 0, 0)
	}

	newStr, ok := args[2].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.replace() third argument must be a string", 0, 0)
	}

	n := -1
	if len(args) == 4 {
		num, ok := args[3].(*objects.Number)
		if !ok {
			return nil, runtime.NewPanic("string.replace() fourth argument must be a number", 0, 0)
		}
		n = int(num.Value)
	}

	result := strings.Replace(str.Value, oldStr.Value, newStr.Value, n)
	return objects.NewString(result), nil
}

func builtinStringIndexFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.index() expects exactly 2 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.index() first argument must be a string", 0, 0)
	}

	substr, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.index() second argument must be a string", 0, 0)
	}

	index := strings.Index(str.Value, substr.Value)
	return objects.NewNumber(float64(index)), nil
}

func builtinStringLastIndexFunc(_ context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) != 2 {
		return nil, runtime.NewPanic("string.lastindex() expects exactly 2 arguments", 0, 0)
	}

	str, ok := args[0].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.lastindex() first argument must be a string", 0, 0)
	}

	substr, ok := args[1].(*objects.String)
	if !ok {
		return nil, runtime.NewPanic("string.lastindex() second argument must be a string", 0, 0)
	}

	index := strings.LastIndex(str.Value, substr.Value)
	return objects.NewNumber(float64(index)), nil
}

func createStringObject() runtime.Object {
	stringObj := objects.NewObject()

	defLiteralBuiltinFunc("upper", stringObj, builtinStringUpperFunc)
	defLiteralBuiltinFunc("lower", stringObj, builtinStringLowerFunc)
	defLiteralBuiltinFunc("trim", stringObj, builtinStringTrimSpaceFunc)
	defLiteralBuiltinFunc("split", stringObj, builtinStringSplitFunc)
	defLiteralBuiltinFunc("join", stringObj, builtinStringJoinFunc)
	defLiteralBuiltinFunc("contains", stringObj, builtinStringContainsFunc)
	defLiteralBuiltinFunc("has_prefix", stringObj, builtinStringHasPrefixFunc)
	defLiteralBuiltinFunc("has_suffix", stringObj, builtinStringHasSuffixFunc)
	defLiteralBuiltinFunc("replace", stringObj, builtinStringReplaceFunc)
	defLiteralBuiltinFunc("index", stringObj, builtinStringIndexFunc)
	defLiteralBuiltinFunc("last_index", stringObj, builtinStringLastIndexFunc)

	return stringObj
}
