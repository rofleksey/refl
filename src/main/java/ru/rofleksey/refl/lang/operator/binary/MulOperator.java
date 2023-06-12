package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;

public class MulOperator implements BinaryOperator {
    public static final MulOperator INSTANCE = new MulOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        return left.multiply(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "*";
    }
}
