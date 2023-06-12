package ru.rofleksey.refl.lang.error;

public class NotReferencableError extends EvalError {
    public NotReferencableError(String str) {
        super("expression '" + str + "' can not be dereferenced");
    }
}
