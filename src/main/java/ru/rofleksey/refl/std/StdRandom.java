package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NumberValue;

import java.util.List;
import java.util.Map;

public final class StdRandom extends FunctionValue {
    public StdRandom() {
        super("random");
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        return new NumberValue(Math.random());
    }
}
