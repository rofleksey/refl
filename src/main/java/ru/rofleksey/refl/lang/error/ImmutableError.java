package ru.rofleksey.refl.lang.error;

public class ImmutableError extends EvalError {
    public ImmutableError(String str) {
        super("expression '" + str + "' is immutable");
    }
}
