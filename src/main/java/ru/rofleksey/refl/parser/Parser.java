package ru.rofleksey.refl.parser;

import ru.rofleksey.refl.lang.node.*;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.Refl;
import ru.rofleksey.refl.lang.value.StringValue;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;
import ru.rofleksey.refl.lexer.lexem.NumberLexem;
import ru.rofleksey.refl.lexer.lexem.StringLexem;
import ru.rofleksey.refl.lexer.lexem.VarLiteral;
import ru.rofleksey.refl.util.Pair;

import java.util.ArrayList;
import java.util.List;

// http://marvin.cs.uidaho.edu/Teaching/CS445/c-Grammar.pdf
// https://cyberzhg.github.io/toolbox/cfg2ll
// https://www.cs.princeton.edu/courses/archive/spring22/cos320/LL1/index.html

/*
  declList ::= decl ; declList'
      decl ::= if and : declList end
        decl    ::= while and : declList end
        decl    ::= s = and
         decl   ::= and
       and ::= orExp and'
     orExp ::= notExp orExp'
    notExp ::= rel
        notExp    ::= not rel
       rel ::= add rel'
       add ::= mul add''
       mul ::= unary mul'
     unary ::= term
        unary    ::= - term
      term ::= const
         term   ::= ( and )
        term    ::= call
      call ::= s ( args )
      args ::= argsList
        args    ::= ϵ
  argsList ::= s argsList''
        argsList    ::= const argsList''
      rel' ::= ϵ
       rel'    ::= < add
      add' ::= + mul
       add'     ::= - mul
 argsList' ::= s
      argsList'      ::= const
 declList' ::= decl ; declList'
       declList'     ::= ϵ
      and' ::= & orExp and'
        and'     ::= ϵ
    orExp' ::= or notExp orExp'
       orExp'     ::= ϵ
     add'' ::= add' add''
       add''     ::= ϵ
      mul' ::= * unary mul'
        mul'    ::= ϵ
argsList'' ::= , argsList' argsList''
       argsList''     ::= ϵ


//

declList -> declList decl ; | decl ;
decl -> if and : declList end | while and : declList end | s = and | and
and -> orExp | and & orExp
orExp -> notExp | orExp or notExp
notExp -> rel | not rel
rel -> add | add < add
add -> mul | add + mul | add - mul
mul -> unary | mul * unary
unary -> term | - term
term -> const | ( and ) | call
call -> s ( args )
args -> argsList | ϵ
argsList -> argsList , s | argsList , const | s | const
 */

public class Parser {

    public List<Node> parse(List<Lexem> input) throws ParserError {
        var ctx = new ParserContext(input);
        return parseStart(ctx);
    }

    private List<Node> parseStart(ParserContext ctx) throws ParserError {
        var result = new ArrayList<Node>();
        var nodeList = parseDeclList(ctx, result);
        ctx.consume(LexemType.EOF);
        return nodeList;
    }

