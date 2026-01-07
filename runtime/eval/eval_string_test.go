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

// TestEvalString verifies that string concatenation works correctly
func TestEvalString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// String + String
		{"two strings", `"hello" + " " + "world"`, "hello world"},
		{"empty strings", `"" + ""`, ""},
		{"string with empty", `"hello" + ""`, "hello"},
		{"multiple strings", `"a" + "b" + "c"`, "abc"},

		// String + Number
		{"string + number", `"a" + 1`, "a1"},
		{"string + negative number", `"value: " + (-5)`, "value: -5"},
		{"string + decimal", `"pi = " + 3.14`, "pi = 3.14"},
		{"string + zero", `"count: " + 0`, "count: 0"},

		// Number + String
		{"number + string", `1 + "a"`, "1a"},
		{"negative + string", `-5 + " items"`, "-5 items"},
		{"decimal + string", `3.14 + " is pi"`, "3.14 is pi"},

		// Complex expressions
		{"chained concatenation", `1 + " + " + 2 + " = " + 3`, "1 + 2 = 3"},
		{"arithmetic in string context", `"result: " + (2 + 3)`, "result: 5"},

		// Unicode
		{"unicode strings", `"café" + " " + "☕"`, "café ☕"},
		{"emoji concatenation", `"👍" + "👎"`, "👍👎"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			program := parseProgram(t, tt.input)
			env := runtime.NewEnvironment(nil)
			result, err := Eval(ctx, program, env)
			require.NoError(t, err)

			assert.IsType(t, &objects.String{}, result)
			str := result.(*objects.String)
			assert.Equal(t, tt.expected, str.Value)
		})
	}
}
