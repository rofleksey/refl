package parser

import (
	"testing"

	"refl/ast"
)

func TestParseEmptyProgram(t *testing.T) {
	p := New()
	program, err := p.Parse("")
	if err != nil {
		t.Fatalf("Failed to parse empty program: %v", err)
	}
	if program == nil {
		t.Fatal("Program is nil")
	}
	if len(program.Statements) != 0 {
		t.Fatalf("Expected 0 statements, got %d", len(program.Statements))
	}
}

func TestParseVarDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		value    interface{}
	}{
		{"var x = 5", "x", 5.0},
		{"var y = nil", "y", nil},
		{"var name = \"hello\"", "name", "hello"},
		{"var pi = 3.14", "pi", 3.14},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.VarDeclaration)
		if !ok {
			t.Fatalf("Expected VarDeclaration, got %T", program.Statements[0])
		}

		if stmt.Name != tt.expected {
			t.Errorf("Expected name %q, got %q", tt.expected, stmt.Name)
		}

		switch expected := tt.value.(type) {
		case float64:
			checkNumberLiteral(t, stmt.Value, expected)
		case string:
			checkStringLiteral(t, stmt.Value, expected)
		case nil:
			_, ok := stmt.Value.(*ast.NilLiteral)
			if !ok {
				t.Errorf("Expected NilLiteral, got %T", stmt.Value)
			}
		}
	}
}

func TestParseAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		value    float64
	}{
		{"x = 5", "x", 5.0},
		{"y = 42", "y", 42.0},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		assign, ok := stmt.Expression.(*ast.Assignment)
		if !ok {
			t.Fatalf("Expected Assignment, got %T", stmt.Expression)
		}

		ident, ok := assign.Left.(*ast.Identifier)
		if !ok {
			t.Fatalf("Expected Identifier on left side, got %T", assign.Left)
		}

		if ident.Name != tt.expected {
			t.Errorf("Expected identifier %q, got %q", tt.expected, ident.Name)
		}

		checkNumberLiteral(t, assign.Right, tt.value)
	}
}

func TestParseBinaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		left     float64
		operator string
		right    float64
	}{
		{"1 + 2", 1.0, "+", 2.0},
		{"3 - 4", 3.0, "-", 4.0},
		{"5 * 6", 5.0, "*", 6.0},
		{"7 / 8", 7.0, "/", 8.0},
		{"9 % 10", 9.0, "%", 10.0},
		{"1 < 2", 1.0, "<", 2.0},
		{"3 > 4", 3.0, ">", 4.0},
		{"5 <= 6", 5.0, "<=", 6.0},
		{"7 >= 8", 7.0, ">=", 8.0},
		{"9 == 10", 9.0, "==", 10.0},
		{"11 != 12", 11.0, "!=", 12.0},
		{"1 && 2", 1.0, "&&", 2.0},
		{"3 || 4", 3.0, "||", 4.0},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		bin, ok := stmt.Expression.(*ast.BinaryExpression)
		if !ok {
			t.Fatalf("Expected BinaryExpression, got %T", stmt.Expression)
		}

		if bin.Operator != tt.operator {
			t.Errorf("Expected operator %q, got %q", tt.operator, bin.Operator)
		}

		checkNumberLiteral(t, bin.Left, tt.left)
		checkNumberLiteral(t, bin.Right, tt.right)
	}
}

func TestParseUnaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    float64
	}{
		{"-5", "-", 5.0},
		{"!10", "!", 10.0},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		unary, ok := stmt.Expression.(*ast.UnaryExpression)
		if !ok {
			t.Fatalf("Expected UnaryExpression, got %T", stmt.Expression)
		}

		if unary.Operator != tt.operator {
			t.Errorf("Expected operator %q, got %q", tt.operator, unary.Operator)
		}

		checkNumberLiteral(t, unary.Right, tt.value)
	}
}

