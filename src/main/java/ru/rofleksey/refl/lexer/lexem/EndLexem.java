package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class EndLexem implements Lexem {
    public static final EndLexem INSTANCE = new EndLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.END;
    }
}
