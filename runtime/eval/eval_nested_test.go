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

// TestEvalNestedScopes verifies that nested scopes work correctly
func TestEvalNestedScopes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic nested blocks
		{"nested blocks", `
			var x = 1
			{
				var x = 2
				{
					var x = 3
				}
			}
			x
		`, 1},

		{"deeply nested blocks", `
			var a = 1
			{
				var b = 2
				{
					var c = 3
					{
						var d = 4
						a + b + c + d
					}
				}
			}
		`, 10},

		// Nested functions
		{"nested function scopes", `
			var global = 1
			var outer = fun () {
				var outerVar = 2
				var inner = fun() {
					var innerVar = 3
					return global + outerVar + innerVar
				}
				return inner()
			}
			outer()
		`, 6},

		{"nested if-else scopes", `
			var result = 0
			if 1 {
				var a = 1
				if 0 {
					var b = 2
				} else {
					var c = 3
					result = a + c
				}
			}
			result
		`, 4},

		// Nested loops
		{"nested while loops", `
			var sum = 0
			var i = 0
			while i < 2 {
				var j = 0
				while j < 2 {
					var product = i * j
					sum = sum + product
					j = j + 1
				}
				i = i + 1
			}
			sum
		`, 1}, // (0*0)+(0*1)+(1*0)+(1*1)=0+0+0+1=1

		{"nested for loops", `
			var total = 0
			for i, x in {1, 2} {
				for j, y in {3, 4} {
					var product = x * y
					total = total + product
				}
			}
			total
		`, 21}, // (1*3)+(1*4)+(2*3)+(2*4)=3+4+6+8=21
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
