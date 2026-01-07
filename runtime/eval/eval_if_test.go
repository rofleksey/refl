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

// TestEvalIfStatements verifies that if statements work correctly
func TestEvalIfStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		// Simple if
		{"if true", "if 1 { 5 }", float64(5)},
		{"if false", "if 0 { 5 }", nil},
		{"if with expression", "if 2 > 1 { 10 }", float64(10)},
		{"if false with expression", "if 1 > 2 { 10 }", nil},

		// If-else
		{"if true else", "if 1 { 5 } else { 10 }", float64(5)},
		{"if false else", "if 0 { 5 } else { 10 }", float64(10)},
		{"complex condition else", "if 1 && 0 { 5 } else { 10 }", float64(10)},

		// If-elif-else
		{"if true elif", "if 1 { 5 } elif 1 { 10 } else { 15 }", float64(5)},
		{"if false elif true", "if 0 { 5 } elif 1 { 10 } else { 15 }", float64(10)},
		{"if false elif false else", "if 0 { 5 } elif 0 { 10 } else { 15 }", float64(15)},
		{"multiple elif", "if 0 { 1 } elif 0 { 2 } elif 1 { 3 } elif 1 { 4 } else { 5 }", float64(3)},

		// String conditions
		{"string true", `if "hello" { 1 }`, float64(1)},
		{"string false", `if "" { 1 }`, nil},
		{"string else", `if "" { 1 } else { 2 }`, float64(2)},

		// Object conditions
		{"object true", "if {} { 1 }", float64(1)},
		{"empty object true", "if {} { 1 } else { 2 }", float64(1)},

		// Nil conditions
		{"nil false", "if nil { 1 }", nil},
		{"nil else", "if nil { 1 } else { 2 }", float64(2)},

		// Function conditions
		{"function true", "if fun() { 1 } { 5 }", float64(5)},

		// Variable assignment in if
		{"assign in if true", "var x = 0\nif 1 { x = 5 }\nx", float64(5)},
		{"assign in if false", "var x = 0\nif 0 { x = 5 }\nx", float64(0)},
		{"assign in else", "var x = 0\nif 0 { x = 1 } else { x = 2 }\nx", float64(2)},
		{"assign in elif", "var x = 0\nif 0 { x = 1 } elif 1 { x = 2 } else { x = 3 }\nx", float64(2)},

		// Return from if
		{"return from if", "fun() { if 1 { return 5 } return 10 }()", float64(5)},
		{"return from else", "fun() { if 0 { return 5 } else { return 10 } }()", float64(10)},
		{"return from elif", "fun() { if 0 { return 1 } elif 1 { return 2 } else { return 3 } }()", float64(2)},

		// Nested if
		{"nested if true true", "if 1 { if 1 { 5 } }", float64(5)},
		{"nested if true false", "if 1 { if 0 { 5 } }", nil},
		{"nested if false", "if 0 { if 1 { 5 } }", nil},
		{"nested with else", "if 1 { if 0 { 1 } else { 2 } } else { 3 }", float64(2)},

		// Complex conditions
		{"and condition true", "if 1 && 1 { 5 }", float64(5)},
		{"and condition false", "if 1 && 0 { 5 }", nil},
		{"or condition true", "if 0 || 1 { 5 }", float64(5)},
		{"or condition false", "if 0 || 0 { 5 }", nil},
		{"not condition", "if !0 { 5 }", float64(5)},
		{"comparison condition", "if 5 > 3 { 1 } else { 0 }", float64(1)},

		// Block scoping
		{"scope in if", "var x = 1\nif 1 { var x = 2 }\nx", float64(1)},
		{"scope in else", "var x = 1\nif 0 { var x = 2 } else { var x = 3 }\nx", float64(1)},
		{"modify in if", "var x = 1\nif 1 { x = 2 }\nx", float64(2)},

		// Empty blocks
		{"empty if true", "if 1 { }", nil},
		{"empty if false", "if 0 { }", nil},
		{"empty if else", "if 0 { } else { }", nil},
		{"empty if with value else", "if 0 { } else { 5 }", float64(5)},

		// Last expression value
		{"last expression in if", "if 1 { 1\n2\n3 }", float64(3)},
		{"last expression in else", "if 0 { 1 } else { 4\n5\n6 }", float64(6)},
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

			if tt.expected == nil {
				assert.Equal(t, runtime.NilType, result.Type())
				assert.Same(t, objects.NilInstance, result)
			} else {
				assert.IsType(t, &objects.Number{}, result)
				num := result.(*objects.Number)
				assert.Equal(t, tt.expected.(float64), num.Value)
			}
		})
	}
}
