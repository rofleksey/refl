package ru.rofleksey.refl.lang.value;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.DivisionByZeroError;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotCallableError;

import java.sql.Ref;
import java.util.List;

public class NumberValue implements Value {
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
    public @NotNull Value add(Value other) {
        return new NumberValue(value + other.asNumber().value);
    }

    @Override
    public @NotNull Value subtract(Value other) {
        return new NumberValue(value - other.asNumber().value);
    }

    @Override
    public @NotNull Value multiply(Value other) {
        return new NumberValue(value * other.asNumber().value);
    }

    @Override
    public @NotNull Value divide(Value other) throws EvalError {
        var otherValue = other.asNumber().value;
        if (otherValue == 0) {
            throw new DivisionByZeroError();
        }
        return new NumberValue(value / otherValue);
    }

    @Override
    public @NotNull Value and(Value other) {
        if (isTruthy() && other.isTruthy()) {
            return TRUE;
        }
        return FALSE;
    }

    @Override
    public @NotNull Value or(Value other) {
        if (isTruthy()|| other.isTruthy()) {
            return TRUE;
        }
        return FALSE;
    }

    @Override
    public @NotNull Value compare(Value other) throws EvalError {
        if (!getType().equals(other.getType())) {
            return Refl.INSTANCE;
        }
        var otherNumber = other.asNumber();
        return new NumberValue(Double.compare(value, otherNumber.value));
    }

    @Override
    public @NotNull Value not() {
        if (isTruthy()) {
            return NumberValue.FALSE;
        }
        return NumberValue.TRUE;
    }

    @Override
    public @NotNull Value call(ReflContext ctx, List<Value> args) throws NotCallableError {
        throw new NotCallableError(toString());
    }

    @Override
    public boolean isTruthy() {
        return false;
    }

    @Override
    public @NotNull StringValue asString() {
        return new StringValue(toString());
    }

    @Override
    public @NotNull NumberValue asNumber() {
        return this;
    }

    @Override
    public @NotNull String getType() {
        return "number";
    }

    @Override
    public String toString() {
        return Double.toString(value);
    }
}
