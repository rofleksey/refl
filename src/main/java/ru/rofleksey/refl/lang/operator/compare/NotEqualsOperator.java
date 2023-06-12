package ru.rofleksey.refl.lang.operator.compare;

import ru.rofleksey.refl.lang.operator.CompareOperator;

public class NotEqualsOperator implements CompareOperator {
    public static final NotEqualsOperator INSTANCE = new NotEqualsOperator();

    @Override
    public boolean test(Double compareResult) {
        return compareResult != 0;
    }

    @Override
    public String toString() {
        return "!=";
    }
}
