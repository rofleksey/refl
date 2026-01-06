package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"refl/ast"
	"strings"

	"refl/parser"
	"refl/runtime"
	"refl/runtime/eval"
	"refl/runtime/objects"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
		os.Exit(1)
	}

	var source string
	var err error

	if len(os.Args) == 2 {
		filename := os.Args[1]
		if !strings.HasSuffix(filename, ".refl") {
			filename += ".refl"
		}
		source, err = readFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		source, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	program, parseErr := parseSource(source)
	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", parseErr)
		os.Exit(1)
	}

	env := createGlobalEnvironment()
	result, runtimeErr := executeProgram(program, env)

	if runtimeErr != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", runtimeErr)
		os.Exit(1)
	}

	if result != objects.NilInstance && result != nil {
		fmt.Println(result.String())
	}
}

func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		// Try to find file with .refl extension if not already added
		if !strings.HasSuffix(filename, ".refl") {
			filenameWithExt := filename + ".refl"
			content, err = os.ReadFile(filenameWithExt)
			if err != nil {
				// Try current directory
				cwd, _ := os.Getwd()
				path := filepath.Join(cwd, filename)
				content, err = os.ReadFile(path)
				if err != nil {
					pathWithExt := filepath.Join(cwd, filename+".refl")
					content, err = os.ReadFile(pathWithExt)
				}
			}
		}
	}

	if err != nil {
		return "", fmt.Errorf("cannot read file '%s': %v", filename, err)
	}
	return string(content), nil
}

func readStdin() (string, error) {
	// Interactive mode
	fmt.Print("refl> ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("no input provided")
}

func parseSource(source string) (*ast.Program, error) {
	p := parser.New()
	program, err := p.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("%v\nSource: %s", err, source)
	}
	return program, nil
}

func createGlobalEnvironment() *runtime.Environment {
	return runtime.NewEnvironment(nil)
}

func executeProgram(program *ast.Program, env *runtime.Environment) (runtime.Object, *runtime.Error) {
	result, err := eval.Eval(program, env)
	if err != nil {
		return nil, err
	}

	switch result.(type) {
	case *objects.ReturnSignal:
		return result.(*objects.ReturnSignal).Value, nil
	case *objects.BreakSignal, *objects.ContinueSignal:
		return nil, runtime.NewError("break/continue outside loop", 0, 0)
	}

	return result, nil
}
