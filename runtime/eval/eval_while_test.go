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

// TestEvalWhileStatements verifies that while statements work correctly
func TestEvalWhileStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic while loops
		{"simple while", "var x = 0\nwhile x < 5 { x = x + 1 }\nx", 5},
		{"while false never runs", "var x = 0\nwhile 0 { x = 1 }\nx", 0},
		{"infinite loop prevention", "var x = 0\nwhile x < 3 { x = x + 1 }\nx", 3},

		// Accumulation
		{"sum with while", "var x = 0\nvar sum = 0\nwhile x < 5 { sum = sum + x\nx = x + 1 }\nsum", 10},
		{"factorial with while", "var n = 5\nvar result = 1\nwhile n > 0 { result = result * n\nn = n - 1 }\nresult", 120},

		// Break statement
		{"break early", "var x = 0\nwhile x < 10 { x = x + 1\nif x == 5 { break } }\nx", 5},
		{"break in nested if", "var x = 0\nwhile x < 10 { x = x + 1\nif x > 3 { if x == 5 { break } } }\nx", 5},
		{"break stops loop", "var x = 0\nvar count = 0\nwhile x < 5 { x = x + 1\nif x == 3 { break }\ncount = count + 1 }\ncount", 2},

		// Continue statement
		{"continue skips", "var x = 0\nvar sum = 0\nwhile x < 5 { x = x + 1\nif x == 3 { continue }\nsum = sum + x }\nsum", 12},                // 1+2+4+5=12
		{"continue in if", "var x = 0\nvar result = 0\nwhile x < 5 { x = x + 1\nif x % 2 == 0 { continue }\nresult = result + x }\nresult", 9}, // 1+3+5=9

		// Return from while
		{"return from while", `
			fun() {
				var x = 0
				while x < 5 {
					x = x + 1
					if x == 3 { return x }
				}
				return 0
			}()
		`, 3},

		// One-time loop
		{"one iteration", "var x = 0\nwhile x < 1 { x = x + 1 }\nx", 1},

		// Modifying loop variable in body
		{"modify counter", "var x = 0\nwhile x < 10 { x = x + 2 }\nx", 10},
		{"decrement", "var x = 10\nwhile x > 0 { x = x - 1 }\nx", 0},
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