    private List<Node> parseDeclList(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case IF, WHILE, STRING, NUMBER, VARIABLE, REFL, NOT, MINUS, BRACKET_OPEN -> {
                var decl = parseDecl(ctx);
                result.add(decl);
                var tailList = parseDeclListSlash(ctx, result);
                result.add(decl);
                result.addAll(tailList);
                return result;
            }
            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseDecl(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case IF, WHILE -> {
                ctx.take();
                var condition = parseAnd(ctx);
                ctx.consume(LexemType.COLON);
                var result = new ArrayList<Node>();
                var body = parseDeclList(ctx, result);
                ctx.consume(LexemType.END);
                if (curLexem.type().equals(LexemType.IF)) {
                    return new IfNode(condition, body);
                } else {
                    return new WhileNode(condition, body);
                }
            }

            case VARIABLE -> {
                if (ctx.lookUp(LexemType.ASSIGN, 1)) {
                    ctx.take();
                    var variable = ((VarLiteral) curLexem);
                    ctx.take();
                    var exp = parseAnd(ctx);
                    return new AssignNode(variable.name(), exp);
                }
                return parseAnd(ctx);
            }

            case NOT, MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                return parseAnd(ctx);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseAnd(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case NOT, MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                var first = parseOrExp(ctx);
                var second = parseAndSlash(ctx);
                if (second == null) {
                    return first;
                }
                return new AndNode(first, second);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseOrExp(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case NOT, MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                var first = parseNotExp(ctx);
                var second = parseOrExpSlash(ctx);
                return new OrNode(first, second);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseNotExp(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                return parseRel(ctx);
            }

            case NOT -> {
                ctx.take();
                var node = parseRel(ctx);
                return new NotNode(node);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseRel(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                var first = parseAdd(ctx);
                var other = parseRelSlash(ctx);
                if (other == null) {
                    return first;
                }
                var predicate = switch (other.second()) {
                    case LT -> CompareNode.LT;
                    case GT -> CompareNode.GT;
                    default -> CompareNode.EQ;
                };
                return new CompareNode(predicate, first, other.first());
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseAdd(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                var first = parseMul(ctx);
                var second = parseAddDoubleSlash(ctx);
                return new AddNode(first, second);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseMul(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS, STRING, NUMBER, REFL, BRACKET_OPEN -> {
                var first = parseUnary(ctx);
                var second = parseMulSlash(ctx);
                if (LexemType.MULTIPLY.equals(second.second())) {
                    return new MultiplyNode(first, second.first());
                }
                return new DivideNode(first, second.first());
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseUnary(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING, NUMBER, REFL, BRACKET_OPEN -> {
                return parseTerm(ctx);
            }

            case MINUS -> {
                ctx.take();
                var node = parseTerm(ctx);
                return new UnaryMinusNode(node);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseTerm(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING -> {
                ctx.take();
                var str = ((StringLexem) curLexem);
                return new ConstNode(new StringValue(str.text()));
            }

            case NUMBER -> {
                ctx.take();
                var num = ((NumberLexem) curLexem);
                return new ConstNode(new NumberValue(num.value()));
            }

            case REFL -> {
                ctx.take();
                return new ConstNode(Refl.INSTANCE);
            }

            case BRACKET_OPEN -> {
                ctx.take();
                var node = parseAnd(ctx);
                ctx.consume(LexemType.BRACKET_CLOSE);
                return node;
            }

            case VARIABLE -> {
                return parseCall(ctx);
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseCall(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.VARIABLE) {
            ctx.take();
            var variable = ((VarLiteral) curLexem);
            ctx.consume(LexemType.BRACKET_OPEN);
            var args = parseArgs(ctx);
            ctx.consume(LexemType.BRACKET_CLOSE);
            return new CallNode(new GetVarNode(variable.name()), args);
        }
        throw new UnexpectedLexemError(curLexem.type());
    }

    private List<Node> parseArgs(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING, NUMBER, REFL, VARIABLE -> {
                var result = new ArrayList<Node>();
                return parseArgsList(ctx, result);
            }

            default -> {
                return List.of();
            }
        }
    }

    private List<Node> parseArgsList(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING -> {
                ctx.take();
                var str = ((StringLexem) curLexem);
                result.add(new ConstNode(new StringValue(str.text())));
                return parseArgsListDoubleSlash(ctx, result);
            }

            case NUMBER -> {
                ctx.take();
                var num = ((NumberLexem) curLexem);
                result.add(new ConstNode(new NumberValue(num.value())));
                return parseArgsListDoubleSlash(ctx, result);
            }

            case REFL -> {
                ctx.take();
                result.add(new ConstNode(Refl.INSTANCE));
                return parseArgsListDoubleSlash(ctx, result);
            }

            case VARIABLE -> {
                ctx.take();
                var variable = ((VarLiteral) curLexem);
                result.add(new GetVarNode(variable.name()));
                return parseArgsListDoubleSlash(ctx, result);
            }

            default -> {
                return List.of();
            }
        }
    }

    private Pair<Node, LexemType> parseRelSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case LT, GT, EQ -> {
                ctx.take();
                return new Pair<>(parseAdd(ctx), curLexem.type());
            }

            default -> {
                return null;
            }
        }
    }

    private Pair<Node, LexemType> parseAddSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS, PLUS -> {
                ctx.take();
                return new Pair<>(parseMul(ctx), curLexem.type());
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseArgsListSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING -> {
                ctx.take();
                var str = ((StringLexem) curLexem);
                return new ConstNode(new StringValue(str.text()));
            }

            case NUMBER -> {
                ctx.take();
                var num = ((NumberLexem) curLexem);
                return new ConstNode(new NumberValue(num.value()));
            }

            case REFL -> {
                ctx.take();
                return new ConstNode(Refl.INSTANCE);
            }

            case VARIABLE -> {
                ctx.take();
                var variable = ((VarLiteral) curLexem);
                return new GetVarNode(variable.name());
            }

            default -> throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private List<Node> parseDeclListSlash(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case IF, WHILE, STRING, NUMBER, VARIABLE, REFL, NOT, MINUS, BRACKET_OPEN -> {
                var decl = parseDecl(ctx);
                result.add(decl);
                var tailList = parseDeclListSlash(ctx, result);
                result.add(decl);
                result.addAll(tailList);
                return result;
            }
            default -> {
                return List.of();
            }
        }
    }

    private Node parseAndSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.AND) {
            ctx.take();
            var right = parseOrExp(ctx);
            var next = parseAndSlash(ctx);
            if (next == null) {
                return right;
            }
            return new AndNode(right, next);
        }
        return null;
    }

    private Node parseOrExpSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.OR) {
            ctx.take();
            var right = parseNotExp(ctx);
            var next = parseOrExpSlash(ctx);
            if (next == null) {
                return right;
            }
            return new OrNode(right, next);
        }
        return null;
    }

    private Node parseAddDoubleSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case PLUS, MINUS -> {
                ctx.take();
                var first = parseAddSlash(ctx);
                var second = parseAddDoubleSlash(ctx);
                if (second == null) {
                    return first.first();
                }
                if (first.second() == LexemType.PLUS) {
                    return new AddNode(first.first(), second);
                }
                return new SubtractNode(first.first(), second);
            }
            default -> {
                return null;
            }
        }
    }

    private Pair<Node, LexemType> parseMulSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MULTIPLY, DIVIDE -> {
                ctx.take();
                var first = parseUnary(ctx);
                var second = parseMulSlash(ctx);
                if (second == null) {
                    return new Pair<>(first, curLexem.type());
                }
                if (second.second() == LexemType.MULTIPLY) {
                    return new Pair<>(new MultiplyNode(first, second.first()), curLexem.type());
                }
                return new Pair<>(new DivideNode(first, second.first()), curLexem.type());
            }

            default -> {
                return null;
            }
        }
    }

    private List<Node> parseArgsListDoubleSlash(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.COMMA) {
            ctx.take();
            var first = parseArgsListSlash(ctx);
            result.add(first);
            return parseArgsListDoubleSlash(ctx, result);
        }
        return List.of();
    }
}
