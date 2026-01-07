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

// TestEvalArgsVariable verifies that the special 'args' variable works correctly
func TestEvalArgsVariable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic args usage
		{"args length with no args", `
			fun() {
				return len(args)
			}()
		`, 0},

		{"args length with args", `
			fun() {
				return len(args)
			}(1, 2, 3)
		`, 3},

		{"sum using args", `
			fun() {
				var total = 0
				for i, val in args {
					total = total + val
				}
				return total
			}(1, 2, 3, 4)
		`, 10},

		{"access args by index", `
			fun() {
				return args[0] + args[1]
			}(5, 3)
		`, 8},

		// Args with named parameters
		{"args with named params", `
			fun(x, y) {
				return x + y + args[2]
			}(1, 2, 3)
		`, 6}, // 1 + 2 + 3 = 6

		{"more args than params", `
			fun(a, b) {
				return a + b + len(args)
			}(1, 2, 3, 4, 5)
		`, 8}, // 1 + 2 + 5 = 8

		{"pass args to inner function", `
			fun() {
				var process = fun(arr) {
					var sum = 0
					for i, val in arr {
						sum = sum + val
					}
					return sum
				}
				return process(args)
			}(1, 2, 3)
		`, 6},

		// Args shadowing
		{"args shadows outer args", `
			var args = 42
			fun() {
				var args = 99
				return args
			}()
		`, 99},

		{"local args variable", `
			fun() {
				var args = 5
				return args
			}(1, 2, 3)
		`, 5},

		// Args in method calls
		{"args in method", `
			var obj = {
				sumAll: fun(self) {
					var total = 0
					for i, val in args {
						if i == 0 {
							continue
						}
						total = total + val
					}
					return total
				}
			}
			obj:sumAll(1, 2, 3)
		`, 6},

		// Args iteration patterns
		{"iterate args with index", `
			fun() {
				var sum = 0
				for i, val in args {
					sum = sum + i * val
				}
				return sum
			}(1, 2, 3)
		`, 8}, // 0*1 + 1*2 + 2*3 = 0+2+6=8
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
