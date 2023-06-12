package ru.rofleksey.refl.lang.error;

public class IndexOutOfBoundsError extends EvalError {
    public IndexOutOfBoundsError(int index) {
        super("index '" + index + "' is out of bounds");
    }
}
