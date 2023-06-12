package ru.rofleksey.refl.lang;

import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.VarUndefinedError;
import ru.rofleksey.refl.lang.value.ValueType;
import ru.rofleksey.refl.util.HandshakeChannel;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

public final class ReflContext implements Value {
    private final Map<Value, Value> vars;
    private final HandshakeChannel<Value> valueChannel;

    private final ReflContext parentCtx;
    private final ReflContext rootCtx;
    private final Value objectCtx;

    private ReflContext(ReflContext parent, Value objectCtx) {
        this.parentCtx = parent;
        this.objectCtx = objectCtx;

        if (parent == null) {
            rootCtx = this;
            valueChannel = new HandshakeChannel<>();
        } else {
            valueChannel = parent.valueChannel;
            rootCtx = Objects.requireNonNullElse(parent.rootCtx, parent);
        }

        vars = new HashMap<>();
    }

    public static ReflContext empty() {
        return new ReflContext(null, null);
    }

    public static ReflContext startScope(Value objectCtx) {
        return new ReflContext(null, objectCtx);
    }

    public ReflContext shallowClone() {
        if (parentCtx == null) {
            return createChild();
        }
        var clone = new ReflContext(parentCtx, objectCtx);
        clone.vars.putAll(vars);
        return clone;
    }

    public ReflContext createChild() {
        return new ReflContext(this, null);
    }

    public ReflContext shallowRootClone() {
        return new ReflContext(rootCtx, null);
    }

//    public ReflContext initStd() {
//        vars.put("wait", new StdWait());
//        vars.put("sleep", new StdSleep());
//        vars.put("exit", new StdExit());
//        vars.put("random", new StdRandom());
//        vars.put("floor", new StdFloor());
//        vars.put("ceil", new StdCeil());
//        vars.put("round", new StdRound());
//        vars.put("string", new StdString());
//        vars.put("number", new StdNumber());
//
//        return this;
//    }

    @Override
    public ValueType getType() {
        return ValueType.OBJECT;
    }

    public void setVar(Value key, Value value) {
        if (objectCtx != null) {
            try {
                objectCtx.setVar(key, value);
            } catch (EvalError ignored) {

            }
            return;
        }

        var curCtx = this;

        while (curCtx != null && !curCtx.vars.containsKey(key)) {
            curCtx = curCtx.parentCtx;
        }

        if (curCtx == null) {
            vars.put(key, value);
        } else {
            curCtx.vars.put(key, value);
        }
    }

    public Value waitCtx() throws InterruptedException {
        return valueChannel.read();
    }

    public boolean notifyCtx(Value value) {
        return valueChannel.write(value);
    }

    public Value getVar(Value key) throws EvalError {
        if (objectCtx != null) {
            return objectCtx.getVar(key);
        }

        var value = vars.get(key);

        if (value == null) {
            if (parentCtx == null) {
                throw new VarUndefinedError(key);
            }

            return parentCtx.getVar(key);
        }

        return value;
    }

    @Override
    public boolean isTruthy() {
        return true;
    }

    @Override
    public String toString() {
        return "scope EvalScope";
    }
}
