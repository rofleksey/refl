package eval

import (
	"context"
	"refl/ast"
	"refl/parser"
	"refl/runtime"
	"refl/runtime/objects"
	"testing"
	"time"
)

func TestEvalNumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42", 42},
		{"3.14", 3.14},
		{"0", 0},
		{"-5", -5},
		{"-3.14", -3.14},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"world"`, "world"},
		{`""`, ""},
		{`"escape \"quote\""`, `escape "quote"`},
		{"`raw string`", "raw string"},
		{"`raw \"string\"`", `raw "string"`},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkStringObject(t, result, tt.expected)
	}
}

func TestEvalNilLiteral(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	program := parseProgram(t, "nil")
	env := runtime.NewEnvironment(nil)
	result, err := Eval(ctx, program, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Type() != runtime.NilType {
		t.Fatalf("expected nil type, got %s", result.Type())
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1", true},
		{"0", false},
		{"!0", true},
		{"!1", false},
		{"!!1", true},
		{"!!0", false},
		{"1 == 1", true},
		{"1 == 0", false},
		{"1 != 0", true},
		{"1 != 1", false},
		{"2 > 1", true},
		{"1 > 2", false},
		{"1 < 2", true},
		{"2 < 1", false},
		{"2 >= 2", true},
		{"2 >= 1", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 <= 2", true},
		{"2 <= 1", false},
		{"1 && 1", true},
		{"1 && 0", false},
		{"0 && 1", false},
		{"0 && 0", false},
		{"1 || 1", true},
		{"1 || 0", true},
		{"0 || 1", true},
		{"0 || 0", false},
		{`"a" == "a"`, true},
		{`"a" != "b"`, true},
		{`"a" < "b"`, true},
		{`"b" > "a"`, true},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkBooleanObject(t, result, tt.expected)
	}
}

func TestEvalArithmeticExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1 + 2", 3},
		{"5 - 3", 2},
		{"2 * 3", 6},
		{"6 / 2", 3},
		{"7 % 3", 1},
		{"-5", -5},
		{"-(-5)", 5},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
		{"10 / 2 * 3", 15},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error evaluating %s: %v", tt.input, err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello" + " " + "world"`, "hello world"},
		{`"a" + 1`, "a1"},
		{`1 + "a"`, "1a"},
		{`"test" * 3`, "testtesttest"},
		{`3 * "x"`, "xxx"},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkStringObject(t, result, tt.expected)
	}
}

func TestEvalVarDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"var x = 5", float64(5)},
		{"var x = 5\nx", float64(5)},
		{"var x = 5\nvar y = x\nx + y", float64(10)},
		{"var x = nil", nil},
		{`var x = "hello"`, "hello"},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		switch expected := tt.expected.(type) {
		case float64:
			checkNumberObject(t, result, expected)
		case string:
			checkStringObject(t, result, expected)
		case nil:
			if result.Type() != runtime.NilType {
				t.Fatalf("expected nil, got %s", result.Type())
			}
		}
	}
}

func TestEvalAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"x = 5", 5},
		{"x = 5\nx = 10\nx", 10},
		{"x = 1\ny = 2\nx + y", 3},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if 1 { 5 }", float64(5)},
		{"if 0 { 5 }", nil},
		{"if 0 { 5 } else { 10 }", float64(10)},
		{"if 0 { 5 } elif 1 { 15 } else { 10 }", float64(15)},
		{"if 0 { 5 } elif 0 { 15 } else { 10 }", float64(10)},
		{"var x = 0\nif 1 { x = 5 }\nx", float64(5)},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tt.expected == nil {
			if result.Type() != runtime.NilType {
				t.Fatalf("expected nil, got %v", result)
			}
		} else {
			checkNumberObject(t, result, tt.expected.(float64))
		}
	}
}

func TestEvalWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var x = 0\nwhile x < 5 { x = x + 1 }\nx", 5},
		{"var x = 0\nvar sum = 0\nwhile x < 5 { sum = sum + x\nx = x + 1 }\nsum", 10},
		{"var x = 10\nwhile x > 0 { x = x - 1 if x == 5 { break } }\nx", 5},
		{"var x = 0\nvar sum = 0\nwhile x < 5 { x = x + 1 if x == 3 { continue } sum = sum + x }\nsum", 12}, // 1 + 2 + 4 + 5 = 12
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalForStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var sum = 0\nfor i, val in {1, 2, 3, 4} { sum = sum + val }\nsum", 10},
		{"var sum = 0\nfor i, val in {1, 2, 3, 4} { if val == 3 { break } sum = sum + val }\nsum", 3},    // 1 + 2
		{"var sum = 0\nfor i, val in {1, 2, 3, 4} { if val == 2 { continue } sum = sum + val }\nsum", 8}, // 1 + 3 + 4
		{"var count = 0\nfor i, _ in {1, 2, 3, 4} { count = count + 1 }\ncount", 4},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalFunctionLiteral(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := "fun (x) { x + 1 }"
	program := parseProgram(t, input)
	env := runtime.NewEnvironment(nil)
	result, err := Eval(ctx, program, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Type() != runtime.FunctionType {
		t.Fatalf("expected function type, got %s", result.Type())
	}
}

func TestEvalFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var add = fun (a, b) { return a + b }\nadd(2, 3)", 5},
		{"var fib = fun (n) { if n == 0 { return 0 } if n == 1 { return 1 } return fib(n-1) + fib(n-2) }\nfib(6)", 8},
		{"var fact = fun (n) { if n == 0 { return 1 } return n * fact(n-1) }\nfact(5)", 120},
		{"var adder = fun (x) { return fun (y) { return x + y } }\nvar add5 = adder(5)\nadd5(3)", 8},
		{"var x = fun () { return 42 }\nx()", 42},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalObjectLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"{}.a", nil},
		{"{a: 1}.a", float64(1)},
		{`{"key": "value"}.key`, "value"},
		{"{1, 2, 3}[1]", float64(2)},
		{"var obj = {x: 5, y: 10}\nobj.x + obj.y", float64(15)},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		switch expected := tt.expected.(type) {
		case float64:
			checkNumberObject(t, result, expected)
		case string:
			checkStringObject(t, result, expected)
		case nil:
			if result.Type() != runtime.NilType {
				t.Fatalf("expected nil, got %v", result)
			}
		}
	}
}

func TestEvalMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`var obj = {add: fun (self, a, b) { return a + b }} obj:add(2, 3)`, 5},
		{`var Counter = {
    new: fun(self, initial) {
        var inst = clone(self)
        inst.value = initial
        return inst
    },
    inc: fun(self) {
        self.value = self.value + 1
        return self.value
    }
}
var c = Counter:new(5)
c:inc()`, 6},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"type(5)", "number"},
		{`type("hello")`, "string"},
		{"type(nil)", "nil"},
		{"type({})", "object"},
		{"type(fun () {})", "function"},
		{`str(5)`, "5"},
		{`str("hello")`, "hello"},
		{`str(nil)`, "nil"},
		{`number("42")`, float64(42)},
		{`number("abc")`, nil},
		{`number(5)`, float64(5)},
		{`len("hello")`, float64(5)},
		{`len({1, 2, 3})`, float64(3)},
		{`len({})`, float64(0)},
		{`var obj = {x: 1, y: 2} var obj2 = clone(obj) obj2.x`, float64(1)},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		switch expected := tt.expected.(type) {
		case float64:
			checkNumberObject(t, result, expected)
		case string:
			checkStringObject(t, result, expected)
		case nil:
			if result.Type() != runtime.NilType {
				t.Fatalf("expected nil, got %v", result)
			}
		}
	}
}

func TestEvalErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{"nil + 1", true},
		{"1 + nil", true},
		{"{} + 5", true},
		{"5 + {}", true},
		{"1 / 0", true},
		{"1 % 0", true},
		{"-nil", true},
		{`"hello" - "world"`, true},
		{`"hello" / 2`, true},
		{`"hello" % 2`, true},
		{"-{}", true},
		{"5()", true},
		{"{}()", true},
		{"exit()", true},
		{"for x, _ in 5 {}", true},
		{"1[0]", true},
		{`"hello"[10]`, true},
		{`"hello"["key"]`, true},
		{`"hello".length`, true},
		{`"hello"[0] = "x"`, true},
		{"var x = {}\nx.y = 5\nx.y.z", true},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		_, err := Eval(ctx, program, env)

		if tt.expectError && err == nil {
			t.Errorf("expected error for input: %s, got nil", tt.input)
		} else if !tt.expectError && err != nil {
			t.Errorf("unexpected error for input: %s, got %v", tt.input, err)
		}
	}
}

func TestEvalGlobalVariable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := `var x = 10
var getGlobal = fun () {
    var global = $
    return global.x
}
getGlobal()`

	program := parseProgram(t, input)
	env := runtime.NewEnvironment(nil)
	result, err := Eval(ctx, program, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkNumberObject(t, result, 10)
}

func TestEvalArgsVariable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := `var sumArgs = fun () {
    var total = 0
    for i, val in args {
        total = total + val
    }
    return total
}
sumArgs(1, 2, 3, 4)`

	program := parseProgram(t, input)
	env := runtime.NewEnvironment(nil)
	result, err := Eval(ctx, program, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkNumberObject(t, result, 10)
}

func TestEvalGlobalVsLocalVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"x = 5\nvar x = 10\nx", 10},
		{"var x = 5\nx = 10\nx", 10},
		{"x = 5\ntest = fun () { var x = 10\nreturn x }\ntest()", 10},
		{"x = 5\ntest = fun () { x = 10 }\ntest()\nx", 10},
		{"var x = 1\nouter = fun () { var x = 2\ninner = fun() { return x }\nreturn inner() }\nouter()", 2},
		{"x = 1\nouter = fun () { x = 2 }\nouter()\nx", 2},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalVariableShadowing(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var x = 1\nif 1 { var x = 2 }\nx", 1},
		{"var x = 1\nif 1 { x = 2 }\nx", 2},
		{"var x = 1\nwhile x < 3 { var x = 5\nbreak }\nx", 1},
		{"var x = 1\nfor i, _ in {1,2} { var x = 2 }\nx", 1},
		{"var x = 1\nfor i, _ in {1,2} { x = 2 }\nx", 2},
		{"var x = 1\ntest = fun () { var x = 2 }\ntest()\nx", 1},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalClosureScoping(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`var makeCounter = fun() {
    var count = 0
    return fun() {
        count = count + 1
        return count
    }
}
var counter = makeCounter()
counter()
counter()`, 2},
		{`var makeAdder = fun(x) {
    return fun(y) { return x + y }
}
var add5 = makeAdder(5)
add5(3)`, 8},
		{`var x = 10
var getX = fun() { return x }
x = 20
getX()`, 20},
		{`var x = 1
var outer = fun() {
    var x = 2
    var inner = fun() { return x }
    return inner
}
var fn = outer()
fn()`, 2},
		{`var createMultiplier = fun(a) {
    return fun(b) { return a * b }
}
var times3 = createMultiplier(3)
var times4 = createMultiplier(4)
times3(2) + times4(2)`, 14},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalNestedScopes(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`var x = 1
if 1 {
    var x = 2
    if 1 {
        var x = 3
    }
    x
}`, 2},
		{`var x = 1
outer = fun () {
    var x = 2
    middle = fun () {
        var x = 3
        inner = fun () {
            return x
        }
        return inner()
    }
    return middle()
}
outer()`, 3},
		{`var x = 1
test = fun () {
    var x = 2
    if 1 {
        var x = 3
        if 1 {
            return x
        }
    }
}
test()`, 3},
		{`var result = 0
for _, val in {1,2,3} {
    var x = val
    if 1 {
        var y = x * 2
        result = result + y
    }
}
result`, 12},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalGlobalVariableAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`x = 5
getGlobal = fun () { return $["x"] }
getGlobal()`, 5},
		{`var x = 10
test = fun () { return $["x"] }
test()`, 10},
		{`test = fun () { x = 5 } 
test()
$["x"]`, 5},
		{`a = 1
b = 2
c = 3
sumGlobals = fun () {
    var total = 0
    for key, val in $ {
        if type(val) == "number" {
            total = total + val
        }
    }
    return total
}
sumGlobals()`, 6},
		{`x = 1
outer = fun () {
    var x = 2
    inner = fun () {
        return $["x"]
    }
    return inner()
}
outer()`, 1},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalFunctionParameterScoping(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`var x = 1
test = fun (x) { return x }
test(5)`, 5},
		{`var x = 1
test = fun (y) { var x = y * 2 return x }
test(3)`, 6},
		{`test = fun(a, b) {
    var result = a + b
    return fun(c) { return result + c }
}
var addToSum = test(1, 2)
addToSum(3)`, 6},
		{`var x = 10
test = fun (x) { 
    var inner = fun () { return x }
    return inner()
}
test(20)`, 20},
		{`recursive = fun (x) {
    if x == 0 { return 0 }
    return x + recursive(x - 1)
}
recursive(5)`, 15},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalArgsVariableScoping(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`sum = fun () {
		   var total = 0
		   for i, val in args {
		       total = total + val
		   }
		   return total
		}
		sum(1, 2, 3, 4)`, 10},
		{`test = fun () {
		   var args = 5
		   return args
		}
		test(1, 2, 3)`, 5},
		{`outer = fun () {
		   var inner = fun () {
		       return len(args)
		   }
		   return inner()
		}
		outer(1, 2, 3)`, 0},
		{`var args = 99
		test = fun () { return args[0] }
		test(1)`, 1},
		{`nested = fun () {
    var inner = fun () {
        return args[0]
    }
    return inner(35)
}
nested(42, 100)`, 35},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalBlockScoping(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`var x = 1
{
    var x = 2
    var y = 3
}
x`, 1},
		{`var x = 1
{
    x = 2
    var y = 3
}
x`, 2},
		{`var result = 0
{
    var x = 5
    {
        var y = 10
        result = x + y
    }
}
result`, 15},
		{`test = fun () {
    var x = 1
    {
        var x = 2
        {
            return x
        }
    }
}
test()`, 2},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalVariableDeletion(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`x = 5
var x = 10
x`, 10},
		{`var x = 5
test = fun () { var x = 10 return x }
test()`, 10},
		{`var x = 1
test = fun () { var x = 2\nvar inner = fun() { return x } return inner() }
test()`, 2},
	}

	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		program := parseProgram(t, tt.input)
		env := runtime.NewEnvironment(nil)
		result, err := Eval(ctx, program, env)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkNumberObject(t, result, tt.expected)
	}
}

func TestEvalComplexClosureChain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := `var makeBankAccount = fun(initial) {
    var balance = initial
    return {
        deposit: fun(self, amount) {
            balance = balance + amount
            return balance
        },
        withdraw: fun(self, amount) {
            if amount > balance { return 0 }
            balance = balance - amount
            return amount
        },
        getBalance: fun(self) { return balance }
    }
}
var account = makeBankAccount(100)
account:deposit(50)
account:withdraw(30)
account:getBalance()`

	program := parseProgram(t, input)
	env := runtime.NewEnvironment(nil)
	result, err := Eval(ctx, program, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkNumberObject(t, result, 120)
}

func TestEvalEnvironmentChain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := `var globalVar = 10
outer = fun () {
    var outerVar = 20
    var middle = fun() {
        var middleVar = 30
				var inner = fun () {
            var innerVar = 40
            return globalVar + outerVar + middleVar + innerVar
        }
        return inner()
    }
    return middle()
}
outer()`

	program := parseProgram(t, input)
	env := runtime.NewEnvironment(nil)
	result, err := Eval(ctx, program, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkNumberObject(t, result, 100)
}

func parseProgram(t *testing.T, input string) *ast.Program {
	p := parser.New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return program
}

func checkNumberObject(t *testing.T, obj runtime.Object, expected float64) {
	t.Helper()

	num, ok := obj.(*objects.Number)
	if !ok {
		t.Fatalf("object is not Number. got=%T (%v)", obj, obj)
	}

	if num.Value != expected {
		t.Fatalf("object has wrong value. got=%v, want=%v", num.Value, expected)
	}
}

func checkStringObject(t *testing.T, obj runtime.Object, expected string) {
	t.Helper()

	str, ok := obj.(*objects.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%v)", obj, obj)
	}

	if str.Value != expected {
		t.Fatalf("object has wrong value. got=%v, want=%v", str.Value, expected)
	}
}

func checkBooleanObject(t *testing.T, obj runtime.Object, expected bool) {
	t.Helper()

	num, ok := obj.(*objects.Number)
	if !ok {
		t.Fatalf("object is not Number (boolean). got=%T (%v)", obj, obj)
	}

	var expectedValue float64
	if expected {
		expectedValue = 1
	}

	if num.Value != expectedValue {
		t.Fatalf("object has wrong value. got=%v, want=%v", num.Value, expectedValue)
	}
}
