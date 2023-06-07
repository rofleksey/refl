package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class ReflLexem implements Lexem {
    public static final ReflLexem INSTANCE = new ReflLexem();

    @Override
    public  LexemType type() {
        return LexemType.REFL;
    }
}
