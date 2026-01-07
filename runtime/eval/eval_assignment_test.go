package eval

import (
	"context"
	"refl/runtime"
	"testing"
	"time"

	"refl/runtime/objects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEvalAssignment verifies that assignment operations work correctly
func TestEvalAssignment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Simple assignments
		{"assign to global", "x = 5", 5},
		{"assign multiple times", "x = 5\nx = 10\nx", 10},
		{"assign to multiple vars", "x = 1\ny = 2\nx + y", 3},

		// Assignment expressions
		{"assign expression result", "x = 2 + 3", 5},
		{"assign function result", "x = len(\"hello\")", 5},
		{"assign with operation", "x = 5\nx = x + 3\nx", 8},

		// Object property assignments
		{"assign object property", "obj = {}\nobj.x = 5\nobj.x", 5},
		{"modify object property", "obj = {x: 1}\nobj.x = 2\nobj.x", 2},
		{"nested assignment", "obj = {}\nobj.inner = {}\nobj.inner.value = 42\nobj.inner.value", 42},

		// Array index assignments
		{"assign array element", "arr = {}\narr[0] = 10\narr[0]", 10},
		{"assign multiple indices", "arr = {}\narr[0] = 1\narr[1] = 2\narr[0] + arr[1]", 3},
		{"modify array element", "arr = {1, 2, 3}\narr[1] = 99\narr[1]", 99},

		// Assignment with different types
		{"change type", "x = 5\nx = \"hello\"\nlen(x)", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			program := parseProgram(t, tt.input)
			env := runtime.NewEnvironment(nil)

			evaluator := New(ctx, program, env)
			result, err := evaluator.Run()
			require.NoError(t, err)

			assert.IsType(t, &objects.Number{}, result)
			num := result.(*objects.Number)
			assert.Equal(t, tt.expected, num.Value)
		})
	}
}
