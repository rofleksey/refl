package ru.rofleksey.refl.lang.error;

public class NoObjectContextError extends EvalError {
    public NoObjectContextError(String str) {
        super("node '" + str + "' doesn't have object context");
    }
}
