package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class DivideLexem implements Lexem {
    public static final DivideLexem INSTANCE = new DivideLexem();

    @Override
    public  LexemType type() {
        return LexemType.DIVIDE;
    }
}
