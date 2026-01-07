package eval

import (
	"context"
	"testing"
	"time"

	"refl/runtime"
	"refl/runtime/objects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEvalVarDeclaration verifies that variable declarations work correctly
func TestEvalVarDeclaration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		// Declarations work
		{"declare local number", "var x = 5", float64(5)},
		{"declare global number", "x = 6", float64(6)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			program := parseProgram(t, tt.input)
			env := runtime.NewEnvironment(nil)
			result, err := Eval(ctx, program, env)
			require.NoError(t, err)

			switch expected := tt.expected.(type) {
			case float64:
				assert.IsType(t, &objects.Number{}, result)
				num := result.(*objects.Number)
				assert.Equal(t, expected, num.Value)
			}
		})
	}
}
