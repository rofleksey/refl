package ru.rofleksey.refl.lang.error;

import ru.rofleksey.refl.lang.Value;

public class ReturnCalledError extends EvalError {
    private final Value value;

    public ReturnCalledError(Value value) {
        super("return called");
        this.value = value;
    }

    public Value getValue() {
        return value;
    }
}
