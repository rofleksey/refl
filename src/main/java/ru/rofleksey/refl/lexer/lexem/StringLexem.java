package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public record StringLexem(String text) implements Lexem {

    @Override
    public @NotNull LexemType type() {
        return LexemType.STRING;
    }
}
