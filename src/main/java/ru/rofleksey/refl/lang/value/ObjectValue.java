package ru.rofleksey.refl.lang.value;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

import java.util.HashMap;
import java.util.Map;

public final class ObjectValue implements Value {
    private final Map<Value, Value> fields = new HashMap<>();

    @Override
    public Value compare(Value other) throws EvalError {
        if (other.getType() != ValueType.OBJECT) {
            return NilValue.INSTANCE;
        }

        var otherObj = ((ObjectValue) other);
        if (fields.equals(otherObj.fields)) {
            return NumberValue.FALSE;
        }

        return NumberValue.TRUE;
    }

    @Override
    public boolean isTruthy() {
        return fields.size() > 0;
    }

    @Override
    public NumberValue asNumber() {
        return isTruthy() ? NumberValue.TRUE : NumberValue.FALSE;
    }

    @Override
    public ValueType getType() {
        return ValueType.OBJECT;
    }

    @Override
    public void setVar(Value key, Value value) throws EvalError {
        fields.put(key, value);
    }

    @Override
    public Value getVar(Value key) throws EvalError {
        var value = fields.get(key);
        if (value == null) {
            return NilValue.INSTANCE;
        }
        return value;
    }
}
