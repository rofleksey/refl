package ru.rofleksey.refl.lang.operator;

import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.node.Node;

public interface BinaryOperator {
    Value apply(ReflContext ctx, Value left, Node right) throws EvalError;
}
