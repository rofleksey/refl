package ru.rofleksey.refl.lang.operator.assign;

import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.operator.AssignOperator;
import ru.rofleksey.refl.lang.operator.GenericOperatorType;
import ru.rofleksey.refl.lang.value.NumberValue;

public class PostfixIncOperator implements AssignOperator {
    public static final PostfixIncOperator INSTANCE = new PostfixIncOperator();

    @Override
    public Value assign(Value leftSide, Value key, Value ignored) throws EvalError {
        var oldValue = leftSide.getVar(key);
        var newValue = oldValue.add(NumberValue.TRUE);
        leftSide.setVar(key, newValue);
        return oldValue;
    }

    @Override
    public GenericOperatorType type() {
        return GenericOperatorType.POSTFIX;
    }

    @Override
    public String toString() {
        return "++";
    }
}
