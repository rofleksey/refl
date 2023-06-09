package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public final class CallNode implements Node {
    private final Node child;
    private final List<Node> args;
    private final Map<String, Node> namedArgs;

    public CallNode(Node child, List<Node> args, Map<String, Node> namedArgs) {
        this.child = child;
        this.args = args;
        this.namedArgs = namedArgs;
    }


    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var argsValueList = new ArrayList<Value>(args.size());
        for (var arg : args) {
            argsValueList.add(arg.evaluate(ctx));
        }
        var namedArgsValueMap = new HashMap<String, Value>();
        for (var entry : namedArgs.entrySet()) {
            namedArgsValueMap.put(entry.getKey(), entry.getValue().evaluate(ctx));
        }
        return child.evaluate(ctx).call(ctx, argsValueList, namedArgsValueMap);
    }

    @Override
    public String toString() {
        return child.toString() + "(" +
                args.stream().map(Object::toString).collect(Collectors.joining(","))
                + ")";
    }
}
