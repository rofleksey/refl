package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;

public class DivOperator implements BinaryOperator {
    public static final DivOperator INSTANCE = new DivOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        return left.divide(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "/";
    }
}
