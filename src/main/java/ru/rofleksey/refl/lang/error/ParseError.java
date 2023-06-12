package ru.rofleksey.refl.lang.error;

public class ParseError extends Exception {
    public ParseError(String message) {
        super(message);
    }

    public ParseError(String message, Exception cause) {
        super(message, cause);
    }
}
