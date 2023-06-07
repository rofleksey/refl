package ru.rofleksey.refl.lexer.lexem;

import org.jetbrains.annotations.NotNull;
import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class WhileLexem implements Lexem {
    public static final WhileLexem INSTANCE = new WhileLexem();

    @Override
    public @NotNull LexemType type() {
        return LexemType.WHILE;
    }
}
