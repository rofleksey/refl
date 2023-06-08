package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class EofLexem implements Lexem {
    public static final EofLexem INSTANCE = new EofLexem();

    @Override
    public LexemType type() {
        return LexemType.EOF;
    }
}
