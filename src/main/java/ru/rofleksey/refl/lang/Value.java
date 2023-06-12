package ru.rofleksey.refl.lang;


import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotCallableError;
import ru.rofleksey.refl.lang.value.NilValue;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.StringValue;
import ru.rofleksey.refl.lang.value.ValueType;

import java.util.List;
import java.util.Map;

public interface Value {

    default Value add(Value other) throws EvalError {
        return NilValue.INSTANCE;
    }

    default Value subtract(Value other) throws EvalError {
        return NilValue.INSTANCE;
    }

    default Value multiply(Value other) throws EvalError {
        return NilValue.INSTANCE;
    }

    default Value divide(Value other) throws EvalError {
        return NilValue.INSTANCE;
    }

    default Value mod(Value other) throws EvalError {
        return NilValue.INSTANCE;
    }

    default Value and(Value other) throws EvalError {
        if (isTruthy() && other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    default Value or(Value other) throws EvalError {
        if (isTruthy() || other.isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    default Value compare(Value other) throws EvalError {
        return NilValue.INSTANCE;
    }

    default Value not() throws EvalError {
        if (isTruthy()) {
            return NumberValue.FALSE;
        }
        return NumberValue.TRUE;
    }

    default Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        throw new NotCallableError(toString());
    }

    default boolean isTruthy() {
        return false;
    }

    default StringValue asString() {
        return new StringValue(toString());
    }

    default NumberValue asNumber() {
        if (isTruthy()) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    ValueType getType();

    void setVar(Value key, Value value) throws EvalError;

    Value getVar(Value key) throws EvalError;

    default void setVar(String key, Value value) throws EvalError {
        setVar(new StringValue(key), value);
    }

    default Value getVar(String key) throws EvalError {
        return getVar(new StringValue("key"));
    }
}
