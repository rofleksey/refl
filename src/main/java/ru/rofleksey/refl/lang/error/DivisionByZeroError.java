package ru.rofleksey.refl.lang.error;

public class DivisionByZeroError extends EvalError {
    public DivisionByZeroError() {
        super("division by zero");
    }
}
