// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		// File mode
		filename := os.Args[1]
		source, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read file: %v\n", err)
			os.Exit(1)
		}

		p := NewParser(string(source))
		prog, err := p.ParseProgram()
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(1)
		}

		env := NewEnvironment()
		DefineBuiltins(env)

		_, err = Interpret(prog, env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// REPL mode
	fmt.Println("EsoLambda REPL  (type :q or :quit or Ctrl+C to exit)")
	fmt.Println("Enter expressions or function definitions. Press Enter twice to evaluate.")
	fmt.Println("------------------------------------------------------------")

	env := NewEnvironment()
	DefineBuiltins(env)

	reader := bufio.NewReader(os.Stdin)
	var buffer strings.Builder

	for {
		fmt.Print("esl> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("\nGoodbye.")
				return
			}
			fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)

		if line == "" {
			// Empty line → try to parse & evaluate accumulated buffer
			code := strings.TrimSpace(buffer.String())
			if code == "" {
				continue
			}

			p := NewParser(code)
			prog, err := p.ParseProgram()
			if err != nil {
				fmt.Printf("parse error: %v\n", err)
				buffer.Reset()
				continue
			}

			result, err := Interpret(prog, env)
			if err != nil {
				fmt.Printf("runtime error: %v\n", err)
			} else if result != nil && result.Type() != Void {
				fmt.Printf("→ %s\n", result.String())
			}

			buffer.Reset()
			continue
		}

		if line == ":q" || line == ":quit" || line == ":exit" {
			fmt.Println("Goodbye.")
			return
		}

		// Accumulate lines (multi-line support)
		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(line)
	}
}
