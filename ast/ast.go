package ast

import (
	"fmt"
	"strings"
)

// Node represents a node in the AST
type Node interface {
	Position() Position
}

// Position represents a position in the source code
type Position struct {
	Line   int
	Column int
}

// Program represents the entire program
type Program struct {
	Pos        Position
	Statements []Statement
}

func (p *Program) Position() Position { return p.Pos }
func (p *Program) String() string {
	var stmts []string
	for _, stmt := range p.Statements {
		stmts = append(stmts, stmt.String())
	}
	return strings.Join(stmts, "\n")
}

// Statement represents a statement
type Statement interface {
	Node
	statementNode()
	String() string
}

// Expression represents an expression
type Expression interface {
	Node
	expressionNode()
	String() string
}

// VarDeclaration represents a variable declaration
type VarDeclaration struct {
	Pos   Position
	Name  string
	Value Expression
}

func (vd *VarDeclaration) Position() Position { return vd.Pos }
func (vd *VarDeclaration) statementNode()     {}
func (vd *VarDeclaration) String() string {
	if vd.Value != nil {
		return fmt.Sprintf("var %s = %v", vd.Name, vd.Value)
	}
	return fmt.Sprintf("var %s", vd.Name)
}

// ExpressionStatement represents an expression as a statement
type ExpressionStatement struct {
	Pos        Position
	Expression Expression
}

func (es *ExpressionStatement) Position() Position { return es.Pos }
func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) String() string     { return es.Expression.String() }

// IfStatement represents an if statement
type IfStatement struct {
	Pos       Position
	Condition Expression
	Then      *BlockStatement
	Elif      []*ElifStatement
	Else      *BlockStatement
}

func (is *IfStatement) Position() Position { return is.Pos }
func (is *IfStatement) statementNode()     {}
func (is *IfStatement) String() string {
	result := fmt.Sprintf("if %v %v", is.Condition, is.Then)
	for _, elif := range is.Elif {
		result += fmt.Sprintf(" %v", elif)
	}
	if is.Else != nil {
		result += fmt.Sprintf(" else %v", is.Else)
	}
	return result
}

type ElifStatement struct {
	Pos       Position
	Condition Expression
	Body      *BlockStatement
}

func (es *ElifStatement) String() string {
	return fmt.Sprintf("elif %v %v", es.Condition, es.Body)
}

// WhileStatement represents a while statement
type WhileStatement struct {
	Pos       Position
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) Position() Position { return ws.Pos }
func (ws *WhileStatement) statementNode()     {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %v %v", ws.Condition, ws.Body)
}

// ForStatement represents a for statement
type ForStatement struct {
	Pos    Position
	Key    string
	Value  string
	Object Expression
	Body   *BlockStatement
}

func (fs *ForStatement) Position() Position { return fs.Pos }
func (fs *ForStatement) statementNode()     {}
func (fs *ForStatement) String() string {
	if fs.Value != "" {
		return fmt.Sprintf("for %s, %s in %v %v", fs.Key, fs.Value, fs.Object, fs.Body)
	}
	return fmt.Sprintf("for %s in %v %v", fs.Key, fs.Object, fs.Body)
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Pos        Position
	Statements []Statement
}

func (bs *BlockStatement) Position() Position { return bs.Pos }
func (bs *BlockStatement) statementNode()     {}
func (bs *BlockStatement) String() string {
	var stmts []string
	for _, stmt := range bs.Statements {
		stmts = append(stmts, stmt.String())
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(stmts, "\n"))
}

// BreakStatement represents a break statement
type BreakStatement struct {
	Pos Position
}

func (bs *BreakStatement) Position() Position { return bs.Pos }
func (bs *BreakStatement) statementNode()     {}
func (bs *BreakStatement) String() string     { return "break" }

// ContinueStatement represents a continue statement
type ContinueStatement struct {
	Pos Position
}

func (cs *ContinueStatement) Position() Position { return cs.Pos }
func (cs *ContinueStatement) statementNode()     {}
func (cs *ContinueStatement) String() string     { return "continue" }

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Pos   Position
	Value Expression
}

func (rs *ReturnStatement) Position() Position { return rs.Pos }
func (rs *ReturnStatement) statementNode()     {}
func (rs *ReturnStatement) String() string {
	if rs.Value != nil {
		return fmt.Sprintf("return %v", rs.Value)
	}
	return "return"
}

// Identifier represents an identifier
type Identifier struct {
	Pos  Position
	Name string
}

func (i *Identifier) Position() Position { return i.Pos }
func (i *Identifier) expressionNode()    {}
func (i *Identifier) String() string     { return i.Name }

// NumberLiteral represents a number literal
type NumberLiteral struct {
	Pos   Position
	Value float64
}

func (nl *NumberLiteral) Position() Position { return nl.Pos }
func (nl *NumberLiteral) expressionNode()    {}
func (nl *NumberLiteral) String() string     { return fmt.Sprintf("%g", nl.Value) }

// StringLiteral represents a string literal
type StringLiteral struct {
	Pos   Position
	Value string
}

