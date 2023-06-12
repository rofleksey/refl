package ru.rofleksey.refl.lang.operator.compare;

import ru.rofleksey.refl.lang.operator.CompareOperator;

public class LteOperator implements CompareOperator {
    public static final LteOperator INSTANCE = new LteOperator();

    @Override
    public boolean test(Double compareResult) {
        return compareResult <= 0;
    }

    @Override
    public String toString() {
        return "<=";
    }
}
