package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class EqLexem implements Lexem {
    public static final EqLexem INSTANCE = new EqLexem();

    @Override
    public  LexemType type() {
        return LexemType.EQ;
    }
}
