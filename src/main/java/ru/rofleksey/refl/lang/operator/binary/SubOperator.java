package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;

public class SubOperator implements BinaryOperator {
    public static final SubOperator INSTANCE = new SubOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        return left.subtract(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "-";
    }
}
