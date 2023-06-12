package ru.rofleksey.refl.lang.operator.compare;

import ru.rofleksey.refl.lang.operator.CompareOperator;

public class LessOperator implements CompareOperator {
    public static final LessOperator INSTANCE = new LessOperator();

    @Override
    public boolean test(Double compareResult) {
        return compareResult < 0;
    }

    @Override
    public String toString() {
        return "<";
    }
}
