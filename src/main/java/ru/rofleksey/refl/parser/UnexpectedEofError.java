package ru.rofleksey.refl.parser;

import ru.rofleksey.refl.lexer.LexemType;

public class UnexpectedEofError extends ParserError {
    public UnexpectedEofError() {
        super("unexpected EOF");
    }
}
