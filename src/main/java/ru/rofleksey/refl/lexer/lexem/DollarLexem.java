package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class DollarLexem implements Lexem {
    public static final DollarLexem INSTANCE = new DollarLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.DOLLAR;
    }
}
