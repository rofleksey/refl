package ru.rofleksey.refl.lang.operator.binary;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;
import ru.rofleksey.refl.lang.operator.BinaryOperator;
import ru.rofleksey.refl.lang.value.NilValue;

public class ElvisOperator implements BinaryOperator {
    public static final ElvisOperator INSTANCE = new ElvisOperator();

    @Override
    public Value apply(ReflContext ctx, Value left, Node right) throws EvalError {
        if (left != NilValue.INSTANCE) {
            return left;
        }
        return right.evaluate(ctx);
    }

    @Override
    public String toString() {
        return "??";
    }
}
