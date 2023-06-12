package ru.rofleksey.refl.lang.operator.assign;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.GenericOperatorType;
import ru.rofleksey.refl.lang.value.NumberValue;

public class PrefixDecOperator implements AssignOperator {
    public static final PrefixDecOperator INSTANCE = new PrefixDecOperator();

    @Override
    public Value assign(Value leftSide, Value key, Value ignored) throws EvalError {
        var newValue = leftSide.getVar(key).subtract(NumberValue.TRUE);
        leftSide.setVar(key, newValue);
        return newValue;
    }

    @Override
    public GenericOperatorType type() {
        return GenericOperatorType.PREFIX;
    }

    @Override
    public String toString() {
        return "--";
    }
}
