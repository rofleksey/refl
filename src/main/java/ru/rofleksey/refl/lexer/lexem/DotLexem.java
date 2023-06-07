package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class DotLexem implements Lexem {
    public static final DotLexem INSTANCE = new DotLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.DOT;
    }
}
