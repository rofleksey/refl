package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class IfLexem implements Lexem {
    public static final IfLexem INSTANCE = new IfLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.IF;
    }
}
