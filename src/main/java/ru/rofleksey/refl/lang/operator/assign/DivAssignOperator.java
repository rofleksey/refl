package ru.rofleksey.refl.lang.operator.assign;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.GenericOperatorType;

public class DivAssignOperator implements AssignOperator {
    public static final DivAssignOperator INSTANCE = new DivAssignOperator();

    @Override
    public Value assign(Value leftSide, Value key, Value value) throws EvalError {
        var newValue = leftSide.getVar(key).divide(value);
        leftSide.setVar(key, newValue);
        return newValue;
    }

    @Override
    public GenericOperatorType type() {
        return GenericOperatorType.INFIX;
    }

    @Override
    public String toString() {
        return "/=";
    }
}
