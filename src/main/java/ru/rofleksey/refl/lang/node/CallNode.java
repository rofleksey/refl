package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

public final class CallNode implements Node {
    private final Node child;
    private final List<Node> args;

    public CallNode(Node child, List<Node> args) {
        this.child = child;
        this.args = args;
    }


    @Override
    public  Value evaluate(ReflContext ctx) throws EvalError {
        var valueList = new ArrayList<Value>(args.size());
        for (var arg : args) {
            valueList.add(arg.evaluate(ctx));
        }
        return child.evaluate(ctx).call(ctx, valueList);
    }

    @Override
    public String toString() {
        return child.toString() + "(" +
                args.stream().map(Object::toString).collect(Collectors.joining(","))
                + ")";
    }
}
