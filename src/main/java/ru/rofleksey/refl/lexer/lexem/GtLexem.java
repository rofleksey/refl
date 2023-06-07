package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class GtLexem implements Lexem {
    public static final GtLexem INSTANCE = new GtLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.GT;
    }
}
