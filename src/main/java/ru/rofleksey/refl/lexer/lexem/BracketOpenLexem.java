package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class BracketOpenLexem implements Lexem {
    public static final BracketOpenLexem INSTANCE = new BracketOpenLexem();

    @Override
    public  LexemType type() {
        return LexemType.BRACKET_OPEN;
    }
}
