package parser

import (
	"errors"
	"fmt"
	"refl/ast"
	"strconv"

	"refl/parser/gen"

	"github.com/antlr4-go/antlr/v4"
)

type Parser struct {
	errors []error
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(code string) (resultP *ast.Program, resultErr error) {
	defer func() {
		if err := recover(); err != nil {
			resultErr = fmt.Errorf("panic: %v, errors: %w", err, errors.Join(p.errors...))
		}
	}()

	p.errors = []error{}

	input := antlr.NewInputStream(code)
	lexer := gen.NewReflLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	antlrParser := gen.NewReflParser(stream)
	antlrParser.RemoveErrorListeners()
	antlrParser.AddErrorListener(p)

	tree := antlrParser.Program()

	if len(p.errors) > 0 {
		return nil, p.errors[0]
	}

	visitor := newReflVisitor(p)
	return visitor.Visit(tree).(*ast.Program), nil
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, errors.New(msg))
}

func (p *Parser) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
	p.error(fmt.Sprintf("%s at line %d, column %d, %s", msg, line, column, offendingSymbol))
}

func (p *Parser) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (p *Parser) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (p *Parser) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}

var _ gen.ReflVisitor = (*ReflVisitor)(nil)

type ReflVisitor struct {
	gen.BaseReflVisitor
	parser *Parser
}

func newReflVisitor(p *Parser) *ReflVisitor {
	return &ReflVisitor{
		parser: p,
	}
}

func (v *ReflVisitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

func (v *ReflVisitor) VisitPrimaryExpr(ctx *gen.PrimaryExprContext) any {
	return ctx.Primary().Accept(v)
}

func (v *ReflVisitor) VisitProgram(ctx *gen.ProgramContext) any {
	program := &ast.Program{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}

	for _, stmtCtx := range ctx.AllStatement() {
		acceptRes := stmtCtx.Accept(v)
		if acceptRes != nil {
			program.Statements = append(program.Statements, acceptRes.(ast.Statement))
		}
	}

	return program
}

func (v *ReflVisitor) VisitVarDeclaration(ctx *gen.VarDeclarationContext) any {
	vd := &ast.VarDeclaration{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Name: ctx.IDENTIFIER().GetText(),
	}

	if ctx.Expression() != nil {
		vd.Value = ctx.Expression().Accept(v).(ast.Expression)
	}

	return vd
}

func (v *ReflVisitor) VisitExpressionStatement(ctx *gen.ExpressionStatementContext) any {
	es := &ast.ExpressionStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}

	if ctx.Expression() != nil {
		es.Expression = ctx.Expression().Accept(v).(ast.Expression)
	}

	return es
}

func (v *ReflVisitor) VisitIfStatement(ctx *gen.IfStatementContext) any {
	is := &ast.IfStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Condition: ctx.Expression(0).Accept(v).(ast.Expression),
		Then:      ctx.Block(0).Accept(v).(*ast.BlockStatement),
	}

	elifs := ctx.AllELIF()

	for i := 0; i < len(elifs); i++ {
		el := &ast.ElifStatement{
			Pos: ast.Position{
				Line:   elifs[i].GetSymbol().GetLine(),
				Column: elifs[i].GetSymbol().GetColumn(),
			},
			Condition: ctx.Expression(i + 1).Accept(v).(ast.Expression),
			Body:      ctx.Block(i + 1).Accept(v).(*ast.BlockStatement),
		}
		is.Elif = append(is.Elif, el)
	}

	if ctx.ELSE() != nil {
		is.Else = ctx.Block(len(ctx.AllBlock()) - 1).Accept(v).(*ast.BlockStatement)
	}

	return is
}

func (v *ReflVisitor) VisitWhileStatement(ctx *gen.WhileStatementContext) any {
	ws := &ast.WhileStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Condition: ctx.Expression().Accept(v).(ast.Expression),
		Body:      ctx.Block().Accept(v).(*ast.BlockStatement),
	}

	return ws
}

func (v *ReflVisitor) VisitForStatement(ctx *gen.ForStatementContext) any {
	fs := &ast.ForStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Key:    ctx.IDENTIFIER(0).GetText(),
		Value:  ctx.IDENTIFIER(1).GetText(),
		Object: ctx.Expression().Accept(v).(ast.Expression),
		Body:   ctx.Block().Accept(v).(*ast.BlockStatement),
	}

	return fs
}