func (sl *StringLiteral) Position() Position { return sl.Pos }
func (sl *StringLiteral) expressionNode()    {}
func (sl *StringLiteral) String() string     { return fmt.Sprintf("%q", sl.Value) }

// RawStringLiteral represents a raw string literal
type RawStringLiteral struct {
	Pos   Position
	Value string
}

func (rsl *RawStringLiteral) Position() Position { return rsl.Pos }
func (rsl *RawStringLiteral) expressionNode()    {}
func (rsl *RawStringLiteral) String() string     { return fmt.Sprintf("`%s`", rsl.Value) }

// NilLiteral represents a nil literal
type NilLiteral struct {
	Pos Position
}

func (nl *NilLiteral) Position() Position { return nl.Pos }
func (nl *NilLiteral) expressionNode()    {}
func (nl *NilLiteral) String() string     { return "nil" }

// ObjectLiteral represents an object literal
type ObjectLiteral struct {
	Pos        Position
	Properties map[string]Expression
}

func (ol *ObjectLiteral) Position() Position { return ol.Pos }
func (ol *ObjectLiteral) expressionNode()    {}
func (ol *ObjectLiteral) String() string {
	if len(ol.Properties) == 0 {
		return "{}"
	}
	var props []string
	for key, value := range ol.Properties {
		props = append(props, fmt.Sprintf("%s: %v", key, value))
	}
	return fmt.Sprintf("{ %s }", strings.Join(props, ", "))
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	Pos      Position
	Elements []Expression
}

func (al *ArrayLiteral) Position() Position { return al.Pos }
func (al *ArrayLiteral) expressionNode()    {}
func (al *ArrayLiteral) String() string {
	var elems []string
	for _, elem := range al.Elements {
		elems = append(elems, elem.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

// FunctionLiteral represents a function literal
type FunctionLiteral struct {
	Pos        Position
	Parameters []string
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Position() Position { return fl.Pos }
func (fl *FunctionLiteral) expressionNode()    {}
func (fl *FunctionLiteral) String() string {
	params := strings.Join(fl.Parameters, ", ")
	return fmt.Sprintf("fun(%s) %v", params, fl.Body)
}

// MemberDot represents a member access using dot notation
type MemberDot struct {
	Pos    Position
	Object Expression
	Member string
}

func (md *MemberDot) Position() Position { return md.Pos }
func (md *MemberDot) expressionNode()    {}
func (md *MemberDot) String() string     { return fmt.Sprintf("%v.%s", md.Object, md.Member) }

// MemberBracket represents a member access using bracket notation
type MemberBracket struct {
	Pos    Position
	Object Expression
	Member Expression
}

func (mb *MemberBracket) Position() Position { return mb.Pos }
func (mb *MemberBracket) expressionNode()    {}
func (mb *MemberBracket) String() string     { return fmt.Sprintf("%v[%v]", mb.Object, mb.Member) }

// FunctionCall represents a function call
type FunctionCall struct {
	Pos       Position
	Function  Expression
	Arguments []Expression
}

func (fc *FunctionCall) Position() Position { return fc.Pos }
func (fc *FunctionCall) expressionNode()    {}
func (fc *FunctionCall) String() string {
	return fmt.Sprintf("%v(%s)", fc.Function, fc.FormatArgs())
}
func (fc *FunctionCall) FormatArgs() string {
	var args []string
	for _, arg := range fc.Arguments {
		args = append(args, arg.String())
	}
	return strings.Join(args, ", ")
}

// MethodCall represents a method call
type MethodCall struct {
	Pos       Position
	Object    Expression
	Method    string
	Arguments []Expression
}

func (mc *MethodCall) Position() Position { return mc.Pos }
func (mc *MethodCall) expressionNode()    {}
func (mc *MethodCall) String() string {
	return fmt.Sprintf("%v:%s(%s)", mc.Object, mc.Method, mc.FormatArgs())
}
func (mc *MethodCall) FormatArgs() string {
	var args []string
	for _, arg := range mc.Arguments {
		args = append(args, arg.String())
	}
	return strings.Join(args, ", ")
}

// UnaryExpression represents a unary expression
type UnaryExpression struct {
	Pos      Position
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) Position() Position { return ue.Pos }
func (ue *UnaryExpression) expressionNode()    {}
func (ue *UnaryExpression) String() string     { return fmt.Sprintf("(%s%v)", ue.Operator, ue.Right) }

// BinaryExpression represents a binary expression
type BinaryExpression struct {
	Pos      Position
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) Position() Position { return be.Pos }
func (be *BinaryExpression) expressionNode()    {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%v %s %v)", be.Left, be.Operator, be.Right)
}

// Assignment represents an assignment expression
type Assignment struct {
	Pos   Position
	Left  Expression
	Right Expression
}

func (a *Assignment) Position() Position { return a.Pos }
func (a *Assignment) expressionNode()    {}
func (a *Assignment) String() string     { return fmt.Sprintf("%v = %v", a.Left, a.Right) }
