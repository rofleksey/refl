package ru.rofleksey.refl.lang.error;

public class NotCallableError extends EvalError {
    public NotCallableError(String str) {
        super("expression '"+ str + "' is not callable");
    }
}
