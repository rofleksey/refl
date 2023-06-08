package ru.rofleksey.refl.parser.error;

import ru.rofleksey.refl.lexer.LexemType;

public class InvalidLexemError extends ParserError {
    public InvalidLexemError(LexemType expected, LexemType found) {
        super(expected + " expected, " + found + " found");
    }
}
