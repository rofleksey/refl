package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class MinusLexem implements Lexem {
    public static final MinusLexem INSTANCE = new MinusLexem();

    @Override
    public  LexemType type() {
        return LexemType.MINUS;
    }
}
