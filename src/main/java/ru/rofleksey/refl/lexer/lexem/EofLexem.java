package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class EofLexem implements Lexem {
    public static final EofLexem INSTANCE = new EofLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.EOF;
    }
}
