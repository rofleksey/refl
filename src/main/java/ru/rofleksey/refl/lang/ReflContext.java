package ru.rofleksey.refl.lang;

import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.VarUndefinedError;
import ru.rofleksey.refl.std.*;
import ru.rofleksey.refl.util.HandshakeChannel;

import java.util.HashMap;
import java.util.Map;

public final class ReflContext {
    private final Map<String, Value> vars;
    private final HandshakeChannel<Value> valueChannel = new HandshakeChannel<>();


    public ReflContext() {
        vars = new HashMap<>();

        vars.put("wait", new StdWait());
        vars.put("exit", new StdExit());
        vars.put("random", new StdRandom());
        vars.put("floor", new StdFloor());
        vars.put("ceil", new StdCeil());
        vars.put("round", new StdRound());
    }

    public void setVar(String name, Value value) {
        vars.put(name, value);
    }

    public Value waitCtx() throws InterruptedException {
        return valueChannel.read();
    }

    public boolean notifyCtx(Value value) {
        return valueChannel.write(value);
    }

    public Value getVar(String name) throws EvalError {
        var value = vars.get(name);
        if (value == null) {
            throw new VarUndefinedError(name);
        }
        return value;
    }
}
