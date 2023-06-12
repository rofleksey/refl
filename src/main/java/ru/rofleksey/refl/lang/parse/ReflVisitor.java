package ru.rofleksey.refl.lang.parse;

import org.apache.commons.text.StringEscapeUtils;
import ru.rofleksey.refl.ReflParser;
import ru.rofleksey.refl.ReflParserBaseVisitor;
import ru.rofleksey.refl.lang.node.*;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.BinaryOperator;
import ru.rofleksey.refl.lang.operator.CompareOperator;
import ru.rofleksey.refl.lang.operator.assign.*;
import ru.rofleksey.refl.lang.operator.binary.*;
import ru.rofleksey.refl.lang.operator.compare.*;
import ru.rofleksey.refl.lang.value.NilValue;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.StringValue;

import java.util.ArrayList;
import java.util.HashMap;

public class ReflVisitor extends ReflParserBaseVisitor<Node> {
    @Override
    public Node visitRoot(ReflParser.RootContext ctx) {
        return visit(ctx.declarationList());
    }

    @Override
    public Node visitDeclarationList(ReflParser.DeclarationListContext ctx) {
        var declarations = ctx.declaration();
        var nodes = new ArrayList<Node>(declarations.size());
        for (var decl : declarations) {
            nodes.add(visit(decl));
        }
        return new ListNode(nodes);
    }

    @Override
    public Node visitIfExpression(ReflParser.IfExpressionContext ctx) {
        var condition = visit(ctx.mainCondition);
        var mainBody = visit(ctx.mainBody);
        var conditions = ctx.simpleExpression();
        var declarationLists = ctx.declarationList();
        var elifNodes = new ArrayList<IfNode>(conditions.size() + 1);
        for (var i = 0; i < conditions.size(); i++) {
            var elifCond = visit(conditions.get(i));
            var elifBody = visit(ctx.declarationList(i));
            elifNodes.add(new IfNode(elifCond, elifBody, null));
        }
        if (declarationLists.size() > conditions.size()) {
            var elseBody = visit(declarationLists.get(declarationLists.size() - 1));
            elifNodes.add(new IfNode(null, elseBody, null));
        }
        return new IfNode(condition, mainBody, elifNodes);
    }

    @Override
    public Node visitWhileExpression(ReflParser.WhileExpressionContext ctx) {
        var condition = visit(ctx.condition);
        var body = visit(ctx.declarationList());
        return new WhileNode(condition, body);
    }

    @Override
    public Node visitNormalFunction(ReflParser.NormalFunctionContext ctx) {
        var name = ctx.name.getText();
        var body = visit(ctx.declarationList());
        return new FunctionDeclarationNode(new StringValue(name), body);
    }

    @Override
    public Node visitScope(ReflParser.ScopeContext ctx) {
        var name = ctx.name.getText();
        var body = visit(ctx.declarationList());
        return new ScopeDeclarationNode(new StringValue(name), body);
    }

    @Override
    public Node visitCompare(ReflParser.CompareContext ctx) {
        var opText = ctx.bop.getText();
        CompareOperator operator;
        switch (opText) {
            case "<":
                operator = LessOperator.INSTANCE;
                break;
            case "<=":
                operator = LteOperator.INSTANCE;
                break;
            case "==":
                operator = EqualsOperator.INSTANCE;
                break;
            case "!=":
                operator = NotEqualsOperator.INSTANCE;
                break;
            case ">":
                operator = GreaterOperator.INSTANCE;
                break;
            case ">=":
                operator = GteOperator.INSTANCE;
                break;
            default:
                throw new IllegalArgumentException("unknown comparison operator " + opText);
        }
        return new CompareNode(operator, visit(ctx.left), visit(ctx.right));
    }

    @Override
    public Node visitAddSub(ReflParser.AddSubContext ctx) {
        var opText = ctx.bop.getText();
        BinaryOperator operator;
        switch (opText) {
            case "+":
                operator = AddOperator.INSTANCE;
                break;
            case "-":
                operator = SubOperator.INSTANCE;
                break;
            default:
                throw new IllegalArgumentException("unknown add/sub operator " + opText);
        }
        return new BinaryNode(visit(ctx.left), visit(ctx.right), operator);
    }

    @Override
    public Node visitMulDiv(ReflParser.MulDivContext ctx) {
        var opText = ctx.bop.getText();
        BinaryOperator operator;
        switch (opText) {
            case "*":
                operator = MulOperator.INSTANCE;
                break;
            case "/":
                operator = DivOperator.INSTANCE;
                break;
            case "%":
                operator = ModOperator.INSTANCE;
                break;
            default:
                throw new IllegalArgumentException("unknown mul/div operator " + opText);
        }
        return new BinaryNode(visit(ctx.left), visit(ctx.right), operator);
    }

    @Override
    public Node visitOr(ReflParser.OrContext ctx) {
        var opText = ctx.bop.getText();
        BinaryOperator operator;
        switch (opText) {
            case "|":
                operator = OrOperator.INSTANCE;
                break;
            case "??":
                operator = ElvisOperator.INSTANCE;
                break;
            default:
                throw new IllegalArgumentException("unknown binary operator " + opText);
        }
        return new BinaryNode(visit(ctx.left), visit(ctx.right), operator);
    }

