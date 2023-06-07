package ru.rofleksey.refl.lang.value;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

import java.util.List;

public class Refl implements Value {
    public static final Refl INSTANCE = new Refl();

    @Override
    public @NotNull Value add(Value other) {
        return INSTANCE;
    }

    @Override
    public @NotNull Value subtract(Value other) {
        return INSTANCE;
    }

    @Override
    public @NotNull Value multiply(Value other) {
        return INSTANCE;
    }

    @Override
    public @NotNull Value divide(Value other) {
        return INSTANCE;
    }

    @Override
    public @NotNull Value and(Value other) {
        return NumberValue.FALSE;
    }

    @Override
    public @NotNull Value or(Value other) {
        if (other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public @NotNull Value compare(Value other) {
        return INSTANCE;
    }

    @Override
    public @NotNull Value not() {
        return INSTANCE;
    }

    @Override
    public @NotNull Value call(ReflContext ctx, List<Value> args) {
        return Refl.INSTANCE;
    }

    @Override
    public boolean isTruthy() {
        return false;
    }

    @Override
    public @NotNull StringValue asString() {
        return new StringValue("refl");
    }

    @Override
    public @NotNull NumberValue asNumber() {
        return new NumberValue(0);
    }

    @Override
    public @NotNull String getType() {
        return "refl";
    }
}
