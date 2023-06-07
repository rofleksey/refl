package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

import java.util.Objects;

public class StringLexem implements Lexem {
    private final String text;

    public StringLexem(String text) {
        this.text = text;
    }

    public String text() {
        return text;
    }

    @Override
    public LexemType type() {
        return LexemType.STRING;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        StringLexem that = (StringLexem) o;
        return Objects.equals(text, that.text);
    }

    @Override
    public int hashCode() {
        return Objects.hash(text);
    }
}
