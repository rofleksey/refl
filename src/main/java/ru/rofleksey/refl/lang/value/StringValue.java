package ru.rofleksey.refl.lang.value;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotCallableError;

import java.util.List;

public class StringValue implements Value {
    private final String value;

    public StringValue(String value) {
        this.value = value;
    }

    @Override
    public @NotNull Value add(Value other) {
        return new StringValue(value + other.asString().value);
    }

    @Override
    public @NotNull Value subtract(Value other) {
        return Refl.INSTANCE;
    }

    @Override
    public @NotNull Value multiply(Value other) {
        return Refl.INSTANCE;
    }

    @Override
    public @NotNull Value divide(Value other) {
        return Refl.INSTANCE;
    }

    @Override
    public @NotNull Value and(Value other) {
        if (isTruthy() && other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public @NotNull Value or(Value other) {
        if (isTruthy()|| other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public @NotNull Value compare(Value other) {
        if (!getType().equals(other.getType())) {
            return Refl.INSTANCE;
        }
        var otherString = other.asString();
        return new NumberValue(value.compareTo(otherString.value));
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
        return !value.isEmpty();
    }

    @Override
    public @NotNull StringValue asString() {
        return this;
    }

    @Override
    public @NotNull NumberValue asNumber() {
        if (isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public @NotNull String getType() {
        return "string";
    }
}
