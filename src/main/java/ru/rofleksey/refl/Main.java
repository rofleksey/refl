package ru.rofleksey.refl;

import ru.rofleksey.refl.lang.Refl;
import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ParseError;
import ru.rofleksey.refl.lang.value.FunctionValue;
import ru.rofleksey.refl.lang.value.NilValue;
import ru.rofleksey.refl.lang.value.NumberValue;

import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) throws EvalError, ParseError {
        var ctx = ReflContext.empty();

        ctx.setVar("x", new NumberValue(5));
        ctx.setVar("print", new FunctionValue("print") {
            @Override
            public Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) {
                var prefix = namedArgs.get("prefix").toString();
                args.forEach(it -> System.out.println(prefix + it.toString()));
                return NilValue.INSTANCE;
            }
        });

        var executor = Refl.compile("while x > 0 \n print(x, prefix ~ '>') \n x-- \n end");
        var result = executor.execute(ReflContext.empty());
    }
}