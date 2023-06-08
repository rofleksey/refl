package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class PlusLexem implements Lexem {
    public static final PlusLexem INSTANCE = new PlusLexem();

    @Override
    public LexemType type() {
        return LexemType.PLUS;
    }
}
