package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class OrLexem implements Lexem {
    public static final OrLexem INSTANCE = new OrLexem();

    @Override
    public LexemType type() {
        return LexemType.OR;
    }
}
