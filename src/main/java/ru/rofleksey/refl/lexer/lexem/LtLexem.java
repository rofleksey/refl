package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class LtLexem implements Lexem {
    public static final LtLexem INSTANCE = new LtLexem();

    @Override
    public LexemType type() {
        return LexemType.LT;
    }
}
