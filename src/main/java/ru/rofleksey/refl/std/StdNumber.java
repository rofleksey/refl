package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NumberValue;

import java.util.List;
import java.util.Map;

public final class StdNumber extends FunctionValue {
    public StdNumber() {
        super("number");
    }

    @Override
    public NumberValue call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        if (args.isEmpty()) {
            return NumberValue.FALSE;
        }

        return args.get(0).asNumber();
    }
}
