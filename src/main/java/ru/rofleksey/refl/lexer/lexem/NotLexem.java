package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class NotLexem implements Lexem {
    public static final NotLexem INSTANCE = new NotLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.NOT;
    }
}
