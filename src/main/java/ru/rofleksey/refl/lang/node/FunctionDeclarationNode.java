package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NoObjectContextError;
import ru.rofleksey.refl.lang.value.*;

import java.util.List;
import java.util.Map;

public final class FunctionDeclarationNode implements Node {
    private final StringValue key;
    private final Node body;

    public FunctionDeclarationNode(StringValue key, Node body) {
        this.key = key;
        this.body = body;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var func = new FunctionValue(key.toString()) {

            @Override
            public Value call(ReflContext ctx, Value thisValue, List<Value> args, Map<String, Value> namedArgs) throws EvalError {
                var newCtx = ctx.shallowClone();

                var argsVar = new ObjectValue();
                for (var i = 0; i < args.size(); i++) {
                    argsVar.setVar(new NumberValue(i), args.get(i));
                }
                for (var entry : namedArgs.entrySet()) {
                    argsVar.setVar(new StringValue(entry.getKey()), entry.getValue());
                }
                argsVar.setVar(new StringValue("length"), new NumberValue(args.size()));

                newCtx.setVar(new StringValue("args"), argsVar);
                if (args.isEmpty()) {
                    newCtx.setVar(new StringValue("it"), NilValue.INSTANCE);
                } else {
                    newCtx.setVar(new StringValue("it"), args.get(0));
                }
                newCtx.setVar(new StringValue("this"), thisValue);

                return body.evaluate(newCtx);
            }
        };

        ctx.setVar(key, func);

        return func;
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws EvalError {
        throw new NoObjectContextError(toString());
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return null;
    }

    @Override
    public String toString() {
        return "fun " + key.toString() + "()";
    }
}