func TestParseObjectLiteral(t *testing.T) {
	input := `{a: 1, "b": 2, c: 3 + 4}`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	obj, ok := stmt.Expression.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("Expected ObjectLiteral, got %T", stmt.Expression)
	}

	if len(obj.Properties) != 3 {
		t.Fatalf("Expected 3 properties, got %d", len(obj.Properties))
	}

	if val, ok := obj.Properties["a"]; !ok {
		t.Error("Missing property 'a'")
	} else {
		checkNumberLiteral(t, val, 1.0)
	}

	if val, ok := obj.Properties["b"]; !ok {
		t.Error("Missing property 'b'")
	} else {
		checkNumberLiteral(t, val, 2.0)
	}

	if val, ok := obj.Properties["c"]; !ok {
		t.Error("Missing property 'c'")
	} else {
		bin, ok := val.(*ast.BinaryExpression)
		if !ok {
			t.Fatalf("Expected BinaryExpression for property 'c', got %T", val)
		}
		if bin.Operator != "+" {
			t.Errorf("Expected operator '+', got %q", bin.Operator)
		}
		checkNumberLiteral(t, bin.Left, 3.0)
		checkNumberLiteral(t, bin.Right, 4.0)
	}
}

func TestParseArrayLiteral(t *testing.T) {
	input := `{1, 2, 3 + 4}`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Expected ArrayLiteral, got %T", stmt.Expression)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(arr.Elements))
	}

	checkNumberLiteral(t, arr.Elements[0], 1.0)
	checkNumberLiteral(t, arr.Elements[1], 2.0)

	bin, ok := arr.Elements[2].(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("Expected BinaryExpression for third element, got %T", arr.Elements[2])
	}
	if bin.Operator != "+" {
		t.Errorf("Expected operator '+', got %q", bin.Operator)
	}
	checkNumberLiteral(t, bin.Left, 3.0)
	checkNumberLiteral(t, bin.Right, 4.0)
}

func TestParseFunctionLiteral(t *testing.T) {
	input := `fun (a, b) { return a + b }`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected FunctionLiteral, got %T", stmt.Expression)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(fn.Parameters))
	}
	if fn.Parameters[0] != "a" {
		t.Errorf("Expected first parameter 'a', got %q", fn.Parameters[0])
	}
	if fn.Parameters[1] != "b" {
		t.Errorf("Expected second parameter 'b', got %q", fn.Parameters[1])
	}

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in body, got %d", len(fn.Body.Statements))
	}

	ret, ok := fn.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Expected ReturnStatement, got %T", fn.Body.Statements[0])
	}

	bin, ok := ret.Value.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("Expected BinaryExpression in return, got %T", ret.Value)
	}
	if bin.Operator != "+" {
		t.Errorf("Expected operator '+', got %q", bin.Operator)
	}

	leftIdent, ok := bin.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier on left, got %T", bin.Left)
	}
	if leftIdent.Name != "a" {
		t.Errorf("Expected identifier 'a', got %q", leftIdent.Name)
	}

	rightIdent, ok := bin.Right.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier on right, got %T", bin.Right)
	}
	if rightIdent.Name != "b" {
		t.Errorf("Expected identifier 'b', got %q", rightIdent.Name)
	}
}

func TestParseFunctionCall(t *testing.T) {
	input := `foo(1, 2)`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	call, ok := stmt.Expression.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("Expected FunctionCall, got %T", stmt.Expression)
	}

	ident, ok := call.Function.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as function, got %T", call.Function)
	}
	if ident.Name != "foo" {
		t.Errorf("Expected function name 'foo', got %q", ident.Name)
	}

	if len(call.Arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(call.Arguments))
	}
	checkNumberLiteral(t, call.Arguments[0], 1.0)
	checkNumberLiteral(t, call.Arguments[1], 2.0)
}

func TestParseMethodCall(t *testing.T) {
	input := `obj:foo(1, 2)`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	methodCall, ok := stmt.Expression.(*ast.MethodCall)
	if !ok {
		t.Fatalf("Expected MethodCall, got %T", stmt.Expression)
	}

	ident, ok := methodCall.Object.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as object, got %T", methodCall.Object)
	}
	if ident.Name != "obj" {
		t.Errorf("Expected object name 'obj', got %q", ident.Name)
	}

	if methodCall.Method != "foo" {
		t.Errorf("Expected method name 'foo', got %q", methodCall.Method)
	}

	if len(methodCall.Arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(methodCall.Arguments))
	}
	checkNumberLiteral(t, methodCall.Arguments[0], 1.0)
	checkNumberLiteral(t, methodCall.Arguments[1], 2.0)
}

