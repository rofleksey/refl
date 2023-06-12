package ru.rofleksey.refl.lang.error;

import ru.rofleksey.refl.lang.Value;

public class VarUndefinedError extends EvalError {
    public VarUndefinedError(Value str) {
        super("var '" + str + "' is not defined");
    }
}
