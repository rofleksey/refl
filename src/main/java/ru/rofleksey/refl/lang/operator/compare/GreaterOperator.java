package ru.rofleksey.refl.lang.operator.compare;

import ru.rofleksey.refl.lang.operator.CompareOperator;

public class GreaterOperator implements CompareOperator {
    public static final GreaterOperator INSTANCE = new GreaterOperator();

    @Override
    public boolean test(Double compareResult) {
        return compareResult > 0;
    }

    @Override
    public String toString() {
        return ">";
    }
}