func TestParseMemberAccess(t *testing.T) {
	tests := []struct {
		input  string
		object string
		member interface{}
		dot    bool
	}{
		{"obj.key", "obj", "key", true},
		{"obj[\"key\"]", "obj", "key", false},
		{"obj[key]", "obj", "key", false},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		if tt.dot {
			md, ok := stmt.Expression.(*ast.MemberDot)
			if !ok {
				t.Fatalf("Expected MemberDot, got %T", stmt.Expression)
			}

			ident, ok := md.Object.(*ast.Identifier)
			if !ok {
				t.Fatalf("Expected Identifier as object, got %T", md.Object)
			}
			if ident.Name != tt.object {
				t.Errorf("Expected object name %q, got %q", tt.object, ident.Name)
			}

			if md.Member != tt.member {
				t.Errorf("Expected member %q, got %q", tt.member, md.Member)
			}
		} else {
			mb, ok := stmt.Expression.(*ast.MemberBracket)
			if !ok {
				t.Fatalf("Expected MemberBracket, got %T", stmt.Expression)
			}

			ident, ok := mb.Object.(*ast.Identifier)
			if !ok {
				t.Fatalf("Expected Identifier as object, got %T", mb.Object)
			}
			if ident.Name != tt.object {
				t.Errorf("Expected object name %q, got %q", tt.object, ident.Name)
			}

			if tt.input == "obj[\"key\"]" {
				strLit, ok := mb.Member.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("Expected StringLiteral as member, got %T", mb.Member)
				}
				if strLit.Value != tt.member {
					t.Errorf("Expected member %q, got %q", tt.member, strLit.Value)
				}
			} else {
				memberIdent, ok := mb.Member.(*ast.Identifier)
				if !ok {
					t.Fatalf("Expected Identifier as member, got %T", mb.Member)
				}
				if memberIdent.Name != tt.member {
					t.Errorf("Expected member %q, got %q", tt.member, memberIdent.Name)
				}
			}
		}
	}
}

func TestParseIfStatement(t *testing.T) {
	input := `if x { return 1 } elif y { return 2 } else { return 3 }`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("Expected IfStatement, got %T", program.Statements[0])
	}

	condIdent, ok := ifStmt.Condition.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as condition, got %T", ifStmt.Condition)
	}
	if condIdent.Name != "x" {
		t.Errorf("Expected condition 'x', got %q", condIdent.Name)
	}

	if len(ifStmt.Then.Statements) != 1 {
		t.Fatalf("Expected 1 statement in then block, got %d", len(ifStmt.Then.Statements))
	}

	if len(ifStmt.Elif) != 1 {
		t.Fatalf("Expected 1 elif, got %d", len(ifStmt.Elif))
	}

	elifCondIdent, ok := ifStmt.Elif[0].Condition.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as elif condition, got %T", ifStmt.Elif[0].Condition)
	}
	if elifCondIdent.Name != "y" {
		t.Errorf("Expected elif condition 'y', got %q", elifCondIdent.Name)
	}

	if ifStmt.Else == nil {
		t.Fatal("Expected else block")
	}

	if len(ifStmt.Else.Statements) != 1 {
		t.Fatalf("Expected 1 statement in else block, got %d", len(ifStmt.Else.Statements))
	}
}

func TestParseWhileStatement(t *testing.T) {
	input := `while x { break }`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	whileStmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("Expected WhileStatement, got %T", program.Statements[0])
	}

	condIdent, ok := whileStmt.Condition.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as condition, got %T", whileStmt.Condition)
	}
	if condIdent.Name != "x" {
		t.Errorf("Expected condition 'x', got %q", condIdent.Name)
	}

	if len(whileStmt.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in body, got %d", len(whileStmt.Body.Statements))
	}

	_, ok = whileStmt.Body.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("Expected BreakStatement in body, got %T", whileStmt.Body.Statements[0])
	}
}

func TestParseForStatement(t *testing.T) {
	input := `for key, value in obj { continue }`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("Expected ForStatement, got %T", program.Statements[0])
	}

	if forStmt.Key != "key" {
		t.Errorf("Expected key 'key', got %q", forStmt.Key)
	}
	if forStmt.Value != "value" {
		t.Errorf("Expected value 'value', got %q", forStmt.Value)
	}

	objIdent, ok := forStmt.Object.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as object, got %T", forStmt.Object)
	}
	if objIdent.Name != "obj" {
		t.Errorf("Expected object 'obj', got %q", objIdent.Name)
	}

	if len(forStmt.Body.Statements) != 1 {
		t.Fatalf("Expected 1 statement in body, got %d", len(forStmt.Body.Statements))
	}

	_, ok = forStmt.Body.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("Expected ContinueStatement in body, got %T", forStmt.Body.Statements[0])
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := `return 42`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	retStmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Expected ReturnStatement, got %T", program.Statements[0])
	}

	if retStmt.Value == nil {
		t.Fatal("Expected return value")
	}

	checkNumberLiteral(t, retStmt.Value, 42.0)
}

func TestParseOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"1 + 2 + 3", "((1 + 2) + 3)"},
		{"1 * 2 * 3", "((1 * 2) * 3)"},
		{"1 + 2 - 3", "((1 + 2) - 3)"},
		{"1 < 2 == 3 > 4", "((1 < 2) == (3 > 4))"},
		{"1 && 2 || 3", "((1 && 2) || 3)"},
		{"!1 + 2", "((!1) + 2)"},
		{"-1 * 2", "((-1) * 2)"},
		{"1 + 2 * 3 == 4", "((1 + (2 * 3)) == 4)"},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		actual := stmt.Expression.String()
		if actual != tt.expected {
			t.Errorf("For input %q, expected %q, got %q", tt.input, tt.expected, actual)
		}
	}
}

func TestParseFibonacciExample(t *testing.T) {
	input := `fib = fun (n) {
    if n == 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    
    return fib(n-1) + fib(n-2)
}

fib(10)`

	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse Fibonacci example: %v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(program.Statements))
	}

	assignStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	assign, ok := assignStmt.Expression.(*ast.Assignment)
	if !ok {
		t.Fatalf("Expected Assignment, got %T", assignStmt.Expression)
	}

	ident, ok := assign.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier on left, got %T", assign.Left)
	}
	if ident.Name != "fib" {
		t.Errorf("Expected identifier 'fib', got %q", ident.Name)
	}

	_, ok = assign.Right.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected FunctionLiteral on right, got %T", assign.Right)
	}

	callStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[1])
	}

	call, ok := callStmt.Expression.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("Expected FunctionCall, got %T", callStmt.Expression)
	}

	callIdent, ok := call.Function.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as function, got %T", call.Function)
	}
	if callIdent.Name != "fib" {
		t.Errorf("Expected function name 'fib', got %q", callIdent.Name)
	}
}

func TestParseStringLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		raw      bool
	}{
		{`"hello"`, "hello", false},
		{`"hello\nworld"`, "hello\nworld", false},
		{`"he said \"hello\""`, "he said \"hello\"", false},
		{"`hello`", "hello", true},
		{"`hello\nworld`", "hello\nworld", true},
		{`"café"`, "café", false},
		{`"👍👎"`, "👍👎", false},
	}

	for _, tt := range tests {
		p := New()
		program, err := p.Parse(tt.input)
		if err != nil {
			t.Fatalf("Failed to parse %q: %v", tt.input, err)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
		}

		if tt.raw {
			rawStr, ok := stmt.Expression.(*ast.RawStringLiteral)
			if !ok {
				t.Fatalf("Expected RawStringLiteral, got %T", stmt.Expression)
			}
			if rawStr.Value != tt.expected {
				t.Errorf("Expected value %q, got %q", tt.expected, rawStr.Value)
			}
		} else {
			str, ok := stmt.Expression.(*ast.StringLiteral)
			if !ok {
				t.Fatalf("Expected StringLiteral, got %T", stmt.Expression)
			}
			if str.Value != tt.expected {
				t.Errorf("Expected value %q, got %q", tt.expected, str.Value)
			}
		}
	}
}

func TestParseComplexExpression(t *testing.T) {
	input := `obj.key:method(arg1, arg2)[index]`
	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse %q: %v", input, err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[0])
	}

	// The expression should be parsed as: ((obj.key):method(arg1, arg2))[index]
	memberBracket, ok := stmt.Expression.(*ast.MemberBracket)
	if !ok {
		t.Fatalf("Expected MemberBracket at top level, got %T", stmt.Expression)
	}

	methodCall, ok := memberBracket.Object.(*ast.MethodCall)
	if !ok {
		t.Fatalf("Expected MethodCall inside MemberBracket, got %T", memberBracket.Object)
	}

	memberDot, ok := methodCall.Object.(*ast.MemberDot)
	if !ok {
		t.Fatalf("Expected MemberDot inside MethodCall, got %T", methodCall.Object)
	}

	objIdent, ok := memberDot.Object.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier inside MemberDot, got %T", memberDot.Object)
	}
	if objIdent.Name != "obj" {
		t.Errorf("Expected object name 'obj', got %q", objIdent.Name)
	}

	if memberDot.Member != "key" {
		t.Errorf("Expected member 'key', got %q", memberDot.Member)
	}

	if methodCall.Method != "method" {
		t.Errorf("Expected method 'method', got %q", methodCall.Method)
	}

	if len(methodCall.Arguments) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(methodCall.Arguments))
	}

	arg1, ok := methodCall.Arguments[0].(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier for arg1, got %T", methodCall.Arguments[0])
	}
	if arg1.Name != "arg1" {
		t.Errorf("Expected argument 'arg1', got %q", arg1.Name)
	}

	arg2, ok := methodCall.Arguments[1].(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier for arg2, got %T", methodCall.Arguments[1])
	}
	if arg2.Name != "arg2" {
		t.Errorf("Expected argument 'arg2', got %q", arg2.Name)
	}

	indexIdent, ok := memberBracket.Member.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier as index, got %T", memberBracket.Member)
	}
	if indexIdent.Name != "index" {
		t.Errorf("Expected index 'index', got %q", indexIdent.Name)
	}
}

