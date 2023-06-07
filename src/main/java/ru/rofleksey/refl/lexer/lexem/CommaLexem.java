package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class CommaLexem implements Lexem {
    public static final CommaLexem INSTANCE = new CommaLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.COMMA;
    }
}
