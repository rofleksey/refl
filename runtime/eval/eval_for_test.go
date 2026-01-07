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

// TestEvalForStatements verifies that for statements work correctly
func TestEvalForStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Basic for loops over arrays
		{"sum array", "var sum = 0\nfor i, val in {1, 2, 3, 4} { sum = sum + val }\nsum", 10},
		{"sum with index", "var sum = 0\nfor i, val in {5, 10, 15} { sum = sum + i }\nsum", 3}, // 0+1+2=3
		{"empty array", "var count = 0\nfor i, val in {} { count = count + 1 }\ncount", 0},
		{"single element", "var sum = 0\nfor i, val in {42} { sum = sum + val }\nsum", 42},

		// For loops over objects
		{"sum object values", "var sum = 0\nfor key, val in {a: 1, b: 2, c: 3} { sum = sum + val }\nsum", 6},
		{"count object keys", "var count = 0\nfor key, val in {x: 1, y: 2, z: 3} { count = count + 1 }\ncount", 3},
		{"empty object", "var count = 0\nfor key, val in {} { count = count + 1 }\ncount", 0},

		// For loops over strings
		{"iterate string", "var result = \"\"\nfor i, ch in \"abc\" { result = result + ch }\nlen(result)", 3},
		{"string indices", "var sum = 0\nfor i, ch in \"hello\" { sum = sum + i }\nsum", 10}, // 0+1+2+3+4=10
		{"empty string", "var count = 0\nfor i, ch in \"\" { count = count + 1 }\ncount", 0},

		// Break statement
		{"break early array", "var sum = 0\nfor i, val in {1, 2, 3, 4} { if val == 3 { break } sum = sum + val }\nsum", 3},              // 1+2
		{"break object", "var sum = 0\nfor _, val in {1, 2, 3, 4} { if val > 2 { break } sum = sum + val }\nsum", 3},                    // 1+2
		{"break string", "var result = \"\"\nfor i, ch in \"hello\" { if ch == \"l\" { break } result = result + ch }\nlen(result)", 2}, // "he"

		// Continue statement
		{"continue array", "var sum = 0\nfor i, val in {1, 2, 3, 4} { if val == 2 { continue } sum = sum + val }\nsum", 8},                    // 1+3+4
		{"continue object", "var sum = 0\nfor _, val in {1, 2, 3} { if val == 2 { continue } sum = sum + val }\nsum", 4},                      // 1+3
		{"continue string", "var result = \"\"\nfor i, ch in \"hello\" { if ch == \"l\" { continue } result = result + ch }\nlen(result)", 3}, // "heo"

		// Return from for
		{"return from for", `
			fun() {
				for i, val in {10, 20, 30} {
					if val == 20 { return val }
				}
				return 0
			}()
		`, 20},
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