    @Override
    public Node visitAnd(ReflParser.AndContext ctx) {
        return new BinaryNode(visit(ctx.left), visit(ctx.right), AndOperator.INSTANCE);
    }

    private MethodCallParseCtx visitMethodCallCustom(ReflParser.MethodCallContext ctx) {
        var methodName = ctx.name.getText();
        var argsListCtx = ctx.argument();
        var posArgs = new ArrayList<Node>();
        var namedArgs = new HashMap<String, Node>();
        for (var argCtx : argsListCtx) {
            var value = visit(argCtx.simpleExpression());
            if (argCtx.name != null) {
                namedArgs.put(argCtx.name.getText(), value);
            } else {
                posArgs.add(value);
            }
        }
        return new MethodCallParseCtx(methodName, posArgs, namedArgs);
    }

    @Override
    public Node visitDot(ReflParser.DotContext ctx) {
        var left = visit(ctx.simpleExpression());
        var methodCall = ctx.methodCall();
        if (methodCall != null) {
            var customCtx = visitMethodCallCustom(methodCall);
            var getMemberNode = new GetMemberNode(left, new StringValue(customCtx.getName()));
            return new FunctionCallNode(getMemberNode, customCtx.getPosArgs(), customCtx.getNamedArgs());
        }
        return new GetMemberNode(left, new StringValue(ctx.name.getText()));
    }

    @Override
    public Node visitArrows(ReflParser.ArrowsContext ctx) {
        return new ReturnNode(visit(ctx.simpleExpression()));
    }

    @Override
    public Node visitUnary(ReflParser.UnaryContext ctx) {
        var opText = ctx.prefix.getText();
        var child = visit(ctx.simpleExpression());
        switch (opText) {
            case "-":
                return new UnaryMinusNode(child);
            case "+":
                return new UnaryPlusNode(child);
            case "--":
                return new UnaryAssignNode(child, PrefixDecOperator.INSTANCE);
            case "++":
                return new UnaryAssignNode(child, PrefixIncOperator.INSTANCE);
            case "!":
                return new NotNode(child);
            default:
                throw new IllegalArgumentException("unknown unary operator " + opText);
        }
    }

    @Override
    public Node visitIncDec(ReflParser.IncDecContext ctx) {
        var opText = ctx.postfix.getText();
        AssignOperator operator;
        switch (opText) {
            case "++":
                operator = PostfixIncOperator.INSTANCE;
                break;
            case "--":
                operator = PostfixDecOperator.INSTANCE;
                break;
            default:
                throw new IllegalArgumentException("unknown postfix assign operator " + opText);
        }
        return new UnaryAssignNode(visit(ctx.simpleExpression()), operator);
    }

    @Override
    public Node visitArrayCall(ReflParser.ArrayCallContext ctx) {
        return new ArrayCallNode(visit(ctx.left), visit(ctx.right));
    }

    @Override
    public Node visitAssign(ReflParser.AssignContext ctx) {
        var opText = ctx.bop.getText();
        AssignOperator operator;
        switch (opText) {
            case "=":
                operator = SimpleAssignOperator.INSTANCE;
                break;
            case "+=":
                operator = AddAssignOperator.INSTANCE;
                break;
            case "-=":
                operator = SubAssignOperator.INSTANCE;
                break;
            case "*=":
                operator = MulAssignOperator.INSTANCE;
                break;
            case "/=":
                operator = DivAssignOperator.INSTANCE;
                break;
            case "%=":
                operator = ModAssignOperator.INSTANCE;
                break;
            case "&=":
                operator = AndAssignOperator.INSTANCE;
                break;
            case "|=":
                operator = OrAssignOperator.INSTANCE;
                break;
            default:
                throw new IllegalArgumentException("unknown infix assign operator " + opText);
        }
        return new AssignNode(visit(ctx.left), visit(ctx.right), operator);
    }

    @Override
    public Node visitParenExpression(ReflParser.ParenExpressionContext ctx) {
        return visit(ctx.simpleExpression());
    }

    @Override
    public Node visitIntegerLiteral(ReflParser.IntegerLiteralContext ctx) {
        var number = Double.parseDouble(ctx.getText());
        return new ConstNode(new NumberValue(number));
    }

    @Override
    public Node visitFloatLiteral(ReflParser.FloatLiteralContext ctx) {
        var number = Double.parseDouble(ctx.getText());
        return new ConstNode(new NumberValue(number));
    }

    @Override
    public Node visitStringLiteral(ReflParser.StringLiteralContext ctx) {
        var str = ctx.getText();
        var actualStr = StringEscapeUtils.unescapeJava(str.substring(1, str.length() - 1));
        return new ConstNode(new StringValue(actualStr));
    }

    @Override
    public Node visitNilLiteral(ReflParser.NilLiteralContext ctx) {
        return new ConstNode(NilValue.INSTANCE);
    }

    @Override
    public Node visitIdentifier(ReflParser.IdentifierContext ctx) {
        return new GetVarNode(new StringValue(ctx.getText()));
    }

    @Override
    public Node visitMethodCall(ReflParser.MethodCallContext ctx) {
        var customCtx = visitMethodCallCustom(ctx);
        var getVar = new GetVarNode(new StringValue(customCtx.getName()));
        return new FunctionCallNode(getVar, customCtx.getPosArgs(), customCtx.getNamedArgs());
    }
}
