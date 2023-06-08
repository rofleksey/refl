package ru.rofleksey.refl.lang.error;

public class ExecutionInterruptedError extends EvalError {
    public ExecutionInterruptedError() {
        super("interrupted");
    }
}