func TestParseInvalidSyntax(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
		errorMsg    string
	}{
		{"var x = ;", true, ""},     // Missing expression
		{"x =", true, ""},           // Missing right side
		{"if x {", true, ""},        // Missing closing brace
		{"while x", true, ""},       // Missing block
		{"for a, b in", true, ""},   // Missing expression
		{"return", false, ""},       // Valid (return without value)
		{"1 +", true, ""},           // Missing right operand
		{"fun () {", true, ""},      // Missing closing brace
		{"{a:}", true, ""},          // Missing value in property
		{"{1,}", true, ""},          // Missing element in array
		{"obj.", true, ""},          // Missing member after dot
		{"obj[", true, ""},          // Missing member and closing bracket
		{"obj:", true, ""},          // Missing method name
		{"obj:method(", true, ""},   // Missing closing paren
		{"obj:method()", false, ""}, // Valid
		{"(1 + 2", true, ""},        // Missing closing paren
	}

	for _, tt := range tests {
		p := New()
		_, err := p.Parse(tt.input)
		if tt.expectError && err == nil {
			t.Errorf("Expected error for input %q, but got none", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Did not expect error for input %q, but got: %v", tt.input, err)
		}
	}
}

func TestParseComments(t *testing.T) {
	input := `# This is a comment
var x = 1  # Another comment
# Yet another comment
x = 2`

	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse code with comments: %v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(program.Statements))
	}

	vd, ok := program.Statements[0].(*ast.VarDeclaration)
	if !ok {
		t.Fatalf("Expected VarDeclaration, got %T", program.Statements[0])
	}
	if vd.Name != "x" {
		t.Errorf("Expected variable name 'x', got %q", vd.Name)
	}

	stmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got %T", program.Statements[1])
	}

	assign, ok := stmt.Expression.(*ast.Assignment)
	if !ok {
		t.Fatalf("Expected Assignment, got %T", stmt.Expression)
	}

	ident, ok := assign.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier on left, got %T", assign.Left)
	}
	if ident.Name != "x" {
		t.Errorf("Expected identifier 'x', got %q", ident.Name)
	}
}

func checkNumberLiteral(t *testing.T, expr ast.Expression, expected float64) {
	t.Helper()
	num, ok := expr.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("Expected NumberLiteral, got %T", expr)
	}
	if num.Value != expected {
		t.Errorf("Expected value %f, got %f", expected, num.Value)
	}
}

func checkStringLiteral(t *testing.T, expr ast.Expression, expected string) {
	t.Helper()
	str, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("Expected StringLiteral, got %T", expr)
	}
	if str.Value != expected {
		t.Errorf("Expected value %q, got %q", expected, str.Value)
	}
}

func TestParseExamplesFromDocs(t *testing.T) {
	examples := []string{
		// Counter example
		`var Counter = {
    new: fun(self, initial) {
        var inst = clone(self)
        inst.value = initial
        return inst
    },
    inc: fun(self) {
        self.value = self.value + 1
        return self.value
    },
    get: fun(self) {
        return self.value
    }
}

var c = Counter:new(5)
c:inc()
c:get()`,

		// Map reduction example
		`map = fun (arr, fn) {
    var result = {}
    for i, val in arr {
        var idx = number(i)
        result[i] = fn(val)
    }
    return result
}

var arr = {1, 2, 3, 4}
var doubled = map(arr, fun(x) { x * 2 })
doubled[2]`,
	}

	for i, example := range examples {
		p := New()
		program, err := p.Parse(example)
		if err != nil {
			t.Fatalf("Failed to parse example %d: %v", i+1, err)
		}
		if program == nil {
			t.Fatalf("Program is nil for example %d", i+1)
		}
		// Just verify it parses without errors
		if len(p.Errors()) > 0 {
			t.Errorf("Unexpected errors parsing example %d: %v", i+1, p.Errors())
		}
	}
}

