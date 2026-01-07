package eval

import (
	"context"
	"refl/runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEvalErrors verifies that error cases are handled properly
func TestEvalErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		errorSubstr string
	}{
		// Type errors
		{"nil addition", "nil + 1", "cannot apply operator + to types nil and number"},
		{"number plus nil", "1 + nil", "cannot add number to nil"},
		{"object plus number", "{} + 5", "cannot apply operator + to types object and number"},
		{"number plus object", "5 + {}", "cannot add number to object"},

		// Division by zero
		{"division by zero", "1 / 0", "division by zero"},
		{"modulo by zero", "1 % 0", "modulo by zero"},

		// Invalid unary operations
		{"negate nil", "-nil", "cannot negate non-number"},
		{"negate object", "-{}", "cannot negate non-number"},
		{"negate string", `-"hello"`, "cannot negate non-number"},

		// Invalid binary operations
		{"string minus string", `"hello" - "world"`, "cannot apply operator - to types string and string"},
		{"string division", `"hello" / 2`, "cannot apply operator / to types string and number"},
		{"string modulo", `"hello" % 2`, "cannot apply operator % to types string and number"},

		// Invalid function calls
		{"call number", "5()", "attempt to call non-function"},
		{"call object", "{}()", "attempt to call non-function"},
		{"call nil", "nil()", "attempt to call non-function"},
		{"call non-existent builtin", "nonexistent()", "attempt to call non-function"},

		// Invalid iteration
		{"iterate number", "for x, _ in 5 {}", "cannot iterate over non-iterable object"},
		{"iterate nil", "for x, _ in nil {}", "cannot iterate over non-iterable object"},

		// Invalid index access
		{"index number", "1[0]", "cannot access member of non-indexable object"},
		{"string index out of bounds", `"hello"[10]`, "string index out of bounds"},
		{"string non-numeric index", `"hello"["key"]`, "string index must be a number"},
		{"nil index", "nil[0]", "cannot access member of non-indexable object"},

		// Invalid property access
		{"property on number", "5.length", "cannot access member of non-indexable object"},
		{"property on nil", "nil.length", "cannot access member of non-indexable object"},

		// Invalid assignments
		{"assign to string index", `"hello"[0] = "x"`, "strings are immutable"},
		{"assign to nil property", "nil.x = 5", "cannot assign to member of non-indexable object"},
		{"assign to number property", "5.x = 10", "cannot assign to member of non-indexable object"},

		// Chained invalid access
		{"chained invalid property", `var x = {}\nx.y = 5\nx.y.z`, "cannot assign to member of non-indexable object"},
		{"property on result of invalid", "nil.x.y", "cannot access member of non-indexable object"},

		// Invalid object operations
		{"call object property that is not function", `var obj = {x: 5}\nobj:x()`, "cannot access method of non-indexable object"},
		{"method not found", `var obj = {}\nobj:nonExistent()`, "cannot access method of non-indexable object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			program := parseProgram(t, tt.input)
			env := runtime.NewEnvironment(nil)

			evaluator := New(ctx, program, env)
			_, err := evaluator.Run()
			require.Error(t, err)

			if tt.errorSubstr != "" {
				assert.Contains(t, err.Error(), tt.errorSubstr)
			}
		})
	}
}
