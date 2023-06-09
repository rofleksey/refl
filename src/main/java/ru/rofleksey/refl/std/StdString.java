package ru.rofleksey.refl.std;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.StringValue;

import java.util.List;
import java.util.Map;

public final class StdString extends FunctionValue {
    public StdString() {
        super("string");
    }

    @Override
    public StringValue call(ReflContext ctx, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
        if (args.isEmpty()) {
            return StringValue.EMPTY;
        }

        return args.get(0).asString();
    }
}
