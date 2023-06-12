package ru.rofleksey.refl.lang.error;

public class NotIndexableError extends EvalError {
    public NotIndexableError(String str) {
        super("expression '" + str + "' is not indexable");
    }
}
