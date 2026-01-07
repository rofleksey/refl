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

// TestEvalArithmeticExpressions verifies that arithmetic expressions work correctly
func TestEvalArithmeticExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Addition
		{"simple addition", "1 + 2", 3},
		{"multiple addition", "1 + 2 + 3", 6},
		{"negative addition", "5 + (-3)", 2},

		// Subtraction
		{"simple subtraction", "5 - 3", 2},
		{"negative result", "3 - 5", -2},
		{"multiple subtraction", "10 - 2 - 3", 5},

		// Multiplication
		{"simple multiplication", "2 * 3", 6},
		{"multiple multiplication", "2 * 3 * 4", 24},
		{"negative multiplication", "2 * (-3)", -6},
		{"zero multiplication", "5 * 0", 0},

		// Division
		{"simple division", "6 / 2", 3},
		{"decimal division", "5 / 2", 2.5},
		{"division by one", "7 / 1", 7},

		// Modulo
		{"simple modulo", "7 % 3", 1},
		{"even modulo", "6 % 3", 0},
		{"modulo with decimals", "5.5 % 2", 1.5},

		// Negation
		{"negate positive", "-5", -5},
		{"negate negative", "-(-5)", 5},
		{"double negation", "-(-(-5))", -5},

		// Operator precedence
		{"multiplication before addition", "2 + 3 * 4", 14},
		{"parentheses override", "(2 + 3) * 4", 20},
		{"division before subtraction", "10 - 6 / 2", 7},
		{"complex precedence", "2 * 3 + 4 / 2 - 1", 7},

		// Mixed operations
		{"all operators", "1 + 2 * 3 - 4 / 2", 5},
		{"with modulo", "10 + 5 % 3 * 2", 14},
		{"nested parentheses", "((1 + 2) * (3 - 1)) / 2", 3},

		// Floating point
		{"float addition", "0.1 + 0.2", 0.3},
		{"float multiplication", "1.5 * 2", 3},
		{"float division", "3.14 / 2", 1.57},
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
			assert.InDelta(t, tt.expected, num.Value, 0.000001)
		})
	}
}
