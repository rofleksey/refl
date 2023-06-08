package ru.rofleksey.refl.parser;

import ru.rofleksey.refl.lang.node.*;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.ReflValue;
import ru.rofleksey.refl.lang.value.StringValue;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;
import ru.rofleksey.refl.lexer.lexem.NumberLexem;
import ru.rofleksey.refl.lexer.lexem.StringLexem;
import ru.rofleksey.refl.lexer.lexem.VarLexem;
import ru.rofleksey.refl.parser.error.ParserError;
import ru.rofleksey.refl.parser.error.UnexpectedLexemError;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Predicate;

// http://marvin.cs.uidaho.edu/Teaching/CS445/c-Grammar.pdf
// https://cyberzhg.github.io/toolbox/cfg2ll
// https://www.cs.princeton.edu/courses/archive/spring22/cos320/LL1/index.html

/*
    declList ::= decl ; declList'
      decl ::= if and : declList end
         decl   ::= while and : declList end
         decl   ::= s = and
         decl   ::= and
       and ::= orExp and'
     orExp ::= notExp orExp'
    notExp ::= rel
         notExp   ::= not rel
       rel ::= add rel'
       add ::= mul add''
       mul ::= unary mul'
     unary ::= term
        unary    ::= - term
      term ::= const
        term    ::= s
        term    ::= ( and )
         term   ::= call
      call ::= s ( args )
      args ::= argsList
        args    ::= ϵ
  argsList ::= s argsList''
        argsList    ::= const argsList''
      rel' ::= ϵ
        rel'    ::= < add
      add' ::= + mul
        add'    ::= - mul
 argsList' ::= s
       argsList'     ::= const
 declList' ::= decl ; declList'
       declList'     ::= ϵ
      and' ::= & orExp and'
        and'    ::= ϵ
    orExp' ::= or notExp orExp'
      orExp'      ::= ϵ
     add'' ::= add' add''
       add''     ::= ϵ
      mul' ::= * unary mul'
        mul'    ::= ϵ
argsList'' ::= , argsList' argsList''
        argsList''    ::= ϵ


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

public final class Parser {

    public List<Node> parse(List<Lexem> input) throws ParserError {
        var ctx = new ParserContext(input);
        return parseStart(ctx);
    }

    private List<Node> parseStart(ParserContext ctx) throws ParserError {
        var result = new ArrayList<Node>();
        parseDeclList(ctx, result);
        ctx.consume(LexemType.EOF);
        return result;
    }

    private void parseDeclList(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case IF:
            case WHILE:
            case STRING:
            case NUMBER:
            case VARIABLE:
            case REFL:
            case NOT:
            case MINUS:
            case BRACKET_OPEN:
                var decl = parseDecl(ctx);
                result.add(decl);
                ctx.consume(LexemType.SEMICOLON);
                parseDeclListSlash(ctx, result);
                return;
            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseDecl(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case IF:
            case WHILE:
                ctx.take();
                var condition = parseAnd(ctx);
                ctx.consume(LexemType.COLON);
                var body = new ArrayList<Node>();
                parseDeclList(ctx, body);
                ctx.consume(LexemType.END);
                if (curLexem.type().equals(LexemType.IF)) {
                    return new IfNode(condition, body);
                } else {
                    return new WhileNode(condition, body);
                }

            case VARIABLE:
                if (ctx.lookUp(LexemType.ASSIGN, 1)) {
                    ctx.take();
                    var variable = ((VarLexem) curLexem);
                    ctx.take();
                    var exp = parseAnd(ctx);
                    return new AssignNode(variable.name(), exp);
                }
                return parseAnd(ctx);

            case NOT:
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case BRACKET_OPEN:
                return parseAnd(ctx);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseAnd(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case NOT:
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                var first = parseOrExp(ctx);
                var queue = new ArrayList<NodeWithLexem>();
                queue.add(new NodeWithLexem(first, null));
                return parseAndSlash(ctx, queue);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseOrExp(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case NOT:
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                var first = parseNotExp(ctx);
                var queue = new ArrayList<NodeWithLexem>();
                queue.add(new NodeWithLexem(first, null));
                return parseOrExpSlash(ctx, queue);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseNotExp(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                return parseRel(ctx);

            case NOT:
                ctx.take();
                var node = parseRel(ctx);
                return new NotNode(node);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseRel(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                var first = parseAdd(ctx);
                var other = parseRelSlash(ctx);
                if (other == null) {
                    return first;
                }

                Predicate<Double> predicate;
                if (other.lexemType() == LexemType.LT) {
                    predicate = CompareNode.LT;
                } else if (other.lexemType() == LexemType.GT) {
                    predicate = CompareNode.GT;
                } else {
                    predicate = CompareNode.EQ;
                }
                return new CompareNode(predicate, first, other.node());

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseAdd(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                var first = parseMul(ctx);
                var queue = new ArrayList<NodeWithLexem>();
                queue.add(new NodeWithLexem(first, null));
                return parseAddDoubleSlash(ctx, queue);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseMul(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS:
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                var first = parseUnary(ctx);
                var queue = new ArrayList<NodeWithLexem>();
                queue.add(new NodeWithLexem(first, null));
                var second = parseMulSlash(ctx, queue);
                if (second == null) {
                    return first;
                }
                return second;

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseUnary(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
            case BRACKET_OPEN:
                return parseTerm(ctx);

            case MINUS:
                ctx.take();
                var node = parseTerm(ctx);
                return new UnaryMinusNode(node);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseTerm(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING:
                ctx.take();
                var str = ((StringLexem) curLexem);
                return new ConstNode(new StringValue(str.text()));

            case NUMBER:
                ctx.take();
                var num = ((NumberLexem) curLexem);
                return new ConstNode(new NumberValue(num.value()));

            case REFL:
                ctx.take();
                return new ConstNode(ReflValue.INSTANCE);

            case BRACKET_OPEN:
                ctx.take();
                var node = parseAnd(ctx);
                ctx.consume(LexemType.BRACKET_CLOSE);
                return node;

            case VARIABLE:
                if (ctx.lookUp(LexemType.BRACKET_OPEN, 1)) {
                    return parseCall(ctx);
                }
                ctx.take();
                var variable = ((VarLexem) curLexem);
                return new GetVarNode(variable.name());

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseCall(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.VARIABLE) {
            ctx.take();
            var variable = ((VarLexem) curLexem);
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
            case STRING:
            case NUMBER:
            case REFL:
            case VARIABLE:
                var result = new ArrayList<Node>();
                return parseArgsList(ctx, result);

            default:
                return List.of();
        }
    }

    private List<Node> parseArgsList(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING:
                ctx.take();
                var str = ((StringLexem) curLexem);
                result.add(new ConstNode(new StringValue(str.text())));
                return parseArgsListDoubleSlash(ctx, result);

            case NUMBER:
                ctx.take();
                var num = ((NumberLexem) curLexem);
                result.add(new ConstNode(new NumberValue(num.value())));
                return parseArgsListDoubleSlash(ctx, result);

            case REFL:
                ctx.take();
                result.add(new ConstNode(ReflValue.INSTANCE));
                return parseArgsListDoubleSlash(ctx, result);

            case VARIABLE:
                ctx.take();
                var variable = ((VarLexem) curLexem);
                result.add(new GetVarNode(variable.name()));
                return parseArgsListDoubleSlash(ctx, result);

            default:
                return List.of();
        }
    }

    private NodeWithLexem parseRelSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case LT:
            case GT:
            case EQ:
                ctx.take();
                return new NodeWithLexem(parseAdd(ctx), curLexem.type());

            default:
                return null;
        }
    }

    private Node parseAddSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MINUS:
            case PLUS:
                ctx.take();
                return parseMul(ctx);

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private Node parseArgsListSlash(ParserContext ctx) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case STRING:
                ctx.take();
                var str = ((StringLexem) curLexem);
                return new ConstNode(new StringValue(str.text()));

            case NUMBER:
                ctx.take();
                var num = ((NumberLexem) curLexem);
                return new ConstNode(new NumberValue(num.value()));

            case REFL:
                ctx.take();
                return new ConstNode(ReflValue.INSTANCE);

            case VARIABLE:
                ctx.take();
                var variable = ((VarLexem) curLexem);
                return new GetVarNode(variable.name());

            default:
                throw new UnexpectedLexemError(curLexem.type());
        }
    }

    private void parseDeclListSlash(ParserContext ctx, List<Node> result) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case IF:
            case WHILE:
            case STRING:
            case NUMBER:
            case VARIABLE:
            case REFL:
            case NOT:
            case MINUS:
            case BRACKET_OPEN:
                var decl = parseDecl(ctx);
                result.add(decl);
                ctx.consume(LexemType.SEMICOLON);
                parseDeclListSlash(ctx, result);
        }
    }

    private Node parseAndSlash(ParserContext ctx, List<NodeWithLexem> queue) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.AND) {
            ctx.take();
            var right = parseOrExp(ctx);
            queue.add(new NodeWithLexem(right, curLexem.type()));
            return parseAndSlash(ctx, queue);
        }
        return combineQueue(queue);
    }

    private Node parseOrExpSlash(ParserContext ctx, List<NodeWithLexem> queue) throws ParserError {
        var curLexem = ctx.peek();
        if (curLexem.type() == LexemType.OR) {
            ctx.take();
            var right = parseNotExp(ctx);
            queue.add(new NodeWithLexem(right, curLexem.type()));
            return parseOrExpSlash(ctx, queue);
        }
        return combineQueue(queue);
    }

    private Node combineQueue(List<NodeWithLexem> queue) {
        if (queue.size() == 1) {
            return queue.get(0).node();
        }
        var result = queue.get(0).node();
        for (var i = 1; i < queue.size(); i++) {
            var item = queue.get(i);
            switch (item.lexemType()) {
                case PLUS:
                    result = new AddNode(result, item.node());
                    break;

                case MINUS:
                    result = new SubtractNode(result, item.node());
                    break;

                case MULTIPLY:
                    result = new MultiplyNode(result, item.node());
                    break;

                case DIVIDE:
                    result = new DivideNode(result, item.node());
                    break;

                case AND:
                    result = new AndNode(result, item.node());
                    break;

                case OR:
                    result = new OrNode(result, item.node());
                    break;
            }
        }
        return result;
    }

    private Node parseAddDoubleSlash(ParserContext ctx, List<NodeWithLexem> queue) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case PLUS:
            case MINUS:
                var first = parseAddSlash(ctx);
                queue.add(new NodeWithLexem(first, curLexem.type()));
                return parseAddDoubleSlash(ctx, queue);

            default:
                return combineQueue(queue);
        }
    }

    private Node parseMulSlash(ParserContext ctx, List<NodeWithLexem> queue) throws ParserError {
        var curLexem = ctx.peek();
        switch (curLexem.type()) {
            case MULTIPLY:
            case DIVIDE:
                ctx.take();
                var first = parseUnary(ctx);
                queue.add(new NodeWithLexem(first, curLexem.type()));
                return parseMulSlash(ctx, queue);

            default:
                return combineQueue(queue);
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
        return result;
    }
}
