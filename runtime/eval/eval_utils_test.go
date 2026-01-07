package eval

import (
	"refl/ast"
	"refl/parser"
	"testing"
)

func parseProgram(t *testing.T, input string) *ast.Program {
	p := parser.New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return program
}
