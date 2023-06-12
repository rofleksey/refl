package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotReferencableError;

public abstract class FunctionValue implements Value {
    private final String name;

    public FunctionValue(String name) {
        this.name = name;
    }

    @Override
    public Value compare(Value other) throws EvalError {
        if (this == other) {
            return NumberValue.TRUE;
        }
        return NilValue.INSTANCE;
    }

    @Override
    public boolean isTruthy() {
        return true;
    }

    @Override
    public ValueType getType() {
        return ValueType.FUNCTION;
    }

    @Override
    public String toString() {
        return "fun " + name;
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
