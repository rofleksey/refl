package ru.rofleksey.refl.lang.operator.assign;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.GenericOperatorType;

public class SimpleAssignOperator implements AssignOperator {
    public static final SimpleAssignOperator INSTANCE = new SimpleAssignOperator();

    @Override
    public Value assign(Value leftSide, Value key, Value value) throws EvalError {
        leftSide.setVar(key, value);
        return value;
    }

    @Override
    public GenericOperatorType type() {
        return GenericOperatorType.INFIX;
    }

    @Override
    public String toString() {
        return "=";
    }
}
