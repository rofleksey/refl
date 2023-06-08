package ru.rofleksey.refl.lang.error;

import ru.rofleksey.refl.lang.Value;

public class ExitCalledError extends EvalError {
    private final Value value;

    public ExitCalledError(Value value) {
        super("exit called");
        this.value = value;
    }

    public Value getValue() {
        return value;
    }
}
