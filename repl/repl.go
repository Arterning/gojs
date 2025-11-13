package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gojs/runtime"
)

const PROMPT = "> "
const CONTINUE_PROMPT = "... "

// Start starts the REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	rt := runtime.New()

	fmt.Fprintf(out, "GoJS REPL v1.0.0\n")
	fmt.Fprintf(out, "Type '.help' for more information\n\n")

	var multilineBuffer strings.Builder
	isMultiline := false

	for {
		if isMultiline {
			fmt.Fprint(out, CONTINUE_PROMPT)
		} else {
			fmt.Fprint(out, PROMPT)
		}

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		// Handle special commands
		if !isMultiline && strings.HasPrefix(line, ".") {
			handleCommand(line, out)
			continue
		}

		// Check for multiline input
		if isMultiline {
			multilineBuffer.WriteString("\n")
			multilineBuffer.WriteString(line)

			// Check if we should end multiline mode
			if shouldEndMultiline(multilineBuffer.String()) {
				code := multilineBuffer.String()
				multilineBuffer.Reset()
				isMultiline = false

				evaluateCode(rt, code, out)
			}
		} else {
			// Check if this starts a multiline input
			if needsMoreInput(line) {
				isMultiline = true
				multilineBuffer.WriteString(line)
			} else {
				evaluateCode(rt, line, out)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(out, "Error reading input: %v\n", err)
	}
}

func handleCommand(cmd string, out io.Writer) {
	switch cmd {
	case ".help":
		fmt.Fprintf(out, "Commands:\n")
		fmt.Fprintf(out, "  .help    Show this help message\n")
		fmt.Fprintf(out, "  .exit    Exit the REPL\n")
		fmt.Fprintf(out, "  .clear   Clear the console\n")
	case ".exit":
		fmt.Fprintf(out, "Goodbye!\n")
		// Exit is handled by returning from Start()
		panic("exit")
	case ".clear":
		fmt.Fprint(out, "\033[H\033[2J")
	default:
		fmt.Fprintf(out, "Unknown command: %s\n", cmd)
		fmt.Fprintf(out, "Type .help for available commands\n")
	}
}

func evaluateCode(rt *runtime.Runtime, code string, out io.Writer) {
	code = strings.TrimSpace(code)
	if code == "" {
		return
	}

	val, err := rt.Eval(code)
	if err != nil {
		fmt.Fprintf(out, "Error: %v\n", err)
		return
	}

	// Print the result (unless it's undefined)
	if val != nil && val.String() != "undefined" {
		fmt.Fprintf(out, "%v\n", val.String())
	}
}

// needsMoreInput checks if the line needs more input (basic heuristic)
func needsMoreInput(line string) bool {
	line = strings.TrimSpace(line)

	// Count braces
	openBraces := strings.Count(line, "{")
	closeBraces := strings.Count(line, "}")
	openParens := strings.Count(line, "(")
	closeParens := strings.Count(line, ")")
	openBrackets := strings.Count(line, "[")
	closeBrackets := strings.Count(line, "]")

	// Check for unclosed strings (basic check)
	singleQuotes := strings.Count(line, "'")
	doubleQuotes := strings.Count(line, "\"")
	backQuotes := strings.Count(line, "`")

	// Need more input if:
	// - Unclosed braces, parentheses, or brackets
	// - Unclosed strings (odd number of quotes)
	// - Line ends with certain keywords or operators
	return openBraces != closeBraces ||
		openParens != closeParens ||
		openBrackets != closeBrackets ||
		singleQuotes%2 != 0 ||
		doubleQuotes%2 != 0 ||
		backQuotes%2 != 0 ||
		strings.HasSuffix(line, "{") ||
		strings.HasSuffix(line, "(") ||
		strings.HasSuffix(line, "[") ||
		strings.HasSuffix(line, ",") ||
		strings.HasSuffix(line, "\\")
}

// shouldEndMultiline checks if multiline input should end
func shouldEndMultiline(code string) bool {
	// Count all braces
	openBraces := strings.Count(code, "{")
	closeBraces := strings.Count(code, "}")
	openParens := strings.Count(code, "(")
	closeParens := strings.Count(code, ")")
	openBrackets := strings.Count(code, "[")
	closeBrackets := strings.Count(code, "]")

	// End multiline if all braces are balanced
	return openBraces == closeBraces &&
		openParens == closeParens &&
		openBrackets == closeBrackets
}
