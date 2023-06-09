package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ExecutionInterruptedError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.ReflValue;

import java.util.List;
import java.util.Map;

public final class StdSleep extends FunctionValue {
    public StdSleep() {
        super("sleep");
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        if (args.isEmpty()) {
            return ReflValue.INSTANCE;
        }

        var time = (long) args.get(0).asNumber().getValue();

        try {
            Thread.sleep(time);
        } catch (InterruptedException e) {
            throw new ExecutionInterruptedError();
        }

        return ReflValue.INSTANCE;
    }
}
