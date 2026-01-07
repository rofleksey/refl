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

// TestEvalClosures verifies that closures work correctly
func TestEvalClosures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic closure
		{"simple closure", `
			var x = 10
			var getX = fun() { return x }
			getX()
		`, 10},

		{"closure captures current value", `
			var x = 1
			var getX = fun() { return x }
			x = 2
			getX()
		`, 2},

		// Closure factory
		{"makeAdder closure", `
			var makeAdder = fun(x) {
				return fun(y) { return x + y }
			}
			var add5 = makeAdder(5)
			add5(3)
		`, 8},

		// Stateful closure
		{"counter closure", `
			var makeCounter = fun() {
				var count = 0
				return fun() {
					count = count + 1
					return count
				}
			}
			var counter = makeCounter()
			counter() + counter() + counter()
		`, 6}, // 1 + 2 + 3 = 6
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
