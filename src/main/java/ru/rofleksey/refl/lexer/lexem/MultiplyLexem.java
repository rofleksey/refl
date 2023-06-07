package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class MultiplyLexem implements Lexem {
    public static final MultiplyLexem INSTANCE = new MultiplyLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.MULTIPLY;
    }
}
