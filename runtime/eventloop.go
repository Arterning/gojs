package runtime

import (
	"container/heap"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// Task represents a task in the event loop
type Task struct {
	Callback func()
	Time     time.Time
	Index    int
}

// TaskQueue is a priority queue for tasks
type TaskQueue []*Task

func (tq TaskQueue) Len() int           { return len(tq) }
func (tq TaskQueue) Less(i, j int) bool { return tq[i].Time.Before(tq[j].Time) }
func (tq TaskQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
	tq[i].Index = i
	tq[j].Index = j
}

func (tq *TaskQueue) Push(x interface{}) {
	n := len(*tq)
	task := x.(*Task)
	task.Index = n
	*tq = append(*tq, task)
}

func (tq *TaskQueue) Pop() interface{} {
	old := *tq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil
	task.Index = -1
	*tq = old[0 : n-1]
	return task
}

// EventLoop represents the JavaScript event loop
type EventLoop struct {
	vm            *goja.Runtime
	macrotasks    TaskQueue
	microtasks    []func()
	timers        map[int]*Task
	intervals     map[int]*Task
	timerID       int
	mutex         sync.Mutex
	running       bool
	stopChan      chan struct{}
	pendingTasks  int
}

// NewEventLoop creates a new event loop
func NewEventLoop(vm *goja.Runtime) *EventLoop {
	el := &EventLoop{
		vm:         vm,
		macrotasks: make(TaskQueue, 0),
		microtasks: make([]func(), 0),
		timers:     make(map[int]*Task),
		intervals:  make(map[int]*Task),
		timerID:    1,
		stopChan:   make(chan struct{}),
	}
	heap.Init(&el.macrotasks)
	return el
}

// QueueMicrotask adds a microtask to the queue
func (el *EventLoop) QueueMicrotask(fn func()) {
	el.mutex.Lock()
	defer el.mutex.Unlock()
	el.microtasks = append(el.microtasks, fn)
}

// SetTimeout schedules a function to run after a delay
func (el *EventLoop) SetTimeout(callback func(), delay time.Duration) int {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	id := el.timerID
	el.timerID++

	task := &Task{
		Callback: callback,
		Time:     time.Now().Add(delay),
	}

	heap.Push(&el.macrotasks, task)
	el.timers[id] = task
	el.pendingTasks++

	return id
}

// ClearTimeout cancels a timeout
func (el *EventLoop) ClearTimeout(id int) {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	if task, exists := el.timers[id]; exists {
		if task.Index >= 0 {
			heap.Remove(&el.macrotasks, task.Index)
			el.pendingTasks--
		}
		delete(el.timers, id)
	}
}

// SetInterval schedules a function to run repeatedly
func (el *EventLoop) SetInterval(callback func(), delay time.Duration) int {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	id := el.timerID
	el.timerID++

	var repeatFunc func()
	repeatFunc = func() {
		callback()

		// Reschedule the interval
		el.mutex.Lock()
		if _, exists := el.intervals[id]; exists {
			task := &Task{
				Callback: repeatFunc,
				Time:     time.Now().Add(delay),
			}
			heap.Push(&el.macrotasks, task)
			el.intervals[id] = task
		}
		el.mutex.Unlock()
	}

	task := &Task{
		Callback: repeatFunc,
		Time:     time.Now().Add(delay),
	}

	heap.Push(&el.macrotasks, task)
	el.intervals[id] = task
	el.pendingTasks++

	return id
}

// ClearInterval cancels an interval
func (el *EventLoop) ClearInterval(id int) {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	if task, exists := el.intervals[id]; exists {
		if task.Index >= 0 {
			heap.Remove(&el.macrotasks, task.Index)
			el.pendingTasks--
		}
		delete(el.intervals, id)
	}
}

// processMicrotasks executes all pending microtasks
func (el *EventLoop) processMicrotasks() {
	for {
		el.mutex.Lock()
		if len(el.microtasks) == 0 {
			el.mutex.Unlock()
			break
		}
		task := el.microtasks[0]
		el.microtasks = el.microtasks[1:]
		el.mutex.Unlock()

		task()
	}
}

// Run starts the event loop
func (el *EventLoop) Run() {
	el.running = true
	defer func() {
		el.running = false
	}()

	for {
		// Process all microtasks first
		el.processMicrotasks()

		// Get next macrotask
		el.mutex.Lock()
		if len(el.macrotasks) == 0 {
			el.mutex.Unlock()
			break
		}

		task := heap.Pop(&el.macrotasks).(*Task)
		el.mutex.Unlock()

		// Wait until it's time to execute
		now := time.Now()
		if task.Time.After(now) {
			time.Sleep(task.Time.Sub(now))
		}

		// Execute the task
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Handle panic in task
					if err, ok := r.(*goja.Exception); ok {
						println("Uncaught exception:", err.String())
					} else {
						println("Panic in task:", r)
					}
				}
			}()
			task.Callback()
		}()

		// Decrement pending tasks if it was a timeout (not interval)
		el.mutex.Lock()
		isInterval := false
		for _, intervalTask := range el.intervals {
			if intervalTask == task {
				isInterval = true
				break
			}
		}
		if !isInterval {
			el.pendingTasks--
		}
		el.mutex.Unlock()

		// Process microtasks after each macrotask
		el.processMicrotasks()
	}
}

// RunUntilIdle runs the event loop until there are no more tasks
func (el *EventLoop) RunUntilIdle() {
	el.Run()
}

// Stop stops the event loop
func (el *EventLoop) Stop() {
	close(el.stopChan)
}
