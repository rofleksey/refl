package ru.rofleksey.refl.lang;

import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.VarUndefinedError;

import java.util.HashMap;
import java.util.Map;

public class ReflContext {
    private final Map<String, Value> vars = new HashMap<>();

    public void setVar(String name, Value value) {
        vars.put(name, value);
    }

    public Value getVar(String name) throws EvalError {
        var value = vars.get(name);
        if (value == null) {
            throw new VarUndefinedError(name);
        }
        return value;
    }
}
