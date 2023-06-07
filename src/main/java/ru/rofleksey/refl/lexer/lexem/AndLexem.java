package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class AndLexem implements Lexem {
    public static final AndLexem INSTANCE = new AndLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.AND;
    }
}
