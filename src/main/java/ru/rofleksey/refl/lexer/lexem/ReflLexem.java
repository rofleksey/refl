package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class ReflLexem implements Lexem {
    public static final ReflLexem INSTANCE = new ReflLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.REFL;
    }
}
