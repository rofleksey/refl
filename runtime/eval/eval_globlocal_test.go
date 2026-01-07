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

// TestEvalGlobalLocalVariables verifies that global and local variables work correctly together
func TestEvalGlobalLocalVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic global vs local
		{"local shadows global in function", `
			x = 5
			var test = fun () {
				var x = 10
				return x
			}
			test()
		`, 10},

		{"global unchanged by local", `
			x = 5
			var test = fun () {
				var x = 10
			}
			test()
			x
		`, 5},

		// Global assignment in local scope
		{"assign to global from local", `
			x = 1
			var test = fun () {
				x = 2
			}
			test()
			x
		`, 2},

		{"create global from local", `
			var test = fun () {
				newGlobal = 42
			}
			test()
			newGlobal
		`, 42},
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
