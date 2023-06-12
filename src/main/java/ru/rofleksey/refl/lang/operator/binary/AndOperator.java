package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;
import ru.rofleksey.refl.lang.value.NumberValue;

public class AndOperator implements BinaryOperator {
    public static final AndOperator INSTANCE = new AndOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        if (!left.isTruthy()) {
            return NumberValue.FALSE;
        }
        return left.and(right.evaluate(ctx));
    }

    @Override
    public String toString() {
        return "&";
    }
}
