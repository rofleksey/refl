package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.NotCallableError;

import java.util.List;
import java.util.Map;

public final class ReflValue implements Value {
    public static final ReflValue INSTANCE = new ReflValue();

    @Override
    public Value add(Value other) {
        return INSTANCE;
    }

    @Override
    public Value subtract(Value other) {
        return INSTANCE;
    }

    @Override
    public Value multiply(Value other) {
        return INSTANCE;
    }

    @Override
    public Value divide(Value other) {
        return INSTANCE;
    }

    @Override
    public Value and(Value other) {
        return NumberValue.FALSE;
    }

    @Override
    public Value or(Value other) {
        if (other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public Value compare(Value other) {
        return INSTANCE;
    }

    @Override
    public Value not() {
        return INSTANCE;
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws NotCallableError {
        throw new NotCallableError(toString());
    }

    @Override
    public boolean isTruthy() {
        return false;
    }

    @Override
    public StringValue asString() {
        return new StringValue(toString());
    }

    @Override
    public NumberValue asNumber() {
        return new NumberValue(0);
    }

    @Override
    public String getType() {
        return "refl";
    }

    @Override
    public String toString() {
        return "refl";
    }
}
