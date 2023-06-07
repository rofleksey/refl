package ru.rofleksey.refl.lexer;

import org.jetbrains.annotations.NotNull;

public interface Lexem {
    @NotNull
    LexemType type();
}
