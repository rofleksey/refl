package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

import java.util.Objects;

public class VarLexem implements Lexem {
    private final String name;

    public VarLexem(String name) {
        this.name = name;
    }

    public String name() {
        return name;
    }

    @Override
    public LexemType type() {
        return LexemType.VARIABLE;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        VarLexem varLexem = (VarLexem) o;
        return Objects.equals(name, varLexem.name);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name);
    }
}
