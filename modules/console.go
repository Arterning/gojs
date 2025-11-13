package modules

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

// SetupConsole sets up the console object
func SetupConsole(vm *goja.Runtime) {
	console := vm.NewObject()

	// Helper function to format values
	format := func(args []goja.Value) string {
		parts := make([]string, len(args))
		for i, arg := range args {
			if arg == nil || goja.IsUndefined(arg) {
				parts[i] = "undefined"
			} else if goja.IsNull(arg) {
				parts[i] = "null"
			} else {
				parts[i] = arg.String()
			}
		}
		return strings.Join(parts, " ")
	}

	// console.log
	console.Set("log", func(call goja.FunctionCall) goja.Value {
		fmt.Println(format(call.Arguments))
		return goja.Undefined()
	})

	// console.info (alias for log)
	console.Set("info", func(call goja.FunctionCall) goja.Value {
		fmt.Println(format(call.Arguments))
		return goja.Undefined()
	})

	// console.warn
	console.Set("warn", func(call goja.FunctionCall) goja.Value {
		fmt.Println("[WARN]", format(call.Arguments))
		return goja.Undefined()
	})

	// console.error
	console.Set("error", func(call goja.FunctionCall) goja.Value {
		fmt.Println("[ERROR]", format(call.Arguments))
		return goja.Undefined()
	})

	// console.debug
	console.Set("debug", func(call goja.FunctionCall) goja.Value {
		fmt.Println("[DEBUG]", format(call.Arguments))
		return goja.Undefined()
	})

	// console.dir
	console.Set("dir", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			fmt.Println(call.Arguments[0].String())
		}
		return goja.Undefined()
	})

	// console.trace
	console.Set("trace", func(call goja.FunctionCall) goja.Value {
		fmt.Println("[TRACE]", format(call.Arguments))
		// TODO: Add stack trace
		return goja.Undefined()
	})

	// console.assert
	console.Set("assert", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}

		assertion := call.Arguments[0].ToBoolean()
		if !assertion {
			message := "Assertion failed"
			if len(call.Arguments) > 1 {
				message += ": " + format(call.Arguments[1:])
			}
			fmt.Println("[ASSERT]", message)
		}
		return goja.Undefined()
	})

	// console.clear
	console.Set("clear", func(call goja.FunctionCall) goja.Value {
		fmt.Print("\033[H\033[2J")
		return goja.Undefined()
	})

	// console.time / console.timeEnd (simple implementation)
	timers := make(map[string]int64)

	console.Set("time", func(call goja.FunctionCall) goja.Value {
		label := "default"
		if len(call.Arguments) > 0 {
			label = call.Arguments[0].String()
		}
		timers[label] = 0 // Simplified - could use actual timing
		return goja.Undefined()
	})

	console.Set("timeEnd", func(call goja.FunctionCall) goja.Value {
		label := "default"
		if len(call.Arguments) > 0 {
			label = call.Arguments[0].String()
		}
		if _, ok := timers[label]; ok {
			fmt.Printf("%s: 0ms\n", label) // Simplified
			delete(timers, label)
		}
		return goja.Undefined()
	})

	vm.Set("console", console)
}
