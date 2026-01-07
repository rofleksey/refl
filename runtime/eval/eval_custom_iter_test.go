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

// TestEvalCustomIterators verifies that custom iterators work correctly
func TestEvalCustomIterators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic custom iterator
		{"custom array iterator", `
			var obj = {"a", "b", "c"}
			obj.__iter = fun(self, yield) {
				yield(0, "a")
				yield(1, "b")
				yield(2, "c")
			}
			var result = ""
			for key, val in obj {
				result = result + val
			}
			len(result)
		`, 3},

		{"custom range iterator", `
			var Range = {
				new: fun(self, start, end) {
					var inst = clone(self)
					inst.start = start
					inst.end = end
					return inst
				},
				__iter: fun(self, yield) {
					var i = 0
					var current = self.start
					while current < self.end {
						if !yield(i, current) { return }
						i = i + 1
						current = current + 1
					}
				}
			}
			var r = Range:new(5, 10)
			var sum = 0
			for i, val in r {
				sum = sum + val
			}
			sum
		`, 35}, // 5+6+7+8+9=35
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
