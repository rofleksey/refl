package eventloop

import (
	"container/heap"
	"context"
	"refl/runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Task func()

type EventCallback func(ctx context.Context, event string, args []runtime.Object)

type EventHandler struct {
	ID       uuid.UUID
	Callback EventCallback
}

var noop = func() {}

type immediateTask struct {
	cancelled *atomic.Bool
	task      Task
}

type delayedTask struct {
	cancelled *atomic.Bool
	task      Task
	executeAt time.Time
	index     int
}

type delayedTaskHeap struct {
	heap []*delayedTask
}

func (h *delayedTaskHeap) Len() int {
	return len(h.heap)
}

func (h *delayedTaskHeap) Less(i, j int) bool {
	return h.heap[i].executeAt.Before(h.heap[j].executeAt)
}

func (h *delayedTaskHeap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
	h.heap[i].index = i
	h.heap[j].index = j
}

func (h *delayedTaskHeap) Push(x any) {
	task := x.(*delayedTask)
	task.index = len(h.heap)
	h.heap = append(h.heap, task)
}

func (h *delayedTaskHeap) Pop() any {
	old := h.heap
	n := len(old)

	task := old[n-1]
	task.index = -1

	h.heap = old[0 : n-1]

	return task
}

type EventLoop struct {
	ctx    context.Context
	cancel context.CancelFunc

	immediateTaskChan chan immediateTask
	delayedTasks      *delayedTaskHeap
	delayedTasksMu    sync.Mutex

	triggerChan chan struct{}

	handlers   map[string][]EventHandler
	handlersMu sync.Mutex

	wg sync.WaitGroup

	loopRunning bool
	loopMu      sync.Mutex

	stopped bool
	stopMu  sync.Mutex

	lastPanic any
	panicMu   sync.Mutex
}

func New(ctx context.Context) *EventLoop {
	ctx, cancel := context.WithCancel(ctx)
	return &EventLoop{
		ctx:               ctx,
		cancel:            cancel,
		immediateTaskChan: make(chan immediateTask, 32),
		delayedTasks:      &delayedTaskHeap{},
		triggerChan:       make(chan struct{}, 1),
		handlers:          make(map[string][]EventHandler),
	}
}

func (e *EventLoop) RegisterCallback(event string, callback EventCallback) func() {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()

	handlerID := uuid.New()

	e.handlers[event] = append(e.handlers[event], EventHandler{
		ID:       handlerID,
		Callback: callback,
	})

	select {
	case e.triggerChan <- struct{}{}:
	default:
	}

	return func() {
		e.unregisterHandler(event, handlerID)
	}
}

func (e *EventLoop) RegisterLock() func() {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()

	handlerID := uuid.New()
	event := "__lock__" + handlerID.String()

	e.handlers[event] = append(e.handlers[event], EventHandler{
		ID:       handlerID,
		Callback: func(ctx context.Context, event string, args []runtime.Object) {},
	})

	select {
	case e.triggerChan <- struct{}{}:
	default:
	}

	return func() {
		e.unregisterHandler(event, handlerID)
	}
}

func (e *EventLoop) Fire(event string, args []runtime.Object) {
	e.Enqueue(func() {
		e.fireEvent(event, args)
	})
}

func (e *EventLoop) fireEvent(event string, args []runtime.Object) {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()

	handlers := e.handlers[event]
	for _, h := range handlers {
		h.Callback(e.ctx, event, args)
	}
}

func (e *EventLoop) unregisterHandler(event string, handlerID uuid.UUID) {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()

	select {
	case e.triggerChan <- struct{}{}:
	default:
	}

	if len(e.handlers[event]) == 1 {
		delete(e.handlers, event)
		return
	}

	newHandlers := make([]EventHandler, 0, len(e.handlers[event]))
	for _, h := range e.handlers[event] {
		if h.ID == handlerID {
			continue
		}

		newHandlers = append(newHandlers, h)
	}

	e.handlers[event] = newHandlers
}

func (e *EventLoop) Enqueue(task Task) func() {
	if task == nil {
		panic("task cannot be nil")
	}

	e.stopMu.Lock()
	if e.stopped {
		e.stopMu.Unlock()
		return noop
	}
	e.stopMu.Unlock()

	imTask := immediateTask{
		cancelled: new(atomic.Bool),
		task:      task,
	}
	cancelFunc := func() {
		imTask.cancelled.Store(true)
	}

	select {
	case e.immediateTaskChan <- imTask:
		return cancelFunc
	case <-e.ctx.Done():
		return cancelFunc
	}
}

func (e *EventLoop) Schedule(task Task, atTime time.Time) func() {
	if task == nil {
		panic("task cannot be nil")
	}

	e.stopMu.Lock()
	if e.stopped {
		e.stopMu.Unlock()
		return func() {}
	}
	e.stopMu.Unlock()

	delTask := &delayedTask{
		cancelled: new(atomic.Bool),
		task:      task,
		executeAt: atTime,
	}
	cancelFunc := func() {
		delTask.cancelled.Store(true)
	}

	select {
	case <-e.ctx.Done():
		return cancelFunc
	default:
		e.delayedTasksMu.Lock()
		heap.Push(e.delayedTasks, delTask)
		e.delayedTasksMu.Unlock()

		select {
		case e.triggerChan <- struct{}{}:
		default:
		}

		return cancelFunc
	}
}

func (e *EventLoop) Start() {
	e.loopMu.Lock()
	if e.loopRunning {
		e.loopMu.Unlock()
		return
	}
	e.loopRunning = true
	e.loopMu.Unlock()

	go func() {
		<-e.ctx.Done()
		e.stop()
	}()

	e.wg.Go(e.runLoop)
}

func (e *EventLoop) stop() {
	e.stopMu.Lock()
	if e.stopped {
		e.stopMu.Unlock()
		return
	}
	e.stopped = true
	e.stopMu.Unlock()

	e.cancel()
	e.wg.Wait()

	e.loopMu.Lock()
	e.loopRunning = false
	e.loopMu.Unlock()

	close(e.immediateTaskChan)
	close(e.triggerChan)
}

func (e *EventLoop) runLoop() {
	var timer *time.Timer
	var timerC <-chan time.Time

	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

outer:
	for {
		select {
		// try immediate tasks first
		case <-e.ctx.Done():
			e.drainTasks()
			return

		case imTask := <-e.immediateTaskChan:
			if imTask.cancelled.Load() {
				continue outer
			}

			if !e.executeTask(imTask.task) {
				return
			}

			continue outer

		default:
		}

		// try scheduled tasks if no immediate tasks found
		nextDelay := e.getNextDelay()

		if timer != nil {
			timer.Stop()
			timer = nil
			timerC = nil
		}

		if nextDelay < 0 {
			// no delayed tasks either - check if event handlers exist
			e.handlersMu.Lock()
			handleCount := len(e.handlers)
			e.handlersMu.Unlock()

			// no scheduled or delayed tasks, no event handlers - nothing will ever happen again
			if handleCount == 0 {
				return
			}
		} else {
			// at least 1 delayed task is pending
			timer = time.NewTimer(nextDelay)
			timerC = timer.C
		}

		select {
		case <-e.ctx.Done():
			e.drainTasks()
			return

		case imTask := <-e.immediateTaskChan:
			if imTask.cancelled.Load() {
				continue
			}

			if !e.executeTask(imTask.task) {
				return
			}

		case <-timerC:
			if !e.executeNextDelayedTask() {
				return
			}

		case <-e.triggerChan:
		}
	}
}

func (e *EventLoop) getNextDelay() time.Duration {
	e.delayedTasksMu.Lock()
	defer e.delayedTasksMu.Unlock()

	if len(e.delayedTasks.heap) == 0 {
		return -1
	}

	next := e.delayedTasks.heap[0]
	delay := time.Until(next.executeAt)
	if delay < 0 {
		delay = 0
	}
	return delay
}

func (e *EventLoop) executeNextDelayedTask() bool {
	e.delayedTasksMu.Lock()
	if e.delayedTasks.heap == nil || len(e.delayedTasks.heap) == 0 {
		e.delayedTasksMu.Unlock()
		return true
	}

	now := time.Now()
	if e.delayedTasks.heap[0].executeAt.After(now) {
		e.delayedTasksMu.Unlock()
		return true
	}

	task := heap.Pop(e.delayedTasks).(*delayedTask)
	e.delayedTasksMu.Unlock()

	if task.cancelled.Load() {
		return true
	}

	return e.executeTask(task.task)
}

func (e *EventLoop) executeTask(task Task) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			e.panicMu.Lock()
			if e.lastPanic == nil {
				e.lastPanic = r
			}
			e.panicMu.Unlock()
		}
	}()

	task()

	return true
}

func (e *EventLoop) drainTasks() {
	for {
		select {
		case imTask := <-e.immediateTaskChan:
			if imTask.cancelled.Load() {
				continue
			}

			if !e.executeTask(imTask.task) {
				return
			}
		default:
			return
		}
	}
}

func (e *EventLoop) IsRunning() bool {
	e.loopMu.Lock()
	defer e.loopMu.Unlock()
	return e.loopRunning
}

func (e *EventLoop) Wait() {
	e.wg.Wait()
}

func (e *EventLoop) LastPanic() any {
	e.panicMu.Lock()
	defer e.panicMu.Unlock()

	return e.lastPanic
}
