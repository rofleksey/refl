package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class OrLexem implements Lexem {
    public static final OrLexem INSTANCE = new OrLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.OR;
    }
}
