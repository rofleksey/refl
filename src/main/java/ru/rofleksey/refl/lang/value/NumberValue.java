package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.DivisionByZeroError;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotCallableError;

import java.util.List;
import java.util.Map;

public final class NumberValue implements Value {
    public static final NumberValue TRUE = new NumberValue(1);
    public static final NumberValue FALSE = new NumberValue(0);
    public static final NumberValue MINUS_ONE = new NumberValue(-1);
    private final double value;

    public NumberValue(double value) {
        this.value = value;
    }

    public double getValue() {
        return value;
    }

    @Override
    public Value add(Value other) {
        return new NumberValue(value + other.asNumber().value);
    }

    @Override
    public Value subtract(Value other) {
        return new NumberValue(value - other.asNumber().value);
    }

    @Override
    public Value multiply(Value other) {
        return new NumberValue(value * other.asNumber().value);
    }

    @Override
    public Value divide(Value other) throws EvalError {
        var otherValue = other.asNumber().value;
        if (otherValue == 0) {
            throw new DivisionByZeroError();
        }
        return new NumberValue(value / otherValue);
    }

    @Override
    public Value and(Value other) {
        if (isTruthy() && other.isTruthy()) {
            return TRUE;
        }
        return FALSE;
    }

    @Override
    public Value or(Value other) {
        if (isTruthy() || other.isTruthy()) {
            return TRUE;
        }
        return FALSE;
    }

    @Override
    public Value compare(Value other) throws EvalError {
        if (!getType().equals(other.getType())) {
            return ReflValue.INSTANCE;
        }
        var otherNumber = other.asNumber();
        return new NumberValue(Double.compare(value, otherNumber.value));
    }

    @Override
    public Value not() {
        if (isTruthy()) {
            return NumberValue.FALSE;
        }
        return NumberValue.TRUE;
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws NotCallableError {
        throw new NotCallableError(toString());
    }

    @Override
    public boolean isTruthy() {
        return value > 0;
    }

    @Override
    public StringValue asString() {
        return new StringValue(toString());
    }

    @Override
    public NumberValue asNumber() {
        return this;
    }

    @Override
    public String getType() {
        return "number";
    }

    @Override
    public String toString() {
        return Double.toString(value);
    }
}
