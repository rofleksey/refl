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

// TestEvalObjectLiterals verifies that object literals work correctly
func TestEvalObjectLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		// Basic object literals
		{"empty object", "{}", runtime.ObjectType_},
		{"object with number property", "{x: 5}.x", float64(5)},
		{"object with string property", `{"key": "value"}.key`, "value"},
		{"object with multiple properties", "{a: 1, b: 2, c: 3}.b", float64(2)},

		// Array literal syntax (sugar for object)
		{"array with elements", "{1, 2, 3}[1]", float64(2)},
		{"mixed array", `{1, "two", 3}[1]`, "two"},
		{"nested array", "{{1, 2}, {3, 4}}[0][1]", float64(2)},

		// Object property access
		{"dot notation", "var obj = {x: 10}\nobj.x", float64(10)},
		{"bracket notation", `var obj = {key: "value"} obj["key"]`, "value"},
		{"bracket with expression", "var obj = {x: 5}\nvar prop = \"x\"\nobj[prop]", float64(5)},

		// Object equality
		{"object reference equality", "var o1 = {}\nvar o2 = o1\no1 == o2", float64(1)},
		{"different objects not equal", "var o1 = {}\nvar o2 = {}\no1 == o2", float64(0)},
		{"objects with same content not equal", "{x: 1} == {x: 1}", float64(0)},

		// Object operations
		{"object length", "len({a: 1, b: 2, c: 3})", float64(3)},
		{"array length", "len({1, 2, 3, 4, 5})", float64(5)},
		{"empty length", "len({})", float64(0)},

		// Object modification
		{"add property", "var obj = {}\nobj.x = 5\nobj.x", float64(5)},
		{"modify property", "var obj = {x: 1}\nobj.x = 2\nobj.x", float64(2)},
		{"nested modification", "var obj = {inner: {x: 1}}\nobj.inner.x = 2\nobj.inner.x", float64(2)},

		// Object truthiness
		{"empty object truthy", "!!{}", float64(1)},
		{"non-empty object truthy", "!!{x: 1}", float64(1)},

		// Object patterns
		{"object as map", `
			var map = {
				"apple": "red",
				"banana": "yellow",
				"grape": "purple"
			}
			map["apple"]
		`, "red"},

		// Array methods
		{"array push pattern", `
			var arr = {}
			arr[len(arr)] = 1
			arr[len(arr)] = 2
			arr[len(arr)] = 3
			len(arr)
		`, float64(3)},
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

			switch expected := tt.expected.(type) {
			case runtime.ObjectType:
				assert.Equal(t, expected, result.Type())
			case float64:
				assert.IsType(t, &objects.Number{}, result)
				num := result.(*objects.Number)
				assert.Equal(t, expected, num.Value)
			case string:
				assert.IsType(t, &objects.String{}, result)
				str := result.(*objects.String)
				assert.Equal(t, expected, str.Value)
			case nil:
				assert.Equal(t, runtime.NilType, result.Type())
				assert.Same(t, objects.NilInstance, result)
			}
		})
	}
}
