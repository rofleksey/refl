package ru.rofleksey.refl.lang;

import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ExitCalledError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.value.ReflValue;

import java.util.Collections;
import java.util.List;

public class ReflExecutor {
    private final List<Node> nodes;

    public ReflExecutor(List<Node> nodes) {
        this.nodes = Collections.unmodifiableList(nodes);
    }

    public Value execute(ReflContext ctx) throws EvalError {
        Value result = ReflValue.INSTANCE;
        for (var node : nodes) {
            try {
                result = node.evaluate(ctx);
            } catch (ExitCalledError e) {
                return e.getValue();
            }
        }
        return result;
    }
}
