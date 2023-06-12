package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NoObjectContextError;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.GenericOperatorType;

public final class UnaryAssignNode implements Node {
    private final Node left;
    private final AssignOperator operator;

    public UnaryAssignNode(Node left, AssignOperator operator) {
        this.left = left;
        this.operator = operator;
    }

    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var objectContext = left.getLeftSide(ctx);
        var key = left.getSetterKey(ctx);
        return operator.assign(objectContext, key, null);
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
        if (operator.type() == GenericOperatorType.PREFIX) {
            return "(" + operator + left.toString() + ")";
        }
        return "(" + left.toString() + operator + ")";
    }
}
