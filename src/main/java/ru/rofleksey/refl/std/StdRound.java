package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NilValue;
import ru.rofleksey.refl.lang.value.NumberValue;

import java.util.List;
import java.util.Map;

public final class StdRound extends FunctionValue {
    public StdRound() {
        super("round");
    }

    @Override
    public Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        if (args.isEmpty()) {
            return NilValue.INSTANCE;
        }
        var value = args.get(0).asNumber().getValue();
        return new NumberValue(Math.round(value));
    }
}
