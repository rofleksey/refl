package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.DivisionByZeroError;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotCallableError;
import ru.rofleksey.refl.lang.error.NotReferencableError;

import java.util.List;
import java.util.Map;
import java.util.Objects;

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
    public Value mod(Value other) throws EvalError {
        var otherValue = other.asNumber().value;
        if (otherValue == 0) {
            throw new DivisionByZeroError();
        }
        return new NumberValue(value % otherValue);
    }

    @Override
    public Value compare(Value other) throws EvalError {
        if (!getType().equals(other.getType())) {
            return NilValue.INSTANCE;
        }
        var otherNumber = other.asNumber();
        return new NumberValue(Double.compare(value, otherNumber.value));
    }

    @Override
    public Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) throws NotCallableError {
        throw new NotCallableError(toString());
    }

    @Override
    public boolean isTruthy() {
        return value > 0;
    }

    @Override
    public NumberValue asNumber() {
        return this;
    }

    @Override
    public ValueType getType() {
        return ValueType.NUMBER;
    }

    @Override
    public void setVar(Value key, Value value) throws EvalError {
        throw new NotReferencableError(toString());
    }

    @Override
    public Value getVar(Value key) throws EvalError {
        throw new NotReferencableError(toString());
    }

    @Override
    public String toString() {
        return Double.toString(value);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        NumberValue that = (NumberValue) o;
        return Double.compare(that.value, value) == 0;
    }

    @Override
    public int hashCode() {
        return Objects.hash(value);
    }
}
