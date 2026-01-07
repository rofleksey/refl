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

// TestEvalLogicalExpressions verifies that logical expressions work correctly
func TestEvalLogicalExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Truthy values
		{"number truthy", "1", true},
		{"string truthy", `!!"hello"`, true},
		{"object truthy", "!!{}", true},

		// Falsy values
		{"zero falsy", "0", false},
		{"empty string falsy", `!!""`, false},
		{"nil falsy", "!!nil", false},

		// NOT operator
		{"not true", "!1", false},
		{"not false", "!0", true},
		{"double not", "!!1", true},
		{"not string", `!"hello"`, false},
		{"not empty string", `!""`, true},

		// Equality
		{"equal numbers", "1 == 1", true},
		{"unequal numbers", "1 == 0", false},
		{"not equal numbers", "1 != 0", true},
		{"equal strings", `"a" == "a"`, true},
		{"unequal strings", `"a" == "b"`, false},
		{"not equal strings", `"a" != "b"`, true},
		{"nil equality", "nil == nil", true},
		{"nil inequality with value", "nil != 1", true},

		// Comparisons
		{"greater than", "2 > 1", true},
		{"greater than false", "1 > 2", false},
		{"less than", "1 < 2", true},
		{"less than false", "2 < 1", false},
		{"greater or equal equal", "2 >= 2", true},
		{"greater or equal greater", "2 >= 1", true},
		{"greater or equal false", "1 >= 2", false},
		{"less or equal equal", "1 <= 1", true},
		{"less or equal less", "1 <= 2", true},
		{"less or equal false", "2 <= 1", false},

		// String comparisons
		{"string less than", `"a" < "b"`, true},
		{"string greater than", `"b" > "a"`, true},
		{"string less or equal", `"a" <= "a"`, true},
		{"string greater or equal", `"b" >= "a"`, true},

		// AND operator
		{"both true", "1 && 1", true},
		{"first false", "0 && 1", false},
		{"second false", "1 && 0", false},
		{"both false", "0 && 0", false},
		{"short-circuit false", "0 && panic()", false},

		// OR operator
		{"both true", "1 || 1", true},
		{"first true", "1 || 0", true},
		{"second true", "0 || 1", true},
		{"both false", "0 || 0", false},
		{"short-circuit true", "1 || panic()", true},

		// Complex expressions
		{"mixed and or", "(1 && 0) || 1", true},
		{"parentheses", "(1 || 0) && 0", false},
		{"chained comparisons", "1 < 2 && 2 < 3", true},
		{"multiple operators", "!(0 || 1) && 1", false},
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
			expectedValue := float64(0)
			if tt.expected {
				expectedValue = 1
			}
			assert.Equal(t, expectedValue, num.Value)
		})
	}
}
