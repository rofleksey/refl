package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class MultiplyLexem implements Lexem {
    public static final MultiplyLexem INSTANCE = new MultiplyLexem();

    @Override
    public LexemType type() {
        return LexemType.MULTIPLY;
    }
}
