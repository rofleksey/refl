package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class AndLexem implements Lexem {
    public static final AndLexem INSTANCE = new AndLexem();

    @Override
    public LexemType type() {
        return LexemType.AND;
    }
}