func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *ast.Program)
	}{
		{
			name:  "empty object",
			input: "var x = {}",
			check: func(t *testing.T, p *ast.Program) {
				vd := p.Statements[0].(*ast.VarDeclaration)
				obj, ok := vd.Value.(*ast.ObjectLiteral)
				if !ok {
					t.Fatalf("Expected ObjectLiteral, got %T", vd.Value)
				}
				if len(obj.Properties) != 0 {
					t.Errorf("Expected empty object, got %d properties", len(obj.Properties))
				}
			},
		},
		{
			name:  "empty array",
			input: "var x = {}",
			check: func(t *testing.T, p *ast.Program) {
				vd := p.Statements[0].(*ast.VarDeclaration)
				_, ok := vd.Value.(*ast.ObjectLiteral)
				if !ok {
					t.Fatalf("Expected ObjectLiteral, got %T", vd.Value)
				}
				// Note: Empty {} is ambiguous between object and array
				// The grammar parses it as object literal
			},
		},
		{
			name:  "function with no parameters",
			input: "var f = fun () { return 1 }",
			check: func(t *testing.T, p *ast.Program) {
				vd := p.Statements[0].(*ast.VarDeclaration)
				fn, ok := vd.Value.(*ast.FunctionLiteral)
				if !ok {
					t.Fatalf("Expected FunctionLiteral, got %T", vd.Value)
				}
				if len(fn.Parameters) != 0 {
					t.Errorf("Expected 0 parameters, got %d", len(fn.Parameters))
				}
			},
		},
		{
			name:  "nested member access",
			input: "a.b.c.d",
			check: func(t *testing.T, p *ast.Program) {
				stmt := p.Statements[0].(*ast.ExpressionStatement)
				// Should be parsed as ((a.b).c).d
				md1, ok := stmt.Expression.(*ast.MemberDot)
				if !ok {
					t.Fatalf("Expected MemberDot, got %T", stmt.Expression)
				}
				if md1.Member != "d" {
					t.Errorf("Expected member 'd', got %q", md1.Member)
				}
				md2, ok := md1.Object.(*ast.MemberDot)
				if !ok {
					t.Fatalf("Expected nested MemberDot, got %T", md1.Object)
				}
				if md2.Member != "c" {
					t.Errorf("Expected member 'c', got %q", md2.Member)
				}
			},
		},
		{
			name:  "chained method calls",
			input: "obj:foo():bar()",
			check: func(t *testing.T, p *ast.Program) {
				stmt := p.Statements[0].(*ast.ExpressionStatement)
				// Should be parsed as (obj:foo()):bar()
				mc1, ok := stmt.Expression.(*ast.MethodCall)
				if !ok {
					t.Fatalf("Expected MethodCall, got %T", stmt.Expression)
				}
				if mc1.Method != "bar" {
					t.Errorf("Expected method 'bar', got %q", mc1.Method)
				}
				mc2, ok := mc1.Object.(*ast.MethodCall)
				if !ok {
					t.Fatalf("Expected nested MethodCall, got %T", mc1.Object)
				}
				if mc2.Method != "foo" {
					t.Errorf("Expected method 'foo', got %q", mc2.Method)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			program, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse %q: %v", tt.input, err)
			}
			tt.check(t, program)
		})
	}
}

func TestParsePositionTracking(t *testing.T) {
	input := `var x = 1
var y = 2
if x {
    return y
}`

	p := New()
	program, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Check positions are set
	for i, stmt := range program.Statements {
		pos := stmt.Position()
		if pos.Line <= 0 {
			t.Errorf("Statement %d has invalid line number: %d", i, pos.Line)
		}
		if pos.Column < 0 {
			t.Errorf("Statement %d has invalid column: %d", i, pos.Column)
		}
	}

	// Specifically check the if statement
	ifStmt, ok := program.Statements[2].(*ast.IfStatement)
	if !ok {
		t.Fatalf("Expected IfStatement, got %T", program.Statements[2])
	}

	if ifStmt.Pos.Line != 3 {
		t.Errorf("Expected if statement at line 3, got %d", ifStmt.Pos.Line)
	}
}
