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

// TestEvalVariableShadowing verifies that variable shadowing works correctly
func TestEvalVariableShadowing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Block shadowing
		{"block shadows outer var", `
			var x = 1
			{
				var x = 2
			}
			x
		`, 1},

		// For loop shadowing
		{"for shadows iteration variable", `
			var x = 1
			for i, x in {10, 20, 30} {
				# x is the iteration value
			}
			x
		`, 1},

		// Function parameter shadowing
		{"parameter shadows global", `
			var x = 1
			test = fun (x) {
				return x
			}
			test(5)
		`, 5},

		{"shadow global with local function", `
			var x = 1
			test = fun () {
				var x = fun() { return 5 }
				return x()
			}
			test()
		`, 5},
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
