package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;

public class AddOperator implements BinaryOperator {
    public static final AddOperator INSTANCE = new AddOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        return left.add(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "+";
    }
}
