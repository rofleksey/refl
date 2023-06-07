package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class LtLexem implements Lexem {
    public static final LtLexem INSTANCE = new LtLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.LT;
    }
}
