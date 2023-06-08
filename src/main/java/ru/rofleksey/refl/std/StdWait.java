package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ExecutionInterruptedError;
import ru.rofleksey.refl.lang.value.FunctionValue;

import java.util.List;

public final class StdWait extends FunctionValue {
    public StdWait() {
        super("wait");
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args) throws EvalError {
        try {
            return ctx.waitCtxInternal();
        } catch (InterruptedException e) {
            throw new ExecutionInterruptedError();
        }
    }
}