func (v *ReflVisitor) VisitBlock(ctx *gen.BlockContext) any {
	bs := &ast.BlockStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}

	for _, stmtCtx := range ctx.AllStatement() {
		acceptRes := stmtCtx.Accept(v)
		if acceptRes != nil {
			bs.Statements = append(bs.Statements, acceptRes.(ast.Statement))
		}
	}

	return bs
}

func (v *ReflVisitor) VisitBlockStatement(ctx *gen.BlockStatementContext) any {
	return ctx.Block().Accept(v)
}

func (v *ReflVisitor) VisitBreakStatement(ctx *gen.BreakStatementContext) any {
	return &ast.BreakStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}
}

func (v *ReflVisitor) VisitContinueStatement(ctx *gen.ContinueStatementContext) any {
	return &ast.ContinueStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}
}

func (v *ReflVisitor) VisitReturnStatement(ctx *gen.ReturnStatementContext) any {
	rs := &ast.ReturnStatement{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}

	if ctx.Expression() != nil {
		rs.Value = ctx.Expression().Accept(v).(ast.Expression)
	}

	return rs
}

func (v *ReflVisitor) VisitMemberDot(ctx *gen.MemberDotContext) any {
	md := &ast.MemberDot{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Object: ctx.Expression().Accept(v).(ast.Expression),
		Member: ctx.IDENTIFIER().GetText(),
	}

	return md
}

func (v *ReflVisitor) VisitMethodCall(ctx *gen.MethodCallContext) any {
	mc := &ast.MethodCall{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Object: ctx.Expression().Accept(v).(ast.Expression),
		Method: ctx.IDENTIFIER().GetText(),
	}

	if ctx.ExpressionList() != nil {
		exprList := ctx.ExpressionList().Accept(v).([]ast.Expression)
		mc.Arguments = append(mc.Arguments, exprList...)
	}

	return mc
}

func (v *ReflVisitor) VisitFunctionCall(ctx *gen.FunctionCallContext) any {
	fc := &ast.FunctionCall{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Function: ctx.Expression().Accept(v).(ast.Expression),
	}

	if ctx.ExpressionList() != nil {
		exprList := ctx.ExpressionList().Accept(v).([]ast.Expression)
		fc.Arguments = append(fc.Arguments, exprList...)
	}

	return fc
}

func (v *ReflVisitor) VisitMemberBracket(ctx *gen.MemberBracketContext) any {
	mb := &ast.MemberBracket{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Object: ctx.Expression(0).Accept(v).(ast.Expression),
		Member: ctx.Expression(1).Accept(v).(ast.Expression),
	}

	return mb
}

func (v *ReflVisitor) VisitUnary(ctx *gen.UnaryContext) any {
	ue := &ast.UnaryExpression{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Operator: ctx.GetOp().GetText(),
		Right:    ctx.Expression().Accept(v).(ast.Expression),
	}

	return ue
}

func (v *ReflVisitor) VisitBinary(ctx *gen.BinaryContext) any {
	be := &ast.BinaryExpression{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Left:     ctx.Expression(0).Accept(v).(ast.Expression),
		Operator: ctx.GetOp().GetText(),
		Right:    ctx.Expression(1).Accept(v).(ast.Expression),
	}

	return be
}

func (v *ReflVisitor) VisitAssignment(ctx *gen.AssignmentContext) any {
	a := &ast.Assignment{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Left:  ctx.Expression(0).Accept(v).(ast.Expression),
		Right: ctx.Expression(1).Accept(v).(ast.Expression),
	}

	return a
}

func (v *ReflVisitor) VisitExpressionList(ctx *gen.ExpressionListContext) any {
	var expressions []ast.Expression

	for _, exprCtx := range ctx.AllExpression() {
		expr := exprCtx.Accept(v).(ast.Expression)
		expressions = append(expressions, expr)
	}

	return expressions
}

func (v *ReflVisitor) VisitLiteralPrimary(ctx *gen.LiteralPrimaryContext) any {
	return ctx.Literal().Accept(v)
}

func (v *ReflVisitor) VisitIdentifierPrimary(ctx *gen.IdentifierPrimaryContext) any {
	return &ast.Identifier{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Name: ctx.IDENTIFIER().GetText(),
	}
}

func (v *ReflVisitor) VisitParenPrimary(ctx *gen.ParenPrimaryContext) any {
	return ctx.Expression().Accept(v)
}

