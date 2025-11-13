package main

import (
	"fmt"
	"os"

	"gojs/repl"
	"gojs/runtime"
)

func main() {
	args := os.Args[1:]

	// If no arguments, start REPL
	if len(args) == 0 {
		defer func() {
			if r := recover(); r != nil {
				if r == "exit" {
					os.Exit(0)
				}
				panic(r)
			}
		}()
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	// Check for flags
	if args[0] == "-h" || args[0] == "--help" {
		printHelp()
		return
	}

	if args[0] == "-v" || args[0] == "--version" {
		fmt.Println("GoJS v1.0.0")
		return
	}

	// Otherwise, treat first argument as a file to execute
	filename := args[0]

	rt := runtime.New()
	if err := rt.RunFile(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("GoJS - A JavaScript runtime written in Go")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gojs [file.js]     Run a JavaScript file")
	fmt.Println("  gojs               Start REPL (interactive mode)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help         Show this help message")
	fmt.Println("  -v, --version      Show version")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - Event loop with macrotasks and microtasks")
	fmt.Println("  - Promise support (async/await)")
	fmt.Println("  - Timers (setTimeout, setInterval)")
	fmt.Println("  - Node.js-style modules (fs, path)")
	fmt.Println("  - Console API")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gojs test.js       # Run test.js")
	fmt.Println("  gojs               # Start REPL")
}
