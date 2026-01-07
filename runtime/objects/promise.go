package objects

import (
	"context"
	"fmt"
	"refl/runtime"
	"sync"
)

type PromiseState string

const (
	PromiseStatePending   PromiseState = "pending"
	PromiseStateFulfilled PromiseState = "fulfilled"
	PromiseStateRejected  PromiseState = "rejected"
)

type Promise struct {
	id string

	state PromiseState

	result runtime.Object
	err    error

	then    []*Function
	catch   []*Function
	finally []*Function

	mu sync.RWMutex
}

func NewPromise() *Promise {
	result := &Promise{
		state:   PromiseStatePending,
		then:    make([]*Function, 0),
		catch:   make([]*Function, 0),
		finally: make([]*Function, 0),
	}

	result.id = fmt.Sprintf("%p", result)

	return result
}

func (p *Promise) Type() runtime.ObjectType { return runtime.ObjectType_ }
func (p *Promise) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Promise<%s>", p.state)
}
func (p *Promise) Truthy() bool { return true }
func (p *Promise) Equal(other runtime.Object) bool {
	return p == other
}
func (p *Promise) Clone() runtime.Object { return p }

func (p *Promise) Get(key runtime.Object) (runtime.Object, error) {
	keyStr := key.String()

	switch keyStr {
	case "then":
		return NewWrapperFunction(func(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
			if len(args) != 1 {
				return nil, runtime.NewPanic("then() expects exactly 1 argument", 0, 0)
			}

			fn, ok := args[0].(*Function)
			if !ok {
				return nil, runtime.NewPanic("then() argument must be a function", 0, 0)
			}

			p.mu.Lock()
			defer p.mu.Unlock()

			if p.state == PromiseStateFulfilled {
				evaluator := ctx.Value("evaluator").(Evaluator)
				evaluator.EnqueueTask(func() {
					_, err := fn.Call(evaluator.Context(), []runtime.Object{p.result})
					if err != nil {
						panic("promise resolve() failed: " + err.Error())
					}
				})
				return p, nil
			}

			p.then = append(p.then, fn)

			return p, nil
		}), nil

	case "catch":
		return NewWrapperFunction(func(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
			if len(args) != 1 {
				return nil, runtime.NewPanic("catch() expects exactly 1 argument", 0, 0)
			}

			fn, ok := args[0].(*Function)
			if !ok {
				return nil, runtime.NewPanic("catch() argument must be a function", 0, 0)
			}

			p.mu.Lock()
			defer p.mu.Unlock()

			if p.state == PromiseStateRejected {
				evaluator := ctx.Value("evaluator").(Evaluator)
				evaluator.EnqueueTask(func() {
					_, err := fn.Call(ctx, []runtime.Object{NewError(p.err.Error())})
					if err != nil {
						panic("promise resolve() failed: " + err.Error())
					}
				})

				return p, nil
			}

			p.catch = append(p.catch, fn)

			return p, nil
		}), nil

	case "finally":
		return NewWrapperFunction(func(ctx context.Context, args []runtime.Object) (runtime.Object, error) {
			if len(args) != 1 {
				return nil, runtime.NewPanic("finally() expects exactly 1 argument", 0, 0)
			}

			fn, ok := args[0].(*Function)
			if !ok {
				return nil, runtime.NewPanic("finally() argument must be a function", 0, 0)
			}

			p.mu.Lock()
			defer p.mu.Unlock()

			if p.state != PromiseStatePending {
				evaluator := ctx.Value("evaluator").(Evaluator)
				evaluator.EnqueueTask(func() {
					_, err := fn.Call(evaluator.Context(), []runtime.Object{})
					if err != nil {
						panic("promise finally() failed: " + err.Error())
					}
				})

				return p, nil
			}

			p.finally = append(p.finally, fn)

			return p, nil
		}), nil
	}

	return nil, nil
}

func (p *Promise) Set(_, _ runtime.Object) error {
	return runtime.NewPanic("cannot modify Promise object", 0, 0)
}

func (p *Promise) Length() int { return 0 }

func (p *Promise) HashKey() runtime.HashKey {
	return runtime.HashKey("promise_" + p.id)
}

func (p *Promise) Not() runtime.Object {
	return NewBoolean(!p.Truthy())
}

func (p *Promise) Resolve(value runtime.Object, evaluator Evaluator) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PromiseStatePending {
		return
	}

	p.state = PromiseStateFulfilled
	p.result = value

	for _, handler := range p.then {
		evaluator.EnqueueTask(func() {
			_, err := handler.Call(evaluator.Context(), []runtime.Object{value})
			if err != nil {
				panic("promise resolve() failed: " + err.Error())
			}
		})
	}

	for _, handler := range p.finally {
		evaluator.EnqueueTask(func() {
			_, err := handler.Call(evaluator.Context(), []runtime.Object{})
			if err != nil {
				panic("promise finally() failed: " + err.Error())
			}
		})
	}
}

func (p *Promise) Reject(errValue error, evaluator Evaluator) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PromiseStatePending {
		return
	}

	p.state = PromiseStateRejected
	p.err = errValue

	for _, handler := range p.catch {
		evaluator.EnqueueTask(func() {
			_, err := handler.Call(evaluator.Context(), []runtime.Object{NewError(errValue.Error())})
			if err != nil {
				panic("promise reject() failed: " + err.Error())
			}
		})
	}

	for _, handler := range p.finally {
		evaluator.EnqueueTask(func() {
			_, err := handler.Call(evaluator.Context(), []runtime.Object{})
			if err != nil {
				panic("promise finally() failed: " + err.Error())
			}
		})
	}
}
