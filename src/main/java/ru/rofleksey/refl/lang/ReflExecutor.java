package ru.rofleksey.refl.lang;

import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.ExitCalledError;
import ru.rofleksey.refl.lang.node.Node;

public class ReflExecutor {
    private final Node rootNode;

    public ReflExecutor(Node rootNode) {
        this.rootNode = rootNode;
    }

    public Value execute(ReflContext ctx) throws EvalError {
        try {
            return rootNode.evaluate(ctx);
        } catch (ExitCalledError e) {
            return e.getValue();
        }
    }
}
