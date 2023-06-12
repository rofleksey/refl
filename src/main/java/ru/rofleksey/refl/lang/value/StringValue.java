package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ImmutableError;
import ru.rofleksey.refl.lang.error.IndexOutOfBoundsError;
import ru.rofleksey.refl.lang.error.NotCallableError;

import java.util.List;
import java.util.Map;
import java.util.Objects;

public final class StringValue implements Value {
    public static final StringValue EMPTY = new StringValue("");
    private final String value;

    public StringValue(String value) {
        this.value = value;
    }

    @Override
    public Value add(Value other) {
        return new StringValue(value + other.asString().value);
    }

    @Override
    public Value subtract(Value other) {
        return NilValue.INSTANCE;
    }

    @Override
    public Value multiply(Value other) {
        return NilValue.INSTANCE;
    }

    @Override
    public Value divide(Value other) {
        return NilValue.INSTANCE;
    }

    @Override
    public Value mod(Value other) {
        return NilValue.INSTANCE;
    }

    @Override
    public Value compare(Value other) {
        if (!getType().equals(other.getType())) {
            return NilValue.INSTANCE;
        }
        var otherString = other.asString();
        return new NumberValue(value.compareTo(otherString.value));
    }

    @Override
    public Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) throws NotCallableError {
        throw new NotCallableError(toString());
    }

    @Override
    public void setVar(Value key, Value value) throws EvalError {
        throw new ImmutableError(StringValue.this.toString());
    }

    @Override
    public Value getVar(Value key) throws EvalError {
        if (key.getType() != ValueType.NUMBER) {
            return NilValue.INSTANCE;
        }

        var index = ((int) key.asNumber().getValue());
        if (index < 0 || index >= value.length()) {
            throw new IndexOutOfBoundsError(index);
        }

        return new StringValue(Character.toString(value.charAt(index)));
    }

    @Override
    public boolean isTruthy() {
        return !value.isEmpty();
    }

    @Override
    public StringValue asString() {
        return this;
    }

    @Override
    public NumberValue asNumber() {
        try {
            return new NumberValue(Double.parseDouble(value));
        } catch (NumberFormatException ignored) {
            if (isTruthy()) {
                return NumberValue.TRUE;
            }
            return NumberValue.FALSE;
        }
    }

    @Override
    public ValueType getType() {
        return ValueType.STRING;
    }

    @Override
    public String toString() {
        return value;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        StringValue that = (StringValue) o;
        return Objects.equals(value, that.value);
    }

    @Override
    public int hashCode() {
        return Objects.hash(value);
    }
}
