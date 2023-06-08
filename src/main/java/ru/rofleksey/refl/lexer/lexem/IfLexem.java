package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public final class IfLexem implements Lexem {
    public static final IfLexem INSTANCE = new IfLexem();

    @Override
    public LexemType type() {
        return LexemType.IF;
    }
}
