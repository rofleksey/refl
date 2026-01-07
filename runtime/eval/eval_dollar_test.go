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

// TestEvalDollarOperator verifies that the $ operator works correctly
func TestEvalDollarOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		// Basic $ usage
		{"$ type", "type($)", "object"},
		{"$ string representation", "str($)", "$"},
		{"$ truthiness", "!!$", float64(1)},

		// Access globals via $
		{"access global via $", "x = 5\n$[\"x\"]", float64(5)},
		{"access multiple globals", "x = 1\ny = 2\n$[\"x\"] + $[\"y\"]", float64(3)},
		{"access non-existent global", "$[\"nonexistent\"]", nil},

		// $ iteration
		{"iterate $", `
			x = 1
			y = 2
			var sum = 0
			for key, val in $ {
				if type(val) == "number" {
					sum = sum + val
				}
			}
			sum >= 3
		`, float64(1)},

		// $ in functions
		{"$ in function", `
			x = 10
			getGlobal = fun (name) {
				return $[name]
			}
			getGlobal("x")
		`, float64(10)},

		{"$ captures current globals", `
			x = 1
			var getX = fun() { return $["x"] }
			x = 2
			getX()
		`, float64(2)},

		{"$ doesn't see locals", `
			test = fun () {
				var local = 42
				return $["local"]
			}
			test()
		`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			program := parseProgram(t, tt.input)
			env := runtime.NewEnvironment(nil)

			result, err := Eval(ctx, program, env)
			require.NoError(t, err)

			switch expected := tt.expected.(type) {
			case string:
				assert.IsType(t, &objects.String{}, result)
				str := result.(*objects.String)
				assert.Equal(t, expected, str.Value)
			case float64:
				assert.IsType(t, &objects.Number{}, result)
				num := result.(*objects.Number)
				assert.Equal(t, expected, num.Value)
			}
		})
	}
}
