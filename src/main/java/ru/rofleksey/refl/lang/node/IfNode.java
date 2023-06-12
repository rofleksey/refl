package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotReferencableError;
import ru.rofleksey.refl.lang.value.NilValue;

import java.util.List;

public final class IfNode implements Node {
    private final Node condition;
    private final Node body;
    private final List<IfNode> elifNodes;

    public IfNode(Node condition, Node body, List<IfNode> elifNodes) {
        this.condition = condition;
        this.body = body;
        this.elifNodes = elifNodes;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        if (condition.evaluate(ctx).isTruthy()) {
            var newCtx = ctx.createChild();
            return body.evaluate(newCtx);
        }

        if (elifNodes != null) {
            for (var node : elifNodes) {
                if (node.condition == null || node.evaluate(ctx).isTruthy()) {
                    var newCtx = ctx.createChild();
                    return node.body.evaluate(newCtx);
                }
            }
        }

        return NilValue.INSTANCE;
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws NotReferencableError {
        throw new NotReferencableError(toString());
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return null;
    }

    @Override
    public String toString() {
        return "if " + condition.toString() + ": " + body.toString();
    }
}
