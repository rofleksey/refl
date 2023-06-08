package ru.rofleksey.refl.lang;

import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.VarUndefinedError;
import ru.rofleksey.refl.std.StdExit;
import ru.rofleksey.refl.std.StdWait;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;

public final class ReflContext {
    private final Map<String, Value> vars;
    private final BlockingQueue<Value> waitQueue;

    public ReflContext() {
        vars = new HashMap<>();
        waitQueue = new ArrayBlockingQueue<>(1);

        vars.put("wait", new StdWait());
        vars.put("exit", new StdExit());
    }

    public void setVar(String name, Value value) {
        vars.put(name, value);
    }

    public Value waitCtxInternal() throws InterruptedException {
        return waitQueue.take();
    }

    public boolean notifyCtx(Value value) {
        return waitQueue.offer(value);
    }

    public Value getVar(String name) throws EvalError {
        var value = vars.get(name);
        if (value == null) {
            throw new VarUndefinedError(name);
        }
        return value;
    }
}
