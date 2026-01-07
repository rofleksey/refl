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

// TestEvalGlobalVariables verifies that global variables work correctly
func TestEvalGlobalVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Global variable assignment
		{"assign global", "x = 5\nx", 5},
		{"reassign global", "x = 1\nx = 2\nx", 2},
		{"multiple globals", "x = 1\ny = 2\nx + y", 3},

		// Global from anywhere
		{"global from function", "fun() { x = 5 }()\nx", 5},
		{"global from if", "if 1 { x = 3 }\nx", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			program := parseProgram(t, tt.input)
			env := runtime.NewEnvironment(nil)

			result, err := Eval(ctx, program, env)
			require.NoError(t, err)

			assert.IsType(t, &objects.Number{}, result)
			num := result.(*objects.Number)
			assert.Equal(t, tt.expected, num.Value)
		})
	}
}
