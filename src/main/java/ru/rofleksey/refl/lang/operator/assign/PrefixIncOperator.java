package ru.rofleksey.refl.lang.operator.assign;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.GenericOperatorType;
import ru.rofleksey.refl.lang.value.NumberValue;

public class PrefixIncOperator implements AssignOperator {
    public static final PrefixIncOperator INSTANCE = new PrefixIncOperator();

    @Override
    public Value assign(Value leftSide, Value key, Value ignored) throws EvalError {
        var newValue = leftSide.getVar(key).add(NumberValue.TRUE);
        leftSide.setVar(key, newValue);
        return newValue;
    }

    @Override
    public GenericOperatorType type() {
        return GenericOperatorType.PREFIX;
    }

    @Override
    public String toString() {
        return "++";
    }
}
