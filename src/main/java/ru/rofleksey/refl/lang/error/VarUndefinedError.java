package ru.rofleksey.refl.lang.error;

public class VarUndefinedError extends EvalError {
    public VarUndefinedError(String str) {
        super("var '"+ str + "' is not defined");
    }
}
