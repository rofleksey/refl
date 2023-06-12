package ru.rofleksey.refl.lang.operator;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;

public interface AssignOperator {
    Value assign(Value leftSide, Value key, Value value) throws EvalError;

    GenericOperatorType type();
}
