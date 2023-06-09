package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.ReflValue;

import java.util.List;
import java.util.Map;

public final class StdCeil extends FunctionValue {
    public StdCeil() {
        super("ceil");
    }

    @Override
    public Value call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        if (args.isEmpty()) {
            return ReflValue.INSTANCE;
        }
        var value = args.get(0).asNumber().getValue();
        return new NumberValue(Math.ceil(value));
    }
}
