package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class SemicolonLexem implements Lexem {
    public static final SemicolonLexem INSTANCE = new SemicolonLexem();

    @Override
    public  LexemType type() {
        return LexemType.SEMICOLON;
    }
}
