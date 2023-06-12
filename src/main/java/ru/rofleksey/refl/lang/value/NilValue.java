package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotReferencableError;

public final class NilValue implements Value {
    public static final NilValue INSTANCE = new NilValue();

    @Override
    public ValueType getType() {
        return ValueType.NIL;
    }

    @Override
    public String toString() {
        return "nil";
    }

    @Override
    public void setVar(Value key, Value value) throws EvalError {
        throw new NotReferencableError(toString());
    }

    @Override
    public Value getVar(Value key) throws EvalError {
        throw new NotReferencableError(toString());
    }
}
