package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NoObjectContextError;
import ru.rofleksey.refl.lang.operator.AssignOperator;

public final class AssignNode implements Node {
    private final Node left;
    private final Node right;
    private final AssignOperator operator;

    public AssignNode(Node left, Node right, AssignOperator operator) {
        this.left = left;
        this.right = right;
        this.operator = operator;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var objectContext = left.getLeftSide(ctx);
        var key = left.getSetterKey(ctx);
        var value = right.evaluate(ctx);
        return operator.assign(objectContext, key, value);
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
