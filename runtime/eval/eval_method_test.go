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

// TestEvalMethodCalls verifies that method calls work correctly
func TestEvalMethodCalls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic method calls
		{"simple method call", `
			var obj = {
				value: 5,
				getValue: fun(self) {
					return self.value
				}
			}
			obj:getValue()
		`, 5},
		{"method with parameters", `
			var obj = {
				add: fun(self, a, b) {
					return a + b
				}
			}
			obj:add(2, 3)
		`, 5},

		// Method modifying object state
		{"method modifies object", `
			var counter = {
				value: 0,
				inc: fun(self) {
					self.value = self.value + 1
					return self.value
				}
			}
			counter:inc()
			counter:inc()
			counter:inc()
		`, 3},

		// Chained method calls
		{"chained methods", `
			var obj = {
				value: 0,
				add: fun(self, x) {
					self.value = self.value + x
					return self
				},
				get: fun(self) {
					return self.value
				}
			}
			obj:add(5):add(3):get()
		`, 8},
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
