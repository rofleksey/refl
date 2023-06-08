package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class DivideLexem implements Lexem {
    public static final DivideLexem INSTANCE = new DivideLexem();

    @Override
    public LexemType type() {
        return LexemType.DIVIDE;
    }
}
