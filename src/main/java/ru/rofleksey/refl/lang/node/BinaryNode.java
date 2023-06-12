package ru.rofleksey.refl.lang.node;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NoObjectContextError;
import ru.rofleksey.refl.lang.operator.BinaryOperator;

public class BinaryNode implements Node {
    private final Node left;
    private final Node right;
    private final BinaryOperator operator;

    public BinaryNode(Node left, Node right, BinaryOperator op) {
        this.left = left;
        this.right = right;
        this.operator = op;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var leftValue = left.evaluate(ctx);
        return operator.apply(ctx, leftValue, right);
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
        return "(" + left.toString() + operator.toString() + right.toString() + ")";
    }
}
