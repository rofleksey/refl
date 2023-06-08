package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class ColonLexem implements Lexem {
    public static final ColonLexem INSTANCE = new ColonLexem();

    @Override
    public LexemType type() {
        return LexemType.COLON;
    }
}
