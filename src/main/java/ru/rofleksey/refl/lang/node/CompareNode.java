package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.error.NoObjectContextError;
import ru.rofleksey.refl.lang.operator.CompareOperator;
import ru.rofleksey.refl.lang.value.NumberValue;
import ru.rofleksey.refl.lang.value.ValueType;

public final class CompareNode implements Node {
    private final CompareOperator operator;
    private final Node left;
    private final Node right;

    public CompareNode(CompareOperator operator, Node left, Node right) {
        this.operator = operator;
        this.left = left;
        this.right = right;
    }


    @Override
    public Value evaluate(ReflContext ctx) throws EvalError {
        var compare = left.evaluate(ctx).compare(right.evaluate(ctx));
        if (compare.getType() == ValueType.NIL) {
            return NumberValue.FALSE;
        }
        var num = compare.asNumber();
        if (operator.test(num.getValue())) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
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
