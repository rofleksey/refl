package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class MinusLexem implements Lexem {
    public static final MinusLexem INSTANCE = new MinusLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.MINUS;
    }
}
