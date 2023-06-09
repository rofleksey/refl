package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

import java.util.Objects;

public final class NumberLexem implements Lexem {
    private final double value;

    public NumberLexem(double value) {
        this.value = value;
    }

    public double value() {
        return value;
    }

    @Override
    public LexemType type() {
        return LexemType.NUMBER;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        NumberLexem that = (NumberLexem) o;
        return Double.compare(that.value, value) == 0;
    }

    @Override
    public int hashCode() {
        return Objects.hash(value);
    }
}
