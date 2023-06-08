package ru.rofleksey.refl.parser.error;

public class UnexpectedEofError extends ParserError {
    public UnexpectedEofError() {
        super("unexpected EOF");
    }
}
