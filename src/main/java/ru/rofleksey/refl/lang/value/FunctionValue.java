package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public abstract class FunctionValue implements Value {
    private final String name;

    public FunctionValue(String name) {
        this.name = name;
    }

    @Override
    public  Value add(Value other) {
        return ReflValue.INSTANCE;
    }

    @Override
    public  Value subtract(Value other) {
        return ReflValue.INSTANCE;
    }

    @Override
    public  Value multiply(Value other) {
        return ReflValue.INSTANCE;
    }

    @Override
    public  Value divide(Value other) throws EvalError {
        return ReflValue.INSTANCE;
    }

    @Override
    public  Value and(Value other) {
        if (isTruthy() && other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public  Value or(Value other) {
        if (isTruthy() || other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public  Value compare(Value other) throws EvalError {
        if (this == other) {
            return NumberValue.TRUE;
        }
        return ReflValue.INSTANCE;
    }

    @Override
    public  Value not() {
        return NumberValue.FALSE;
    }

    @Override
    public boolean isTruthy() {
        return true;
    }

    @Override
    public  StringValue asString() {
        return new StringValue(toString());
    }

    @Override
    public  NumberValue asNumber() {
        return NumberValue.TRUE;
    }

    @Override
    public  String getType() {
        return "function";
    }

    @Override
    public String toString() {
        return "function "+name + "{}";
    }
}
