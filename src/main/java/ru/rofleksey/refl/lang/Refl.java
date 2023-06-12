package ru.rofleksey.refl.lang;

import org.antlr.v4.runtime.BailErrorStrategy;
import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;
import ru.rofleksey.refl.ReflLexer;
import ru.rofleksey.refl.ReflParser;
import ru.rofleksey.refl.lang.error.ParseError;
import ru.rofleksey.refl.lang.parse.ReflVisitor;

public class Refl {
    public static ReflExecutor compile(String code) throws ParseError {
        try {
            var inputStream = CharStreams.fromString(code);
            var reflLexer = new ReflLexer(inputStream);
            var tokenStream = new CommonTokenStream(reflLexer);
            var reflParser = new ReflParser(tokenStream);
            reflParser.setErrorHandler(new BailErrorStrategy());
            var antlrContext = reflParser.root();
            var reflVisitor = new ReflVisitor();
            var rootNode = reflVisitor.visitRoot(antlrContext);
            return new ReflExecutor(rootNode);
        } catch (Exception e) {
            throw new ParseError("parsing failed", e);
        }
    }
}
