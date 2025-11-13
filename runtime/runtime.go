package runtime

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/dop251/goja"
	"gojs/modules"
)

// Runtime represents the JavaScript runtime environment
type Runtime struct {
	VM        *goja.Runtime
	EventLoop *EventLoop
}

// New creates a new JavaScript runtime
func New() *Runtime {
	vm := goja.New()
	loop := NewEventLoop(vm)

	rt := &Runtime{
		VM:        vm,
		EventLoop: loop,
	}

	// Setup global functions
	rt.setupGlobals()

	// Setup Promise
	if err := SetupPromise(vm, loop); err != nil {
		panic(err)
	}

	// Setup console
	modules.SetupConsole(vm)

	// Setup require function first (defaults to current directory)
	if err := modules.SetupRequire(vm, "."); err != nil {
		panic(err)
	}

	// Setup built-in modules (fs, path) - these will be registered in the require cache
	if err := modules.SetupFS(vm); err != nil {
		panic(err)
	}
	if err := modules.SetupPath(vm); err != nil {
		panic(err)
	}

	return rt
}

// setupGlobals sets up global functions like setTimeout, setInterval, etc.
func (rt *Runtime) setupGlobals() {
	vm := rt.VM
	loop := rt.EventLoop

	// queueMicrotask
	vm.Set("queueMicrotask", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("queueMicrotask requires a function argument"))
		}

		fn, ok := goja.AssertFunction(call.Arguments[0])
		if !ok {
			panic(vm.ToValue("queueMicrotask argument must be a function"))
		}

		loop.QueueMicrotask(func() {
			fn(goja.Undefined(), goja.Undefined())
		})

		return goja.Undefined()
	})

	// setTimeout
	vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("setTimeout requires at least 1 argument"))
		}

		fn, ok := goja.AssertFunction(call.Arguments[0])
		if !ok {
			panic(vm.ToValue("setTimeout first argument must be a function"))
		}

		delay := int64(0)
		if len(call.Arguments) > 1 {
			delay = call.Arguments[1].ToInteger()
		}

		// Capture additional arguments
		args := make([]goja.Value, 0)
		if len(call.Arguments) > 2 {
			args = call.Arguments[2:]
		}

		id := loop.SetTimeout(func() {
			fn(goja.Undefined(), args...)
		}, time.Duration(delay)*time.Millisecond)

		return vm.ToValue(id)
	})

	// clearTimeout
	vm.Set("clearTimeout", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}

		id := int(call.Arguments[0].ToInteger())
		loop.ClearTimeout(id)

		return goja.Undefined()
	})

	// setInterval
	vm.Set("setInterval", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("setInterval requires at least 1 argument"))
		}

		fn, ok := goja.AssertFunction(call.Arguments[0])
		if !ok {
			panic(vm.ToValue("setInterval first argument must be a function"))
		}

		delay := int64(0)
		if len(call.Arguments) > 1 {
			delay = call.Arguments[1].ToInteger()
		}

		// Capture additional arguments
		args := make([]goja.Value, 0)
		if len(call.Arguments) > 2 {
			args = call.Arguments[2:]
		}

		id := loop.SetInterval(func() {
			fn(goja.Undefined(), args...)
		}, time.Duration(delay)*time.Millisecond)

		return vm.ToValue(id)
	})

	// clearInterval
	vm.Set("clearInterval", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}

		id := int(call.Arguments[0].ToInteger())
		loop.ClearInterval(id)

		return goja.Undefined()
	})

	// setImmediate (executes in next macrotask)
	vm.Set("setImmediate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("setImmediate requires a function argument"))
		}

		fn, ok := goja.AssertFunction(call.Arguments[0])
		if !ok {
			panic(vm.ToValue("setImmediate argument must be a function"))
		}

		args := make([]goja.Value, 0)
		if len(call.Arguments) > 1 {
			args = call.Arguments[1:]
		}

		id := loop.SetTimeout(func() {
			fn(goja.Undefined(), args...)
		}, 0)

		return vm.ToValue(id)
	})

	// clearImmediate
	vm.Set("clearImmediate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}

		id := int(call.Arguments[0].ToInteger())
		loop.ClearTimeout(id)

		return goja.Undefined()
	})

	// global object
	vm.Set("global", vm.GlobalObject())
}

// RunScript runs a JavaScript script
func (rt *Runtime) RunScript(script string, filename string) (goja.Value, error) {
	// Compile and run the script
	prg, err := goja.Compile(filename, script, false)
	if err != nil {
		return nil, err
	}

	val, err := rt.VM.RunProgram(prg)
	if err != nil {
		return nil, err
	}

	// Run the event loop
	rt.EventLoop.Run()

	return val, nil
}

// RunFile runs a JavaScript file
func (rt *Runtime) RunFile(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Update require's base directory to the file's directory
	fileDir := filepath.Dir(absPath)
	if err := modules.SetupRequire(rt.VM, fileDir); err != nil {
		return err
	}

	_, err = rt.RunScript(string(content), absPath)
	return err
}

// Eval evaluates JavaScript code (for REPL)
func (rt *Runtime) Eval(code string) (goja.Value, error) {
	val, err := rt.VM.RunString(code)
	if err != nil {
		return nil, err
	}

	// Process any microtasks that were queued
	rt.EventLoop.processMicrotasks()

	return val, nil
}