func (v *ReflVisitor) VisitFunctionLiteral(ctx *gen.FunctionLiteralContext) any {
	fl := &ast.FunctionLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Body: ctx.Block().Accept(v).(*ast.BlockStatement),
	}

	if ctx.Parameters() != nil {
		paramsCtx := ctx.Parameters().(*gen.ParametersContext)
		for _, id := range paramsCtx.AllIDENTIFIER() {
			fl.Parameters = append(fl.Parameters, id.GetText())
		}
	}

	return fl
}

func (v *ReflVisitor) VisitObjectLiteralPrimary(ctx *gen.ObjectLiteralPrimaryContext) any {
	return ctx.ObjectLiteral().Accept(v)
}

func (v *ReflVisitor) VisitArrayLiteralPrimary(ctx *gen.ArrayLiteralPrimaryContext) any {
	return ctx.ArrayLiteral().Accept(v)
}

func (v *ReflVisitor) VisitParameters(ctx *gen.ParametersContext) any {
	var params []string
	for _, id := range ctx.AllIDENTIFIER() {
		params = append(params, id.GetText())
	}
	return params
}

func (v *ReflVisitor) VisitObjectLiteral(ctx *gen.ObjectLiteralContext) any {
	ol := &ast.ObjectLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Properties: make(map[string]ast.Expression),
	}

	if ctx.AllProperty() != nil {
		for _, propCtx := range ctx.AllProperty() {
			prop := propCtx.(*gen.PropertyContext)
			var key string
			if prop.STRING() != nil {
				key = parseString(prop.STRING().GetText())
			} else {
				key = prop.IDENTIFIER().GetText()
			}
			value := prop.Expression().Accept(v).(ast.Expression)
			ol.Properties[key] = value
		}
	}

	return ol
}

func (v *ReflVisitor) VisitProperty(ctx *gen.PropertyContext) any {
	return ctx
}

func (v *ReflVisitor) VisitArrayLiteral(ctx *gen.ArrayLiteralContext) any {
	al := &ast.ArrayLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}

	for _, exprCtx := range ctx.AllExpression() {
		expr := exprCtx.Accept(v).(ast.Expression)
		al.Elements = append(al.Elements, expr)
	}

	return al
}

func (v *ReflVisitor) VisitNumberLiteral(ctx *gen.NumberLiteralContext) any {
	value, err := parseNumber(ctx.NUMBER().GetText())
	if err != nil {
		v.parser.error(err.Error())
		return &ast.NumberLiteral{Value: 0}
	}

	return &ast.NumberLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Value: value,
	}
}

func (v *ReflVisitor) VisitStringLiteral(ctx *gen.StringLiteralContext) any {
	return &ast.StringLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Value: parseString(ctx.STRING().GetText()),
	}
}

func (v *ReflVisitor) VisitRawStringLiteral(ctx *gen.RawStringLiteralContext) any {
	return &ast.RawStringLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
		Value: parseRawString(ctx.RAW_STRING().GetText()),
	}
}

func (v *ReflVisitor) VisitNilLiteral(ctx *gen.NilLiteralContext) any {
	return &ast.NilLiteral{
		Pos: ast.Position{
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		},
	}
}

func (v *ReflVisitor) VisitStatement(ctx *gen.StatementContext) any {
	if ctx.VarDeclaration() != nil {
		return ctx.VarDeclaration().Accept(v)
	}
	if ctx.ExpressionStatement() != nil {
		return ctx.ExpressionStatement().Accept(v)
	}
	if ctx.IfStatement() != nil {
		return ctx.IfStatement().Accept(v)
	}
	if ctx.WhileStatement() != nil {
		return ctx.WhileStatement().Accept(v)
	}
	if ctx.ForStatement() != nil {
		return ctx.ForStatement().Accept(v)
	}
	if ctx.BlockStatement() != nil {
		return ctx.BlockStatement().Accept(v)
	}
	if ctx.BreakStatement() != nil {
		return ctx.BreakStatement().Accept(v)
	}
	if ctx.ContinueStatement() != nil {
		return ctx.ContinueStatement().Accept(v)
	}
	if ctx.ReturnStatement() != nil {
		return ctx.ReturnStatement().Accept(v)
	}
	return nil
}

func parseNumber(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseString(s string) string {
	s = s[1 : len(s)-1]
	result := ""
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '\\':
				result += "\\"
				i += 2
			case '"':
				result += "\""
				i += 2
			case 'n':
				result += "\n"
				i += 2
			case 'r':
				result += "\r"
				i += 2
			case 't':
				result += "\t"
				i += 2
			default:
				result += string(s[i])
				i++
			}
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

func parseRawString(s string) string {
	return s[1 : len(s)-1]
}
