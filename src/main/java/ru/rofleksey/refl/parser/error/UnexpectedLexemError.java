package ru.rofleksey.refl.parser.error;

import ru.rofleksey.refl.lexer.LexemType;

public class UnexpectedLexemError extends ParserError {
    public UnexpectedLexemError(LexemType lexemType) {
        super("unexpected lexem " + lexemType);
    }
}
