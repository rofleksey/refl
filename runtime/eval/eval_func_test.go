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

// TestEvalFunction verifies that function arguments work correctly
func TestEvalFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic argument passing
		{"single argument", "fun(x) { return x }(42)", 42},
		{"multiple arguments", "fun(a, b, c) { return a + b + c }(1, 2, 3)", 6},
		{"argument order matters", "fun(a, b) { return a - b }(10, 3)", 7},

		// Argument expressions
		{"expression as argument", "fun(x) { return x * 2 }(5 + 3)", 16},
		{"function call as argument", `
			var double = fun(x) { return x * 2 }
			fun(y) { return y + 1 }(double(4))
		`, 9}, // double(4)=8, 8+1=9

		// Default argument values (nil when not provided)
		{"missing arguments become nil", "fun(a, b) { if b == nil { return a } else { return a + b } }(5)", 5},
		{"all missing arguments", "fun(a, b, c) { return len(type(a)) }\n()", 3},

		// Extra arguments are ignored
		{"extra arguments ignored", "fun(x) { return x }(1, 2, 3, 4, 5)", 1},

		// Object arguments (passed by reference)
		{"modify object argument", `
			var f = fun(obj) {
				obj.value = 99
				return obj.value
			}
			var o = {value: 1}
			f(o)
			o.value
		`, 99}, // object modified

		// Complex argument patterns
		{"nested function arguments", `
			var apply = fun(f, x) { return f(x) }
			apply(fun(y) { return y * 3 }, 4)
		`, 12},
		{"function returning function with arguments", `
			var makeMultiplier = fun(factor) {
				return fun(x) { return x * factor }
			}
			makeMultiplier(3)(4)
		`, 12},

		// Arguments in recursive calls
		{"recursive with arguments", `
			var sumTo = fun(n) {
				if n == 0 { return 0 }
				return n + sumTo(n - 1)
			}
			sumTo(5)
		`, 15}, // 1+2+3+4+5=15

		// Variable number of arguments using args
		{"variable arguments via args", `
			var sumAll = fun() {
				var total = 0
				for i, val in args {
					total = total + val
				}
				return total
			}
			sumAll(1, 2, 3, 4, 5)
		`, 15},
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
