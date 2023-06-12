package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;
import ru.rofleksey.refl.lang.value.NumberValue;

public class OrOperator implements BinaryOperator {
    public static final OrOperator INSTANCE = new OrOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        if (left.isTruthy()) {
            return NumberValue.TRUE;
        }
        return left.or(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "|";
    }
}
