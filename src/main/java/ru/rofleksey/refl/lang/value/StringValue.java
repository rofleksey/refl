package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.NotCallableError;

import java.util.List;

public final class StringValue implements Value {
    private final String value;

    public StringValue(String value) {
        this.value = value;
    }

    @Override
    public Value add(Value other) {
        return new StringValue(value + other.asString().value);
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
    public  Value divide(Value other) {
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
    public  Value compare(Value other) {
        if (!getType().equals(other.getType())) {
            return ReflValue.INSTANCE;
        }
        var otherString = other.asString();
        return new NumberValue(value.compareTo(otherString.value));
    }

    @Override
    public  Value not() {
        if (isTruthy()) {
            return NumberValue.FALSE;
        }
        return NumberValue.TRUE;
    }

    @Override
    public  Value call(ReflContext ctx, List<Value> args) throws NotCallableError {
        throw new NotCallableError(toString());
    }

    @Override
    public boolean isTruthy() {
        return !value.isEmpty();
    }

    @Override
    public  StringValue asString() {
        return this;
    }

    @Override
    public  NumberValue asNumber() {
        if (isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public  String getType() {
        return "string";
    }

    @Override
    public String toString() {
        return value;
    }
}
