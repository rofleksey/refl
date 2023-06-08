package ru.rofleksey.refl;

import ru.rofleksey.refl.lang.ReflExecutor;
import ru.rofleksey.refl.lexer.Lexer;
import ru.rofleksey.refl.lexer.LexerError;
import ru.rofleksey.refl.parser.Parser;
import ru.rofleksey.refl.parser.error.ParserError;

public final class ReflCompiler {
    private final Lexer lexer = new Lexer();
    private final Parser parser = new Parser();

    public ReflExecutor compile(String code) throws LexerError, ParserError {
        var lexems = lexer.process(code);
        var nodes = parser.parse(lexems);
        return new ReflExecutor(nodes);
    }
}
