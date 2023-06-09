package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ExitCalledError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.ReflValue;

import java.util.List;
import java.util.Map;

public final class StdExit extends FunctionValue {
    public StdExit() {
        super("exit");
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        Value returnValue = ReflValue.INSTANCE;
        if (!args.isEmpty()) {
            returnValue = args.get(0);
        }
        throw new ExitCalledError(returnValue);
    }
}
