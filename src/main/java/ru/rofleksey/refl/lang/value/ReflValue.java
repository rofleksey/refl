package ru.rofleksey.refl.lang.value;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;

import java.util.List;

public class ReflValue implements Value {
    public static final ReflValue INSTANCE = new ReflValue();

    @Override
    public  Value add(Value other) {
        return INSTANCE;
    }

    @Override
    public  Value subtract(Value other) {
        return INSTANCE;
    }

    @Override
    public  Value multiply(Value other) {
        return INSTANCE;
    }

    @Override
    public  Value divide(Value other) {
        return INSTANCE;
    }

    @Override
    public  Value and(Value other) {
        return NumberValue.FALSE;
    }

    @Override
    public  Value or(Value other) {
        if (other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public  Value compare(Value other) {
        return INSTANCE;
    }

    @Override
    public  Value not() {
        return INSTANCE;
    }

    @Override
    public  Value call(ReflContext ctx, List<Value> args) {
        return ReflValue.INSTANCE;
    }

    @Override
    public boolean isTruthy() {
        return false;
    }

    @Override
    public  StringValue asString() {
        return new StringValue("refl");
    }

    @Override
    public  NumberValue asNumber() {
        return new NumberValue(0);
    }

    @Override
    public  String getType() {
        return "refl";
    }
}
