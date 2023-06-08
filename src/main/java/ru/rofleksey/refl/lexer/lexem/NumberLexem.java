package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

import java.util.Objects;

public final class NumberLexem implements Lexem {
    private final Integer value;

    public NumberLexem(Integer value) {
        this.value = value;
    }

    public Integer value() {
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
        return Objects.equals(value, that.value);
    }

    @Override
    public int hashCode() {
        return Objects.hash(value);
    }
}
