package ru.rofleksey.refl.lang.node;


import ru.rofleksey.refl.lang.ReflContext;
import ru.rofleksey.refl.lang.Value;
import ru.rofleksey.refl.lang.error.EvalError;
import ru.rofleksey.refl.lang.value.NumberValue;

import java.util.Comparator;
import java.util.function.Predicate;

public class CompareNode implements Node {
    public static final Predicate<Double> LT = value -> value < 0;
    public static final Predicate<Double> EQ = value -> value == 0;
    public static final Predicate<Double> GT = value -> value > 0;
    private final Predicate<Double> predicate;
    private final Node left;
    private final Node right;

    public CompareNode(Predicate<Double> predicate, Node left, Node right) {
        this.predicate = predicate;
        this.left = left;
        this.right = right;
    }


    @Override
    public  Value evaluate(ReflContext ctx) throws EvalError {
        var compare = left.evaluate(ctx).compare(right.evaluate(ctx));
        if (compare.getType().equals("refl")) {
            return NumberValue.FALSE;
        }
        var num = compare.asNumber();
        if (predicate.test(num.getValue())) {
            return NumberValue.TRUE;
        }
        return NumberValue.FALSE;
    }

    @Override
    public String toString() {
        if (predicate == LT) {
            return "(" + left.toString() + "<" + right.toString() + ")";
        }
        if (predicate == GT) {
            return "(" + left.toString() + ">" + right.toString() + ")";
        }
        return "(" + left.toString() + "==" + right.toString() + ")";
    }
}
