package ru.rofleksey.refl.lang.node;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NotReferencableError;
import ru.rofleksey.refl.lang.value.NilValue;

import java.util.List;

public final class ListNode implements Node {
    private final List<Node> children;

    public ListNode(List<Node> children) {
        this.children = children;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        Value result = NilValue.INSTANCE;

        for (var node : children) {
            result = node.evaluate(ctx);
        }

        return result;
    }

    @Override
    public Value getLeftSide(ReflContext ctx) throws NotReferencableError {
        throw new NotReferencableError(toString());
    }

    @Override
    public Value getSetterKey(ReflContext ctx) throws EvalError {
        return null;
    }
}
