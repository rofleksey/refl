package ru.rofleksey.refl.lexer.lexem;


import ru.rofleksey.refl.lexer.Lexem;
import ru.rofleksey.refl.lexer.LexemType;

public class CommaLexem implements Lexem {
    public static final CommaLexem INSTANCE = new CommaLexem();

    @Override
    public  LexemType type() {
        return LexemType.COMMA;
    }
}
