package ru.rofleksey.refl.lang.operator.compare;

import ru.rofleksey.refl.lang.operator.CompareOperator;

public class GteOperator implements CompareOperator {
    public static final GteOperator INSTANCE = new GteOperator();

    @Override
    public boolean test(Double compareResult) {
        return compareResult >= 0;
    }

    @Override
    public String toString() {
        return ">=";
    }
}
