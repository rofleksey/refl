package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class PlusLexem implements Lexem {
    public static final PlusLexem INSTANCE = new PlusLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.PLUS;
    }
}
