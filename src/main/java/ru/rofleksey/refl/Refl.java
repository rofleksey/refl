package ru.rofleksey.refl;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.ReflValue;
import ru.rofleksey.refl.lexer.Lexer;
import ru.rofleksey.refl.lexer.LexerError;
import ru.rofleksey.refl.parser.Parser;
import ru.rofleksey.refl.parser.ParserError;

public class Refl {
    public static Value eval(ReflContext ctx, String code) throws LexerError, ParserError, EvalError {
        var lexer = new Lexer();
        var lexems = lexer.process(code);
        var parser = new Parser();
        var nodes = parser.parse(lexems);
        Value result = ReflValue.INSTANCE;
        for (var node : nodes) {
            result = node.evaluate(ctx);
        }
        return result;
    }
}
