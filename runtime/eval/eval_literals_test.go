package eval

import (
	"context"
	"refl/runtime"
	"refl/runtime/objects"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEvalNumberLiterals verifies that number literals are properly evaluated
func TestEvalNumberLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"positive integer", "42", 42},
		{"decimal", "3.14", 3.14},
		{"zero", "0", 0},
		{"negative integer", "-5", -5},
		{"negative decimal", "-3.14", -3.14},
		{"large number", "1000000", 1000000},
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

// TestEvalStringLiterals verifies that string literals are properly evaluated
func TestEvalStringLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", `"hello"`, "hello"},
		{"empty string", `""`, ""},
		{"string with spaces", `"hello world"`, "hello world"},
		{"escaped quotes", `"escape \"quote\""`, `escape "quote"`},
		{"raw string", "`raw string`", "raw string"},
		{"raw string with quotes", "`raw \"string\"`", `raw "string"`},
		{"mixed quotes in raw", "`it's raw`", "it's raw"},
		{"multiline raw", "`line1\nline2`", "line1\nline2"},
		{"unicode string", `"café"`, "café"},
		{"special characters", `"\t\n\r"`, "\t\n\r"},
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

			assert.IsType(t, &objects.String{}, result)
			str := result.(*objects.String)
			assert.Equal(t, tt.expected, str.Value)
		})
	}
}

// TestEvalNilLiteral verifies that nil literal is properly evaluated
func TestEvalNilLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected runtime.Object
	}{
		{"nil literal", "nil", objects.NilInstance},
		{"nil comparison", "nil == nil", objects.NewBoolean(true)},
		{"nil inequality", "nil != 5", objects.NewBoolean(true)},
		{"nil truthiness", "!nil", objects.NewBoolean(true)},
		{"nil in object", "{key: nil}.key", objects.NilInstance},
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

			if tt.expected == objects.NilInstance {
				assert.Equal(t, runtime.NilType, result.Type())
				assert.Same(t, objects.NilInstance, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
