package ru.rofleksey.refl.lang.operator.compare;

import ru.rofleksey.refl.lang.operator.CompareOperator;

public class EqualsOperator implements CompareOperator {
    public static final EqualsOperator INSTANCE = new EqualsOperator();

    @Override
    public boolean test(Double compareResult) {
        return compareResult == 0;
    }

    @Override
    public String toString() {
        return "==";
    }
}
