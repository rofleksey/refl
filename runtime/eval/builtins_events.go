package eval

import (
	"context"
	"refl/runtime"
	"refl/runtime/eventloop"
	"refl/runtime/objects"
	"time"
)

func builtinScheduleFunc(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
	if len(args) < 2 {
		return nil, runtime.NewPanic("schedule() expects at least 2 argument", 0, 0)
	}

	fn, ok := args[0].(*objects.Function)
	if !ok {
		return nil, runtime.NewPanic("schedule() first argument must be a function", 0, 0)
	}

	millis, ok := args[1].(*objects.Number)
	if !ok {
		return nil, runtime.NewPanic("schedule() second argument must be a number", 0, 0)
	}

	otherArgs := args[2:]

	t := time.UnixMilli(int64(millis.Value))

	eventLoop := ctx.Value("event_loop").(*eventloop.EventLoop)
	cancelFunc := eventLoop.Schedule(func() {
		_, err := fn.Call(ctx, otherArgs)
		if err != nil {
			panic("schedule call failed: " + err.Error())
		}
	}, t)

	return objects.NewWrapperFunction(func(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
		cancelFunc()

		return objects.NilInstance, nil
	}), nil
}

func createEventsObject() runtime.Object {
	obj := objects.NewObject()

	defLiteralBuiltinFunc("schedule", obj, builtinScheduleFunc)

	return obj
}
